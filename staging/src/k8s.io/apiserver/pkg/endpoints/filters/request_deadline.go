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

package filters

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	auditinternal "k8s.io/apiserver/pkg/apis/audit"
	"k8s.io/apiserver/pkg/audit"
	"k8s.io/apiserver/pkg/endpoints/handlers/responsewriters"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/klog/v2"
	"k8s.io/utils/clock"
)

const (
	// The 'timeout' query parameter in the request URL has an invalid duration specifier
	invalidTimeoutInURL = "invalid timeout specified in the request URL"
)

// WithRequestDeadline determines the timeout duration applicable to the given request and sets a new context
// with the appropriate deadline.
// auditWrapper provides an http.Handler that audits a failed request.
// longRunning returns true if he given request is a long running request.
// requestTimeoutMaximum specifies the default request timeout value.
// WithRequestDeadline 确定适用于给定请求的超时持续时间，并设置具有适当截止日期的新上下文。 auditWrapper 提供了一个 http.Handler 来审计失败的请求。如果他给定的请求是长时间运行的请求，则 longRunning 返回 true。 requestTimeoutMaximum 指定默认的请求超时值。
func WithRequestDeadline(handler http.Handler, sink audit.Sink, policy audit.PolicyRuleEvaluator, longRunning request.LongRunningRequestCheck,
	negotiatedSerializer runtime.NegotiatedSerializer, requestTimeoutMaximum time.Duration) http.Handler {
	return withRequestDeadline(handler, sink, policy, longRunning, negotiatedSerializer, requestTimeoutMaximum, clock.RealClock{})
}

func withRequestDeadline(handler http.Handler, sink audit.Sink, policy audit.PolicyRuleEvaluator, longRunning request.LongRunningRequestCheck,
	negotiatedSerializer runtime.NegotiatedSerializer, requestTimeoutMaximum time.Duration, clock clock.PassiveClock) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		requestInfo, ok := request.RequestInfoFrom(ctx)
		if !ok {
			handleError(w, req, http.StatusInternalServerError, fmt.Errorf("no RequestInfo found in context, handler chain must be wrong"))
			return
		}
		if longRunning(req, requestInfo) {
			handler.ServeHTTP(w, req)
			return
		}

		userSpecifiedTimeout, ok, err := parseTimeout(req)
		if err != nil {
			statusErr := apierrors.NewBadRequest(err.Error())

			klog.Errorf("Error - %s: %#v", err.Error(), req.RequestURI)

			failed := failedErrorHandler(negotiatedSerializer, statusErr)
			failWithAudit := withFailedRequestAudit(failed, statusErr, sink, policy)
			failWithAudit.ServeHTTP(w, req)
			return
		}

		timeout := requestTimeoutMaximum
		if ok {
			// we use the default timeout enforced by the apiserver:
			// - if the user has specified a timeout of 0s, this implies no timeout on the user's part.
			// - if the user has specified a timeout that exceeds the maximum deadline allowed by the apiserver.
			if userSpecifiedTimeout > 0 && userSpecifiedTimeout < requestTimeoutMaximum {
				timeout = userSpecifiedTimeout
			}
		}

		started := clock.Now()
		if requestStartedTimestamp, ok := request.ReceivedTimestampFrom(ctx); ok {
			started = requestStartedTimestamp
		}

		ctx, cancel := context.WithDeadline(ctx, started.Add(timeout))
		defer cancel()

		req = req.WithContext(ctx)
		handler.ServeHTTP(w, req)
	})
}

// withFailedRequestAudit decorates a failed http.Handler and is used to audit a failed request.
// statusErr is used to populate the Message property of ResponseStatus.
func withFailedRequestAudit(failedHandler http.Handler, statusErr *apierrors.StatusError, sink audit.Sink, policy audit.PolicyRuleEvaluator) http.Handler {
	if sink == nil || policy == nil {
		return failedHandler
	}
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ac, err := evaluatePolicyAndCreateAuditEvent(req, policy)
		if err != nil {
			utilruntime.HandleError(fmt.Errorf("failed to create audit event: %v", err))
			responsewriters.InternalError(w, req, errors.New("failed to create audit event"))
			return
		}

		if ac == nil || ac.Event == nil {
			failedHandler.ServeHTTP(w, req)
			return
		}
		ev := ac.Event

		ev.ResponseStatus = &metav1.Status{}
		ev.Stage = auditinternal.StageResponseStarted
		if statusErr != nil {
			ev.ResponseStatus.Message = statusErr.Error()
		}

		rw := decorateResponseWriter(req.Context(), w, ev, sink, ac.RequestAuditConfig.OmitStages)
		failedHandler.ServeHTTP(rw, req)
	})
}

// failedErrorHandler returns an http.Handler that uses the specified StatusError object
// to render an error response to the request.
func failedErrorHandler(s runtime.NegotiatedSerializer, statusError *apierrors.StatusError) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		requestInfo, found := request.RequestInfoFrom(ctx)
		if !found {
			responsewriters.InternalError(w, req, errors.New("no RequestInfo found in the context"))
			return
		}

		gv := schema.GroupVersion{Group: requestInfo.APIGroup, Version: requestInfo.APIVersion}
		responsewriters.ErrorNegotiated(statusError, s, gv, w, req)
	})
}

// parseTimeout parses the given HTTP request URL and extracts the timeout query parameter
// value if specified by the user.
// If a timeout is not specified the function returns false and err is set to nil
// If the value specified is malformed then the function returns false and err is set
func parseTimeout(req *http.Request) (time.Duration, bool, error) {
	value := req.URL.Query().Get("timeout")
	if value == "" {
		return 0, false, nil
	}

	timeout, err := time.ParseDuration(value)
	if err != nil {
		return 0, false, fmt.Errorf("%s - %s", invalidTimeoutInURL, err.Error())
	}

	return timeout, true, nil
}

func handleError(w http.ResponseWriter, r *http.Request, code int, err error) {
	errorMsg := fmt.Sprintf("Error - %s: %#v", err.Error(), r.RequestURI)
	http.Error(w, errorMsg, code)
	klog.Errorf(errorMsg)
}
