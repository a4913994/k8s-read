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

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Indexer extends Store with multiple indices and restricts each
// accumulator to simply hold the current object (and be empty after
// Delete).
//
// There are three kinds of strings here:
//  1. a storage key, as defined in the Store interface,
//  2. a name of an index, and
//  3. an "indexed value", which is produced by an IndexFunc and
//     can be a field value or any other string computed from the object.
//
// Indexer 扩展了 Store，具有多个索引，并且将每个累加器限制为仅保存当前对象（在 Delete 后为空）。
// 这里有三种字符串：
// 1. 存储键，如 Store 接口中所定义的那样，
// 2. 索引的名称，以及
// 3. “索引值”，由 IndexFunc 生成，可以是字段值或从对象计算出的任何其他字符串。
type Indexer interface {
	Store
	// Index returns the stored objects whose set of indexed values
	// intersects the set of indexed values of the given object, for
	// the named index
	// Index 返回存储的对象，其索引值的集合与给定对象的索引值的集合相交，对于命名的索引
	Index(indexName string, obj interface{}) ([]interface{}, error)
	// IndexKeys returns the storage keys of the stored objects whose
	// set of indexed values for the named index includes the given
	// indexed value
	// IndexKeys 返回存储对象的存储键，其命名索引的索引值的集合包含给定的索引值
	IndexKeys(indexName, indexedValue string) ([]string, error)
	// ListIndexFuncValues returns all the indexed values of the given index
	// ListIndexFuncValues 返回给定索引的所有索引值
	ListIndexFuncValues(indexName string) []string
	// ByIndex returns the stored objects whose set of indexed values
	// for the named index includes the given indexed value
	// ByIndex 返回存储的对象，其命名索引的索引值的集合包含给定的索引值
	ByIndex(indexName, indexedValue string) ([]interface{}, error)
	// GetIndexers return the indexers
	// GetIndexers 返回索引器
	GetIndexers() Indexers

	// AddIndexers adds more indexers to this store.  If you call this after you already have data
	// in the store, the results are undefined.
	// AddIndexers 向此存储添加更多索引器。 如果在存储中已经有数据之后调用此方法，则结果是未定义的。
	AddIndexers(newIndexers Indexers) error
}

// IndexFunc knows how to compute the set of indexed values for an object.
// IndexFunc 知道如何计算对象的索引值集。
type IndexFunc func(obj interface{}) ([]string, error)

// IndexFuncToKeyFuncAdapter adapts an indexFunc to a keyFunc.  This is only useful if your index function returns
// unique values for every object.  This conversion can create errors when more than one key is found.  You
// should prefer to make proper key and index functions.
// IndexFuncToKeyFuncAdapter 将 indexFunc 适配为 keyFunc。 仅当您的索引函数为每个对象返回唯一值时，此转换才有用。
// 当找到多个键时，此转换可能会创建错误。 您应该首选正确的键和索引函数。
func IndexFuncToKeyFuncAdapter(indexFunc IndexFunc) KeyFunc {
	return func(obj interface{}) (string, error) {
		indexKeys, err := indexFunc(obj)
		if err != nil {
			return "", err
		}
		if len(indexKeys) > 1 {
			return "", fmt.Errorf("too many keys: %v", indexKeys)
		}
		if len(indexKeys) == 0 {
			return "", fmt.Errorf("unexpected empty indexKeys")
		}
		return indexKeys[0], nil
	}
}

const (
	// NamespaceIndex is the lookup name for the most common index function, which is to index by the namespace field.
	// NamespaceIndex 是最常见的索引函数的查找名称，即通过命名空间字段进行索引。
	NamespaceIndex string = "namespace"
)

// MetaNamespaceIndexFunc is a default index function that indexes based on an object's namespace
// MetaNamespaceIndexFunc 是一个默认的索引函数，它基于对象的命名空间进行索引
func MetaNamespaceIndexFunc(obj interface{}) ([]string, error) {
	meta, err := meta.Accessor(obj)
	if err != nil {
		return []string{""}, fmt.Errorf("object has no meta: %v", err)
	}
	return []string{meta.GetNamespace()}, nil
}

// Index maps the indexed value to a set of keys in the store that match on that value
// Index 将索引值映射到存储中与该值匹配的键集
type Index map[string]sets.String

// Indexers maps a name to an IndexFunc
// Indexers 将名称映射到 IndexFunc
type Indexers map[string]IndexFunc

// Indices maps a name to an Index
// Indices 将名称映射到 Index
type Indices map[string]Index
