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

package storage

import (
	"fmt"
	"strconv"

	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// APIObjectVersioner implements versioning and extracting etcd node information
// for objects that have an embedded ObjectMeta or ListMeta field.
// APIObjectVersioner 实现了对象的版本控制和提取 etcd 节点信息的功能，这些对象都有一个内嵌的 ObjectMeta 或者 ListMeta 字段。
type APIObjectVersioner struct{}

// UpdateObject implements Versioner
// UpdateObject 实现了 Versioner 接口
func (a APIObjectVersioner) UpdateObject(obj runtime.Object, resourceVersion uint64) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	versionString := ""
	if resourceVersion != 0 {
		versionString = strconv.FormatUint(resourceVersion, 10)
	}
	accessor.SetResourceVersion(versionString)
	return nil
}

// UpdateList implements Versioner
// UpdateList 实现了 Versioner 接口
func (a APIObjectVersioner) UpdateList(obj runtime.Object, resourceVersion uint64, nextKey string, count *int64) error {
	if resourceVersion == 0 {
		return fmt.Errorf("illegal resource version from storage: %d", resourceVersion)
	}
	listAccessor, err := meta.ListAccessor(obj)
	if err != nil || listAccessor == nil {
		return err
	}
	versionString := strconv.FormatUint(resourceVersion, 10)
	listAccessor.SetResourceVersion(versionString)
	listAccessor.SetContinue(nextKey)
	listAccessor.SetRemainingItemCount(count)
	return nil
}

// PrepareObjectForStorage clears resourceVersion and selfLink prior to writing to etcd.
// PrepareObjectForStorage 在写入 etcd 之前清除 resourceVersion 和 selfLink。
func (a APIObjectVersioner) PrepareObjectForStorage(obj runtime.Object) error {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return err
	}
	accessor.SetResourceVersion("")
	accessor.SetSelfLink("")
	return nil
}

// ObjectResourceVersion implements Versioner
// ObjectResourceVersion 实现了 Versioner 接口
func (a APIObjectVersioner) ObjectResourceVersion(obj runtime.Object) (uint64, error) {
	accessor, err := meta.Accessor(obj)
	if err != nil {
		return 0, err
	}
	version := accessor.GetResourceVersion()
	if len(version) == 0 {
		return 0, nil
	}
	return strconv.ParseUint(version, 10, 64)
}

// ParseResourceVersion takes a resource version argument and converts it to
// the etcd version. For watch we should pass to helper.Watch(). Because resourceVersion is
// an opaque value, the default watch behavior for non-zero watch is to watch
// the next value (if you pass "1", you will see updates from "2" onwards).
// ParseResourceVersion 将一个资源版本参数转换为 etcd 版本。对于 watch，我们应该传递给 helper.Watch()。因为 resourceVersion 是一个不透明的值，
// 对于非零 watch，默认的 watch 行为是观察下一个值（如果你传递 "1"，你将从 "2" 开始看到更新）。
func (a APIObjectVersioner) ParseResourceVersion(resourceVersion string) (uint64, error) {
	if resourceVersion == "" || resourceVersion == "0" {
		return 0, nil
	}
	version, err := strconv.ParseUint(resourceVersion, 10, 64)
	if err != nil {
		return 0, NewInvalidError(field.ErrorList{
			// Validation errors are supposed to return version-specific field
			// paths, but this is probably close enough.
			field.Invalid(field.NewPath("resourceVersion"), resourceVersion, err.Error()),
		})
	}
	return version, nil
}

// Versioner implements Versioner
// Versioner 实现了 Versioner 接口
var _ Versioner = APIObjectVersioner{}

// CompareResourceVersion compares etcd resource versions.  Outside this API they are all strings,
// but etcd resource versions are special, they're actually ints, so we can easily compare them.
// CompareResourceVersion 比较 etcd 资源版本。在这个 API 之外，它们都是字符串，但是 etcd 资源版本是特殊的，它们实际上是整数，所以我们可以轻松地比较它们。
func (a APIObjectVersioner) CompareResourceVersion(lhs, rhs runtime.Object) int {
	lhsVersion, err := a.ObjectResourceVersion(lhs)
	if err != nil {
		// coder error
		panic(err)
	}
	rhsVersion, err := a.ObjectResourceVersion(rhs)
	if err != nil {
		// coder error
		panic(err)
	}

	if lhsVersion == rhsVersion {
		return 0
	}
	if lhsVersion < rhsVersion {
		return -1
	}

	return 1
}
