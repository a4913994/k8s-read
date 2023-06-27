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
	authenticationv1 "k8s.io/api/authentication/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AdmissionReview describes an admission review request/response.
// AdmissionReview 描述了一个准入审查请求/响应
type AdmissionReview struct {
	metav1.TypeMeta `json:",inline"`
	// Request describes the attributes for the admission request.
	// +optional
	// Request 描述了准入请求的属性
	Request *AdmissionRequest `json:"request,omitempty" protobuf:"bytes,1,opt,name=request"`
	// Response describes the attributes for the admission response.
	// +optional
	// Response 描述了准入响应的属性
	Response *AdmissionResponse `json:"response,omitempty" protobuf:"bytes,2,opt,name=response"`
}

// AdmissionRequest describes the admission.Attributes for the admission request.
// AdmissionRequest 描述了准入请求的准入属性
type AdmissionRequest struct {
	// UID is an identifier for the individual request/response. It allows us to distinguish instances of requests which are
	// otherwise identical (parallel requests, requests when earlier requests did not modify etc)
	// The UID is meant to track the round trip (request/response) between the KAS and the WebHook, not the user request.
	// It is suitable for correlating log entries between the webhook and apiserver, for either auditing or debugging.
	// UID 是单个请求/响应的标识符。它允许我们区分其他相同的请求实例（并行请求，请求在之前的请求未修改等）
	// UID 旨在跟踪 KAS 和 WebHook 之间的往返（请求/响应），而不是用户请求。
	// 它适合在 Webhook 和 apiserver 之间关联日志条目，用于审计或调试。
	UID types.UID `json:"uid" protobuf:"bytes,1,opt,name=uid"`
	// Kind is the fully-qualified type of object being submitted (for example, v1.Pod or autoscaling.v1.Scale)
	// Kind 是正在提交的对象的完全限定类型（例如，v1.Pod 或 autoscaling.v1.Scale）
	Kind metav1.GroupVersionKind `json:"kind" protobuf:"bytes,2,opt,name=kind"`
	// Resource is the fully-qualified resource being requested (for example, v1.pods)
	// Resource 是正在请求的完全限定资源（例如，v1.pods）
	Resource metav1.GroupVersionResource `json:"resource" protobuf:"bytes,3,opt,name=resource"`
	// SubResource is the subresource being requested, if any (for example, "status" or "scale")
	// +optional
	// SubResource 是正在请求的子资源（如果有的话，例如“status”或“scale”）
	SubResource string `json:"subResource,omitempty" protobuf:"bytes,4,opt,name=subResource"`

	// RequestKind is the fully-qualified type of the original API request (for example, v1.Pod or autoscaling.v1.Scale).
	// If this is specified and differs from the value in "kind", an equivalent match and conversion was performed.
	//
	// For example, if deployments can be modified via apps/v1 and apps/v1beta1, and a webhook registered a rule of
	// `apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]` and `matchPolicy: Equivalent`,
	// an API request to apps/v1beta1 deployments would be converted and sent to the webhook
	// with `kind: {group:"apps", version:"v1", kind:"Deployment"}` (matching the rule the webhook registered for),
	// and `requestKind: {group:"apps", version:"v1beta1", kind:"Deployment"}` (indicating the kind of the original API request).
	//
	// See documentation for the "matchPolicy" field in the webhook configuration type for more details.
	// +optional
	// RequestKind 是原始 API 请求的完全限定类型（例如，v1.Pod 或 autoscaling.v1.Scale）。 如果指定了此值并且与“kind”中的值不同，则执行等效的匹配和转换。
	// 例如，如果可以通过 apps/v1 和 apps/v1beta1 修改部署，并且 Webhook 注册了一个规则 `apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]` 和 `matchPolicy: Equivalent`，
	// 则对 apps/v1beta1 部署的 API 请求将被转换并发送到 Webhook，其中 `kind: {group:"apps", version:"v1", kind:"Deployment"}`（与 Webhook 注册的规则匹配），
	// 并且 `requestKind: {group:"apps", version:"v1beta1", kind:"Deployment"}`（指示原始 API 请求的类型）。
	// 有关详细信息，请参阅 Webhook 配置类型中的“matchPolicy”字段的文档。
	RequestKind *metav1.GroupVersionKind `json:"requestKind,omitempty" protobuf:"bytes,13,opt,name=requestKind"`
	// RequestResource is the fully-qualified resource of the original API request (for example, v1.pods).
	// If this is specified and differs from the value in "resource", an equivalent match and conversion was performed.
	//
	// For example, if deployments can be modified via apps/v1 and apps/v1beta1, and a webhook registered a rule of
	// `apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]` and `matchPolicy: Equivalent`,
	// an API request to apps/v1beta1 deployments would be converted and sent to the webhook
	// with `resource: {group:"apps", version:"v1", resource:"deployments"}` (matching the resource the webhook registered for),
	// and `requestResource: {group:"apps", version:"v1beta1", resource:"deployments"}` (indicating the resource of the original API request).
	//
	// See documentation for the "matchPolicy" field in the webhook configuration type.
	// +optional
	// RequestResource 是原始 API 请求的完全限定资源（例如，v1.pods）。 如果指定了此值并且与“resource”中的值不同，则执行等效的匹配和转换。
	// 例如，如果可以通过 apps/v1 和 apps/v1beta1 修改部署，并且 Webhook 注册了一个规则 `apiGroups:["apps"], apiVersions:["v1"], resources: ["deployments"]` 和 `matchPolicy: Equivalent`，
	// 则对 apps/v1beta1 部署的 API 请求将被转换并发送到 Webhook，其中 `resource: {group:"apps", version:"v1", resource:"deployments"}`（与 Webhook 注册的资源匹配），
	// 并且 `requestResource: {group:"apps", version:"v1beta1", resource:"deployments"}`（指示原始 API 请求的资源）。
	// 有关详细信息，请参阅 Webhook 配置类型中的“matchPolicy”字段的文档。
	RequestResource *metav1.GroupVersionResource `json:"requestResource,omitempty" protobuf:"bytes,14,opt,name=requestResource"`
	// RequestSubResource is the name of the subresource of the original API request, if any (for example, "status" or "scale")
	// If this is specified and differs from the value in "subResource", an equivalent match and conversion was performed.
	// See documentation for the "matchPolicy" field in the webhook configuration type.
	// +optional
	// RequestSubResource 是原始 API 请求的子资源的名称（如果有的话，例如“status”或“scale”）。 如果指定了此值并且与“subResource”中的值不同，则执行等效的匹配和转换。
	// 有关详细信息，请参阅 Webhook 配置类型中的“matchPolicy”字段的文档。
	RequestSubResource string `json:"requestSubResource,omitempty" protobuf:"bytes,15,opt,name=requestSubResource"`

	// Name is the name of the object as presented in the request.  On a CREATE operation, the client may omit name and
	// rely on the server to generate the name.  If that is the case, this field will contain an empty string.
	// +optional
	// Name 是以请求中呈现的方式提供的对象的名称。 在创建操作中，客户端可以省略名称并依赖服务器生成名称。 如果是这种情况，此字段将包含一个空字符串。
	Name string `json:"name,omitempty" protobuf:"bytes,5,opt,name=name"`
	// Namespace is the namespace associated with the request (if any).
	// +optional
	// Namespace 与请求相关联的名称空间（如果有的话）。
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,6,opt,name=namespace"`
	// Operation is the operation being performed. This may be different than the operation
	// requested. e.g. a patch can result in either a CREATE or UPDATE Operation.
	// Operation 是正在执行的操作。 这可能与请求的操作不同。 例如，补丁可能导致 CREATE 或 UPDATE 操作。
	Operation Operation `json:"operation" protobuf:"bytes,7,opt,name=operation"`
	// UserInfo is information about the requesting user
	// UserInfo 是有关请求用户的信息
	UserInfo authenticationv1.UserInfo `json:"userInfo" protobuf:"bytes,8,opt,name=userInfo"`
	// Object is the object from the incoming request.
	// +optional
	// Object 是传入请求的对象。
	Object runtime.RawExtension `json:"object,omitempty" protobuf:"bytes,9,opt,name=object"`
	// OldObject is the existing object. Only populated for DELETE and UPDATE requests.
	// +optional
	// OldObject 是现有对象。 仅适用于 DELETE 和 UPDATE 请求。
	OldObject runtime.RawExtension `json:"oldObject,omitempty" protobuf:"bytes,10,opt,name=oldObject"`
	// DryRun indicates that modifications will definitely not be persisted for this request.
	// Defaults to false.
	// +optional
	// DryRun 表示此请求的修改绝对不会被持久化。 默认为 false。
	DryRun *bool `json:"dryRun,omitempty" protobuf:"varint,11,opt,name=dryRun"`
	// Options is the operation option structure of the operation being performed.
	// e.g. `meta.k8s.io/v1.DeleteOptions` or `meta.k8s.io/v1.CreateOptions`. This may be
	// different than the options the caller provided. e.g. for a patch request the performed
	// Operation might be a CREATE, in which case the Options will a
	// `meta.k8s.io/v1.CreateOptions` even though the caller provided `meta.k8s.io/v1.PatchOptions`.
	// +optional
	// Options 是正在执行的操作的操作选项结构。 例如，`meta.k8s.io/v1.DeleteOptions` 或 `meta.k8s.io/v1.CreateOptions`。
	// 这可能与调用者提供的选项不同。 例如，对于补丁请求，执行的操作可能是 CREATE，在这种情况下，选项将是 `meta.k8s.io/v1.CreateOptions`，即使调用者提供了 `meta.k8s.io/v1.PatchOptions`。
	Options runtime.RawExtension `json:"options,omitempty" protobuf:"bytes,12,opt,name=options"`
}

