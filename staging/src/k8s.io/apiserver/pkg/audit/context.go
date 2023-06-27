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

package audit

import (
	"context"
	"sync"

	"k8s.io/apimachinery/pkg/types"
	auditinternal "k8s.io/apiserver/pkg/apis/audit"
	genericapirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/klog/v2"
)

// The key type is unexported to prevent collisions
type key int

// auditKey is the context key for storing the audit context that is being
// captured and the evaluated policy that applies to the given request.
// auditKey 是上下文键，用于存储正在捕获的审计上下文和适用于给定请求的评估策略。
const auditKey key = iota

// AuditContext holds the information for constructing the audit events for the current request.
// AuditContext 包含为当前请求构造审计事件的信息。
type AuditContext struct {
	// RequestAuditConfig is the audit configuration that applies to the request
	// RequestAuditConfig 是应用于请求的审计配置
	RequestAuditConfig RequestAuditConfig

	// Event is the audit Event object that is being captured to be written in
	// the API audit log. It is set to nil when the request is not being audited.
	// Event 是被捕获并写入 API 审计日志的审计事件对象。当请求未被审计时，它被设置为 nil。
	Event *auditinternal.Event

	// annotations holds audit annotations that are recorded before the event has been initialized.
	// This is represented as a slice rather than a map to preserve order.
	// annotations 保存在事件初始化之前记录的审计注释。这表示为切片而不是映射以保持顺序。
	annotations []annotation
	// annotationMutex guards annotations AND event.Annotations
	// annotationMutex 保护注释和事件。注释
	annotationMutex sync.Mutex

	// auditID is the Audit ID associated with this request.
	// auditID 是与此请求关联的审核 ID。
	auditID types.UID
}

type annotation struct {
	key, value string
}

// AddAuditAnnotation sets the audit annotation for the given key, value pair.
// It is safe to call at most parts of request flow that come after WithAuditAnnotations.
// The notable exception being that this function must not be called via a
// defer statement (i.e. after ServeHTTP) in a handler that runs before WithAudit
// as at that point the audit event has already been sent to the audit sink.
// Handlers that are unaware of their position in the overall request flow should
// prefer AddAuditAnnotation over LogAnnotation to avoid dropping annotations.
// AddAuditAnnotation 为给定的键值对设置审计注解。调用 WithAuditAnnotations 之后的请求流的大部分部分是安全的。
// 值得注意的例外是，在 WithAudit 之前运行的处理程序中，不得通过延迟语句（即在 ServeHTTP 之后）调用此函数，因为此时审计事件已经发送到审计接收器。
// 不知道自己在整个请求流中的位置的处理程序应该更喜欢 AddAuditAnnotation 而不是 LogAnnotation 以避免删除注释。
func AddAuditAnnotation(ctx context.Context, key, value string) {
	ac := AuditContextFrom(ctx)
	if ac == nil {
		// auditing is not enabled
		return
	}

	ac.annotationMutex.Lock()
	defer ac.annotationMutex.Unlock()

	addAuditAnnotationLocked(ac, key, value)
}

// AddAuditAnnotations is a bulk version of AddAuditAnnotation. Refer to AddAuditAnnotation for
// restrictions on when this can be called.
// keysAndValues are the key-value pairs to add, and must have an even number of items.
// AddAuditAnnotations 是 AddAuditAnnotation 的批量版本。有关何时可以调用的限制，请参阅 AddAuditAnnotation。
// keysAndValues 是要添加的键值对，并且必须具有偶数个项目。
func AddAuditAnnotations(ctx context.Context, keysAndValues ...string) {
	ac := AuditContextFrom(ctx)
	if ac == nil {
		// auditing is not enabled
		return
	}

	ac.annotationMutex.Lock()
	defer ac.annotationMutex.Unlock()

	if len(keysAndValues)%2 != 0 {
		klog.Errorf("Dropping mismatched audit annotation %q", keysAndValues[len(keysAndValues)-1])
	}
	for i := 0; i < len(keysAndValues); i += 2 {
		addAuditAnnotationLocked(ac, keysAndValues[i], keysAndValues[i+1])
	}
}

