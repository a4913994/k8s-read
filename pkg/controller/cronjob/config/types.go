/*
Copyright 2020 The Kubernetes Authors.

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

// CronJobControllerConfiguration contains elements describing the
// CronJobControllerV2.
type CronJobControllerConfiguration struct {
	// ConcurrentCronJobSyncs is the number of cron job objects that are
	// allowed to sync concurrently. Larger number = more responsive jobs,
	// but more CPU (and network) load.
	// Concurrent CronJob Syncs 是允许同时同步的 cron 作业对象的数量。大数字 = 更多的作业，但 CPU（和网络）负载更多。
	ConcurrentCronJobSyncs int32
}
