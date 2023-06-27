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

package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Event is a report of an event somewhere in the cluster. It generally denotes some state change in the system.
// Events have a limited retention time and triggers and messages may evolve
// with time.  Event consumers should not rely on the timing of an event
// with a given Reason reflecting a consistent underlying trigger, or the
// continued existence of events with that Reason.  Events should be
// treated as informative, best-effort, supplemental data.
// Event是集群中某处的事件报告。它通常表示系统中某种状态的变化。事件具有有限的保留时间，触发器和消息可能随时间而演变。事件消费者不应该依赖于具有给定原因的事件的时间反映一致的基础触发器，或者该原因的事件的持续存在。事件应该被视为信息性的，尽力而为的，补充数据。
type Event struct {
	metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	// eventTime is the time when this Event was first observed. It is required.
	// eventTime是第一次观察到此事件的时间。这是必需的。
	EventTime metav1.MicroTime `json:"eventTime" protobuf:"bytes,2,opt,name=eventTime"`

	// series is data about the Event series this event represents or nil if it's a singleton Event.
	// +optional
	// series是有关此事件表示的事件系列的数据，如果它是单个事件，则为零。
	Series *EventSeries `json:"series,omitempty" protobuf:"bytes,3,opt,name=series"`

	// reportingController is the name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`.
	// This field cannot be empty for new Events.
	// reportingController是发出此事件的控制器的名称，例如`kubernetes.io/kubelet`。对于新事件，此字段不能为空。
	ReportingController string `json:"reportingController,omitempty" protobuf:"bytes,4,opt,name=reportingController"`

	// reportingInstance is the ID of the controller instance, e.g. `kubelet-xyzf`.
	// This field cannot be empty for new Events and it can have at most 128 characters.
	// reportingInstance是控制器实例的ID，例如`kubelet-xyzf`。对于新事件，此字段不能为空，最多可以有128个字符。
	ReportingInstance string `json:"reportingInstance,omitempty" protobuf:"bytes,5,opt,name=reportingInstance"`

	// action is what action was taken/failed regarding to the regarding object. It is machine-readable.
	// This field cannot be empty for new Events and it can have at most 128 characters.
	// action是采取了什么行动/失败了关于regarding对象。这是机器可读的。对于新事件，此字段不能为空，最多可以有128个字符。
	Action string `json:"action,omitempty" protobuf:"bytes,6,name=action"`

	// reason is why the action was taken. It is human-readable.
	// This field cannot be empty for new Events and it can have at most 128 characters.
	// reason是为什么采取了这个行动。这是人类可读的。对于新事件，此字段不能为空，最多可以有128个字符。
	Reason string `json:"reason,omitempty" protobuf:"bytes,7,name=reason"`

	// regarding contains the object this Event is about. In most cases it's an Object reporting controller
	// implements, e.g. ReplicaSetController implements ReplicaSets and this event is emitted because
	// it acts on some changes in a ReplicaSet object.
	// +optional
	// regarding包含此事件所涉及的对象。在大多数情况下，它是一个对象报告控制器实现，例如ReplicaSetControlle
	// r实现了ReplicaSets，因为它会对ReplicaSet对象中的某些更改采取行动，所以会发出此事件。
	Regarding corev1.ObjectReference `json:"regarding,omitempty" protobuf:"bytes,8,opt,name=regarding"`

	// related is the optional secondary object for more complex actions. E.g. when regarding object triggers
	// a creation or deletion of related object.
	// +optional
	// related是可选的次要对象，用于更复杂的操作。例如，当regarding对象触发创建或删除相关对象时。
	Related *corev1.ObjectReference `json:"related,omitempty" protobuf:"bytes,9,opt,name=related"`

	// note is a human-readable description of the status of this operation.
	// Maximal length of the note is 1kB, but libraries should be prepared to
	// handle values up to 64kB.
	// +optional
	// 注释是对该操作状态的可读描述。注释的最大长度是1kB，但是库应该准备好处理高达64kB的值。
	Note string `json:"note,omitempty" protobuf:"bytes,10,opt,name=note"`

	// type is the type of this event (Normal, Warning), new types could be added in the future.
	// It is machine-readable.
	// This field cannot be empty for new Events.
	// type是此事件的类型（Normal，Warning），以后可能会添加新类型。这是机器可读的。对于新事件，此字段不能为空。
	Type string `json:"type,omitempty" protobuf:"bytes,11,opt,name=type"`

	// deprecatedSource is the deprecated field assuring backward compatibility with core.v1 Event type.
	// +optional
	// deprecatedSource是确保与core.v1 Event类型向后兼容的已弃用字段。
	DeprecatedSource corev1.EventSource `json:"deprecatedSource,omitempty" protobuf:"bytes,12,opt,name=deprecatedSource"`
	// deprecatedFirstTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
	// +optional
	DeprecatedFirstTimestamp metav1.Time `json:"deprecatedFirstTimestamp,omitempty" protobuf:"bytes,13,opt,name=deprecatedFirstTimestamp"`
	// deprecatedLastTimestamp is the deprecated field assuring backward compatibility with core.v1 Event type.
	// +optional
	DeprecatedLastTimestamp metav1.Time `json:"deprecatedLastTimestamp,omitempty" protobuf:"bytes,14,opt,name=deprecatedLastTimestamp"`
	// deprecatedCount is the deprecated field assuring backward compatibility with core.v1 Event type.
	// +optional
	DeprecatedCount int32 `json:"deprecatedCount,omitempty" protobuf:"varint,15,opt,name=deprecatedCount"`
}

// EventSeries contain information on series of events, i.e. thing that was/is happening
// continuously for some time. How often to update the EventSeries is up to the event reporters.
// The default event reporter in "k8s.io/client-go/tools/events/event_broadcaster.go" shows
// how this struct is updated on heartbeats and can guide customized reporter implementations.
// EventSeries 包含一系列事件的信息, 即一直发生的事情。如何更新 EventSeries 取决于事件报告者。默认事件报告者在
// "k8s.io/client-go/tools/events/event_broadcaster.go" 中显示了如何在心跳上更新此结构，并可以指导自定义报告程序实现。
type EventSeries struct {
	// count is the number of occurrences in this series up to the last heartbeat time.
	// count 是此系列中直到最后一次心跳时间为止的发生次数。
	Count int32 `json:"count" protobuf:"varint,1,opt,name=count"`
	// lastObservedTime is the time when last Event from the series was seen before last heartbeat.
	// lastObservedTime 是在最后一次心跳之前看到系列中最后一个事件的时间。
	LastObservedTime metav1.MicroTime `json:"lastObservedTime" protobuf:"bytes,2,opt,name=lastObservedTime"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventList is a list of Event objects.
type EventList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// items is a list of schema objects.
	Items []Event `json:"items" protobuf:"bytes,2,rep,name=items"`
}
