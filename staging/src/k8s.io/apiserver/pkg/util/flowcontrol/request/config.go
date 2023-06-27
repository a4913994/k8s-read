/*
Copyright 2022 The Kubernetes Authors.

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

package request

import (
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	// minimumSeats 是请求必须占用的最少席位数。
	minimumSeats = 1
	// maximumSeats 是请求可以占用的最大席位数。
	maximumSeats = 10
	// objectsPerSeat 是每个席位的对象数。
	objectsPerSeat = 100.0
	// watchesPerSeat 是每个席位的监视器数。
	watchesPerSeat = 10.0
	// enableMutatingWorkEstimator 指示是否启用可变工作估算器。
	enableMutatingWorkEstimator = true
)

var eventAdditionalDuration = 5 * time.Millisecond

// WorkEstimatorConfig holds work estimator parameters.
// WorkEstimatorConfig 保存工作估计器参数。
type WorkEstimatorConfig struct {
	*ListWorkEstimatorConfig     `json:"listWorkEstimatorConfig,omitempty"`
	*MutatingWorkEstimatorConfig `json:"mutatingWorkEstimatorConfig,omitempty"`

	// MinimumSeats is the minimum number of seats a request must occupy.
	MinimumSeats uint64 `json:"minimumSeats,omitempty"`
	// MaximumSeats is the maximum number of seats a request can occupy
	//
	// NOTE: work_estimate_seats_samples metric uses the value of maximumSeats
	// as the upper bound, so when we change maximumSeats we should also
	// update the buckets of the metric.
	// 注意：work_estimate_seats_samples 指标使用 maximumSeats 的值作为上限，因此当我们更改 maximumSeats 时，我们还应该更新指标的桶。
	MaximumSeats uint64 `json:"maximumSeats,omitempty"`
}

// ListWorkEstimatorConfig holds work estimator parameters related to list requests.
// ListWorkEstimatorConfig 保存与列表请求相关的工作估计器参数。
type ListWorkEstimatorConfig struct {
	ObjectsPerSeat float64 `json:"objectsPerSeat,omitempty"`
}

// MutatingWorkEstimatorConfig holds work estimator
// parameters related to watches of mutating objects.
// MutatingWorkEstimatorConfig 保存与变异对象监视相关的工作估计器参数。
type MutatingWorkEstimatorConfig struct {
	// TODO(wojtekt): Remove it once we tune the algorithm to not fail
	// scalability tests.
	Enabled                 bool            `json:"enable,omitempty"`
	EventAdditionalDuration metav1.Duration `json:"eventAdditionalDurationMs,omitempty"`
	WatchesPerSeat          float64         `json:"watchesPerSeat,omitempty"`
}

// DefaultWorkEstimatorConfig creates a new WorkEstimatorConfig with default values.
// DefaultWorkEstimatorConfig 使用默认值创建一个新的 WorkEstimatorConfig。
func DefaultWorkEstimatorConfig() *WorkEstimatorConfig {
	return &WorkEstimatorConfig{
		MinimumSeats:                minimumSeats,
		MaximumSeats:                maximumSeats,
		ListWorkEstimatorConfig:     defaultListWorkEstimatorConfig(),
		MutatingWorkEstimatorConfig: defaultMutatingWorkEstimatorConfig(),
	}
}

// defaultListWorkEstimatorConfig creates a new ListWorkEstimatorConfig with default values.
// defaultListWorkEstimatorConfig 使用默认值创建一个新的 ListWorkEstimatorConfig。
func defaultListWorkEstimatorConfig() *ListWorkEstimatorConfig {
	return &ListWorkEstimatorConfig{ObjectsPerSeat: objectsPerSeat}
}

// defaultMutatingWorkEstimatorConfig creates a new MutatingWorkEstimatorConfig with default values.
// defaultMutatingWorkEstimatorConfig 使用默认值创建一个新的 MutatingWorkEstimatorConfig。
func defaultMutatingWorkEstimatorConfig() *MutatingWorkEstimatorConfig {
	return &MutatingWorkEstimatorConfig{
		Enabled:                 enableMutatingWorkEstimator,
		EventAdditionalDuration: metav1.Duration{Duration: eventAdditionalDuration},
		WatchesPerSeat:          watchesPerSeat,
	}
}

// eventAdditionalDuration converts eventAdditionalDurationMs to a time.Duration type.
// eventAdditionalDuration 将 eventAdditionalDurationMs 转换为 time.Duration 类型。
func (c *MutatingWorkEstimatorConfig) eventAdditionalDuration() time.Duration {
	return c.EventAdditionalDuration.Duration
}
