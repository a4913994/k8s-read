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

// EndpointControllerConfiguration contains elements describing EndpointController.
type EndpointControllerConfiguration struct {
	// concurrentEndpointSyncs is the number of endpoint syncing operations
	// that will be done concurrently. Larger number = faster endpoint updating,
	// but more CPU (and network) load.
	// concurrentEndpointSyncs 是将同时完成的端点同步操作的数量。更大的数字 = 更快的端点更新，但更多的 CPU（和网络）负载。
	ConcurrentEndpointSyncs int32

	// EndpointUpdatesBatchPeriod can be used to batch endpoint updates.
	// All updates of endpoint triggered by pod change will be delayed by up to
	// 'EndpointUpdatesBatchPeriod'. If other pods in the same endpoint change
	// in that period, they will be batched to a single endpoint update.
	// Default 0 value means that each pod update triggers an endpoint update.
	// EndpointUpdatesBatchPeriod 可用于批量端点更新。由 pod 更改触发的所有端点更新将延迟最多“EndpointUpdatesBatchPeriod”。
	// 如果同一端点中的其他 pod 在此期间发生更改，它们将被批处理到单个端点更新。默认 0 值意味着每个 pod 更新都会触发一个端点更新。
	EndpointUpdatesBatchPeriod metav1.Duration
}
