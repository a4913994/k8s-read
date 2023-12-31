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

// DaemonSetControllerConfiguration contains elements describing DaemonSetController.
type DaemonSetControllerConfiguration struct {
	// concurrentDaemonSetSyncs is the number of daemonset objects that are
	// allowed to sync concurrently. Larger number = more responsive daemonset,
	// but more CPU (and network) load.
	// concurrentDaemonSetSyncs 是允许并发同步的 daemonset 对象的数量。更大的数字 = 更多的守护进程，但有更多的 CPU（和网络）负载
	ConcurrentDaemonSetSyncs int32
}
