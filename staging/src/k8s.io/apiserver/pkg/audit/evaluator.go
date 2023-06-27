/*
Copyright 2021 The Kubernetes Authors.

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
	"k8s.io/apiserver/pkg/apis/audit"
	"k8s.io/apiserver/pkg/authorization/authorizer"
)

// RequestAuditConfig is the evaluated audit configuration that is applicable to
// a given request. PolicyRuleEvaluator evaluates the audit policy against the
// authorizer attributes and returns a RequestAuditConfig that applies to the request.
// RequestAuditConfig 包含了针对给定请求的已评估的审计配置。
// PolicyRuleEvaluator 评估审计策略与授权方属性，并返回适用于请求的 RequestAuditConfig。
type RequestAuditConfig struct {
	// Level at which the request is being audited at
	// 请求被审计的级别
	Level audit.Level

	// OmitStages is the stages that need to be omitted from being audited.
	// 需要从审计中省略的阶段
	OmitStages []audit.Stage

	// OmitManagedFields indicates whether to omit the managed fields of the request
	// and response bodies from being written to the API audit log.
	// OmitManagedFields 表示是否从写入 API 审计日志中省略请求和响应体的管理字段。
	OmitManagedFields bool
}

// PolicyRuleEvaluator exposes methods for evaluating the policy rules.
// PolicyRuleEvaluator 公开了评估策略规则的方法。
type PolicyRuleEvaluator interface {
	// EvaluatePolicyRule evaluates the audit policy of the apiserver against
	// the given authorizer attributes and returns the audit configuration that
	// is applicable to the given equest.
	// EvaluatePolicyRule 评估 apiserver 的审计策略与给定的授权方属性，并返回适用于给定请求的审计配置。
	EvaluatePolicyRule(authorizer.Attributes) RequestAuditConfig
}
