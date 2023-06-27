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

package resourcequota

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/storage/names"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	api "k8s.io/kubernetes/pkg/apis/core"
	"k8s.io/kubernetes/pkg/apis/core/validation"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

// resourcequotaStrategy implements behavior for ResourceQuota objects
// resourcequotaStrategy 实现了 ResourceQuota 对象的行为。
type resourcequotaStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

// Strategy is the default logic that applies when creating and updating ResourceQuota
// objects via the REST API.
// Strategy 是在通过 REST API 创建和更新 ResourceQuota 对象时默认的逻辑。
var Strategy = resourcequotaStrategy{legacyscheme.Scheme, names.SimpleNameGenerator}

// NamespaceScoped is true for resourcequotas.
// NamespaceScoped 为 true 表示 resourcequotas 是命名空间范围的。
func (resourcequotaStrategy) NamespaceScoped() bool {
	return true
}

// GetResetFields returns the set of fields that get reset by the strategy
// and should not be modified by the user.
// GetResetFields 返回策略重置的字段集，用户不应该修改这些字段。
func (resourcequotaStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	fields := map[fieldpath.APIVersion]*fieldpath.Set{
		"v1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("status"),
		),
	}

	return fields
}

// PrepareForCreate clears fields that are not allowed to be set by end users on creation.
// PrepareForCreate 清除创建时不允许由最终用户设置的字段。
func (resourcequotaStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	resourcequota := obj.(*api.ResourceQuota)
	resourcequota.Status = api.ResourceQuotaStatus{}
}

// PrepareForUpdate clears fields that are not allowed to be set by end users on update.
// PrepareForUpdate 清除更新时不允许由最终用户设置的字段。
func (resourcequotaStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newResourcequota := obj.(*api.ResourceQuota)
	oldResourcequota := old.(*api.ResourceQuota)
	newResourcequota.Status = oldResourcequota.Status
}

// Validate validates a new resourcequota.
// Validate 验证新的 resourcequota。
func (resourcequotaStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	resourcequota := obj.(*api.ResourceQuota)
	return validation.ValidateResourceQuota(resourcequota)
}

// WarningsOnCreate returns warnings for the creation of the given object.
// WarningsOnCreate 返回给定对象创建时的警告。
func (resourcequotaStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

// Canonicalize normalizes the object after validation.
// Canonicalize 在验证后规范化对象。
func (resourcequotaStrategy) Canonicalize(obj runtime.Object) {
}

// AllowCreateOnUpdate is false for resourcequotas.
// AllowCreateOnUpdate 为 resourcequotas 返回 false。
func (resourcequotaStrategy) AllowCreateOnUpdate() bool {
	return false
}

// ValidateUpdate is the default update validation for an end user.
// ValidateUpdate 是最终用户更新验证的默认值。
func (resourcequotaStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	newObj, oldObj := obj.(*api.ResourceQuota), old.(*api.ResourceQuota)
	return validation.ValidateResourceQuotaUpdate(newObj, oldObj)
}

// WarningsOnUpdate returns warnings for the given update.
// WarningsOnUpdate 返回给定更新时的警告。
func (resourcequotaStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}

func (resourcequotaStrategy) AllowUnconditionalUpdate() bool {
	return true
}

type resourcequotaStatusStrategy struct {
	resourcequotaStrategy
}

// StatusStrategy is the default logic invoked when updating object status.
// StatusStrategy 是更新对象状态时默认调用的逻辑。
var StatusStrategy = resourcequotaStatusStrategy{Strategy}

// GetResetFields returns the set of fields that get reset by the strategy
// and should not be modified by the user.
// GetResetFields 返回由策略重置的字段集，用户不应该修改它们。
func (resourcequotaStatusStrategy) GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set {
	fields := map[fieldpath.APIVersion]*fieldpath.Set{
		"v1": fieldpath.NewSet(
			fieldpath.MakePathOrDie("spec"),
		),
	}

	return fields
}

func (resourcequotaStatusStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
	newResourcequota := obj.(*api.ResourceQuota)
	oldResourcequota := old.(*api.ResourceQuota)
	newResourcequota.Spec = oldResourcequota.Spec
}

func (resourcequotaStatusStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return validation.ValidateResourceQuotaStatusUpdate(obj.(*api.ResourceQuota), old.(*api.ResourceQuota))
}

// WarningsOnUpdate returns warnings for the given update.
func (resourcequotaStatusStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}
