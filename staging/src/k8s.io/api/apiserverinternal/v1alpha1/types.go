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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Storage version of a specific resource.
// 特定资源的存储版本。
type StorageVersion struct {
	metav1.TypeMeta `json:",inline"`
	// The name is <group>.<resource>.
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec is an empty spec. It is here to comply with Kubernetes API style.
	// Spec 是否为空。这里是为了符合 Kubernetes API 风格。
	Spec StorageVersionSpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`

	// API server instances report the version they can decode and the version they
	// encode objects to when persisting objects in the backend.
	// API 服务器实例报告它们可以解码的版本以及在后端持久化对象时它们编码对象的版本。
	Status StorageVersionStatus `json:"status" protobuf:"bytes,3,opt,name=status"`
}

// StorageVersionSpec is an empty spec.
type StorageVersionSpec struct{}

// API server instances report the versions they can decode and the version they
// encode objects to when persisting objects in the backend.
// API 服务器实例报告它们可以解码的版本以及在后端持久化对象时它们编码对象的版本。
type StorageVersionStatus struct {
	// The reported versions per API server instance.
	// +optional
	// +listType=map
	// +listMapKey=apiServerID
	// 每个API服务器实例报告的版本。
	StorageVersions []ServerStorageVersion `json:"storageVersions,omitempty" protobuf:"bytes,1,opt,name=storageVersions"`
	// If all API server instances agree on the same encoding storage version,
	// then this field is set to that version. Otherwise this field is left empty.
	// API servers should finish updating its storageVersionStatus entry before
	// serving write operations, so that this field will be in sync with the reality.
	// +optional
	// 如果所有API服务器实例都同意使用相同的编码存储版本，则将此字段设置为该版本。否则，此字段将保持为空。
	// API服务器应在提供写操作之前完成其storageVersionStatus条目的更新，以便此字段与现实保持同步。
	CommonEncodingVersion *string `json:"commonEncodingVersion,omitempty" protobuf:"bytes,2,opt,name=commonEncodingVersion"`

	// The latest available observations of the storageVersion's state.
	// +optional
	// +listType=map
	// +listMapKey=type
	// 存储版本状态的最新可用观察值。
	Conditions []StorageVersionCondition `json:"conditions,omitempty" protobuf:"bytes,3,opt,name=conditions"`
}

// An API server instance reports the version it can decode and the version it
// encodes objects to when persisting objects in the backend.
// API 服务器实例报告它可以解码的版本以及在后端持久化对象时它编码对象的版本。
type ServerStorageVersion struct {
	// The ID of the reporting API server.
	APIServerID string `json:"apiServerID,omitempty" protobuf:"bytes,1,opt,name=apiServerID"`

	// The API server encodes the object to this version when persisting it in
	// the backend (e.g., etcd).
	// API服务器在后台持久化对象时(例如，etcd)将对象编码为这个版本。
	EncodingVersion string `json:"encodingVersion,omitempty" protobuf:"bytes,2,opt,name=encodingVersion"`

	// The API server can decode objects encoded in these versions.
	// The encodingVersion must be included in the decodableVersions.
	// +listType=set
	// API服务器可以解码在这些版本中编码的对象。encodingVersion必须包含在decodableVersions中。
	DecodableVersions []string `json:"decodableVersions,omitempty" protobuf:"bytes,3,opt,name=decodableVersions"`
}

type StorageVersionConditionType string

const (
	// Indicates that encoding storage versions reported by all servers are equal.
	AllEncodingVersionsEqual StorageVersionConditionType = "AllEncodingVersionsEqual"
)

type ConditionStatus string

const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

// Describes the state of the storageVersion at a certain point.
// 描述存储版本在某个时间点的状态。
type StorageVersionCondition struct {
	// Type of the condition.
	// +required
	Type StorageVersionConditionType `json:"type" protobuf:"bytes,1,opt,name=type"`
	// Status of the condition, one of True, False, Unknown.
	// +required
	Status ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status"`
	// If set, this represents the .metadata.generation that the condition was set based upon.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,3,opt,name=observedGeneration"`
	// Last time the condition transitioned from one status to another.
	// +required
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// The reason for the condition's last transition.
	// +required
	Reason string `json:"reason" protobuf:"bytes,5,opt,name=reason"`
	// A human readable message indicating details about the transition.
	// +required
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// A list of StorageVersions.
type StorageVersionList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Items holds a list of StorageVersion
	Items []StorageVersion `json:"items" protobuf:"bytes,2,rep,name=items"`
}
