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

package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NamespaceControllerConfiguration contains elements describing NamespaceController.
type NamespaceControllerConfiguration struct {
	// NamespaceSyncPeriod is the period for syncing namespace life-cycle
	// updates.
	// NamespaceSyncPeriod是同步命名空间生命周期更新的周期。
	NamespaceSyncPeriod metav1.Duration
	// ConcurrentNamespaceSyncs is the number of namespace objects that are
	// allowed to sync concurrently.
	// ConcurrentNamespaceSyncs是允许并发同步的命名空间对象的数量。
	ConcurrentNamespaceSyncs int32
}