// AddAuditAnnotationsMap is a bulk version of AddAuditAnnotation. Refer to AddAuditAnnotation for
// restrictions on when this can be called.
// AddAuditAnnotationsMap 是 AddAuditAnnotation 的批量版本。有关何时可以调用的限制，请参阅 AddAuditAnnotation。
func AddAuditAnnotationsMap(ctx context.Context, annotations map[string]string) {
	ac := AuditContextFrom(ctx)
	if ac == nil {
		// auditing is not enabled
		return
	}

	ac.annotationMutex.Lock()
	defer ac.annotationMutex.Unlock()

	for k, v := range annotations {
		addAuditAnnotationLocked(ac, k, v)
	}
}

// addAuditAnnotationLocked is the shared code for recording an audit annotation. This method should
// only be called while the auditAnnotationsMutex is locked.
// addAuditAnnotationLocked 是记录审计注解的共享代码。只有在锁定 auditAnnotationsMutex 时才应调用此方法。
func addAuditAnnotationLocked(ac *AuditContext, key, value string) {
	if ac.Event != nil {
		logAnnotation(ac.Event, key, value)
	} else {
		ac.annotations = append(ac.annotations, annotation{key: key, value: value})
	}
}

// This is private to prevent reads/write to the slice from outside of this package.
// The audit event should be directly read to get access to the annotations.
// 这是私有的，以防止从这个包的外部读写切片。应直接读取审计事件以访问注释。
func addAuditAnnotationsFrom(ctx context.Context, ev *auditinternal.Event) {
	ac := AuditContextFrom(ctx)
	if ac == nil {
		// auditing is not enabled
		return
	}

	ac.annotationMutex.Lock()
	defer ac.annotationMutex.Unlock()

	for _, kv := range ac.annotations {
		logAnnotation(ev, kv.key, kv.value)
	}
}

// LogAnnotation fills in the Annotations according to the key value pair.
// LogAnnotation 根据键值对填写Annotations。
func logAnnotation(ae *auditinternal.Event, key, value string) {
	if ae == nil || ae.Level.Less(auditinternal.LevelMetadata) {
		return
	}
	if ae.Annotations == nil {
		ae.Annotations = make(map[string]string)
	}
	if v, ok := ae.Annotations[key]; ok && v != value {
		klog.Warningf("Failed to set annotations[%q] to %q for audit:%q, it has already been set to %q", key, value, ae.AuditID, ae.Annotations[key])
		return
	}
	ae.Annotations[key] = value
}

// WithAuditContext returns a new context that stores the AuditContext.
func WithAuditContext(parent context.Context) context.Context {
	if AuditContextFrom(parent) != nil {
		return parent // Avoid double registering.
	}

	return genericapirequest.WithValue(parent, auditKey, &AuditContext{})
}

// AuditEventFrom returns the audit event struct on the ctx
func AuditEventFrom(ctx context.Context) *auditinternal.Event {
	if o := AuditContextFrom(ctx); o != nil {
		return o.Event
	}
	return nil
}

// AuditContextFrom returns the pair of the audit configuration object
// that applies to the given request and the audit event that is going to
// be written to the API audit log.
func AuditContextFrom(ctx context.Context) *AuditContext {
	ev, _ := ctx.Value(auditKey).(*AuditContext)
	return ev
}

// WithAuditID sets the AuditID on the AuditContext. The AuditContext must already be present in the
// request context. If the specified auditID is empty, no value is set.
func WithAuditID(ctx context.Context, auditID types.UID) {
	if auditID == "" {
		return
	}
	ac := AuditContextFrom(ctx)
	if ac == nil {
		return
	}
	ac.auditID = auditID
	if ac.Event != nil {
		ac.Event.AuditID = auditID
	}
}

// AuditIDFrom returns the value of the audit ID from the request context.
func AuditIDFrom(ctx context.Context) (types.UID, bool) {
	if ac := AuditContextFrom(ctx); ac != nil {
		return ac.auditID, ac.auditID != ""
	}
	return "", false
}

// GetAuditIDTruncated returns the audit ID (truncated) from the request context.
// If the length of the Audit-ID value exceeds the limit, we truncate it to keep
// the first N (maxAuditIDLength) characters.
// This is intended to be used in logging only.
func GetAuditIDTruncated(ctx context.Context) string {
	auditID, ok := AuditIDFrom(ctx)
	if !ok {
		return ""
	}

	// if the user has specified a very long audit ID then we will use the first N characters
	// Note: assuming Audit-ID header is in ASCII
	const maxAuditIDLength = 64
	if len(auditID) > maxAuditIDLength {
		auditID = auditID[:maxAuditIDLength]
	}

	return string(auditID)
}
