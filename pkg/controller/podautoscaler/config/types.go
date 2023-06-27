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

// HPAControllerConfiguration contains elements describing HPAController.
type HPAControllerConfiguration struct {
	// ConcurrentHorizontalPodAutoscalerSyncs is the number of HPA objects that are allowed to sync concurrently.
	// Larger number = more responsive HPA processing, but more CPU (and network) load.
	// ConcurrentHorizontalPodAutoscalerSyncs 是允许并发同步的 HPA 对象的数量。更大的数字 = 更灵敏的 HPA 处理，但更多的 CPU（和网络）负载。
	ConcurrentHorizontalPodAutoscalerSyncs int32
	// horizontalPodAutoscalerSyncPeriod is the period for syncing the number of
	// pods in horizontal pod autoscaler.
	// horizontalPodAutoscalerSyncPeriod 是 horizontal Pod autoscaler 同步 Pod 数量的周期。
	HorizontalPodAutoscalerSyncPeriod metav1.Duration
	// horizontalPodAutoscalerUpscaleForbiddenWindow is a period after which next upscale allowed.
	// horizontalPodAutoscalerUpscaleForbiddenWindow 是允许下一次升级的时间段。
	HorizontalPodAutoscalerUpscaleForbiddenWindow metav1.Duration
	// horizontalPodAutoscalerDownscaleForbiddenWindow is a period after which next downscale allowed.
	// horizontalPodAutoscalerDownscaleForbiddenWindow 是允许下一次缩小的时间段。
	HorizontalPodAutoscalerDownscaleForbiddenWindow metav1.Duration
	// HorizontalPodAutoscalerDowncaleStabilizationWindow is a period for which autoscaler will look
	// backwards and not scale down below any recommendation it made during that period.
	// HorizontalPodAutoscalerDowncaleStabilizationWindow 是一个期间，autoscaler 将向后看，并且不会缩小到低于它在此期间提出的任何建议。
	HorizontalPodAutoscalerDownscaleStabilizationWindow metav1.Duration
	// horizontalPodAutoscalerTolerance is the tolerance for when
	// resource usage suggests upscaling/downscaling
	// horizontalPodAutoscalerTolerance是指当资源使用量建议向上扩容或向下缩放时的容忍度。
	HorizontalPodAutoscalerTolerance float64
	// HorizontalPodAutoscalerCPUInitializationPeriod is the period after pod start when CPU samples
	// might be skipped.
	// HorizontalPodAutoscalerCPUInitializationPeriod是pod启动后CPU采样可能被跳过的时期。
	HorizontalPodAutoscalerCPUInitializationPeriod metav1.Duration
	// HorizontalPodAutoscalerInitialReadinessDelay is period after pod start during which readiness
	// changes are treated as readiness being set for the first time. The only effect of this is that
	// HPA will disregard CPU samples from unready pods that had last readiness change during that
	// period.
	// HorizontalPodAutoscalerInitialReadinessDelay是指POD启动后的一段时间，在此期间，准备状态的变化被视为首次设置准备状态。
	// 这样做的唯一影响是，HPA将不考虑来自未准备好的POD的CPU样本，这些POD的最后一次准备状态变化是在这段时间内。
	HorizontalPodAutoscalerInitialReadinessDelay metav1.Duration
}
