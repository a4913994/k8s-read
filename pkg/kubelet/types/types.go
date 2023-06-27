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

package types

import (
	"net/http"
	"time"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

// TODO: Reconcile custom types in kubelet/types and this subpackage

// HTTPDoer encapsulates http.Do functionality
// HTTPDoer封装了http.Do功能
type HTTPDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// Timestamp wraps around time.Time and offers utilities to format and parse
// the time using RFC3339Nano
// Timestamp包装了time.Time，并提供了格式化和解析时间的实用程序，使用RFC3339Nano
type Timestamp struct {
	time time.Time
}

// NewTimestamp returns a Timestamp object using the current time.
// NewTimestamp使用当前时间返回一个Timestamp对象。
func NewTimestamp() *Timestamp {
	return &Timestamp{time.Now()}
}

// ConvertToTimestamp takes a string, parses it using the RFC3339NanoLenient layout,
// and converts it to a Timestamp object.
// ConvertToTimestamp接受一个字符串，使用RFC3339NanoLenient布局对其进行解析，并将其转换为Timestamp对象。
func ConvertToTimestamp(timeString string) *Timestamp {
	parsed, _ := time.Parse(RFC3339NanoLenient, timeString)
	return &Timestamp{parsed}
}

// Get returns the time as time.Time.
// Get返回时间作为time.Time。
func (t *Timestamp) Get() time.Time {
	return t.time
}

// GetString returns the time in the string format using the RFC3339NanoFixed
// layout.
// GetString使用RFC3339NanoFixed布局返回字符串格式的时间。
func (t *Timestamp) GetString() string {
	return t.time.Format(RFC3339NanoFixed)
}

// SortedContainerStatuses is a type to help sort container statuses based on container names.
// SortedContainerStatuses是一种类型，可帮助根据容器名称对容器状态进行排序。
type SortedContainerStatuses []v1.ContainerStatus

func (s SortedContainerStatuses) Len() int      { return len(s) }
func (s SortedContainerStatuses) Swap(i, j int) { s[i], s[j] = s[j], s[i] }

func (s SortedContainerStatuses) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}

// SortInitContainerStatuses ensures that statuses are in the order that their
// init container appears in the pod spec
// SortInitContainerStatuses确保状态按其初始化容器在pod spec中出现的顺序排列
func SortInitContainerStatuses(p *v1.Pod, statuses []v1.ContainerStatus) {
	containers := p.Spec.InitContainers
	current := 0
	for _, container := range containers {
		for j := current; j < len(statuses); j++ {
			if container.Name == statuses[j].Name {
				statuses[current], statuses[j] = statuses[j], statuses[current]
				current++
				break
			}
		}
	}
}

// SortStatusesOfInitContainers returns the statuses of InitContainers of pod p,
// in the order that they appear in its spec.
// SortStatusesOfInitContainers按照它们在其规范中出现的顺序返回pod p的InitContainers的状态。
func SortStatusesOfInitContainers(p *v1.Pod, statusMap map[string]*v1.ContainerStatus) []v1.ContainerStatus {
	containers := p.Spec.InitContainers
	statuses := []v1.ContainerStatus{}
	for _, container := range containers {
		if status, found := statusMap[container.Name]; found {
			statuses = append(statuses, *status)
		}
	}
	return statuses
}

// Reservation represents reserved resources for non-pod components.
// Reservation表示为非pod组件保留的资源。
type Reservation struct {
	// System represents resources reserved for non-kubernetes components.
	// System表示为非kubernetes组件保留的资源。
	System v1.ResourceList
	// Kubernetes represents resources reserved for kubernetes system components.
	// Kubernetes表示为kubernetes系统组件保留的资源。
	Kubernetes v1.ResourceList
}

// ResolvedPodUID is a pod UID which has been translated/resolved to the representation known to kubelets.
// ResolvedPodUID是一个pod UID，它已被转换/解析为已知的kubelet表示形式。
type ResolvedPodUID types.UID

// MirrorPodUID is a pod UID for a mirror pod.
// MirrorPodUID是镜像pod的pod UID。
type MirrorPodUID types.UID
