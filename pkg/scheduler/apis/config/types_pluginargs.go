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

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// DefaultPreemptionArgs holds arguments used to configure the
// DefaultPreemption plugin.
// DefaultPreemptionArgs保存用于配置DefaultPreemption插件的参数。
type DefaultPreemptionArgs struct {
	metav1.TypeMeta

	// MinCandidateNodesPercentage is the minimum number of candidates to
	// shortlist when dry running preemption as a percentage of number of nodes.
	// Must be in the range [0, 100]. Defaults to 10% of the cluster size if
	// unspecified.
	// MinCandidateNodesPercentage是干运行抢占时候选候选列表的最小数量占节点数的百分比。必须在[0,100]范围内。如果未指定，默认为集群大小的10%。
	MinCandidateNodesPercentage int32
	// MinCandidateNodesAbsolute is the absolute minimum number of candidates to
	// shortlist. The likely number of candidates enumerated for dry running
	// preemption is given by the formula:
	// numCandidates = max(numNodes * minCandidateNodesPercentage, minCandidateNodesAbsolute)
	// We say "likely" because there are other factors such as PDB violations
	// that play a role in the number of candidates shortlisted. Must be at least
	// 0 nodes. Defaults to 100 nodes if unspecified.
	// MinCandidateNodesAbsolute是候选名单的绝对最小数量。
	// 我们说“可能”是因为还有其他因素，比如PDB违规，这些因素会对候选候选的数量产生影响。至少为0个节点。如果未指定，默认为100个节点。
	MinCandidateNodesAbsolute int32
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InterPodAffinityArgs holds arguments used to configure the InterPodAffinity plugin.
// InterPodAffinityArgs 保存用于配置InterPodAffinity插件的参数。
type InterPodAffinityArgs struct {
	metav1.TypeMeta

	// HardPodAffinityWeight is the scoring weight for existing pods with a
	// matching hard affinity to the incoming pod.
	// HardPodAffinityWeight 是现有Pods与传入pod具有匹配硬亲和力的评分权重。
	HardPodAffinityWeight int32
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeResourcesFitArgs holds arguments used to configure the NodeResourcesFit plugin.
// NodeResourcesFitArgs 保存用于配置NodeResourcesFit插件的参数。
type NodeResourcesFitArgs struct {
	metav1.TypeMeta

	// IgnoredResources is the list of resources that NodeResources fit filter
	// should ignore.
	// IgnoredResources是NodeResources fit过滤器应该忽略的资源列表。
	IgnoredResources []string
	// IgnoredResourceGroups defines the list of resource groups that NodeResources fit filter should ignore.
	// e.g. if group is ["example.com"], it will ignore all resource names that begin
	// with "example.com", such as "example.com/aaa" and "example.com/bbb".
	// A resource group name can't contain '/'.
	// IgnoredResourceGroups定义了NodeResources fit过滤器应该忽略的资源组列表。
	// 例如，如果group是["example.com"]，它将忽略所有以"example.com"开头的资源名，
	// 例如"example.com/aaa"和"example.com/bbb"。资源组名称中不能包含“/”。
	IgnoredResourceGroups []string

	// ScoringStrategy selects the node resource scoring strategy.
	// ScoringStrategy 选择节点资源评分策略。
	ScoringStrategy *ScoringStrategy
}

// PodTopologySpreadConstraintsDefaulting defines how to set default constraints
// for the PodTopologySpread plugin.
// PodTopologySpreadConstraintsDefaulting 定义如何为PodTopologySpread插件设置默认约束。
type PodTopologySpreadConstraintsDefaulting string

const (
	// SystemDefaulting instructs to use the kubernetes defined default.
	// SystemDefaulting指示使用kubernetes定义的default。
	SystemDefaulting PodTopologySpreadConstraintsDefaulting = "System"
	// ListDefaulting instructs to use the config provided default.
	ListDefaulting PodTopologySpreadConstraintsDefaulting = "List"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodTopologySpreadArgs holds arguments used to configure the PodTopologySpread plugin.
// PodTopologySpreadArgs保存用于配置PodTopologySpread插件的参数。
type PodTopologySpreadArgs struct {
	metav1.TypeMeta

	// DefaultConstraints defines topology spread constraints to be applied to
	// Pods that don't define any in `pod.spec.topologySpreadConstraints`.
	// `.defaultConstraints[*].labelSelectors` must be empty, as they are
	// deduced from the Pod's membership to Services, ReplicationControllers,
	// ReplicaSets or StatefulSets.
	// When not empty, .defaultingType must be "List".
	// DefaultConstraints定义了拓扑扩展约束，应用于没有在' pod.spec.topologySpreadConstraints '中定义任何约束的Pods。
	// ' . defaultconstraints []. labelselectors '必须为空，因为它们是从Pod的成员到Services, ReplicationControllers, ReplicaSets或StatefulSets推断出来的。当不为空时，. defaultingtype必须为"List"。
	DefaultConstraints []v1.TopologySpreadConstraint

	// DefaultingType determines how .defaultConstraints are deduced. Can be one
	// of "System" or "List".
	//
	// - "System": Use kubernetes defined constraints that spread Pods among
	//   Nodes and Zones.
	// - "List": Use constraints defined in .defaultConstraints.
	//
	// Defaults to "System".
	// +optional
	// DefaultingType决定如何推导. defaultconstraints。可以是“系统”或“列表”之一。-“系统”:使用kubernetes定义的约束，在节点和区域之间传播pod。-“List”:使用. defaultconstraints中定义的约束。默认为“System”。+可选
	DefaultingType PodTopologySpreadConstraintsDefaulting
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeResourcesBalancedAllocationArgs holds arguments used to configure NodeResourcesBalancedAllocation plugin.
// NodeResourcesBalancedAllocationArgs保存用于配置NodeResourcesBalancedAllocation插件的参数。
type NodeResourcesBalancedAllocationArgs struct {
	metav1.TypeMeta

	// Resources to be considered when scoring.
	// The default resource set includes "cpu" and "memory", only valid weight is 1.
	// 评分时要考虑的资源。默认资源集包括“cpu”和“memory”，只有有效的权重为1。
	Resources []ResourceSpec
}

// UtilizationShapePoint represents a single point of a priority function shape.
// UtilizationShapePoint 表示优先级函数形状的单个点。
type UtilizationShapePoint struct {
	// Utilization (x axis). Valid values are 0 to 100. Fully utilized node maps to 100.
	// 利用率(x轴)。取值范围为0 ~ 100。充分利用的节点映射到100。
	Utilization int32
	// Score assigned to a given utilization (y axis). Valid values are 0 to 10.
	// 分配给给定利用率的分数(y轴)。取值范围为0 ~ 10。
	Score int32
}

// ResourceSpec represents single resource.
type ResourceSpec struct {
	// Name of the resource.
	Name string
	// Weight of the resource.
	Weight int64
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// VolumeBindingArgs holds arguments used to configure the VolumeBinding plugin.
type VolumeBindingArgs struct {
	metav1.TypeMeta

	// BindTimeoutSeconds is the timeout in seconds in volume binding operation.
	// Value must be non-negative integer. The value zero indicates no waiting.
	// If this value is nil, the default value will be used.
	// BindTimeoutSeconds:卷绑定操作超时时间，单位为秒。值必须是非负整数。0表示不等待。如果该值为nil，将使用默认值
	BindTimeoutSeconds int64

	// Shape specifies the points defining the score function shape, which is
	// used to score nodes based on the utilization of statically provisioned
	// PVs. The utilization is calculated by dividing the total requested
	// storage of the pod by the total capacity of feasible PVs on each node.
	// Each point contains utilization (ranges from 0 to 100) and its
	// associated score (ranges from 0 to 10). You can turn the priority by
	// specifying different scores for different utilization numbers.
	// The default shape points are:
	// 1) 0 for 0 utilization
	// 2) 10 for 100 utilization
	// All points must be sorted in increasing order by utilization.
	// +featureGate=VolumeCapacityPriority
	// +optional
	// Shape指定定义评分函数形状的点，该函数形状用于根据静态配置pv的使用情况对节点进行评分。通过将pod请求的总存储容量除以每个节点上可行pv的总容量来计算利用率。
	// 每个点包含利用率(范围从0到100)及其关联的得分(范围从0到10)。您可以通过为不同的利用率指定不同的分数来改变优先级。
	Shape []UtilizationShapePoint
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeAffinityArgs holds arguments to configure the NodeAffinity plugin.
type NodeAffinityArgs struct {
	metav1.TypeMeta

	// AddedAffinity is applied to all Pods additionally to the NodeAffinity
	// specified in the PodSpec. That is, Nodes need to satisfy AddedAffinity
	// AND .spec.NodeAffinity. AddedAffinity is empty by default (all Nodes
	// match).
	// When AddedAffinity is used, some Pods with affinity requirements that match
	// a specific Node (such as Daemonset Pods) might remain unschedulable.
	// 除了PodSpec中指定的NodeAffinity外，adddedaffinity还应用于所有pod。
	// 也就是说，节点需要满足adddedaffinity和.spec. nodeaffinity。AddedAffinity默认为空(所有节点匹配)。
	// 当使用adddedaffinity时，一些具有匹配特定节点的亲和需求的Pods(例如Daemonset Pods)可能仍然不可调度。
	AddedAffinity *v1.NodeAffinity
}

// ScoringStrategyType the type of scoring strategy used in NodeResourcesFit plugin.
type ScoringStrategyType string

const (
	// LeastAllocated strategy prioritizes nodes with least allocated resources.
	LeastAllocated ScoringStrategyType = "LeastAllocated"
	// MostAllocated strategy prioritizes nodes with most allocated resources.
	MostAllocated ScoringStrategyType = "MostAllocated"
	// RequestedToCapacityRatio strategy allows specifying a custom shape function
	// to score nodes based on the request to capacity ratio.
	RequestedToCapacityRatio ScoringStrategyType = "RequestedToCapacityRatio"
)

// ScoringStrategy define ScoringStrategyType for node resource plugin
// 为节点资源插件定义ScoringStrategyType
type ScoringStrategy struct {
	// Type selects which strategy to run.
	Type ScoringStrategyType

	// Resources to consider when scoring.
	// The default resource set includes "cpu" and "memory" with an equal weight.
	// Allowed weights go from 1 to 100.
	// Weight defaults to 1 if not specified or explicitly set to 0.
	// 评分时需要考虑的资源。默认资源集包括“cpu”和“memory”，两者权重相等。允许重量从1到100。Weight如果没有指定或显式设置为0，则默认为1。
	Resources []ResourceSpec

	// Arguments specific to RequestedToCapacityRatio strategy.
	RequestedToCapacityRatio *RequestedToCapacityRatioParam
}

// RequestedToCapacityRatioParam define RequestedToCapacityRatio parameters
type RequestedToCapacityRatioParam struct {
	// Shape is a list of points defining the scoring function shape.
	// 形状是定义计分函数形状的点列表。
	Shape []UtilizationShapePoint
}
