/*
Copyright 2017 The Kubernetes Authors.

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

package audit

import (
	authnv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// Header keys used by the audit system.
const (
	// Header to hold the audit ID as the request is propagated through the serving hierarchy. The
	// Audit-ID header should be set by the first server to receive the request (e.g. the federation
	// server or kube-aggregator).
	//
	// Audit ID is also returned to client by http response header.
	// It's not guaranteed Audit-Id http header is sent for all requests. When kube-apiserver didn't
	// audit the events according to the audit policy, no Audit-ID is returned. Also, for request to
	// pods/exec, pods/attach, pods/proxy, kube-apiserver works like a proxy and redirect the request
	// to kubelet node, users will only get http headers sent from kubelet node, so no Audit-ID is
	// sent when users run command like "kubectl exec" or "kubectl attach".
	HeaderAuditID = "Audit-ID"
)

// Level defines the amount of information logged during auditing
// Level 定义审计期间记录的信息量
type Level string

// Valid audit levels
const (
	// LevelNone disables auditing
	LevelNone Level = "None"
	// LevelMetadata provides the basic level of auditing.
	LevelMetadata Level = "Metadata"
	// LevelRequest provides Metadata level of auditing, and additionally
	// logs the request object (does not apply for non-resource requests).
	LevelRequest Level = "Request"
	// LevelRequestResponse provides Request level of auditing, and additionally
	// logs the response object (does not apply for non-resource requests).
	LevelRequestResponse Level = "RequestResponse"
)

// Stage defines the stages in request handling that audit events may be generated.
// Stage 定义了可能生成审计事件的请求处理阶段。
type Stage string

// Valid audit stages.
const (
	// The stage for events generated as soon as the audit handler receives the request, and before it
	// is delegated down the handler chain.
	StageRequestReceived Stage = "RequestReceived"
	// The stage for events generated once the response headers are sent, but before the response body
	// is sent. This stage is only generated for long-running requests (e.g. watch).
	StageResponseStarted Stage = "ResponseStarted"
	// The stage for events generated once the response body has been completed, and no more bytes
	// will be sent.
	StageResponseComplete Stage = "ResponseComplete"
	// The stage for events generated when a panic occurred.
	StagePanic Stage = "Panic"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Event captures all the information that can be included in an API audit log.
// Event 捕获可以包含在 API 审计日志中的所有信息。
type Event struct {
	metav1.TypeMeta

	// AuditLevel at which event was generated
	// 生成事件的 AuditLevel
	Level Level

	// Unique audit ID, generated for each request.
	AuditID types.UID
	// Stage of the request handling when this event instance was generated.
	// 生成此事件实例时请求处理的阶段。
	Stage Stage

	// RequestURI is the request URI as sent by the client to a server.
	// RequestURI 是客户端发送给服务器的请求 URI。
	RequestURI string
	// Verb is the kubernetes verb associated with the request.
	// For non-resource requests, this is the lower-cased HTTP method.
	// Verb 是与请求关联的 kubernetes 动词。对于非资源请求，这是小写的 HTTP 方法。
	Verb string
	// Authenticated user information.
	// 经过身份验证的用户信息。
	User authnv1.UserInfo
	// Impersonated user information.
	// +optional
	// 模拟用户信息。
	ImpersonatedUser *authnv1.UserInfo
	// Source IPs, from where the request originated and intermediate proxies.
	// The source IPs are listed from (in order):
	// 1. X-Forwarded-For request header IPs
	// 2. X-Real-Ip header, if not present in the X-Forwarded-For list
	// 3. The remote address for the connection, if it doesn't match the last
	//    IP in the list up to here (X-Forwarded-For or X-Real-Ip).
	// Note: All but the last IP can be arbitrarily set by the client.
	// +optional
	SourceIPs []string
	// UserAgent records the user agent string reported by the client.
	// Note that the UserAgent is provided by the client, and must not be trusted.
	// +optional
	UserAgent string
	// Object reference this request is targeted at.
	// Does not apply for List-type requests, or non-resource requests.
	// +optional
	ObjectRef *ObjectReference
	// The response status, populated even when the ResponseObject is not a Status type.
	// For successful responses, this will only include the Code. For non-status type
	// error responses, this will be auto-populated with the error Message.
	// +optional
	ResponseStatus *metav1.Status

	// API object from the request, in JSON format. The RequestObject is recorded as-is in the request
	// (possibly re-encoded as JSON), prior to version conversion, defaulting, admission or
	// merging. It is an external versioned object type, and may not be a valid object on its own.
	// Omitted for non-resource requests.  Only logged at Request Level and higher.
	// +optional
	// 来自请求的 API 对象，采用 JSON 格式。
	// RequestObject 在版本转换、默认、接纳或合并之前按原样记录在请求中（可能重新编码为 JSON）。
	// 它是一个外部版本控制的对象类型，并且它本身可能不是一个有效的对象。为非资源请求省略。仅在请求级别及更高级别记录。
	RequestObject *runtime.Unknown
	// API object returned in the response, in JSON. The ResponseObject is recorded after conversion
	// to the external type, and serialized as JSON. Omitted for non-resource requests.  Only logged
	// at Response Level.
	// +optional
	// 响应中以 JSON 格式返回的 API 对象。
	// ResponseObject 转换为外部类型后记录，并序列化为 JSON。为非资源请求省略。仅在响应级别记录。
	ResponseObject *runtime.Unknown

	// Time the request reached the apiserver.
	// 请求到达 apiserver 的时间。
	RequestReceivedTimestamp metav1.MicroTime
	// Time the request reached current audit stage.
	// 请求到达当前审核阶段的时间。
	StageTimestamp metav1.MicroTime

	// Annotations is an unstructured key value map stored with an audit event that may be set by
	// plugins invoked in the request serving chain, including authentication, authorization and
	// admission plugins. Note that these annotations are for the audit event, and do not correspond
	// to the metadata.annotations of the submitted object. Keys should uniquely identify the informing
	// component to avoid name collisions (e.g. podsecuritypolicy.admission.k8s.io/policy). Values
	// should be short. Annotations are included in the Metadata level.
	// +optional
	Annotations map[string]string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventList is a list of audit Events.
type EventList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	Items []Event
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Policy defines the configuration of audit logging, and the rules for how different request
// categories are logged.
// Policy 定义审计日志的配置，以及不同请求类别如何记录的规则
type Policy struct {
	metav1.TypeMeta
	// ObjectMeta is included for interoperability with API infrastructure.
	// +optional
	metav1.ObjectMeta

	// Rules specify the audit Level a request should be recorded at.
	// A request may match multiple rules, in which case the FIRST matching rule is used.
	// The default audit level is None, but can be overridden by a catch-all rule at the end of the list.
	// PolicyRules are strictly ordered.
	// 规则指定应记录请求的审核级别。一个请求可能匹配多个规则，在这种情况下使用第一个匹配规则。
	// 默认审计级别是无，但可以被列表末尾的包罗万象的规则覆盖。 PolicyRules 是严格有序的。
	Rules []PolicyRule

	// OmitStages is a list of stages for which no events are created. Note that this can also
	// be specified per rule in which case the union of both are omitted.
	// OmitStages 是未创建事件的阶段列表。请注意，这也可以按规则指定，在这种情况下，两者的并集将被省略。
	// +optional
	OmitStages []Stage

	// OmitManagedFields indicates whether to omit the managed fields of the request
	// and response bodies from being written to the API audit log.
	// This is used as a global default - a value of 'true' will omit the managed fileds,
	// otherwise the managed fields will be included in the API audit log.
	// Note that this can also be specified per rule in which case the value specified
	// in a rule will override the global default.
	// OmitManagedFields 指示是否忽略请求和响应主体的托管字段写入 API 审计日志。这用作全局默认值——“true”值将忽略托管字段，
	// 否则托管字段将包含在 API 审计日志中。请注意，这也可以按规则指定，在这种情况下，规则中指定的值将覆盖全局默认值。
	// +optional
	OmitManagedFields bool
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PolicyList is a list of audit Policies.
type PolicyList struct {
	metav1.TypeMeta
	// +optional
	metav1.ListMeta

	Items []Policy
}

// PolicyRule maps requests based off metadata to an audit Level.
// Requests must match the rules of every field (an intersection of rules).
// PolicyRule 将基于元数据的请求映射到审计级别。请求必须匹配每个字段的规则（规则的交集）。
type PolicyRule struct {
	// The Level that requests matching this rule are recorded at.
	Level Level

	// The users (by authenticated user name) this rule applies to.
	// An empty list implies every user.
	// +optional
	Users []string
	// The user groups this rule applies to. A user is considered matching
	// if it is a member of any of the UserGroups.
	// An empty list implies every user group.
	// +optional
	UserGroups []string

	// The verbs that match this rule.
	// An empty list implies every verb.
	// +optional
	Verbs []string

	// Rules can apply to API resources (such as "pods" or "secrets"),
	// non-resource URL paths (such as "/api"), or neither, but not both.
	// If neither is specified, the rule is treated as a default for all URLs.

	// Resources that this rule matches. An empty list implies all kinds in all API groups.
	// +optional
	Resources []GroupResources
	// Namespaces that this rule matches.
	// The empty string "" matches non-namespaced resources.
	// An empty list implies every namespace.
	// +optional
	Namespaces []string

	// NonResourceURLs is a set of URL paths that should be audited.
	// *s are allowed, but only as the full, final step in the path.
	// Examples:
	//  "/metrics" - Log requests for apiserver metrics
	//  "/healthz*" - Log all health checks
	// +optional
	NonResourceURLs []string

	// OmitStages is a list of stages for which no events are created. Note that this can also
	// be specified policy wide in which case the union of both are omitted.
	// An empty list means no restrictions will apply.
	// +optional
	OmitStages []Stage

	// OmitManagedFields indicates whether to omit the managed fields of the request
	// and response bodies from being written to the API audit log.
	// - a value of 'true' will drop the managed fields from the API audit log
	// - a value of 'false' indicates that the managed fileds should be included
	//   in the API audit log
	// Note that the value, if specified, in this rule will override the global default
	// If a value is not specified then the global default specified in
	// Policy.OmitManagedFields will stand.
	// +optional
	OmitManagedFields *bool
}

// GroupResources represents resource kinds in an API group.
type GroupResources struct {
	// Group is the name of the API group that contains the resources.
	// The empty string represents the core API group.
	// +optional
	Group string
	// Resources is a list of resources this rule applies to.
	//
	// For example:
	// 'pods' matches pods.
	// 'pods/log' matches the log subresource of pods.
	// '*' matches all resources and their subresources.
	// 'pods/*' matches all subresources of pods.
	// '*/scale' matches all scale subresources.
	//
	// If wildcard is present, the validation rule will ensure resources do not
	// overlap with each other.
	//
	// An empty list implies all resources and subresources in this API groups apply.
	// +optional
	Resources []string
	// ResourceNames is a list of resource instance names that the policy matches.
	// Using this field requires Resources to be specified.
	// An empty list implies that every instance of the resource is matched.
	// +optional
	ResourceNames []string
}

// ObjectReference contains enough information to let you inspect or modify the referred object.
type ObjectReference struct {
	// +optional
	Resource string
	// +optional
	Namespace string
	// +optional
	Name string
	// +optional
	UID types.UID
	// APIGroup is the name of the API group that contains the referred object.
	// The empty string represents the core API group.
	// +optional
	APIGroup string
	// APIVersion is the version of the API group that contains the referred object.
	// +optional
	APIVersion string
	// +optional
	ResourceVersion string
	// +optional
	Subresource string
}
