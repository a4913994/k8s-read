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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RuntimeClass defines a class of container runtime supported in the cluster.
// The RuntimeClass is used to determine which container runtime is used to run
// all containers in a pod. RuntimeClasses are manually defined by a
// user or cluster provisioner, and referenced in the PodSpec. The Kubelet is
// responsible for resolving the RuntimeClassName reference before running the
// pod.  For more details, see
// https://kubernetes.io/docs/concepts/containers/runtime-class/
// RuntimeClass 定义了集群支持的容器运行时类。
// RuntimeClass 用于确定用于运行 Pod 中的所有容器的容器运行时。 RuntimeClasses 是由用户或集群供应商手动定义的，并在 PodSpec 中引用。
// Kubelet 负责在运行 Pod 之前解析 RuntimeClassName 引用。
type RuntimeClass struct {
	metav1.TypeMeta `json:",inline"`
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Handler specifies the underlying runtime and configuration that the CRI
	// implementation will use to handle pods of this class. The possible values
	// are specific to the node & CRI configuration.  It is assumed that all
	// handlers are available on every node, and handlers of the same name are
	// equivalent on every node.
	// For example, a handler called "runc" might specify that the runc OCI
	// runtime (using native Linux containers) will be used to run the containers
	// in a pod.
	// The Handler must be lowercase, conform to the DNS Label (RFC 1123) requirements,
	// and is immutable.
	// Handler 指定CRI实现将用于处理该类pod的底层运行时和配置. 可能的值是特定于节点和CRI配置的. 假定所有处理程序都可用于每个节点, 并且同名处理程序在每个节点上是等效的.
	// 例如, 名为 "runc" 的处理程序可能会指定 runc OCI 运行时(使用本机 Linux 容器)将用于在 pod 中运行容器.
	// Handler 必须小写, 符合 DNS 标签(RFC 1123)要求, 并且是不可变的.
	Handler string `json:"handler" protobuf:"bytes,2,opt,name=handler"`

	// Overhead represents the resource overhead associated with running a pod for a
	// given RuntimeClass. For more details, see
	//  https://kubernetes.io/docs/concepts/scheduling-eviction/pod-overhead/
	// +optional
	// Overhead 表示与运行给定 RuntimeClass 的 pod 关联的资源开销.
	Overhead *Overhead `json:"overhead,omitempty" protobuf:"bytes,3,opt,name=overhead"`

	// Scheduling holds the scheduling constraints to ensure that pods running
	// with this RuntimeClass are scheduled to nodes that support it.
	// If scheduling is nil, this RuntimeClass is assumed to be supported by all
	// nodes.
	// +optional
	// Scheduling 保存调度约束, 以确保使用此 RuntimeClass 运行的 pod 被调度到支持它的节点.
	Scheduling *Scheduling `json:"scheduling,omitempty" protobuf:"bytes,4,opt,name=scheduling"`
}

// Overhead structure represents the resource overhead associated with running a pod.
// Overhead 结构表示与运行 pod 关联的资源开销.
type Overhead struct {
	// PodFixed represents the fixed resource overhead associated with running a pod.
	// +optional
	// PodFixed 表示与运行 pod 关联的固定资源开销.
	PodFixed corev1.ResourceList `json:"podFixed,omitempty" protobuf:"bytes,1,opt,name=podFixed,casttype=k8s.io/api/core/v1.ResourceList,castkey=k8s.io/api/core/v1.ResourceName,castvalue=k8s.io/apimachinery/pkg/api/resource.Quantity"`
}

// Scheduling specifies the scheduling constraints for nodes supporting a
// RuntimeClass.
// Scheduling 指定支持 RuntimeClass 的节点的调度约束.
type Scheduling struct {
	// nodeSelector lists labels that must be present on nodes that support this
	// RuntimeClass. Pods using this RuntimeClass can only be scheduled to a
	// node matched by this selector. The RuntimeClass nodeSelector is merged
	// with a pod's existing nodeSelector. Any conflicts will cause the pod to
	// be rejected in admission.
	// +optional
	// +mapType=atomic
	// nodeSelector 列出必须出现在支持此 RuntimeClass 的节点上的标签. 使用此 RuntimeClass 的 pod 只能调度到与此选择器匹配的节点上.
	// RuntimeClass nodeSelector 与 pod 的现有 nodeSelector 合并. 任何冲突都会导致 pod 在入站时被拒绝.
	NodeSelector map[string]string `json:"nodeSelector,omitempty" protobuf:"bytes,1,opt,name=nodeSelector"`

	// tolerations are appended (excluding duplicates) to pods running with this
	// RuntimeClass during admission, effectively unioning the set of nodes
	// tolerated by the pod and the RuntimeClass.
	// +optional
	// +listType=atomic
	// tolerations 在入站期间附加(排除重复项)到使用此 RuntimeClass 运行的 pod, 有效地将 pod 和 RuntimeClass 容忍的节点集合联合起来.
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,2,rep,name=tolerations"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RuntimeClassList is a list of RuntimeClass objects.
type RuntimeClassList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is a list of schema objects.
	Items []RuntimeClass `json:"items" protobuf:"bytes,2,rep,name=items"`
}
