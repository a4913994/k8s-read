/*
Copyright 2018 The Kubernetes Authors.

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

package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClientConnectionConfiguration contains details for constructing a client.
// ClientConnectionConfiguration 包含构建客户端的详细信息。
type ClientConnectionConfiguration struct {
	// kubeconfig is the path to a KubeConfig file.
	// kubeconfig 是 KubeConfig 文件的路径。
	Kubeconfig string
	// acceptContentTypes defines the Accept header sent by clients when connecting to a server, overriding the
	// default value of 'application/json'. This field will control all connections to the server used by a particular
	// client.
	// acceptContentTypes 定义客户端在连接到服务器时发送的 Accept 标头，覆盖“application/json”的默认值。该字段将控制特定客户端使用的与服务器的所有连接。
	AcceptContentTypes string
	// contentType is the content type used when sending data to the server from this client.
	// contentType 是从该客户端向服务器发送数据时使用的内容类型。
	ContentType string
	// qps controls the number of queries per second allowed for this connection.
	// qps 控制此连接允许的每秒查询数。
	QPS float32
	// burst allows extra queries to accumulate when a client is exceeding its rate.
	// Burst允许在客户机超过其速率时累积额外的查询。
	Burst int32
}

// LeaderElectionConfiguration defines the configuration of leader election
// clients for components that can run with leader election enabled.
// LeaderElection Configuration 为可以在启用领导者选举的情况下运行的组件定义领导者选举客户端的配置
type LeaderElectionConfiguration struct {
	// leaderElect enables a leader election client to gain leadership
	// before executing the main loop. Enable this when running replicated
	// components for high availability.
	// LeaderElect 使领导者选举客户端能够在执行主循环之前获得领导权。在运行复制组件以实现高可用性时启用此功能。
	LeaderElect bool
	// leaseDuration is the duration that non-leader candidates will wait
	// after observing a leadership renewal until attempting to acquire
	// leadership of a led but unrenewed leader slot. This is effectively the
	// maximum duration that a leader can be stopped before it is replaced
	// by another candidate. This is only applicable if leader election is
	// enabled.
	// leaseDuration 是非领导候选人在观察到领导更新后等待的持续时间，直到尝试获得领导但未更新的领导位置的领导。
	// 这实际上是领导者在被另一个候选人取代之前可以停止的最长持续时间。这仅适用于启用领导者选举的情况。
	LeaseDuration metav1.Duration
	// renewDeadline is the interval between attempts by the acting master to
	// renew a leadership slot before it stops leading. This must be less
	// than or equal to the lease duration. This is only applicable if leader
	// election is enabled.
	// renewDeadline 是代理主机在停止领导之前尝试更新领导位置之间的间隔。这必须小于或等于租用期限。这仅适用于启用领导者选举的情况。
	RenewDeadline metav1.Duration
	// retryPeriod is the duration the clients should wait between attempting
	// acquisition and renewal of a leadership. This is only applicable if
	// leader election is enabled.
	// retryPeriod 是客户在尝试获取领导权和更新领导权之间应等待的时间。这仅适用于启用领导者选举的情况。
	RetryPeriod metav1.Duration
	// resourceLock indicates the resource object type that will be used to lock
	// during leader election cycles.
	// resourceLock 指示将在领导者选举周期中用于锁定的资源对象类型。
	ResourceLock string
	// resourceName indicates the name of resource object that will be used to lock
	// during leader election cycles.
	// resourceName 指示将在领导者选举周期中用于锁定的资源对象的名称。
	ResourceName string
	// resourceNamespace indicates the namespace of resource object that will be used to lock
	// during leader election cycles.
	// resourceNamespace 指示将用于在领导者选举周期中锁定的资源对象的名称空间。
	ResourceNamespace string
}

// DebuggingConfiguration holds configuration for Debugging related features.
// DebuggingConfiguration 保存调试相关功能的配置。
type DebuggingConfiguration struct {
	// enableProfiling enables profiling via web interface host:port/debug/pprof/
	EnableProfiling bool
	// enableContentionProfiling enables lock contention profiling, if
	// enableProfiling is true.
	EnableContentionProfiling bool
}
