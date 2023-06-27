/*
Copyright 2014 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cache

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/api/meta"
)

// Store is a generic object storage and processing interface.  A
// Store holds a map from string keys to accumulators, and has
// operations to add, update, and delete a given object to/from the
// accumulator currently associated with a given key.  A Store also
// knows how to extract the key from a given object, so many operations
// are given only the object.
//
// In the simplest Store implementations each accumulator is simply
// the last given object, or empty after Delete, and thus the Store's
// behavior is simple storage.
//
// Reflector knows how to watch a server and update a Store.  This
// package provides a variety of implementations of Store.
// Store是一个通用的对象存储和处理接口。Store保存一个从字符串键到累加器的映射，并且有操作来添加、更新和删除给定对象到/从给定键关联的累加器。
// Store还知道如何从给定对象中提取键，因此许多操作只给出对象。
// 在最简单的Store实现中，每个累加器都是最后给出的对象，或者在Delete后为空，因此Store的行为是简单的存储。
// Reflector知道如何监视服务器并更新Store。该包提供了Store的各种实现。
type Store interface {

	// Add adds the given object to the accumulator associated with the given object's key
	// Add 增加给定对象到与给定对象的键关联的累加器
	Add(obj interface{}) error

	// Update updates the given object in the accumulator associated with the given object's key
	// Update 更新与给定对象的键关联的累加器中的给定对象
	Update(obj interface{}) error

	// Delete deletes the given object from the accumulator associated with the given object's key
	// Delete 从与给定对象的键关联的累加器中删除给定对象
	Delete(obj interface{}) error

	// List returns a list of all the currently non-empty accumulators
	// List 返回当前所有非空累加器的列表
	List() []interface{}

	// ListKeys returns a list of all the keys currently associated with non-empty accumulators
	// ListKeys 返回当前与非空累加器关联的所有键的列表
	ListKeys() []string

	// Get returns the accumulator associated with the given object's key
	// Get 返回与给定对象的键关联的累加器
	Get(obj interface{}) (item interface{}, exists bool, err error)

	// GetByKey returns the accumulator associated with the given key
	// GetByKey 返回与给定键关联的累加器
	GetByKey(key string) (item interface{}, exists bool, err error)

	// Replace will delete the contents of the store, using instead the
	// given list. Store takes ownership of the list, you should not reference
	// it after calling this function.
	// Replace 将删除存储的内容，而不是使用给定的列表。Store接管列表，您不应该在调用此函数后引用它。
	Replace([]interface{}, string) error

	// Resync is meaningless in the terms appearing here but has
	// meaning in some implementations that have non-trivial
	// additional behavior (e.g., DeltaFIFO).
	// Resync 在这里出现的术语中没有意义，但在一些实现中具有非平凡的附加行为（例如DeltaFIFO）时具有意义。
	Resync() error
}

// KeyFunc knows how to make a key from an object. Implementations should be deterministic.
// KeyFunc知道如何从对象中制作键。实现应该是确定性的。
type KeyFunc func(obj interface{}) (string, error)

// KeyError will be returned any time a KeyFunc gives an error; it includes the object
// at fault.
// KeyError将在任何时候返回KeyFunc给出错误时;它包括有问题的对象。
type KeyError struct {
	Obj interface{}
	Err error
}

// Error gives a human-readable description of the error.
// Error 提供了错误的人类可读描述。
func (k KeyError) Error() string {
	return fmt.Sprintf("couldn't create key for object %+v: %v", k.Obj, k.Err)
}

// Unwrap implements errors.Unwrap
// Unwrap 实现了 errors.Unwrap
func (k KeyError) Unwrap() error {
	return k.Err
}

// ExplicitKey can be passed to MetaNamespaceKeyFunc if you have the key for
// the object but not the object itself.
// ExplicitKey 可以传递给 MetaNamespaceKeyFunc，如果您有对象的键但没有对象本身。
type ExplicitKey string

// MetaNamespaceKeyFunc is a convenient default KeyFunc which knows how to make
// keys for API objects which implement meta.Interface.
// The key uses the format <namespace>/<name> unless <namespace> is empty, then
// it's just <name>.
//
// MetaNamespaceKeyFunc 是一个方便的默认 KeyFunc，它知道如何为实现 meta.Interface 的 API 对象制作键。
// 该键使用格式 <namespace>/<name>，除非 <namespace> 为空，然后它只是 <name>。
// TODO: replace key-as-string with a key-as-struct so that this
// packing/unpacking won't be necessary.
func MetaNamespaceKeyFunc(obj interface{}) (string, error) {
	if key, ok := obj.(ExplicitKey); ok {
		return string(key), nil
	}
	meta, err := meta.Accessor(obj)
	if err != nil {
		return "", fmt.Errorf("object has no meta: %v", err)
	}
	if len(meta.GetNamespace()) > 0 {
		return meta.GetNamespace() + "/" + meta.GetName(), nil
	}
	return meta.GetName(), nil
}

// SplitMetaNamespaceKey returns the namespace and name that
// MetaNamespaceKeyFunc encoded into key.
//
// SplitMetaNamespaceKey 返回 MetaNamespaceKeyFunc 编码到 key 中的命名空间和名称。
// TODO: replace key-as-string with a key-as-struct so that this
// packing/unpacking won't be necessary.
func SplitMetaNamespaceKey(key string) (namespace, name string, err error) {
	parts := strings.Split(key, "/")
	switch len(parts) {
	case 1:
		// name only, no namespace
		return "", parts[0], nil
	case 2:
		// namespace and name
		return parts[0], parts[1], nil
	}

	return "", "", fmt.Errorf("unexpected key format: %q", key)
}

// `*cache` implements Indexer in terms of a ThreadSafeStore and an
// associated KeyFunc.
// `*cache` 在 ThreadSafeStore 和相关的 KeyFunc 的条件下实现了 Indexer。
type cache struct {
	// cacheStorage bears the burden of thread safety for the cache
	// cacheStorage 承担缓存的线程安全性的负担
	cacheStorage ThreadSafeStore
	// keyFunc is used to make the key for objects stored in and retrieved from items, and
	// should be deterministic.
	// keyFunc 用于制作存储在 items 中并从 items 中检索的对象的键，并且应该是确定性的。
	keyFunc KeyFunc
}

var _ Store = &cache{}

// Add inserts an item into the cache.
// Add 将项目插入缓存。
func (c *cache) Add(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return KeyError{obj, err}
	}
	c.cacheStorage.Add(key, obj)
	return nil
}

// Update sets an item in the cache to its updated state.
func (c *cache) Update(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return KeyError{obj, err}
	}
	c.cacheStorage.Update(key, obj)
	return nil
}

// Delete removes an item from the cache.
func (c *cache) Delete(obj interface{}) error {
	key, err := c.keyFunc(obj)
	if err != nil {
		return KeyError{obj, err}
	}
	c.cacheStorage.Delete(key)
	return nil
}

// List returns a list of all the items.
// List is completely threadsafe as long as you treat all items as immutable.
func (c *cache) List() []interface{} {
	return c.cacheStorage.List()
}

// ListKeys returns a list of all the keys of the objects currently
// in the cache.
func (c *cache) ListKeys() []string {
	return c.cacheStorage.ListKeys()
}

// GetIndexers returns the indexers of cache
func (c *cache) GetIndexers() Indexers {
	return c.cacheStorage.GetIndexers()
}

// Index returns a list of items that match on the index function
// Index is thread-safe so long as you treat all items as immutable
func (c *cache) Index(indexName string, obj interface{}) ([]interface{}, error) {
	return c.cacheStorage.Index(indexName, obj)
}

// IndexKeys returns the storage keys of the stored objects whose set of
// indexed values for the named index includes the given indexed value.
// The returned keys are suitable to pass to GetByKey().
func (c *cache) IndexKeys(indexName, indexedValue string) ([]string, error) {
	return c.cacheStorage.IndexKeys(indexName, indexedValue)
}

// ListIndexFuncValues returns the list of generated values of an Index func
func (c *cache) ListIndexFuncValues(indexName string) []string {
	return c.cacheStorage.ListIndexFuncValues(indexName)
}

// ByIndex returns the stored objects whose set of indexed values
// for the named index includes the given indexed value.
func (c *cache) ByIndex(indexName, indexedValue string) ([]interface{}, error) {
	return c.cacheStorage.ByIndex(indexName, indexedValue)
}

func (c *cache) AddIndexers(newIndexers Indexers) error {
	return c.cacheStorage.AddIndexers(newIndexers)
}

// Get returns the requested item, or sets exists=false.
// Get is completely threadsafe as long as you treat all items as immutable.
func (c *cache) Get(obj interface{}) (item interface{}, exists bool, err error) {
	key, err := c.keyFunc(obj)
	if err != nil {
		return nil, false, KeyError{obj, err}
	}
	return c.GetByKey(key)
}

// GetByKey returns the request item, or exists=false.
// GetByKey is completely threadsafe as long as you treat all items as immutable.
func (c *cache) GetByKey(key string) (item interface{}, exists bool, err error) {
	item, exists = c.cacheStorage.Get(key)
	return item, exists, nil
}

// Replace will delete the contents of 'c', using instead the given list.
// 'c' takes ownership of the list, you should not reference the list again
// after calling this function.
func (c *cache) Replace(list []interface{}, resourceVersion string) error {
	items := make(map[string]interface{}, len(list))
	for _, item := range list {
		key, err := c.keyFunc(item)
		if err != nil {
			return KeyError{item, err}
		}
		items[key] = item
	}
	c.cacheStorage.Replace(items, resourceVersion)
	return nil
}

// Resync is meaningless for one of these
// Resync 对这些是没有意义的
func (c *cache) Resync() error {
	return nil
}

// NewStore returns a Store implemented simply with a map and a lock.
// NewStore 返回一个简单的 Store，它使用 map 和锁实现。
func NewStore(keyFunc KeyFunc) Store {
	return &cache{
		cacheStorage: NewThreadSafeStore(Indexers{}, Indices{}),
		keyFunc:      keyFunc,
	}
}

// NewIndexer returns an Indexer implemented simply with a map and a lock.
// NewIndexer 返回一个简单的 Indexer，它使用 map 和锁实现。
func NewIndexer(keyFunc KeyFunc, indexers Indexers) Indexer {
	return &cache{
		cacheStorage: NewThreadSafeStore(indexers, Indices{}),
		keyFunc:      keyFunc,
	}
}
