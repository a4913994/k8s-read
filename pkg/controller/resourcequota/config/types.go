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

// ResourceQuotaControllerConfiguration contains elements describing ResourceQuotaController.
type ResourceQuotaControllerConfiguration struct {
	// resourceQuotaSyncPeriod is the period for syncing quota usage status
	// in the system.
	// resourceQuotaSyncPeriod是系统中同步配额使用状态的周期。
	ResourceQuotaSyncPeriod metav1.Duration
	// concurrentResourceQuotaSyncs is the number of resource quotas that are
	// allowed to sync concurrently. Larger number = more responsive quota
	// management, but more CPU (and network) load.
	// concurrentResourceQuotaSyncs是允许并发同步的资源配额的数量。更大的数量=更灵敏的配额管理，但更多的CPU（和网络）负载。
	ConcurrentResourceQuotaSyncs int32
}
