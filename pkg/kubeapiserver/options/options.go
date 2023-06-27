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

package options

import (
	"net"

	utilnet "k8s.io/apimachinery/pkg/util/net"
	netutils "k8s.io/utils/net"
)

// DefaultServiceNodePortRange is the default port range for NodePort services.
// DefaultServiceNodePortRange 是 NodePort 服务的默认端口范围
var DefaultServiceNodePortRange = utilnet.PortRange{Base: 30000, Size: 2768}

// DefaultServiceIPCIDR is a CIDR notation of IP range from which to allocate service cluster IPs
// DefaultServiceIPCIDR 是从中分配服务集群 IP 的 IP 范围的 CIDR 表示法
var DefaultServiceIPCIDR = net.IPNet{IP: netutils.ParseIPSloppy("10.0.0.0"), Mask: net.CIDRMask(24, 32)}

// DefaultEtcdPathPrefix is the default key prefix of etcd for API Server
// DefaultEtcdPathPrefix 是 API 服务器的 etcd 的默认键前缀
const DefaultEtcdPathPrefix = "/registry"
