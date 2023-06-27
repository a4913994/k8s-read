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

package v1

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

const (
	// ImpersonateUserHeader is used to impersonate a particular user during an API server request
	// ImpersonateUserHeader用于在API服务器请求期间模拟特定用户
	ImpersonateUserHeader = "Impersonate-User"

	// ImpersonateGroupHeader is used to impersonate a particular group during an API server request.
	// It can be repeated multiplied times for multiple groups.
	// ImpersonateGroupHeader用于在API服务器请求期间模拟一个特定的组。可以对多个组重复多次。
	ImpersonateGroupHeader = "Impersonate-Group"

	// ImpersonateUIDHeader is used to impersonate a particular UID during an API server request
	// ImpersonateUIDHeader用于在API服务器请求期间模拟特定的UID
	ImpersonateUIDHeader = "Impersonate-Uid"

	// ImpersonateUserExtraHeaderPrefix is a prefix for any header used to impersonate an entry in the
	// extra map[string][]string for user.Info.  The key will be every after the prefix.
	// It can be repeated multiplied times for multiple map keys and the same key can be repeated multiple
	// times to have multiple elements in the slice under a single key
	// ImpersonateUserExtraHeaderPrefix是用于模拟user.Info中extra map[string] [] string中条目的任何标题的前缀。键将在前缀之后。可以对多个映射键重复多次，同一键可以多次重复以在单个键下的切片中具有多个元素
	ImpersonateUserExtraHeaderPrefix = "Impersonate-Extra-"
)

// +genclient
// +genclient:nonNamespaced
// +genclient:onlyVerbs=create
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TokenReview attempts to authenticate a token to a known user.
// Note: TokenReview requests may be cached by the webhook token authenticator
// plugin in the kube-apiserver.
// TokenReview 尝试将令牌验证为已知用户。
// 注意：TokenReview请求可能会被kube-apiserver中的webhook令牌验证器插件缓存。
type TokenReview struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	// 对象的标准元数据。
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec holds information about the request being evaluated
	// Spec 处理有关正在评估的请求的信息
	Spec TokenReviewSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`

	// Status is filled in by the server and indicates whether the request can be authenticated.
	// +optional
	// Status由服务器填写，并指示请求是否可以进行身份验证。
	Status TokenReviewStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// TokenReviewSpec is a description of the token authentication request.
// TokenReviewSpec是令牌身份验证请求的描述。
type TokenReviewSpec struct {
	// Token is the opaque bearer token.
	// +optional
	// Token是不透明的bearer令牌。
	Token string `json:"token,omitempty" protobuf:"bytes,1,opt,name=token"`
	// Audiences is a list of the identifiers that the resource server presented
	// with the token identifies as. Audience-aware token authenticators will
	// verify that the token was intended for at least one of the audiences in
	// this list. If no audiences are provided, the audience will default to the
	// audience of the Kubernetes apiserver.
	// +optional
	// Audiences是资源服务器提供的令牌标识的标识符列表。受众感知的令牌验证器将验证该令牌是否适用于此列表中的至少一个受众。如果未提供受众，则受众将默认为Kubernetes apiserver的受众。
	Audiences []string `json:"audiences,omitempty" protobuf:"bytes,2,rep,name=audiences"`
}

// TokenReviewStatus is the result of the token authentication request.
// TokenReviewStatus是令牌身份验证请求的结果。
type TokenReviewStatus struct {
	// Authenticated indicates that the token was associated with a known user.
	// +optional
	// Authenticated表示令牌与已知用户相关联。
	Authenticated bool `json:"authenticated,omitempty" protobuf:"varint,1,opt,name=authenticated"`
	// User is the UserInfo associated with the provided token.
	// +optional
	// User与提供的令牌相关联的UserInfo。
	User UserInfo `json:"user,omitempty" protobuf:"bytes,2,opt,name=user"`
	// Audiences are audience identifiers chosen by the authenticator that are
	// compatible with both the TokenReview and token. An identifier is any
	// identifier in the intersection of the TokenReviewSpec audiences and the
	// token's audiences. A client of the TokenReview API that sets the
	// spec.audiences field should validate that a compatible audience identifier
	// is returned in the status.audiences field to ensure that the TokenReview
	// server is audience aware. If a TokenReview returns an empty
	// status.audience field where status.authenticated is "true", the token is
	// valid against the audience of the Kubernetes API server.
	// +optional
	// 受众是验证者选择的受众标识符，它们与TokenReview和令牌兼容。标识符是TokenReviewSpec受众和令牌受众交集中的任何标识符。
	// TokenReview API的客户端如果设置了spec.audiences字段，则应验证status.audiences字段中返回的兼容受众标识符，以确保TokenReview服务器是受众感知的。
	Audiences []string `json:"audiences,omitempty" protobuf:"bytes,4,rep,name=audiences"`
	// Error indicates that the token couldn't be checked
	// +optional
	// Error表示无法检查令牌
	Error string `json:"error,omitempty" protobuf:"bytes,3,opt,name=error"`
}

// UserInfo holds the information about the user needed to implement the
// user.Info interface.
// UserInfo包含实现user.Info接口所需的用户信息。
type UserInfo struct {
	// The name that uniquely identifies this user among all active users.
	// +optional
	// 该名称在所有活动用户中唯一标识该用户。
	Username string `json:"username,omitempty" protobuf:"bytes,1,opt,name=username"`
	// A unique value that identifies this user across time. If this user is
	// deleted and another user by the same name is added, they will have
	// different UIDs.
	// +optional
	// 该值在时间跨度内唯一标识该用户。如果删除了该用户，并且又添加了同名的另一个用户，则它们将具有不同的UID。
	UID string `json:"uid,omitempty" protobuf:"bytes,2,opt,name=uid"`
	// The names of groups this user is a part of.
	// +optional
	// 该用户所属组的名称。
	Groups []string `json:"groups,omitempty" protobuf:"bytes,3,rep,name=groups"`
	// Any additional information provided by the authenticator.
	// +optional
	// 由验证器提供的任何其他信息。
	Extra map[string]ExtraValue `json:"extra,omitempty" protobuf:"bytes,4,rep,name=extra"`
}

// ExtraValue masks the value so protobuf can generate
// +protobuf.nullable=true
// +protobuf.options.(gogoproto.goproto_stringer)=false
type ExtraValue []string

func (t ExtraValue) String() string {
	return fmt.Sprintf("%v", []string(t))
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// TokenRequest requests a token for a given service account.
// TokenRequest请求给定服务帐户的令牌。
type TokenRequest struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	//
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec holds information about the request being evaluated
	// Spec处理正在评估的请求的信息
	Spec TokenRequestSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`

	// Status is filled in by the server and indicates whether the token can be authenticated.
	// +optional
	// Status由服务器填写，并指示令牌是否可以进行身份验证。
	Status TokenRequestStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// TokenRequestSpec contains client provided parameters of a token request.
