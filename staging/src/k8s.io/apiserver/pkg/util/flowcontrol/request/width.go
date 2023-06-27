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

package request

import (
	"fmt"
	"net/http"
	"time"

	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/klog/v2"
)

// WorkEstimate carries three of the four parameters that determine the work in a request.
// The fourth parameter is the duration of the initial phase of execution.
// WorkEstimate 携带确定请求工作的四个参数中的三个。第四个参数是执行初始阶段的持续时间。
type WorkEstimate struct {
	// InitialSeats is the number of seats occupied while the server is
	// executing this request.
	// InitialSeats 是服务器执行此请求时占用的席位数。
	InitialSeats uint64

	// FinalSeats is the number of seats occupied at the end,
	// during the AdditionalLatency.
	// FinalSeats 是在 AdditionalLatency 期间最后占用的席位数。
	FinalSeats uint64

	// AdditionalLatency specifies the additional duration the seats allocated
	// to this request must be reserved after the given request had finished.
	// AdditionalLatency should not have any impact on the user experience, the
	// caller must not experience this additional latency.
	// AdditionalLatency 指定分配给此请求的席位必须在给定请求完成后保留的额外持续时间。
	// AdditionalLatency 不应对用户体验产生任何影响，调用者不得体验这种额外的延迟。
	AdditionalLatency time.Duration
}

// MaxSeats returns the maximum number of seats the request occupies over the
// phases of being served.
// MaxSeats 返回请求在服务阶段占据的最大席位数。
func (we *WorkEstimate) MaxSeats() int {
	if we.InitialSeats >= we.FinalSeats {
		return int(we.InitialSeats)
	}

	return int(we.FinalSeats)
}

// objectCountGetterFunc represents a function that gets the total
// number of objects for a given resource.
// objectCountGetterFunc 表示获取给定资源的对象总数的函数。
type objectCountGetterFunc func(string) (int64, error)

// watchCountGetterFunc represents a function that gets the total
// number of watchers potentially interested in a given request.
// watchCountGetterFunc 表示获取可能对给定请求感兴趣的观察者总数的函数。
type watchCountGetterFunc func(*apirequest.RequestInfo) int

// NewWorkEstimator estimates the work that will be done by a given request,
// if no WorkEstimatorFunc matches the given request then the default
// work estimate of 1 seat is allocated to the request.
// NewWorkEstimator 估计给定请求将完成的工作，如果没有 WorkEstimatorFunc 与给定请求匹配，则将 1 个席位的默认工作估计分配给该请求。
func NewWorkEstimator(objectCountFn objectCountGetterFunc, watchCountFn watchCountGetterFunc, config *WorkEstimatorConfig) WorkEstimatorFunc {
	estimator := &workEstimator{
		minimumSeats:          config.MinimumSeats,
		maximumSeats:          config.MaximumSeats,
		listWorkEstimator:     newListWorkEstimator(objectCountFn, config),
		mutatingWorkEstimator: newMutatingWorkEstimator(watchCountFn, config),
	}
	return estimator.estimate
}

// WorkEstimatorFunc returns the estimated work of a given request.
// This function will be used by the Priority & Fairness filter to
// estimate the work of of incoming requests.
// WorkEstimatorFunc 返回给定请求的估计工作。 Priority & Fairness 过滤器将使用此函数来估计传入请求的工作。
type WorkEstimatorFunc func(request *http.Request, flowSchemaName, priorityLevelName string) WorkEstimate

func (e WorkEstimatorFunc) EstimateWork(r *http.Request, flowSchemaName, priorityLevelName string) WorkEstimate {
	return e(r, flowSchemaName, priorityLevelName)
}

type workEstimator struct {
	// the minimum number of seats a request must occupy
	minimumSeats uint64
	// the maximum number of seats a request can occupy
	maximumSeats uint64
	// listWorkEstimator estimates work for list request(s)
	listWorkEstimator WorkEstimatorFunc
	// mutatingWorkEstimator calculates the width of mutating request(s)
	mutatingWorkEstimator WorkEstimatorFunc
}

func (e *workEstimator) estimate(r *http.Request, flowSchemaName, priorityLevelName string) WorkEstimate {
	requestInfo, ok := apirequest.RequestInfoFrom(r.Context())
	if !ok {
		klog.ErrorS(fmt.Errorf("no RequestInfo found in context"), "Failed to estimate work for the request", "URI", r.RequestURI)
		// no RequestInfo should never happen, but to be on the safe side let's return maximumSeats
		return WorkEstimate{InitialSeats: e.maximumSeats}
	}

	switch requestInfo.Verb {
	case "list":
		return e.listWorkEstimator.EstimateWork(r, flowSchemaName, priorityLevelName)
	case "create", "update", "patch", "delete":
		return e.mutatingWorkEstimator.EstimateWork(r, flowSchemaName, priorityLevelName)
	}

	return WorkEstimate{InitialSeats: e.minimumSeats}
}
