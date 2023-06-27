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

package handlers

import (
	"fmt"
	"net/http"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/endpoints/request"
)

// ScopeNamer handles accessing names from requests and objects
// ScopeNamer 处理来自请求和对象的访问名称
type ScopeNamer interface {
	// Namespace returns the appropriate namespace value from the request (may be empty) or an
	// error.
	// Namespace 从请求（可能为空）或错误中返回适当的命名空间值。
	Namespace(req *http.Request) (namespace string, err error)
	// Name returns the name from the request, and an optional namespace value if this is a namespace
	// scoped call. An error is returned if the name is not available.
	// Name 从请求返回名称，并返回可选的命名空间值（如果这是命名空间范围的调用）。 如果名称不可用，则返回错误。
	Name(req *http.Request) (namespace, name string, err error)
	// ObjectName returns the namespace and name from an object if they exist, or an error if the object
	// does not support names.
	// ObjectName 从对象返回命名空间和名称（如果它们存在），如果对象不支持名称，则返回错误。
	ObjectName(obj runtime.Object) (namespace, name string, err error)
}

type ContextBasedNaming struct {
	Namer         runtime.Namer
	ClusterScoped bool
}

// ContextBasedNaming implements ScopeNamer
// ContextBasedNaming 实现 ScopeNamer
var _ ScopeNamer = ContextBasedNaming{}

func (n ContextBasedNaming) Namespace(req *http.Request) (namespace string, err error) {
	requestInfo, ok := request.RequestInfoFrom(req.Context())
	if !ok {
		return "", fmt.Errorf("missing requestInfo")
	}
	return requestInfo.Namespace, nil
}

func (n ContextBasedNaming) Name(req *http.Request) (namespace, name string, err error) {
	requestInfo, ok := request.RequestInfoFrom(req.Context())
	if !ok {
		return "", "", fmt.Errorf("missing requestInfo")
	}

	if len(requestInfo.Name) == 0 {
		return "", "", errEmptyName
	}
	return requestInfo.Namespace, requestInfo.Name, nil
}

func (n ContextBasedNaming) ObjectName(obj runtime.Object) (namespace, name string, err error) {
	name, err = n.Namer.Name(obj)
	if err != nil {
		return "", "", err
	}
	if len(name) == 0 {
		return "", "", errEmptyName
	}
	namespace, err = n.Namer.Namespace(obj)
	if err != nil {
		return "", "", err
	}
	return namespace, name, err
}

// errEmptyName is returned when API requests do not fill the name section of the path.
var errEmptyName = errors.NewBadRequest("name must be provided")
