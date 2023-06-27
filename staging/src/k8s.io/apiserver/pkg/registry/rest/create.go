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

package rest

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	genericvalidation "k8s.io/apimachinery/pkg/api/validation"
	"k8s.io/apimachinery/pkg/api/validation/path"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/admission"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/apiserver/pkg/warning"
)

// RESTCreateStrategy defines the minimum validation, accepted input, and
// name generation behavior to create an object that follows Kubernetes
// API conventions.
// RESTCreateStrategy 定义了创建对象时的最小验证、接受的输入和名称生成行为，以遵循 Kubernetes API 约定。
type RESTCreateStrategy interface {
	runtime.ObjectTyper
	// The name generator is used when the standard GenerateName field is set.
	// The NameGenerator will be invoked prior to validation.
	// 设置标准 GenerateName 字段时使用名称生成器。 NameGenerator 将在验证之前被调用。
	names.NameGenerator

	// NamespaceScoped returns true if the object must be within a namespace.
	// NamespaceScoped 如果对象必须在命名空间中，则返回 true。
	NamespaceScoped() bool
	// PrepareForCreate is invoked on create before validation to normalize
	// the object.  For example: remove fields that are not to be persisted,
	// sort order-insensitive list fields, etc.  This should not remove fields
	// whose presence would be considered a validation error.
	//
	// Often implemented as a type check and an initailization or clearing of
	// status. Clear the status because status changes are internal. External
	// callers of an api (users) should not be setting an initial status on
	// newly created objects.
	// PrepareForCreate 在创建之前调用验证以规范化对象。 例如：删除不要持久化的字段、排序不敏感的列表字段等。 这不应该删除其存在将被视为验证错误的字段。
	// 通常作为类型检查和初始化或清除状态的实现。 清除状态是因为状态更改是内部的。 外部调用者的 api（用户）不应该在新创建的对象上设置初始状态。
	PrepareForCreate(ctx context.Context, obj runtime.Object)
	// Validate returns an ErrorList with validation errors or nil.  Validate
	// is invoked after default fields in the object have been filled in
	// before the object is persisted.  This method should not mutate the
	// object.
	// Validate 返回带有验证错误或 nil 的 ErrorList。 在对象持久化之前，Validate 在对象中填充默认字段之后被调用。 此方法不应该改变对象。
	Validate(ctx context.Context, obj runtime.Object) field.ErrorList
	// WarningsOnCreate returns warnings to the client performing a create.
	// WarningsOnCreate is invoked after default fields in the object have been filled in
	// and after Validate has passed, before Canonicalize is called, and the object is persisted.
	// This method must not mutate the object.
	//
	// Be brief; limit warnings to 120 characters if possible.
	// Don't include a "Warning:" prefix in the message (that is added by clients on output).
	// Warnings returned about a specific field should be formatted as "path.to.field: message".
	// For example: `spec.imagePullSecrets[0].name: invalid empty name ""`
	//
	// Use warning messages to describe problems the client making the API request should correct or be aware of.
	// For example:
	// - use of deprecated fields/labels/annotations that will stop working in a future release
	// - use of obsolete fields/labels/annotations that are non-functional
	// - malformed or invalid specifications that prevent successful handling of the submitted object,
	//   but are not rejected by validation for compatibility reasons
	//
	// Warnings should not be returned for fields which cannot be resolved by the caller.
	// For example, do not warn about spec fields in a subresource creation request.
	// WarningsOnCreate 在对象中填充默认字段之后、在 Validate 通过之后、在 Canonicalize 被调用之前、在对象被持久化之前被调用。
	// 此方法不应该改变对象。
	// 请简洁; 如果可能，请将警告限制在 120 个字符以内。
	// 不要在消息中包含“警告：”前缀（客户端在输出时会添加该前缀）。
	// 关于特定字段的警告应以“path.to.field: message”格式返回。
	// 例如：`spec.imagePullSecrets[0].name: invalid empty name ""`
	// 使用警告消息来描述客户端应该更正或注意的 API 请求的问题。
	// 例如：
	// - 将在未来版本中停止工作的已弃用字段/标签/注释
	// - 不可用的过时字段/标签/注释
	// - 阻止成功处理提交的对象的格式错误或无效的规范，但由于兼容性原因而不被验证拒绝
	// 不应返回无法由调用者解决的字段的警告。
	// 例如，不要在子资源创建请求中警告 spec 字段。
	WarningsOnCreate(ctx context.Context, obj runtime.Object) []string
	// Canonicalize allows an object to be mutated into a canonical form. This
	// ensures that code that operates on these objects can rely on the common
	// form for things like comparison.  Canonicalize is invoked after
	// validation has succeeded but before the object has been persisted.
	// This method may mutate the object. Often implemented as a type check or
	// empty method.
	// Canonicalize 允许对象被转换为规范形式。 这确保了操作这些对象的代码可以依赖于比较等内容的公共形式。
	// Canonicalize 在验证成功之后、对象被持久化之前被调用。 此方法可以改变对象。 通常作为类型检查或空方法的实现。
	Canonicalize(obj runtime.Object)
}

