/*
Copyright 2014 The Kubernetes Authors.

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

package admission

import (
	"context"
	"io"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	auditinternal "k8s.io/apiserver/pkg/apis/audit"
	"k8s.io/apiserver/pkg/authentication/user"
)

// Attributes is an interface used by AdmissionController to get information about a request
// that is used to make an admission decision.
// Attributes 是 AdmissionController 用来获取有关用于做出授权决定的请求的信息的接口。
type Attributes interface {
	// GetName returns the name of the object as presented in the request.  On a CREATE operation, the client
	// may omit name and rely on the server to generate the name.  If that is the case, this method will return
	// the empty string
	// GetName 返回请求中对象的名称。在 CREATE 操作中，客户端可能会省略名称，并依赖服务器生成名称。如果是这种情况，此方法将返回空字符串。
	GetName() string
	// GetNamespace is the namespace associated with the request (if any)
	// GetNamespace 返回与请求关联的命名空间（如果有的话）。
	GetNamespace() string
	// GetResource is the name of the resource being requested.  This is not the kind.  For example: pods
	// GetResource 返回正在请求的资源的名称。这不是种类。例如：pods。
	GetResource() schema.GroupVersionResource
	// GetSubresource is the name of the subresource being requested.  This is a different resource, scoped to the parent resource, but it may have a different kind.
	// For instance, /pods has the resource "pods" and the kind "Pod", while /pods/foo/status has the resource "pods", the sub resource "status", and the kind "Pod"
	// (because status operates on pods). The binding resource for a pod though may be /pods/foo/binding, which has resource "pods", subresource "binding", and kind "Binding".
	// GetSubresource 返回正在请求的子资源的名称。这是一个不同的资源，范围限于父资源，但它可能具有不同的种类。
	// 例如，/pods 具有资源“pods”和种类“Pod”，而 /pods/foo/status 具有资源“pods”，子资源“status”和种类“Pod”
	// （因为状态在pods上操作）。但是，pod的绑定资源可能是 /pods/foo/binding，它具有资源“pods”，子资源“binding”和种类“Binding”。
	GetSubresource() string
	// GetOperation is the operation being performed
	// GetOperation 返回正在执行的操作。
	GetOperation() Operation
	// GetOperationOptions is the options for the operation being performed
	// GetOperationOptions 返回正在执行的操作的选项。
	GetOperationOptions() runtime.Object
	// IsDryRun indicates that modifications will definitely not be persisted for this request. This is to prevent
	// admission controllers with side effects and a method of reconciliation from being overwhelmed.
	// However, a value of false for this does not mean that the modification will be persisted, because it
	// could still be rejected by a subsequent validation step.
	// IsDryRun 表示此请求的修改绝对不会被持久化。这是为了防止具有副作用和一种协调方法的准入控制器被压倒。
	// 但是，对于这个值的 false 并不意味着修改将被持久化，因为它仍然可能被后续的验证步骤拒绝。
	IsDryRun() bool
	// GetObject is the object from the incoming request prior to default values being applied
	// GetObject 返回应用默认值之前从传入请求中获取的对象。
	GetObject() runtime.Object
	// GetOldObject is the existing object. Only populated for UPDATE and DELETE requests.
	// GetOldObject 返回现有对象。仅用于 UPDATE 和 DELETE 请求。
	GetOldObject() runtime.Object
	// GetKind is the type of object being manipulated.  For example: Pod
	// GetKind 返回正在操作的对象的类型。例如：Pod。
	GetKind() schema.GroupVersionKind
	// GetUserInfo is information about the requesting user
	// GetUserInfo 返回有关请求用户的信息。
	GetUserInfo() user.Info

	// AddAnnotation sets annotation according to key-value pair. The key should be qualified, e.g., podsecuritypolicy.admission.k8s.io/admit-policy, where
	// "podsecuritypolicy" is the name of the plugin, "admission.k8s.io" is the name of the organization, "admit-policy" is the key name.
	// An error is returned if the format of key is invalid. When trying to overwrite annotation with a new value, an error is returned.
	// Both ValidationInterface and MutationInterface are allowed to add Annotations.
	// By default, an annotation gets logged into audit event if the request's audit level is greater or
	// equal to Metadata.
	// AddAnnotation 根据键值对设置注释。键应该是合格的，例如，podsecuritypolicy.admission.k8s.io/admit-policy，其中
	// “podsecuritypolicy”是插件的名称，“admission.k8s.io”是组织的名称，“admit-policy”是键名。
	// 如果键的格式无效，则返回错误。当尝试使用新值覆盖注释时，将返回错误。
	// ValidationInterface 和 MutationInterface 都允许添加注释。
	// 默认情况下，如果请求的审计级别大于或等于元数据，则注释将记录到审计事件中。
	AddAnnotation(key, value string) error

	// AddAnnotationWithLevel sets annotation according to key-value pair with additional intended audit level.
	// An Annotation gets logged into audit event if the request's audit level is greater or equal to the
	// intended audit level.
	// AddAnnotationWithLevel 根据键值对设置注释，并附加预期的审计级别。
	// 如果请求的审计级别大于或等于预期的审计级别，则注释将记录到审计事件中。
	AddAnnotationWithLevel(key, value string, level auditinternal.Level) error

	// GetReinvocationContext tracks the admission request information relevant to the re-invocation policy.
	// GetReinvocationContext 跟踪与重新调用策略相关的准入请求信息。
	GetReinvocationContext() ReinvocationContext
}

// ObjectInterfaces is an interface used by AdmissionController to get object interfaces
// such as Converter or Defaulter. These interfaces are normally coming from Request Scope
// to handle special cases like CRDs.
type ObjectInterfaces interface {
	// GetObjectCreater is the ObjectCreator appropriate for the requested object.
	GetObjectCreater() runtime.ObjectCreater
	// GetObjectTyper is the ObjectTyper appropriate for the requested object.
	GetObjectTyper() runtime.ObjectTyper
	// GetObjectDefaulter is the ObjectDefaulter appropriate for the requested object.
	GetObjectDefaulter() runtime.ObjectDefaulter
	// GetObjectConvertor is the ObjectConvertor appropriate for the requested object.
	GetObjectConvertor() runtime.ObjectConvertor
	// GetEquivalentResourceMapper is the EquivalentResourceMapper appropriate for finding equivalent resources and expected kind for the requested object.
	GetEquivalentResourceMapper() runtime.EquivalentResourceMapper
}

// privateAnnotationsGetter is a private interface which allows users to get annotations from Attributes.
type privateAnnotationsGetter interface {
	getAnnotations(maxLevel auditinternal.Level) map[string]string
}

// AnnotationsGetter allows users to get annotations from Attributes. An alternate Attribute should implement
// this interface.
type AnnotationsGetter interface {
	GetAnnotations(maxLevel auditinternal.Level) map[string]string
}

// ReinvocationContext provides access to the admission related state required to implement the re-invocation policy.
// ReinvocationContext 提供了访问实现重新调用策略所需的准入相关状态的访问。
type ReinvocationContext interface {
	// IsReinvoke returns true if the current admission check is a re-invocation.
	IsReinvoke() bool
	// SetIsReinvoke sets the current admission check as a re-invocation.
	SetIsReinvoke()
	// ShouldReinvoke returns true if any plugin has requested a re-invocation.
	ShouldReinvoke() bool
	// SetShouldReinvoke signals that a re-invocation is desired.
	SetShouldReinvoke()
	// AddValue set a value for a plugin name, possibly overriding a previous value.
	SetValue(plugin string, v interface{})
	// Value reads a value for a webhook.
	Value(plugin string) interface{}
}

// Interface is an abstract, pluggable interface for Admission Control decisions.
// Interface 是一个抽象的、可插拔的接口，用于 Admission Control 决策。
type Interface interface {
	// Handles returns true if this admission controller can handle the given operation
	// where operation can be one of CREATE, UPDATE, DELETE, or CONNECT
	// Handles 返回 true 表示该 admission controller 可以处理给定的操作。
	// operation 可以是 CREATE、UPDATE、DELETE 或 CONNECT。
	Handles(operation Operation) bool
}

type MutationInterface interface {
	Interface

	// Admit makes an admission decision based on the request attributes.
	// Context is used only for timeout/deadline/cancellation and tracing information.
	Admit(ctx context.Context, a Attributes, o ObjectInterfaces) (err error)
}

// ValidationInterface is an abstract, pluggable interface for Admission Control decisions.
type ValidationInterface interface {
	Interface

	// Validate makes an admission decision based on the request attributes.  It is NOT allowed to mutate
	// Context is used only for timeout/deadline/cancellation and tracing information.
	Validate(ctx context.Context, a Attributes, o ObjectInterfaces) (err error)
}

// Operation is the type of resource operation being checked for admission control
type Operation string

// Operation constants
const (
	Create  Operation = "CREATE"
	Update  Operation = "UPDATE"
	Delete  Operation = "DELETE"
	Connect Operation = "CONNECT"
)

// PluginInitializer is used for initialization of shareable resources between admission plugins.
// After initialization the resources have to be set separately
type PluginInitializer interface {
	Initialize(plugin Interface)
}

// InitializationValidator holds ValidateInitialization functions, which are responsible for validation of initialized
// shared resources and should be implemented on admission plugins
type InitializationValidator interface {
	ValidateInitialization() error
}

// ConfigProvider provides a way to get configuration for an admission plugin based on its name
type ConfigProvider interface {
	ConfigFor(pluginName string) (io.Reader, error)
}
