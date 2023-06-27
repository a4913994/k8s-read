/*
Copyright 2015 The Kubernetes Authors.

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

package generic

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	"k8s.io/apiserver/pkg/storage/storagebackend/factory"
	"k8s.io/client-go/tools/cache"
)

// StorageDecorator is a function signature for producing a storage.Interface
// and an associated DestroyFunc from given parameters.
// StorageDecorator 是从给定参数生成 storage.Interface 和相关 DestroyFunc 的函数签名。
type StorageDecorator func(
	config *storagebackend.ConfigForResource,
	resourcePrefix string,
	keyFunc func(obj runtime.Object) (string, error),
	newFunc func() runtime.Object,
	newListFunc func() runtime.Object,
	getAttrsFunc storage.AttrFunc,
	trigger storage.IndexerFuncs,
	indexers *cache.Indexers) (storage.Interface, factory.DestroyFunc, error)

// UndecoratedStorage returns the given a new storage from the given config
// without any decoration.
// UndecoratedStorage 从给定的配置返回给定的新存储，而不进行任何装饰。
func UndecoratedStorage(
	config *storagebackend.ConfigForResource,
	resourcePrefix string,
	keyFunc func(obj runtime.Object) (string, error),
	newFunc func() runtime.Object,
	newListFunc func() runtime.Object,
	getAttrsFunc storage.AttrFunc,
	trigger storage.IndexerFuncs,
	indexers *cache.Indexers) (storage.Interface, factory.DestroyFunc, error) {
	return NewRawStorage(config, newFunc)
}

// NewRawStorage creates the low level kv storage. This is a work-around for current
// two layer of same storage interface.
// NewRawStorage 创建底层 kv 存储。这是当前两层相同存储接口的解决方法。
// TODO: Once cacher is enabled on all registries (event registry is special), we will remove this method.
func NewRawStorage(config *storagebackend.ConfigForResource, newFunc func() runtime.Object) (storage.Interface, factory.DestroyFunc, error) {
	return factory.Create(*config, newFunc)
}
