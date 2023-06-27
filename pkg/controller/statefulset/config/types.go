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

// StatefulSetControllerConfiguration contains elements describing StatefulSetController.
type StatefulSetControllerConfiguration struct {
	// concurrentStatefulSetSyncs is the number of statefulset objects that are
	// allowed to sync concurrently. Larger number = more responsive statefulsets,
	// but more CPU (and network) load.
	// concurrentStatefulSetSyncs是允许并发同步的状态集对象的数量。更大的数量=更多的状态集响应，但更多的CPU（和网络）负载。
	ConcurrentStatefulSetSyncs int32
}
