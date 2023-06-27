/*
Copyright 2019 The Kubernetes Authors.

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

package cacher

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
)

var _ runtime.CacheableObject = &cachingObject{}

// metaRuntimeInterface implements runtime.Object and
// metav1.Object interfaces.
// metaRuntimeInterface 实现了 runtime.Object 和 metav1.Object 接口
type metaRuntimeInterface interface {
	runtime.Object
	metav1.Object
}

// serializationResult captures a result of serialization.
// serializationResult 捕获了序列化的结果
type serializationResult struct {
	// once should be used to ensure serialization is computed once.
	once sync.Once

	// raw is serialized object.
	raw []byte
	// err is error from serialization.
	err error
}

// serializationsCache is a type for caching serialization results.
// serializationsCache 是缓存序列化结果的类型
type serializationsCache map[runtime.Identifier]*serializationResult

// cachingObject is an object that is able to cache its serializations
// so that each of those is computed exactly once.
//
// cachingObject implements the metav1.Object interface (accessors for
// all metadata fields).
// cachingObject 实现了 metav1.Object 接口（所有元数据字段的访问器）
// cachingObject 实现了 runtime.CacheableObject 接口
type cachingObject struct {
	lock sync.RWMutex

	// deepCopied defines whether the object below has already been
	// deep copied. The operation is performed lazily on the first
	// setXxx operation.
	//
	// The lazy deep-copy make is useful, as effectively the only
	// case when we are setting some fields are ResourceVersion for
	// DELETE events, so in all other cases we can effectively avoid
	// performing any deep copies.
	// deepCopied 定义了对象下面是否已经被深度复制过了
	// 操作是在第一次 setXxx 操作时延迟执行的
	// 延迟深度复制对于我们来说是有用的，因为实际上我们设置字段的唯一情况是 DELETE 事件的 ResourceVersion
	// 所以在所有其他情况下，我们可以有效地避免执行任何深度复制
	deepCopied bool

	// Object for which serializations are cached.
	// 用于缓存序列化的对象
	object metaRuntimeInterface

	// serializations is a cache containing object`s serializations.
	// The value stored in atomic.Value is of type serializationsCache.
	// The atomic.Value type is used to allow fast-path.
	// serializations 是一个包含对象序列化的缓存
	// 存储在 atomic.Value 中的值的类型是 serializationsCache
	// 使用 atomic.Value 类型可以实现快速路径
	serializations atomic.Value
}

// newCachingObject performs a deep copy of the given object and wraps it
// into a cachingObject.
// An error is returned if it's not possible to cast the object to
// metav1.Object type.
// newCachingObject 对给定对象执行深度复制，并将其包装到 cachingObject 中
// 如果无法将对象转换为 metav1.Object 类型，则返回错误
func newCachingObject(object runtime.Object) (*cachingObject, error) {
	if obj, ok := object.(metaRuntimeInterface); ok {
		result := &cachingObject{
			object:     obj,
			deepCopied: false,
		}
		result.serializations.Store(make(serializationsCache))
		return result, nil
	}
	return nil, fmt.Errorf("can't cast object to metav1.Object: %#v", object)
}

func (o *cachingObject) getSerializationResult(id runtime.Identifier) *serializationResult {
	// Fast-path for getting from cache.
	serializations := o.serializations.Load().(serializationsCache)
	if result, exists := serializations[id]; exists {
		return result
	}

	// Slow-path (that may require insert).
	o.lock.Lock()
	defer o.lock.Unlock()

	serializations = o.serializations.Load().(serializationsCache)
	// Check if in the meantime it wasn't inserted.
	if result, exists := serializations[id]; exists {
		return result
	}

	// Insert an entry for <id>. This requires copy of existing map.
	newSerializations := make(serializationsCache)
	for k, v := range serializations {
		newSerializations[k] = v
	}
	result := &serializationResult{}
	newSerializations[id] = result
	o.serializations.Store(newSerializations)
	return result
}

// CacheEncode implements runtime.CacheableObject interface.
// It serializes the object and writes the result to given io.Writer trying
// to first use the already cached result and falls back to a given encode
// function in case of cache miss.
// It assumes that for a given identifier, the encode function always encodes
// each input object into the same output format.
// CacheEncode 实现了 runtime.CacheableObject 接口
// 它序列化对象并将结果写入给定的 io.Writer，尝试首先使用已缓存的结果，如果缓存未命中，则回退到给定的 encode 函数
// 它假定对于给定的标识符，encode 函数总是将每个输入对象编码为相同的输出格式
func (o *cachingObject) CacheEncode(id runtime.Identifier, encode func(runtime.Object, io.Writer) error, w io.Writer) error {
	result := o.getSerializationResult(id)
	result.once.Do(func() {
		buffer := bytes.NewBuffer(nil)
		// TODO(wojtek-t): This is currently making a copy to avoid races
		//   in cases where encoding is making subtle object modifications,
		//   e.g. #82497
		//   Figure out if we can somehow avoid this under some conditions.
		result.err = encode(o.GetObject(), buffer)
		result.raw = buffer.Bytes()
	})
	// Once invoked, fields of serialization will not change.
	if result.err != nil {
		return result.err
	}
	_, err := w.Write(result.raw)
	return err
}

// GetObject implements runtime.CacheableObject interface.
// It returns deep-copy of the wrapped object to return ownership of it
// to the called according to the contract of the interface.
// GetObject 实现了 runtime.CacheableObject 接口
// 它返回包装对象的深度副本，以便根据接口的约定将其所有权返回给调用者
func (o *cachingObject) GetObject() runtime.Object {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.DeepCopyObject().(metaRuntimeInterface)
}

// GetObjectKind implements runtime.Object interface.
// GetObjectKind 实现了 runtime.Object 接口
func (o *cachingObject) GetObjectKind() schema.ObjectKind {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetObjectKind()
}

// DeepCopyObject implements runtime.Object interface.
// DeepCopyObject 实现了 runtime.Object 接口
func (o *cachingObject) DeepCopyObject() runtime.Object {
	// DeepCopyObject on cachingObject is not expected to be called anywhere.
	// However, to be on the safe-side, we implement it, though given the
	// cache is only an optimization we ignore copying it.
	result := &cachingObject{
		deepCopied: true,
	}
	result.serializations.Store(make(serializationsCache))

	o.lock.RLock()
	defer o.lock.RUnlock()
	result.object = o.object.DeepCopyObject().(metaRuntimeInterface)
	return result
}

var (
	invalidationCacheTimestampLock sync.Mutex
	invalidationCacheTimestamp     time.Time
)

// shouldLogCacheInvalidation allows for logging cache-invalidation
// at most once per second (to avoid spamming logs in case of issues).
// shouldLogCacheInvalidation 允许每秒最多记录一次缓存无效（以避免由于问题而导致日志记录过多）
func shouldLogCacheInvalidation(now time.Time) bool {
	invalidationCacheTimestampLock.Lock()
	defer invalidationCacheTimestampLock.Unlock()
	if invalidationCacheTimestamp.Add(time.Second).Before(now) {
		invalidationCacheTimestamp = now
		return true
	}
	return false
}

func (o *cachingObject) invalidateCacheLocked() {
	if cache, ok := o.serializations.Load().(serializationsCache); ok && len(cache) == 0 {
		return
	}
	// We don't expect cache invalidation to happen - so we want
	// to log the stacktrace to allow debugging if that will happen.
	// OTOH, we don't want to spam logs with it.
	// So we try to log it at most once per second.
	if shouldLogCacheInvalidation(time.Now()) {
		klog.Warningf("Unexpected cache invalidation for %#v\n%s", o.object, string(debug.Stack()))
	}
	o.serializations.Store(make(serializationsCache))
}

// The following functions implement metav1.Object interface:
//   - getters simply delegate for the underlying object
//   - setters check if operations isn't noop and if so,
//     invalidate the cache and delegate for the underlying object
//
// 下面的函数实现了 metav1.Object 接口：
// - getters 只是委托给底层对象
// - setters 检查操作是否不是 noop，如果是，则使缓存无效并委托给底层对象
func (o *cachingObject) conditionalSet(isNoop func() bool, set func()) {
	if fastPath := func() bool {
		o.lock.RLock()
		defer o.lock.RUnlock()
		return isNoop()
	}(); fastPath {
		return
	}
	o.lock.Lock()
	defer o.lock.Unlock()
	if isNoop() {
		return
	}
	if !o.deepCopied {
		o.object = o.object.DeepCopyObject().(metaRuntimeInterface)
		o.deepCopied = true
	}
	o.invalidateCacheLocked()
	set()
}

func (o *cachingObject) GetNamespace() string {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetNamespace()
}
func (o *cachingObject) SetNamespace(namespace string) {
	o.conditionalSet(
		func() bool { return o.object.GetNamespace() == namespace },
		func() { o.object.SetNamespace(namespace) },
	)
}
func (o *cachingObject) GetName() string {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetName()
}
func (o *cachingObject) SetName(name string) {
	o.conditionalSet(
		func() bool { return o.object.GetName() == name },
		func() { o.object.SetName(name) },
	)
}
func (o *cachingObject) GetGenerateName() string {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetGenerateName()
}
func (o *cachingObject) SetGenerateName(name string) {
	o.conditionalSet(
		func() bool { return o.object.GetGenerateName() == name },
		func() { o.object.SetGenerateName(name) },
	)
}
func (o *cachingObject) GetUID() types.UID {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetUID()
}
func (o *cachingObject) SetUID(uid types.UID) {
	o.conditionalSet(
		func() bool { return o.object.GetUID() == uid },
		func() { o.object.SetUID(uid) },
	)
}
func (o *cachingObject) GetResourceVersion() string {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetResourceVersion()
}
func (o *cachingObject) SetResourceVersion(version string) {
	o.conditionalSet(
		func() bool { return o.object.GetResourceVersion() == version },
		func() { o.object.SetResourceVersion(version) },
	)
}
func (o *cachingObject) GetGeneration() int64 {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetGeneration()
}
func (o *cachingObject) SetGeneration(generation int64) {
	o.conditionalSet(
		func() bool { return o.object.GetGeneration() == generation },
		func() { o.object.SetGeneration(generation) },
	)
}
func (o *cachingObject) GetSelfLink() string {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetSelfLink()
}
func (o *cachingObject) SetSelfLink(selfLink string) {
	o.conditionalSet(
		func() bool { return o.object.GetSelfLink() == selfLink },
		func() { o.object.SetSelfLink(selfLink) },
	)
}
func (o *cachingObject) GetCreationTimestamp() metav1.Time {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetCreationTimestamp()
}
func (o *cachingObject) SetCreationTimestamp(timestamp metav1.Time) {
	o.conditionalSet(
		func() bool { return o.object.GetCreationTimestamp() == timestamp },
		func() { o.object.SetCreationTimestamp(timestamp) },
	)
}
func (o *cachingObject) GetDeletionTimestamp() *metav1.Time {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetDeletionTimestamp()
}
func (o *cachingObject) SetDeletionTimestamp(timestamp *metav1.Time) {
	o.conditionalSet(
		func() bool { return o.object.GetDeletionTimestamp() == timestamp },
		func() { o.object.SetDeletionTimestamp(timestamp) },
	)
}
func (o *cachingObject) GetDeletionGracePeriodSeconds() *int64 {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetDeletionGracePeriodSeconds()
}
func (o *cachingObject) SetDeletionGracePeriodSeconds(gracePeriodSeconds *int64) {
	o.conditionalSet(
		func() bool { return o.object.GetDeletionGracePeriodSeconds() == gracePeriodSeconds },
		func() { o.object.SetDeletionGracePeriodSeconds(gracePeriodSeconds) },
	)
}
func (o *cachingObject) GetLabels() map[string]string {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetLabels()
}
func (o *cachingObject) SetLabels(labels map[string]string) {
	o.conditionalSet(
		func() bool { return reflect.DeepEqual(o.object.GetLabels(), labels) },
		func() { o.object.SetLabels(labels) },
	)
}
func (o *cachingObject) GetAnnotations() map[string]string {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetAnnotations()
}
func (o *cachingObject) SetAnnotations(annotations map[string]string) {
	o.conditionalSet(
		func() bool { return reflect.DeepEqual(o.object.GetAnnotations(), annotations) },
		func() { o.object.SetAnnotations(annotations) },
	)
}
func (o *cachingObject) GetFinalizers() []string {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetFinalizers()
}
func (o *cachingObject) SetFinalizers(finalizers []string) {
	o.conditionalSet(
		func() bool { return reflect.DeepEqual(o.object.GetFinalizers(), finalizers) },
		func() { o.object.SetFinalizers(finalizers) },
	)
}
func (o *cachingObject) GetOwnerReferences() []metav1.OwnerReference {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetOwnerReferences()
}
func (o *cachingObject) SetOwnerReferences(references []metav1.OwnerReference) {
	o.conditionalSet(
		func() bool { return reflect.DeepEqual(o.object.GetOwnerReferences(), references) },
		func() { o.object.SetOwnerReferences(references) },
	)
}
func (o *cachingObject) GetManagedFields() []metav1.ManagedFieldsEntry {
	o.lock.RLock()
	defer o.lock.RUnlock()
	return o.object.GetManagedFields()
}
func (o *cachingObject) SetManagedFields(managedFields []metav1.ManagedFieldsEntry) {
	o.conditionalSet(
		func() bool { return reflect.DeepEqual(o.object.GetManagedFields(), managedFields) },
		func() { o.object.SetManagedFields(managedFields) },
	)
}
