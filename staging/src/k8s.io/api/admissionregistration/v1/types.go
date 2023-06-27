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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Rule is a tuple of APIGroups, APIVersion, and Resources.It is recommended
// to make sure that all the tuple expansions are valid.
// 规则是APIGroups，APIVersion和Resources的元组。建议确保所有元组扩展都是有效的。
type Rule struct {
	// APIGroups is the API groups the resources belong to. '*' is all groups.
	// If '*' is present, the length of the slice must be one.
	// Required.
	// +listType=atomic
	// APIGroups 是资源所属的API组。'*'是所有组。如果存在'*'，则切片的长度必须为1。必填。
	APIGroups []string `json:"apiGroups,omitempty" protobuf:"bytes,1,rep,name=apiGroups"`

	// APIVersions is the API versions the resources belong to. '*' is all versions.
	// If '*' is present, the length of the slice must be one.
	// Required.
	// +listType=atomic
	// APIVersions 是资源所属的API版本。'*'是所有版本。如果存在'*'，则切片的长度必须为1。必填。
	APIVersions []string `json:"apiVersions,omitempty" protobuf:"bytes,2,rep,name=apiVersions"`

	// Resources is a list of resources this rule applies to.
	//
	// For example:
	// 'pods' means pods.
	// 'pods/log' means the log subresource of pods.
	// '*' means all resources, but not subresources.
	// 'pods/*' means all subresources of pods.
	// '*/scale' means all scale subresources.
	// '*/*' means all resources and their subresources.
	//
	// If wildcard is present, the validation rule will ensure resources do not
	// overlap with each other.
	//
	// Depending on the enclosing object, subresources might not be allowed.
	// Required.
	// +listType=atomic
	// Resources 是此规则适用的资源列表。
	// 例如：
	// 'pods'表示pods。
	// 'pods/log'表示pods的日志子资源。
	// '*'表示所有资源，但不包括子资源。
	// 'pods/*'表示pods的所有子资源。
	// '*/scale'表示所有规模子资源。
	// '*/*'表示所有资源及其子资源。
	// 如果存在通配符，则验证规则将确保资源不会相互重叠。
	// 根据封闭对象，可能不允许子资源。
	Resources []string `json:"resources,omitempty" protobuf:"bytes,3,rep,name=resources"`

	// scope specifies the scope of this rule.
	// Valid values are "Cluster", "Namespaced", and "*"
	// "Cluster" means that only cluster-scoped resources will match this rule.
	// Namespace API objects are cluster-scoped.
	// "Namespaced" means that only namespaced resources will match this rule.
	// "*" means that there are no scope restrictions.
	// Subresources match the scope of their parent resource.
	// Default is "*".
	//
	// +optional
	// scope 指定此规则的范围。有效值为“Cluster”，“Namespaced”和“*”“Cluster”表示只有集群范围的资源才能与此规则匹配。命名空间API对象是集群范围的。“Namespaced”表示只有命名空间资源才能与此规则匹配。“*”表示没有范围限制。子资源与其父资源的范围匹配。默认值为“*”。
	Scope *ScopeType `json:"scope,omitempty" protobuf:"bytes,4,rep,name=scope"`
}

// ScopeType specifies a scope for a Rule.
// +enum
// ScopeType 指定规则的范围。
type ScopeType string

const (
	// ClusterScope means that scope is limited to cluster-scoped objects.
	// Namespace objects are cluster-scoped.
	// ClusterScope 表示范围仅限于集群范围的对象。命名空间对象是集群范围的。
	ClusterScope ScopeType = "Cluster"
	// NamespacedScope means that scope is limited to namespaced objects.
	// NamespacedScope 表示范围仅限于命名空间对象。
	NamespacedScope ScopeType = "Namespaced"
	// AllScopes means that all scopes are included.
	// AllScopes 表示包含所有范围。
	AllScopes ScopeType = "*"
)

// FailurePolicyType specifies a failure policy that defines how unrecognized errors from the admission endpoint are handled.
// +enum
// FailurePolicyType 指定定义如何处理来自admission端点的未识别错误的失败策略。
type FailurePolicyType string

const (
	// Ignore means that an error calling the webhook is ignored.
	// Ignore 表示忽略调用webhook的错误。
	Ignore FailurePolicyType = "Ignore"
	// Fail means that an error calling the webhook causes the admission to fail.
	// Fail 表示调用webhook的错误导致admission失败。
	Fail FailurePolicyType = "Fail"
)

