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

package rest

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// RESTCreateUpdateStrategy is a union of RESTUpdateStrategy and RESTCreateStrategy,
// and it defines the minimum validation, accepted input, and name generation
// behavior to create and update an object that follows Kubernetes API conventions.
// RESTCreateUpdateStrategy 是 RESTUpdateStrategy 和 RESTCreateStrategy 的联合，
// 它定义了创建和更新遵循 Kubernetes API 约定的对象的最小验证、接受的输入和名称生成行为。
type RESTCreateUpdateStrategy interface {
	RESTCreateStrategy
	// AllowCreateOnUpdate returns true if the object can be created by a PUT.
	// AllowCreateOnUpdate 如果对象可以通过 PUT 创建，则返回 true。
	AllowCreateOnUpdate() bool
	// PrepareForUpdate is invoked on update before validation to normalize
	// the object.  For example: remove fields that are not to be persisted,
	// sort order-insensitive list fields, etc.  This should not remove fields
	// whose presence would be considered a validation error.
	// PrepareForUpdate 在验证之前调用以对对象进行规范化。 例如：删除不要持久化的字段、
	// 对无序列表字段进行排序等。 这不应该删除其存在将被视为验证错误的字段。
	PrepareForUpdate(ctx context.Context, obj, old runtime.Object)
	// ValidateUpdate is invoked after default fields in the object have been
	// filled in before the object is persisted.  This method should not mutate
	// the object.
	// ValidateUpdate 在对象持久化之前填充对象中的默认字段之后调用。 此方法不应该改变对象。
	ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList
	// WarningsOnUpdate returns warnings to the client performing the update.
	// WarningsOnUpdate is invoked after default fields in the object have been filled in
	// and after ValidateUpdate has passed, before Canonicalize is called, and before the object is persisted.
	// This method must not mutate either object.
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
	// For example, do not warn about spec fields in a status update.
	WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string
	// AllowUnconditionalUpdate returns true if the object can be updated
	// unconditionally (irrespective of the latest resource version), when
	// there is no resource version specified in the object.
	// AllowUnconditionalUpdate 如果对象可以在没有指定资源版本的情况下无条件更新（无论最新的资源版本如何），
	// 则返回 true。
	AllowUnconditionalUpdate() bool
}

// Ensure that RESTCreateUpdateStrategy extends RESTCreateStrategy
// Ensure that RESTCreateUpdateStrategy extends RESTCreateStrategy
var _ RESTCreateStrategy = (RESTCreateUpdateStrategy)(nil)

// Ensure that RESTCreateUpdateStrategy extends RESTUpdateStrategy
// Ensure that RESTCreateUpdateStrategy extends RESTUpdateStrategy
var _ RESTUpdateStrategy = (RESTCreateUpdateStrategy)(nil)
