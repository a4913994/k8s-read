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

// EndpointSliceControllerConfiguration contains elements describing
// EndpointSliceController.
type EndpointSliceControllerConfiguration struct {
	// concurrentServiceEndpointSyncs is the number of service endpoint syncing
	// operations that will be done concurrently. Larger number = faster
	// endpoint slice updating, but more CPU (and network) load.
	// concurrentServiceEndpointSyncs 是将同时完成的服务端点同步操作的数量。更大的数字 = 更快的端点切片更新，但更多的 CPU（和网络）负载。
	ConcurrentServiceEndpointSyncs int32

	// maxEndpointsPerSlice is the maximum number of endpoints that will be
	// added to an EndpointSlice. More endpoints per slice will result in fewer
	// and larger endpoint slices, but larger resources.
	// maxEndpointsPerSlice 是将添加到 EndpointSlice 的端点的最大数量。每个切片更多的端点将导致更少和更大的端点切片，但资源更大。
	MaxEndpointsPerSlice int32

	// EndpointUpdatesBatchPeriod can be used to batch endpoint updates.
	// All updates of endpoint triggered by pod change will be delayed by up to
	// 'EndpointUpdatesBatchPeriod'. If other pods in the same endpoint change
	// in that period, they will be batched to a single endpoint update.
	// Default 0 value means that each pod update triggers an endpoint update.
	// EndpointUpdatesBatchPeriod 可用于批量端点更新。由 pod 更改触发的所有端点更新将延迟最多“EndpointUpdatesBatchPeriod”。
	// 如果同一端点中的其他 pod 在此期间发生更改，它们将被批处理到单个端点更新。默认 0 值意味着每个 pod 更新都会触发一个端点更新。
	EndpointUpdatesBatchPeriod metav1.Duration
}
