/*
Copyright 2015 The Kubernetes Authors.

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

package qos

import (
	v1 "k8s.io/api/core/v1"
	v1qos "k8s.io/kubernetes/pkg/apis/core/v1/helper/qos"
	"k8s.io/kubernetes/pkg/kubelet/types"
)

const (
	// KubeletOOMScoreAdj is the OOM score adjustment for Kubelet
	// KubeletOOMScoreAdj是Kubelet的OOM分数调整
	KubeletOOMScoreAdj int = -999
	// KubeProxyOOMScoreAdj is the OOM score adjustment for kube-proxy
	// KubeProxyOOMScoreAdj是kube-proxy的OOM分数调整
	KubeProxyOOMScoreAdj  int = -999
	guaranteedOOMScoreAdj int = -997
	besteffortOOMScoreAdj int = 1000
)

// GetContainerOOMScoreAdjust returns the amount by which the OOM score of all processes in the
// container should be adjusted.
// The OOM score of a process is the percentage of memory it consumes
// multiplied by 10 (barring exceptional cases) + a configurable quantity which is between -1000
// and 1000. Containers with higher OOM scores are killed if the system runs out of memory.
// See https://lwn.net/Articles/391222/ for more information.
// GetContainerOOMScoreAdjust返回应调整容器中所有进程的OOM分数的数量。
// 进程的OOM分数是它消耗的内存的百分比乘以10（除非有特殊情况）+可配置的数量，该数量介于-1000和1000之间。
// 如果系统内存不足，将杀死具有更高OOM分数的容器。 有关更多信息，请参见https://lwn.net/Articles/391222/。
func GetContainerOOMScoreAdjust(pod *v1.Pod, container *v1.Container, memoryCapacity int64) int {
	if types.IsNodeCriticalPod(pod) {
		// Only node critical pod should be the last to get killed.
		// Only node critical pod应该是最后一个被杀死的。
		return guaranteedOOMScoreAdj
	}

	switch v1qos.GetPodQOS(pod) {
	case v1.PodQOSGuaranteed:
		// Guaranteed containers should be the last to get killed.
		// Guaranteed容器应该是最后一个被杀死的。
		return guaranteedOOMScoreAdj
	case v1.PodQOSBestEffort:
		return besteffortOOMScoreAdj
	}

	// Burstable containers are a middle tier, between Guaranteed and Best-Effort. Ideally,
	// we want to protect Burstable containers that consume less memory than requested.
	// The formula below is a heuristic. A container requesting for 10% of a system's
	// memory will have an OOM score adjust of 900. If a process in container Y
	// uses over 10% of memory, its OOM score will be 1000. The idea is that containers
	// which use more than their request will have an OOM score of 1000 and will be prime
	// targets for OOM kills.
	// Note that this is a heuristic, it won't work if a container has many small processes.
	// Burstable容器是保证和最佳效果之间的中间层。 理想情况下，我们希望保护消耗小于请求的内存的可变容器。
	// 下面的公式是一种启发式方法。 请求系统内存的10％的容器将具有调整为900的OOM分数。
	// 如果容器Y中的进程使用超过10％的内存，则其OOM分数将为1000。 想法是使用超过其请求的容器将具有1000的OOM分数，并将成为OOM杀死的首要目标。
	// 请注意，这是一种启发式方法，如果容器有许多小进程，则不起作用。
	memoryRequest := container.Resources.Requests.Memory().Value()
	oomScoreAdjust := 1000 - (1000*memoryRequest)/memoryCapacity
	// A guaranteed pod using 100% of memory can have an OOM score of 10. Ensure
	// that burstable pods have a higher OOM score adjustment.
	// 使用100％内存的保证pod可以具有10的OOM分数。 确保可变pod具有更高的OOM分数调整。
	if int(oomScoreAdjust) < (1000 + guaranteedOOMScoreAdj) {
		return (1000 + guaranteedOOMScoreAdj)
	}
	// Give burstable pods a higher chance of survival over besteffort pods.
	// 给可变pod一个比最佳效果pod更高的生存机会。
	if int(oomScoreAdjust) == besteffortOOMScoreAdj {
		return int(oomScoreAdjust - 1)
	}
	return int(oomScoreAdjust)
}
