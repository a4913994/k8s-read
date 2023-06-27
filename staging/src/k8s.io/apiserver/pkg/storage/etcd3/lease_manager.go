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

package etcd3

import (
	"context"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"k8s.io/apiserver/pkg/storage/etcd3/metrics"
)

const (
	defaultLeaseReuseDurationSeconds = 60
	defaultLeaseMaxObjectCount       = 1000
)

// LeaseManagerConfig is configuration for creating a lease manager.
// LeaseManagerConfig 是用于创建租赁管理器的配置。
type LeaseManagerConfig struct {
	// ReuseDurationSeconds specifies time in seconds that each lease is reused
	// ReuseDurationSeconds 指定每个租约被重用的时间（以秒为单位）
	ReuseDurationSeconds int64
	// MaxObjectCount specifies how many objects that a lease can attach
	// MaxObjectCount 指定租约可以附加多少个对象
	MaxObjectCount int64
}

// NewDefaultLeaseManagerConfig creates a LeaseManagerConfig with default values
// NewDefaultLeaseManagerConfig 创建具有默认值的 LeaseManagerConfig
func NewDefaultLeaseManagerConfig() LeaseManagerConfig {
	return LeaseManagerConfig{
		ReuseDurationSeconds: defaultLeaseReuseDurationSeconds,
		MaxObjectCount:       defaultLeaseMaxObjectCount,
	}
}

// leaseManager is used to manage leases requested from etcd. If a new write
// needs a lease that has similar expiration time to the previous one, the old
// lease will be reused to reduce the overhead of etcd, since lease operations
// are expensive. In the implementation, we only store one previous lease,
// since all the events have the same ttl.
// leaseManager 用于管理从 etcd 请求的租约。如果新的写入需要一个与前一个租约具有相似到期时间的租约，
// 旧租约将被重用以减少 etcd 的开销，因为租约操作是昂贵的。在实现中，我们只存储一个以前的租约，因为所有事件都有相同的 ttl。
type leaseManager struct {
	client                  *clientv3.Client // etcd client used to grant leases
	leaseMu                 sync.Mutex
	prevLeaseID             clientv3.LeaseID
	prevLeaseExpirationTime time.Time
	// The period of time in seconds and percent of TTL that each lease is
	// reused. The minimum of them is used to avoid unreasonably large
	// numbers.
	leaseReuseDurationSeconds   int64
	leaseReuseDurationPercent   float64
	leaseMaxAttachedObjectCount int64
	leaseAttachedObjectCount    int64
}

// newDefaultLeaseManager creates a new lease manager using default setting.
// newDefaultLeaseManager 使用默认设置创建一个新的租赁管理器。
func newDefaultLeaseManager(client *clientv3.Client, config LeaseManagerConfig) *leaseManager {
	if config.MaxObjectCount <= 0 {
		config.MaxObjectCount = defaultLeaseMaxObjectCount
	}
	return newLeaseManager(client, config.ReuseDurationSeconds, 0.05, config.MaxObjectCount)
}

// newLeaseManager creates a new lease manager with the number of buffered
// leases, lease reuse duration in seconds and percentage. The percentage
// value x means x*100%.
// newLeaseManager 创建一个新的租赁管理器，其中包含缓冲租赁的数量、以秒为单位的租赁重用持续时间和百分比。百分比值 x 表示 x100%。
func newLeaseManager(client *clientv3.Client, leaseReuseDurationSeconds int64, leaseReuseDurationPercent float64, maxObjectCount int64) *leaseManager {
	return &leaseManager{
		client:                      client,
		leaseReuseDurationSeconds:   leaseReuseDurationSeconds,
		leaseReuseDurationPercent:   leaseReuseDurationPercent,
		leaseMaxAttachedObjectCount: maxObjectCount,
	}
}

// GetLease returns a lease based on requested ttl: if the cached previous
// lease can be reused, reuse it; otherwise request a new one from etcd.
// GetLease 根据请求的 ttl 返回一个租约：如果之前缓存的租约可以重用，则重用；否则从 etcd 请求一个新的。
func (l *leaseManager) GetLease(ctx context.Context, ttl int64) (clientv3.LeaseID, error) {
	now := time.Now()
	l.leaseMu.Lock()
	defer l.leaseMu.Unlock()
	// check if previous lease can be reused
	reuseDurationSeconds := l.getReuseDurationSecondsLocked(ttl)
	valid := now.Add(time.Duration(ttl) * time.Second).Before(l.prevLeaseExpirationTime)
	sufficient := now.Add(time.Duration(ttl+reuseDurationSeconds) * time.Second).After(l.prevLeaseExpirationTime)

	// We count all operations that happened in the same lease, regardless of success or failure.
	// Currently each GetLease call only attach 1 object
	l.leaseAttachedObjectCount++

	if valid && sufficient && l.leaseAttachedObjectCount <= l.leaseMaxAttachedObjectCount {
		return l.prevLeaseID, nil
	}

	// request a lease with a little extra ttl from etcd
	ttl += reuseDurationSeconds
	lcr, err := l.client.Lease.Grant(ctx, ttl)
	if err != nil {
		return clientv3.LeaseID(0), err
	}
	// cache the new lease id
	l.prevLeaseID = lcr.ID
	l.prevLeaseExpirationTime = now.Add(time.Duration(ttl) * time.Second)
	// refresh count
	metrics.UpdateLeaseObjectCount(l.leaseAttachedObjectCount)
	l.leaseAttachedObjectCount = 1
	return lcr.ID, nil
}

// getReuseDurationSecondsLocked returns the reusable duration in seconds
// based on the configuration. Lock has to be acquired before calling this
// function.
// getReuseDurationSecondsLocked 根据配置返回可重复使用的持续时间（以秒为单位）。在调用此函数之前必须获取锁。
func (l *leaseManager) getReuseDurationSecondsLocked(ttl int64) int64 {
	reuseDurationSeconds := int64(l.leaseReuseDurationPercent * float64(ttl))
	if reuseDurationSeconds > l.leaseReuseDurationSeconds {
		reuseDurationSeconds = l.leaseReuseDurationSeconds
	}
	return reuseDurationSeconds
}
