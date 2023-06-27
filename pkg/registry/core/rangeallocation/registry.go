/*
Copyright 2016 The Kubernetes Authors.

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

package rangeallocation

import (
	api "k8s.io/kubernetes/pkg/apis/core"
)

// RangeRegistry is a registry that can retrieve or persist a RangeAllocation object.
// RangeRegistry 是一个注册表，可以检索或持久化 RangeAllocation 对象。
type RangeRegistry interface {
	// Get returns the latest allocation, an empty object if no allocation has been made,
	// or an error if the allocation could not be retrieved.
	// Get 返回最新的分配，如果没有分配，则返回空对象，如果无法检索分配，则返回错误。
	Get() (*api.RangeAllocation, error)
	// CreateOrUpdate should create or update the provide allocation, unless a conflict
	// has occurred since the item was last created.
	// CreateOrUpdate 应该创建或更新提供的分配，除非自上次创建以来发生冲突。
	CreateOrUpdate(*api.RangeAllocation) error
}
