/*
Copyright 2016 The Kubernetes Authors.

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

package storagebackend

import (
	"time"

	oteltrace "go.opentelemetry.io/otel/trace"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/server/egressselector"
	"k8s.io/apiserver/pkg/storage/etcd3"
	"k8s.io/apiserver/pkg/storage/value"
	flowcontrolrequest "k8s.io/apiserver/pkg/util/flowcontrol/request"
)

const (
	StorageTypeUnset = ""
	StorageTypeETCD2 = "etcd2"
	StorageTypeETCD3 = "etcd3"

	DefaultCompactInterval      = 5 * time.Minute
	DefaultDBMetricPollInterval = 30 * time.Second
	DefaultHealthcheckTimeout   = 2 * time.Second
	DefaultReadinessTimeout     = 2 * time.Second
)

// TransportConfig holds all connection related info,  i.e. equal TransportConfig means equal servers we talk to.
type TransportConfig struct {
	// ServerList is the list of storage servers to connect with.
	// ServerList 是要连接的存储服务器列表。
	ServerList []string
	// TLS credentials
	KeyFile       string
	CertFile      string
	TrustedCAFile string
	// function to determine the egress dialer. (i.e. konnectivity server dialer)
	// 函数来确定出口拨号器。 （即 konnectivity 服务器拨号器）
	EgressLookup egressselector.Lookup
	// The TracerProvider can add tracing the connection
	// TracerProvider 可以添加跟踪连接
	TracerProvider oteltrace.TracerProvider
}

// Config is configuration for creating a storage backend.
// Config 是创建存储后端的配置。
type Config struct {
	// Type defines the type of storage backend. Default ("") is "etcd3".
	// Type 定义存储后端的类型。默认（“”）是“etcd3”。
	Type string
	// Prefix is the prefix to all keys passed to storage.Interface methods.
	// Prefix 是传递给 storage.Interface 方法的所有键的前缀。
	Prefix string
	// Transport holds all connection related info, i.e. equal TransportConfig means equal servers we talk to.
	// Transport 包含所有与连接相关的信息，即相等的 TransportConfig 意味着我们与之交谈的服务器相等。
	Transport TransportConfig
	// Paging indicates whether the server implementation should allow paging (if it is
	// supported). This is generally configured by feature gating, or by a specific
	// resource type not wishing to allow paging, and is not intended for end users to
	// set.
	// 分页指示服务器实现是否应允许分页（如果支持）。这通常由功能门控配置，或者由不希望允许分页的特定资源类型配置，并且不适合最终用户设置。
	Paging bool

	Codec runtime.Codec
	// EncodeVersioner is the same groupVersioner used to build the
	// storage encoder. Given a list of kinds the input object might belong
	// to, the EncodeVersioner outputs the gvk the object will be
	// converted to before persisted in etcd.
	// EncodeVersioner 与用于构建存储编码器的 groupVersioner 相同。给定输入对象可能属于的种类列表，EncodeVersioner 输出对象在持久化到 etcd 之前将转换成的 gvk。
	EncodeVersioner runtime.GroupVersioner
	// Transformer allows the value to be transformed prior to persisting into etcd.
	// Transformer 允许在持久化到 etcd 之前转换值。
	Transformer value.Transformer

	// CompactionInterval is an interval of requesting compaction from apiserver.
	// If the value is 0, no compaction will be issued.
	CompactionInterval time.Duration
	// CountMetricPollPeriod specifies how often should count metric be updated
	CountMetricPollPeriod time.Duration
	// DBMetricPollInterval specifies how often should storage backend metric be updated.
	DBMetricPollInterval time.Duration
	// HealthcheckTimeout specifies the timeout used when checking health
	HealthcheckTimeout time.Duration
	// ReadycheckTimeout specifies the timeout used when checking readiness
	ReadycheckTimeout time.Duration

	LeaseManagerConfig etcd3.LeaseManagerConfig

	// StorageObjectCountTracker is used to keep track of the total
	// number of objects in the storage per resource.
	// StorageObjectCountTracker 用于跟踪每个资源存储中的对象总数。
	StorageObjectCountTracker flowcontrolrequest.StorageObjectCountTracker
}

// ConfigForResource is a Config specialized to a particular `schema.GroupResource`
// ConfigForResource 是专门针对特定 `schema.GroupResource` 的配置
type ConfigForResource struct {
	// Config is the resource-independent configuration
	// config是资源无关的配置
	Config

	// GroupResource is the relevant one
	// GroupResource 是相关的
	GroupResource schema.GroupResource
}

// ForResource specializes to the given resource
func (config *Config) ForResource(resource schema.GroupResource) *ConfigForResource {
	return &ConfigForResource{
		Config:        *config,
		GroupResource: resource,
	}
}

func NewDefaultConfig(prefix string, codec runtime.Codec) *Config {
	return &Config{
		Paging:               true,
		Prefix:               prefix,
		Codec:                codec,
		CompactionInterval:   DefaultCompactInterval,
		DBMetricPollInterval: DefaultDBMetricPollInterval,
		HealthcheckTimeout:   DefaultHealthcheckTimeout,
		ReadycheckTimeout:    DefaultReadinessTimeout,
		LeaseManagerConfig:   etcd3.NewDefaultLeaseManagerConfig(),
		Transport:            TransportConfig{TracerProvider: oteltrace.NewNoopTracerProvider()},
	}
}
