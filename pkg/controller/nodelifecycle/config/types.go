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

// NodeLifecycleControllerConfiguration contains elements describing NodeLifecycleController.
type NodeLifecycleControllerConfiguration struct {
	// If set to true enables NoExecute Taints and will evict all not-tolerating
	// Pod running on Nodes tainted with this kind of Taints.
	// 如果设置为 "true"，将启用NoExecute Taints，并将驱逐所有运行在被这种Taints污染的节点上的不容忍的Pod。
	EnableTaintManager bool
	// nodeEvictionRate is the number of nodes per second on which pods are deleted in case of node failure when a zone is healthy
	// nodeEvictionRate是指当一个区域是健康的时候，在节点故障的情况下，每秒删除pod的节点数。
	NodeEvictionRate float32
	// secondaryNodeEvictionRate is the number of nodes per second on which pods are deleted in case of node failure when a zone is unhealthy
	// secondaryNodeEvictionRate是指当一个区域不健康时，在节点故障的情况下，每秒删除pod的节点数量。
	SecondaryNodeEvictionRate float32
	// nodeStartupGracePeriod is the amount of time which we allow starting a node to
	// be unresponsive before marking it unhealthy.
	// nodeStartupGracePeriod是指在标记为不健康之前，我们允许启动一个节点不响应的时间。
	NodeStartupGracePeriod metav1.Duration
	// NodeMonitorGracePeriod is the amount of time which we allow a running node to be
	// unresponsive before marking it unhealthy. Must be N times more than kubelet's
	// nodeStatusUpdateFrequency, where N means number of retries allowed for kubelet
	// to post node status.
	// NodeMonitorGracePeriod是我们允许一个运行中的节点在标记为不健康之前无响应的时间。
	// 必须是kubelet的nodeStatusUpdateFrequency的N倍，其中N是指kubelet发布节点状态时允许的重试次数。
	NodeMonitorGracePeriod metav1.Duration
	// podEvictionTimeout is the grace period for deleting pods on failed nodes.
	// podEvictionTimeout是在失败的节点上删除pod的宽限期
	PodEvictionTimeout metav1.Duration
	// secondaryNodeEvictionRate is implicitly overridden to 0 for clusters smaller than or equal to largeClusterSizeThreshold
	// 对于小于或等于largeClusterSizeThreshold的集群，secondaryNodeEvictionRate被隐式地覆盖为0。
	LargeClusterSizeThreshold int32
	// Zone is treated as unhealthy in nodeEvictionRate and secondaryNodeEvictionRate when at least
	// unhealthyZoneThreshold (no less than 3) of Nodes in the zone are NotReady
	// 当区域内至少有unhealthyZoneThreshold（不少于3个）的节点不准备时，
	// 在nodeEvictionRate和secondaryNodeEvictionRate中，区域被视为不健康的。
	UnhealthyZoneThreshold float32
}
