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

// ReplicationControllerConfiguration contains elements describing ReplicationController.
type ReplicationControllerConfiguration struct {
	// concurrentRCSyncs is the number of replication controllers that are
	// allowed to sync concurrently. Larger number = more responsive replica
	// management, but more CPU (and network) load.
	// concurrentRCSyncs是允许同步的复制控制器的数量。更大的数量=更多的复制管理响应，但更多的CPU（和网络）负载。
	ConcurrentRCSyncs int32
}
