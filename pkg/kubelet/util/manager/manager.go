/*
Copyright 2018 The Kubernetes Authors.

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

package manager

import (
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Manager is the interface for registering and unregistering
// objects referenced by pods in the underlying cache and
// extracting those from that cache if needed.
// Manager 是用于在底层缓存中注册和注销 Pod 引用的对象的接口，并在需要时从该缓存中提取这些对象。
type Manager interface {
	// GetObject by its namespace and name.
	// 通过命名空间和名称获取对象。
	GetObject(namespace, name string) (runtime.Object, error)

	// WARNING: Register/UnregisterPod functions should be efficient,
	// i.e. should not block on network operations.

	// RegisterPod registers all objects referenced from a given pod.
	//
	// NOTE: All implementations of RegisterPod should be idempotent.
	// 注册 Pod 中引用的所有对象。
	// 注意：RegisterPod 的所有实现都应该是幂等的。
	RegisterPod(pod *v1.Pod)

	// UnregisterPod unregisters objects referenced from a given pod that are not
	// used by any other registered pod.
	//
	// NOTE: All implementations of UnregisterPod should be idempotent.
	// 注销 Pod 中引用的所有未被其他已注册 Pod 使用的对象。
	// 注意：UnregisterPod 的所有实现都应该是幂等的。
	UnregisterPod(pod *v1.Pod)
}

// Store is the interface for a object cache that
// can be used by cacheBasedManager.
// Store 是一个对象缓存的接口，可以由 cacheBasedManager 使用。
type Store interface {
	// AddReference adds a reference to the object to the store.
	// Note that multiple additions to the store has to be allowed
	// in the implementations and effectively treated as refcounted.
	// AddReference 将对对象的引用添加到存储中。 注意，实现中必须允许对存储的多次添加，并有效地将其视为引用计数。
	AddReference(namespace, name string)
	// DeleteReference deletes reference to the object from the store.
	// Note that object should be deleted only when there was a
	// corresponding Delete call for each of Add calls (effectively
	// when refcount was reduced to zero).
	// DeleteReference 从存储中删除对对象的引用。 注意，只有在每个 Add 调用都有相应的 Delete 调用时（即当引用计数减少到零时），才应删除对象。
	DeleteReference(namespace, name string)
	// Get an object from a store.
	// 从存储中获取对象。
	Get(namespace, name string) (runtime.Object, error)
}