// BeforeCreate ensures that common operations for all resources are performed on creation. It only returns
// errors that can be converted to api.Status. It invokes PrepareForCreate, then Validate.
// It returns nil if the object should be created.
// BeforeCreate 确保所有资源的常见操作在创建时执行。 它只返回可以转换为 api.Status 的错误。 它调用 PrepareForCreate，然后调用 Validate。
// 如果对象应该被创建，则返回 nil。
func BeforeCreate(strategy RESTCreateStrategy, ctx context.Context, obj runtime.Object) error {
	objectMeta, kind, kerr := objectMetaAndKind(strategy, obj)
	if kerr != nil {
		return kerr
	}

	// ensure that system-critical metadata has been populated
	if !metav1.HasObjectMetaSystemFieldValues(objectMeta) {
		return errors.NewInternalError(fmt.Errorf("system metadata was not initialized"))
	}

	// ensure the name has been generated
	if len(objectMeta.GetGenerateName()) > 0 && len(objectMeta.GetName()) == 0 {
		return errors.NewInternalError(fmt.Errorf("metadata.name was not generated"))
	}

	// ensure namespace on the object is correct, or error if a conflicting namespace was set in the object
	requestNamespace, ok := genericapirequest.NamespaceFrom(ctx)
	if !ok {
		return errors.NewInternalError(fmt.Errorf("no namespace information found in request context"))
	}
	if err := EnsureObjectNamespaceMatchesRequestNamespace(ExpectedNamespaceForScope(requestNamespace, strategy.NamespaceScoped()), objectMeta); err != nil {
		return err
	}

	strategy.PrepareForCreate(ctx, obj)

	if errs := strategy.Validate(ctx, obj); len(errs) > 0 {
		return errors.NewInvalid(kind.GroupKind(), objectMeta.GetName(), errs)
	}

	// Custom validation (including name validation) passed
	// Now run common validation on object meta
	// Do this *after* custom validation so that specific error messages are shown whenever possible
	if errs := genericvalidation.ValidateObjectMetaAccessor(objectMeta, strategy.NamespaceScoped(), path.ValidatePathSegmentName, field.NewPath("metadata")); len(errs) > 0 {
		return errors.NewInvalid(kind.GroupKind(), objectMeta.GetName(), errs)
	}

	for _, w := range strategy.WarningsOnCreate(ctx, obj) {
		warning.AddWarning(ctx, "", w)
	}

	strategy.Canonicalize(obj)

	return nil
}

// CheckGeneratedNameError checks whether an error that occurred creating a resource is due
// to generation being unable to pick a valid name.
// CheckGeneratedNameError  检查创建资源时发生的错误是否是由于生成无法选择有效名称导致的。
func CheckGeneratedNameError(ctx context.Context, strategy RESTCreateStrategy, err error, obj runtime.Object) error {
	if !errors.IsAlreadyExists(err) {
		return err
	}

	objectMeta, gvk, kerr := objectMetaAndKind(strategy, obj)
	if kerr != nil {
		return kerr
	}

	if len(objectMeta.GetGenerateName()) == 0 {
		// If we don't have a generated name, return the original error (AlreadyExists).
		// When we're here, the user picked a name that is causing a conflict.
		return err
	}

	// Get the group resource information from the context, if populated.
	gr := schema.GroupResource{}
	if requestInfo, found := genericapirequest.RequestInfoFrom(ctx); found {
		gr = schema.GroupResource{Group: gvk.Group, Resource: requestInfo.Resource}
	}

	// If we have a name and generated name, the server picked a name
	// that already exists.
	return errors.NewGenerateNameConflict(gr, objectMeta.GetName(), 1)
}

// objectMetaAndKind retrieves kind and ObjectMeta from a runtime object, or returns an error.
// objectMetaAndKind 从运行时对象中检索 kind 和 ObjectMeta，或返回错误。
func objectMetaAndKind(typer runtime.ObjectTyper, obj runtime.Object) (metav1.Object, schema.GroupVersionKind, error) {
	objectMeta, err := meta.Accessor(obj)
	if err != nil {
		return nil, schema.GroupVersionKind{}, errors.NewInternalError(err)
	}
	kinds, _, err := typer.ObjectKinds(obj)
	if err != nil {
		return nil, schema.GroupVersionKind{}, errors.NewInternalError(err)
	}
	return objectMeta, kinds[0], nil
}

// NamespaceScopedStrategy has a method to tell if the object must be in a namespace.
// NamespaceScopedStrategy 有一个方法来告诉对象是否必须在命名空间中。
type NamespaceScopedStrategy interface {
	// NamespaceScoped returns if the object must be in a namespace.
	// NamespaceScoped 返回对象是否必须在命名空间中。
	NamespaceScoped() bool
}

// AdmissionToValidateObjectFunc converts validating admission to a rest validate object func
// AdmissionToValidateObjectFunc 将验证准入转换为 rest 验证对象函数
func AdmissionToValidateObjectFunc(admit admission.Interface, staticAttributes admission.Attributes, o admission.ObjectInterfaces) ValidateObjectFunc {
	validatingAdmission, ok := admit.(admission.ValidationInterface)
	if !ok {
		return func(ctx context.Context, obj runtime.Object) error { return nil }
	}
	return func(ctx context.Context, obj runtime.Object) error {
		name := staticAttributes.GetName()
		// in case the generated name is populated
		if len(name) == 0 {
			if metadata, err := meta.Accessor(obj); err == nil {
				name = metadata.GetName()
			}
		}

		finalAttributes := admission.NewAttributesRecord(
			obj,
			staticAttributes.GetOldObject(),
			staticAttributes.GetKind(),
			staticAttributes.GetNamespace(),
			name,
			staticAttributes.GetResource(),
			staticAttributes.GetSubresource(),
			staticAttributes.GetOperation(),
			staticAttributes.GetOperationOptions(),
			staticAttributes.IsDryRun(),
			staticAttributes.GetUserInfo(),
		)
		if !validatingAdmission.Handles(finalAttributes.GetOperation()) {
			return nil
		}
		return validatingAdmission.Validate(ctx, finalAttributes, o)
	}
}