// TokenRequestSpec包含客户端提供的令牌请求参数。
type TokenRequestSpec struct {
	// Audiences are the intendend audiences of the token. A recipient of a
	// token must identify themself with an identifier in the list of
	// audiences of the token, and otherwise should reject the token. A
	// token issued for multiple audiences may be used to authenticate
	// against any of the audiences listed but implies a high degree of
	// trust between the target audiences.
	// Audiences 是令牌的目标受众。令牌的接收者必须使用令牌的受众列表中的标识符对自己进行身份验证，否则应拒绝令牌。
	// 为多个受众颁发的令牌可以用于对列出的任何受众进行身份验证，但这意味着目标受众之间的信任程度很高。
	Audiences []string `json:"audiences" protobuf:"bytes,1,rep,name=audiences"`

	// ExpirationSeconds is the requested duration of validity of the request. The
	// token issuer may return a token with a different validity duration so a
	// client needs to check the 'expiration' field in a response.
	// +optional
	// ExpirationSeconds是请求的有效期。令牌颁发者可能会返回具有不同有效期的令牌，因此客户端需要检查响应中的“expiration”字段。
	ExpirationSeconds *int64 `json:"expirationSeconds" protobuf:"varint,4,opt,name=expirationSeconds"`

	// BoundObjectRef is a reference to an object that the token will be bound to.
	// The token will only be valid for as long as the bound object exists.
	// NOTE: The API server's TokenReview endpoint will validate the
	// BoundObjectRef, but other audiences may not. Keep ExpirationSeconds
	// small if you want prompt revocation.
	// +optional
	// BoundObjectRef是令牌将绑定到的对象的引用。令牌仅在绑定对象存在时才有效。
	// 注意：API服务器的TokenReview端点将验证BoundObjectRef，但其他受众可能不会。如果要迅速撤销，请保持ExpirationSeconds较小。
	BoundObjectRef *BoundObjectReference `json:"boundObjectRef" protobuf:"bytes,3,opt,name=boundObjectRef"`
}

// TokenRequestStatus is the result of a token request.
// TokenRequestStatus是令牌请求的结果。
type TokenRequestStatus struct {
	// Token is the opaque bearer token.
	// Token是不透明的持有者令牌。
	Token string `json:"token" protobuf:"bytes,1,opt,name=token"`
	// ExpirationTimestamp is the time of expiration of the returned token.
	// ExpirationTimestamp是返回的令牌的过期时间。
	ExpirationTimestamp metav1.Time `json:"expirationTimestamp" protobuf:"bytes,2,opt,name=expirationTimestamp"`
}

// BoundObjectReference is a reference to an object that a token is bound to.
// BoundObjectReference是令牌绑定到的对象的引用。
type BoundObjectReference struct {
	// Kind of the referent. Valid kinds are 'Pod' and 'Secret'.
	// +optional
	// 引用的类型。有效的类型是“Pod”和“Secret”。
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
	// API version of the referent.
	// +optional
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,2,opt,name=apiVersion"`

	// Name of the referent.
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,3,opt,name=name"`
	// UID of the referent.
	// +optional
	UID types.UID `json:"uid,omitempty" protobuf:"bytes,4,opt,name=uID,casttype=k8s.io/apimachinery/pkg/types.UID"`
}