// MatchPolicyType specifies the type of match policy.
// +enum
// MatchPolicyType 指定匹配策略的类型。
type MatchPolicyType string

const (
	// Exact means requests should only be sent to the webhook if they exactly match a given rule.
	// Exact 表示只有在请求与给定规则完全匹配时，才将请求发送到webhook。
	Exact MatchPolicyType = "Exact"
	// Equivalent means requests should be sent to the webhook if they modify a resource listed in rules via another API group or version.
	// Equivalent 表示如果请求通过另一个API组或版本修改规则中列出的资源，则应将请求发送到webhook。
	Equivalent MatchPolicyType = "Equivalent"
)

// SideEffectClass specifies the types of side effects a webhook may have.
// +enum
// SideEffectClass 指定webhook可能具有的副作用类型。
type SideEffectClass string

const (
	// SideEffectClassUnknown means that no information is known about the side effects of calling the webhook.
	// If a request with the dry-run attribute would trigger a call to this webhook, the request will instead fail.
	// SideEffectClassUnknown 表示关于调用webhook的副作用的信息是未知的。如果带有dry-run属性的请求将触发对此webhook的调用，则该请求将失败。
	SideEffectClassUnknown SideEffectClass = "Unknown"
	// SideEffectClassNone means that calling the webhook will have no side effects.
	// SideEffectClassNone 表示调用webhook不会产生任何副作用。
	SideEffectClassNone SideEffectClass = "None"
	// SideEffectClassSome means that calling the webhook will possibly have side effects.
	// If a request with the dry-run attribute would trigger a call to this webhook, the request will instead fail.
	// SideEffectClassSome 表示调用webhook可能会产生副作用。如果带有dry-run属性的请求将触发对此webhook的调用，则该请求将失败。
	SideEffectClassSome SideEffectClass = "Some"
	// SideEffectClassNoneOnDryRun means that calling the webhook will possibly have side effects, but if the
	// request being reviewed has the dry-run attribute, the side effects will be suppressed.
	// SideEffectClassNoneOnDryRun 表示调用webhook可能会产生副作用，但如果正在审查的请求具有dry-run属性，则将抑制副作用。
	SideEffectClassNoneOnDryRun SideEffectClass = "NoneOnDryRun"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ValidatingWebhookConfiguration describes the configuration of and admission webhook that accept or reject and object without changing it.
// ValidatingWebhookConfiguration 描述接受或拒绝对象而不更改它的admission webhook的配置。
type ValidatingWebhookConfiguration struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Webhooks is a list of webhooks and the affected resources and operations.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// Webhooks 是webhook的列表，以及受影响的资源和操作。
	Webhooks []ValidatingWebhook `json:"webhooks,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=Webhooks"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ValidatingWebhookConfigurationList is a list of ValidatingWebhookConfiguration.
// ValidatingWebhookConfigurationList 是 ValidatingWebhookConfiguration 的列表。
type ValidatingWebhookConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// List of ValidatingWebhookConfiguration.
	Items []ValidatingWebhookConfiguration `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MutatingWebhookConfiguration describes the configuration of and admission webhook that accept or reject and may change the object.
// MutatingWebhookConfiguration 描述接受或拒绝并可能更改对象的admission webhook的配置。
type MutatingWebhookConfiguration struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object metadata; More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Webhooks is a list of webhooks and the affected resources and operations.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// Webhooks 是webhook的列表，以及受影响的资源和操作。
	Webhooks []MutatingWebhook `json:"webhooks,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=Webhooks"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MutatingWebhookConfigurationList is a list of MutatingWebhookConfiguration.
// MutatingWebhookConfigurationList 是 MutatingWebhookConfiguration 的列表。
type MutatingWebhookConfigurationList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// List of MutatingWebhookConfiguration.
	Items []MutatingWebhookConfiguration `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// ValidatingWebhook describes an admission webhook and the resources and operations it applies to.
// ValidatingWebhook 描述了admission webhook以及它适用于的资源和操作。
type ValidatingWebhook struct {
	// The name of the admission webhook.
	// Name should be fully qualified, e.g., imagepolicy.kubernetes.io, where
	// "imagepolicy" is the name of the webhook, and kubernetes.io is the name
	// of the organization.
	// Required.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// ClientConfig defines how to communicate with the hook.
	// Required
	// ClientConfig定义了如何与钩子通信。
	ClientConfig WebhookClientConfig `json:"clientConfig" protobuf:"bytes,2,opt,name=clientConfig"`

	// Rules describes what operations on what resources/subresources the webhook cares about.
	// The webhook cares about an operation if it matches _any_ Rule.
	// However, in order to prevent ValidatingAdmissionWebhooks and MutatingAdmissionWebhooks
	// from putting the cluster in a state which cannot be recovered from without completely
	// disabling the plugin, ValidatingAdmissionWebhooks and MutatingAdmissionWebhooks are never called
	// on admission requests for ValidatingWebhookConfiguration and MutatingWebhookConfiguration objects.
	// Rules 描述了webhook关心的资源/子资源上的操作。 如果它与任何规则匹配，则webhook关心操作。 但是，为了防止ValidatingAdmissionWebhooks和MutatingAdmissionWebhooks将集群置于无法从中恢复的状态，而不完全禁用插件，ValidatingAdmissionWebhooks和MutatingAdmissionWebhooks永远不会在对ValidatingWebhookConfiguration和MutatingWebhookConfiguration对象的准入请求上调用。
	Rules []RuleWithOperations `json:"rules,omitempty" protobuf:"bytes,3,rep,name=rules"`

	// FailurePolicy defines how unrecognized errors from the admission endpoint are handled -
	// allowed values are Ignore or Fail. Defaults to Fail.
	// +optional
	// FailurePolicy定义了如何处理来自admission端点的未识别的错误-允许的值为Ignore或Fail。 默认为Fail。
	FailurePolicy *FailurePolicyType `json:"failurePolicy,omitempty" protobuf:"bytes,4,opt,name=failurePolicy,casttype=FailurePolicyType"`

	// matchPolicy defines how the "rules" list is used to match incoming requests.
	// Allowed values are "Exact" or "Equivalent".
	//
	// - Exact: match a request only if it exactly matches a specified rule.
	// For example, if deployments can be modified via apps/v1, apps/v1beta1, and extensions/v1beta1,
	// but "rules" only included `apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]`,
	// a request to apps/v1beta1 or extensions/v1beta1 would not be sent to the webhook.
	//
	// - Equivalent: match a request if modifies a resource listed in rules, even via another API group or version.
	// For example, if deployments can be modified via apps/v1, apps/v1beta1, and extensions/v1beta1,
	// and "rules" only included `apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]`,
	// a request to apps/v1beta1 or extensions/v1beta1 would be converted to apps/v1 and sent to the webhook.
	//
	// Defaults to "Equivalent"
	// +optional
	// matchPolicy定义了如何使用“rules”列表来匹配传入请求。 允许的值为“Exact”或“Equivalent”。
	// - Exact: 仅当请求与指定规则完全匹配时才匹配请求。 例如，如果部署可以通过apps/v1，apps/v1beta1和extensions/v1beta1进行修改，但“rules”仅包括`apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]`，则对apps/v1beta1或extensions/v1beta1的请求将不会发送到webhook。
	// - Equivalent: 如果通过另一个API组或版本修改“rules”中列出的资源，则匹配请求。 例如，如果部署可以通过apps/v1，apps/v1beta1和extensions/v1beta1进行修改，并且“rules”仅包括`apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]`，则对apps/v1beta1或extensions/v1beta1的请求将转换为apps/v1并发送到webhook。
	// 默认为“Equivalent”
	MatchPolicy *MatchPolicyType `json:"matchPolicy,omitempty" protobuf:"bytes,9,opt,name=matchPolicy,casttype=MatchPolicyType"`

	// NamespaceSelector decides whether to run the webhook on an object based
	// on whether the namespace for that object matches the selector. If the
	// object itself is a namespace, the matching is performed on
	// object.metadata.labels. If the object is another cluster scoped resource,
	// it never skips the webhook.
	//
	// For example, to run the webhook on any objects whose namespace is not
	// associated with "runlevel" of "0" or "1";  you will set the selector as
	// follows:
	// "namespaceSelector": {
	//   "matchExpressions": [
	//     {
	//       "key": "runlevel",
	//       "operator": "NotIn",
	//       "values": [
	//         "0",
	//         "1"
	//       ]
	//     }
	//   ]
	// }
	//
	// If instead you want to only run the webhook on any objects whose
	// namespace is associated with the "environment" of "prod" or "staging";
	// you will set the selector as follows:
	// "namespaceSelector": {
	//   "matchExpressions": [
	//     {
	//       "key": "environment",
	//       "operator": "In",
	//       "values": [
	//         "prod",
	//         "staging"
	//       ]
	//     }
	//   ]
	// }
	//
	// See
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/labels
	// for more examples of label selectors.
	//
	// Default to the empty LabelSelector, which matches everything.
	// +optional
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty" protobuf:"bytes,5,opt,name=namespaceSelector"`

	// ObjectSelector decides whether to run the webhook based on if the
	// object has matching labels. objectSelector is evaluated against both
	// the oldObject and newObject that would be sent to the webhook, and
	// is considered to match if either object matches the selector. A null
	// object (oldObject in the case of create, or newObject in the case of
	// delete) or an object that cannot have labels (like a
	// DeploymentRollback or a PodProxyOptions object) is not considered to
	// match.
	// Use the object selector only if the webhook is opt-in, because end
	// users may skip the admission webhook by setting the labels.
	// Default to the empty LabelSelector, which matches everything.
	// +optional
	// ObjectSelector决定是否根据对象是否具有匹配标签来运行webhook。 objectSelector针对将发送到webhook的oldObject和newObject进行评估，并且如果任一对象与选择器匹配，则被视为匹配。 空对象（在创建的情况下为oldObject，或在删除的情况下为newObject）或不能具有标签的对象（例如DeploymentRollback或PodProxyOptions对象）不被视为匹配。
	ObjectSelector *metav1.LabelSelector `json:"objectSelector,omitempty" protobuf:"bytes,10,opt,name=objectSelector"`

	// SideEffects states whether this webhook has side effects.
	// Acceptable values are: None, NoneOnDryRun (webhooks created via v1beta1 may also specify Some or Unknown).
	// Webhooks with side effects MUST implement a reconciliation system, since a request may be
	// rejected by a future step in the admission chain and the side effects therefore need to be undone.
	// Requests with the dryRun attribute will be auto-rejected if they match a webhook with
	// sideEffects == Unknown or Some.
	// SideEffects 说明这个网络钩子是否有副作用。
	// 可接受的值为：None，NoneOnDryRun（通过v1beta1创建的webhook也可以指定Some或Unknown）。
	// 具有副作用的Webhook必须实现一种协调系统，因为请求可能会被拒绝，从而需要撤消后续步骤中的副作用。
	// 如果请求与sideEffects == Unknown或Some的网络钩子匹配，则具有dryRun属性的请求将自动被拒绝。
	SideEffects *SideEffectClass `json:"sideEffects" protobuf:"bytes,6,opt,name=sideEffects,casttype=SideEffectClass"`

	// TimeoutSeconds specifies the timeout for this webhook. After the timeout passes,
	// the webhook call will be ignored or the API call will fail based on the
	// failure policy.
	// The timeout value must be between 1 and 30 seconds.
	// Default to 10 seconds.
	// +optional
	// TimeoutSeconds指定此网络钩子的超时时间。 超时过后，将忽略网络钩子调用或根据失败策略使API调用失败。
	TimeoutSeconds *int32 `json:"timeoutSeconds,omitempty" protobuf:"varint,7,opt,name=timeoutSeconds"`

	// AdmissionReviewVersions is an ordered list of preferred `AdmissionReview`
	// versions the Webhook expects. API server will try to use first version in
	// the list which it supports. If none of the versions specified in this list
	// supported by API server, validation will fail for this object.
	// If a persisted webhook configuration specifies allowed versions and does not
	// include any versions known to the API Server, calls to the webhook will fail
	// and be subject to the failure policy.
	// AdmissionReviewVersions 是首选的“AdmissionReview”版本列表，Webhook期望API服务器尝试使用列表中的第一个版本。 如果此列表中指定的版本中的任何一个版本都不受API服务器支持，则此对象的验证将失败。 如果持久化的网络钩子配置指定了允许的版本并且不包含任何已知的版本，则对网络钩子的调用将失败并受到失败策略的约束。
	AdmissionReviewVersions []string `json:"admissionReviewVersions" protobuf:"bytes,8,rep,name=admissionReviewVersions"`
}

// MutatingWebhook describes an admission webhook and the resources and operations it applies to.
// MutatingWebhook 描述了一个admission webhook，以及它适用的资源和操作。
type MutatingWebhook struct {
	// The name of the admission webhook.
	// Name should be fully qualified, e.g., imagepolicy.kubernetes.io, where
	// "imagepolicy" is the name of the webhook, and kubernetes.io is the name
	// of the organization.
	// Required.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// ClientConfig defines how to communicate with the hook.
	// Required
	// ClientConfig 定义如何与hook进行通信。
	ClientConfig WebhookClientConfig `json:"clientConfig" protobuf:"bytes,2,opt,name=clientConfig"`

	// Rules describes what operations on what resources/subresources the webhook cares about.
	// The webhook cares about an operation if it matches _any_ Rule.
	// However, in order to prevent ValidatingAdmissionWebhooks and MutatingAdmissionWebhooks
	// from putting the cluster in a state which cannot be recovered from without completely
	// disabling the plugin, ValidatingAdmissionWebhooks and MutatingAdmissionWebhooks are never called
	// on admission requests for ValidatingWebhookConfiguration and MutatingWebhookConfiguration objects.
	// Rules 描述了webhook关心的资源/子资源上的操作。 如果匹配任何规则，则webhook关心操作。 但是，为了防止ValidatingAdmissionWebhooks和MutatingAdmissionWebhooks将集群置于无法从中恢复的状态，而不完全禁用插件，ValidatingAdmissionWebhooks和MutatingAdmissionWebhooks永远不会在针对ValidatingWebhookConfiguration和MutatingWebhookConfiguration对象的准入请求上调用。
	Rules []RuleWithOperations `json:"rules,omitempty" protobuf:"bytes,3,rep,name=rules"`

	// FailurePolicy defines how unrecognized errors from the admission endpoint are handled -
	// allowed values are Ignore or Fail. Defaults to Fail.
	// +optional
	// FailurePolicy 定义了从admission端点未识别的错误如何处理-允许的值为Ignore或Fail。 默认为Fail。
	FailurePolicy *FailurePolicyType `json:"failurePolicy,omitempty" protobuf:"bytes,4,opt,name=failurePolicy,casttype=FailurePolicyType"`

	// matchPolicy defines how the "rules" list is used to match incoming requests.
	// Allowed values are "Exact" or "Equivalent".
	//
	// - Exact: match a request only if it exactly matches a specified rule.
	// For example, if deployments can be modified via apps/v1, apps/v1beta1, and extensions/v1beta1,
	// but "rules" only included `apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]`,
	// a request to apps/v1beta1 or extensions/v1beta1 would not be sent to the webhook.
	//
	// - Equivalent: match a request if modifies a resource listed in rules, even via another API group or version.
	// For example, if deployments can be modified via apps/v1, apps/v1beta1, and extensions/v1beta1,
	// and "rules" only included `apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]`,
	// a request to apps/v1beta1 or extensions/v1beta1 would be converted to apps/v1 and sent to the webhook.
	//
	// Defaults to "Equivalent"
	// +optional
	// matchPolicy 定义了“rules”列表如何用于匹配传入请求。 允许的值为“Exact”或“Equivalent”。
	// -Exact：仅当请求与指定规则完全匹配时才匹配请求。 例如，如果部署可以通过apps/v1，apps/v1beta1和extensions/v1beta1进行修改，但“rules”仅包括`apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]`，则对apps/v1beta1或extensions/v1beta1的请求将不会发送到网络钩子。
	// -Equivalent：如果通过另一个API组或版本修改“rules”中列出的资源，则匹配请求。 例如，如果部署可以通过apps/v1，apps/v1beta1和extensions/v1beta1进行修改，并且“rules”仅包括`apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]`，则对apps/v1beta1或extensions/v1beta1的请求将转换为apps/v1并发送到网络钩子。
	MatchPolicy *MatchPolicyType `json:"matchPolicy,omitempty" protobuf:"bytes,9,opt,name=matchPolicy,casttype=MatchPolicyType"`

	// NamespaceSelector decides whether to run the webhook on an object based
	// on whether the namespace for that object matches the selector. If the
	// object itself is a namespace, the matching is performed on
	// object.metadata.labels. If the object is another cluster scoped resource,
	// it never skips the webhook.
	//
	// For example, to run the webhook on any objects whose namespace is not
	// associated with "runlevel" of "0" or "1";  you will set the selector as
	// follows:
	// "namespaceSelector": {
	//   "matchExpressions": [
	//     {
	//       "key": "runlevel",
	//       "operator": "NotIn",
	//       "values": [
	//         "0",
	//         "1"
	//       ]
	//     }
	//   ]
	// }
	//
	// If instead you want to only run the webhook on any objects whose
	// namespace is associated with the "environment" of "prod" or "staging";
	// you will set the selector as follows:
	// "namespaceSelector": {
	//   "matchExpressions": [
	//     {
	//       "key": "environment",
	//       "operator": "In",
	//       "values": [
	//         "prod",
	//         "staging"
	//       ]
	//     }
	//   ]
	// }
	//
	// See
	// https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/
	// for more examples of label selectors.
	//
	// Default to the empty LabelSelector, which matches everything.
	// +optional
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty" protobuf:"bytes,5,opt,name=namespaceSelector"`

	// ObjectSelector decides whether to run the webhook based on if the
	// object has matching labels. objectSelector is evaluated against both
	// the oldObject and newObject that would be sent to the webhook, and
	// is considered to match if either object matches the selector. A null
	// object (oldObject in the case of create, or newObject in the case of
	// delete) or an object that cannot have labels (like a
	// DeploymentRollback or a PodProxyOptions object) is not considered to
	// match.
	// Use the object selector only if the webhook is opt-in, because end
	// users may skip the admission webhook by setting the labels.
	// Default to the empty LabelSelector, which matches everything.
	// +optional
	ObjectSelector *metav1.LabelSelector `json:"objectSelector,omitempty" protobuf:"bytes,11,opt,name=objectSelector"`

	// SideEffects states whether this webhook has side effects.
	// Acceptable values are: None, NoneOnDryRun (webhooks created via v1beta1 may also specify Some or Unknown).
	// Webhooks with side effects MUST implement a reconciliation system, since a request may be
	// rejected by a future step in the admission chain and the side effects therefore need to be undone.
	// Requests with the dryRun attribute will be auto-rejected if they match a webhook with
	// sideEffects == Unknown or Some.
	SideEffects *SideEffectClass `json:"sideEffects" protobuf:"bytes,6,opt,name=sideEffects,casttype=SideEffectClass"`

	// TimeoutSeconds specifies the timeout for this webhook. After the timeout passes,
	// the webhook call will be ignored or the API call will fail based on the
	// failure policy.
	// The timeout value must be between 1 and 30 seconds.
	// Default to 10 seconds.
	// +optional
	TimeoutSeconds *int32 `json:"timeoutSeconds,omitempty" protobuf:"varint,7,opt,name=timeoutSeconds"`

	// AdmissionReviewVersions is an ordered list of preferred `AdmissionReview`
	// versions the Webhook expects. API server will try to use first version in
	// the list which it supports. If none of the versions specified in this list
	// supported by API server, validation will fail for this object.
	// If a persisted webhook configuration specifies allowed versions and does not
	// include any versions known to the API Server, calls to the webhook will fail
	// and be subject to the failure policy.
	AdmissionReviewVersions []string `json:"admissionReviewVersions" protobuf:"bytes,8,rep,name=admissionReviewVersions"`

	// reinvocationPolicy indicates whether this webhook should be called multiple times as part of a single admission evaluation.
	// Allowed values are "Never" and "IfNeeded".
	//
	// Never: the webhook will not be called more than once in a single admission evaluation.
	//
	// IfNeeded: the webhook will be called at least one additional time as part of the admission evaluation
	// if the object being admitted is modified by other admission plugins after the initial webhook call.
	// Webhooks that specify this option *must* be idempotent, able to process objects they previously admitted.
	// Note:
	// * the number of additional invocations is not guaranteed to be exactly one.
	// * if additional invocations result in further modifications to the object, webhooks are not guaranteed to be invoked again.
	// * webhooks that use this option may be reordered to minimize the number of additional invocations.
	// * to validate an object after all mutations are guaranteed complete, use a validating admission webhook instead.
	//
	// Defaults to "Never".
	// +optional
	ReinvocationPolicy *ReinvocationPolicyType `json:"reinvocationPolicy,omitempty" protobuf:"bytes,10,opt,name=reinvocationPolicy,casttype=ReinvocationPolicyType"`
}

// ReinvocationPolicyType specifies what type of policy the admission hook uses.
// +enum
type ReinvocationPolicyType string

const (
	// NeverReinvocationPolicy indicates that the webhook must not be called more than once in a
	// single admission evaluation.
	NeverReinvocationPolicy ReinvocationPolicyType = "Never"
	// IfNeededReinvocationPolicy indicates that the webhook may be called at least one
	// additional time as part of the admission evaluation if the object being admitted is
	// modified by other admission plugins after the initial webhook call.
	IfNeededReinvocationPolicy ReinvocationPolicyType = "IfNeeded"
)

// RuleWithOperations is a tuple of Operations and Resources. It is recommended to make
// sure that all the tuple expansions are valid.
type RuleWithOperations struct {
	// Operations is the operations the admission hook cares about - CREATE, UPDATE, DELETE, CONNECT or *
	// for all of those operations and any future admission operations that are added.
	// If '*' is present, the length of the slice must be one.
	// Required.
	// +listType=atomic
	Operations []OperationType `json:"operations,omitempty" protobuf:"bytes,1,rep,name=operations,casttype=OperationType"`
	// Rule is embedded, it describes other criteria of the rule, like
	// APIGroups, APIVersions, Resources, etc.
	Rule `json:",inline" protobuf:"bytes,2,opt,name=rule"`
}

// OperationType specifies an operation for a request.
// +enum
type OperationType string

// The constants should be kept in sync with those defined in k8s.io/kubernetes/pkg/admission/interface.go.
const (
	OperationAll OperationType = "*"
	Create       OperationType = "CREATE"
	Update       OperationType = "UPDATE"
	Delete       OperationType = "DELETE"
	Connect      OperationType = "CONNECT"
)

// WebhookClientConfig contains the information to make a TLS
// connection with the webhook
// WebhookClientConfig 包含与webhook建立TLS连接的信息
type WebhookClientConfig struct {
	// `url` gives the location of the webhook, in standard URL form
	// (`scheme://host:port/path`). Exactly one of `url` or `service`
	// must be specified.
	//
	// The `host` should not refer to a service running in the cluster; use
	// the `service` field instead. The host might be resolved via external
	// DNS in some apiservers (e.g., `kube-apiserver` cannot resolve
	// in-cluster DNS as that would be a layering violation). `host` may
	// also be an IP address.
	//
	// Please note that using `localhost` or `127.0.0.1` as a `host` is
	// risky unless you take great care to run this webhook on all hosts
	// which run an apiserver which might need to make calls to this
	// webhook. Such installs are likely to be non-portable, i.e., not easy
	// to turn up in a new cluster.
	//
	// The scheme must be "https"; the URL must begin with "https://".
	//
	// A path is optional, and if present may be any string permissible in
	// a URL. You may use the path to pass an arbitrary string to the
	// webhook, for example, a cluster identifier.
	//
	// Attempting to use a user or basic auth e.g. "user:password@" is not
	// allowed. Fragments ("#...") and query parameters ("?...") are not
	// allowed, either.
	//
	// +optional
	URL *string `json:"url,omitempty" protobuf:"bytes,3,opt,name=url"`

	// `service` is a reference to the service for this webhook. Either
	// `service` or `url` must be specified.
	//
	// If the webhook is running within the cluster, then you should use `service`.
	//
	// +optional
	Service *ServiceReference `json:"service,omitempty" protobuf:"bytes,1,opt,name=service"`

	// `caBundle` is a PEM encoded CA bundle which will be used to validate the webhook's server certificate.
	// If unspecified, system trust roots on the apiserver are used.
	// +optional
	CABundle []byte `json:"caBundle,omitempty" protobuf:"bytes,2,opt,name=caBundle"`
}

// ServiceReference holds a reference to Service.legacy.k8s.io
type ServiceReference struct {
	// `namespace` is the namespace of the service.
	// Required
	Namespace string `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`
	// `name` is the name of the service.
	// Required
	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`

	// `path` is an optional URL path which will be sent in any request to
	// this service.
	// +optional
	Path *string `json:"path,omitempty" protobuf:"bytes,3,opt,name=path"`

	// If specified, the port on the service that hosting webhook.
	// Default to 443 for backward compatibility.
	// `port` should be a valid port number (1-65535, inclusive).
	// +optional
	Port *int32 `json:"port,omitempty" protobuf:"varint,4,opt,name=port"`
}