// AdmissionResponse describes an admission response.
// AdmissionResponse描述了一个admission响应。
type AdmissionResponse struct {
	// UID is an identifier for the individual request/response.
	// This must be copied over from the corresponding AdmissionRequest.
	// UID 是单个请求/响应的标识符。 必须从相应的 AdmissionRequest 中复制。
	UID types.UID `json:"uid" protobuf:"bytes,1,opt,name=uid"`

	// Allowed indicates whether or not the admission request was permitted.
	// Allowed 表示是否允许admission请求。
	Allowed bool `json:"allowed" protobuf:"varint,2,opt,name=allowed"`

	// Result contains extra details into why an admission request was denied.
	// This field IS NOT consulted in any way if "Allowed" is "true".
	// +optional
	// Result 包含有关为什么拒绝admission请求的额外详细信息。 如果“Allowed”为“true”，则不会以任何方式考虑此字段。
	Result *metav1.Status `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`

	// The patch body. Currently we only support "JSONPatch" which implements RFC 6902.
	// +optional
	// 补丁正文。 目前我们只支持“JSONPatch”，它实现了 RFC 6902。
	Patch []byte `json:"patch,omitempty" protobuf:"bytes,4,opt,name=patch"`

	// The type of Patch. Currently we only allow "JSONPatch".
	// +optional
	// 补丁的类型。 目前我们只允许“JSONPatch”。
	PatchType *PatchType `json:"patchType,omitempty" protobuf:"bytes,5,opt,name=patchType"`

	// AuditAnnotations is an unstructured key value map set by remote admission controller (e.g. error=image-blacklisted).
	// MutatingAdmissionWebhook and ValidatingAdmissionWebhook admission controller will prefix the keys with
	// admission webhook name (e.g. imagepolicy.example.com/error=image-blacklisted). AuditAnnotations will be provided by
	// the admission webhook to add additional context to the audit log for this request.
	// +optional
	// AuditAnnotations 是由远程admission控制器（例如error = image-blacklisted）设置的非结构化键值映射。
	// MutatingAdmissionWebhook 和 ValidatingAdmissionWebhook admission 控制器将使用 admission webhook 名称为键添加前缀（例如 imagepolicy.example.com/error = image-blacklisted）。
	// AuditAnnotations 将由 admission webhook 提供，以便为此请求的审计日志添加其他上下文。
	AuditAnnotations map[string]string `json:"auditAnnotations,omitempty" protobuf:"bytes,6,opt,name=auditAnnotations"`

	// warnings is a list of warning messages to return to the requesting API client.
	// Warning messages describe a problem the client making the API request should correct or be aware of.
	// Limit warnings to 120 characters if possible.
	// Warnings over 256 characters and large numbers of warnings may be truncated.
	// +optional
	// warnings 是要返回给请求 API 客户端的警告消息列表。 警告消息描述了 API 请求的客户端应该纠正或注意的问题。 如果可能，请将警告限制为 120 个字符。
	// 超过 256 个字符和大量警告的警告可能会被截断。
	Warnings []string `json:"warnings,omitempty" protobuf:"bytes,7,rep,name=warnings"`
}

// PatchType is the type of patch being used to represent the mutated object
// PatchType 是用于表示变异对象的补丁类型
type PatchType string

// PatchType constants.
const (
	PatchTypeJSONPatch PatchType = "JSONPatch"
)

// Operation is the type of resource operation being checked for admission control
type Operation string

// Operation constants
const (
	Create  Operation = "CREATE"
	Update  Operation = "UPDATE"
	Delete  Operation = "DELETE"
	Connect Operation = "CONNECT"
)
