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

package config

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PersistentVolumeBinderControllerConfiguration contains elements describing
// PersistentVolumeBinderController.
type PersistentVolumeBinderControllerConfiguration struct {
	// pvClaimBinderSyncPeriod is the period for syncing persistent volumes
	// and persistent volume claims.
	// pvClaimBinderSyncPeriod是同步持久化卷和持久化卷要求的周期。
	PVClaimBinderSyncPeriod metav1.Duration
	// volumeConfiguration holds configuration for volume related features.
	// volumeConfiguration持有与卷有关的功能的配置。
	VolumeConfiguration VolumeConfiguration
	// VolumeHostCIDRDenylist is a list of CIDRs that should not be reachable by the
	// controller from plugins.
	// VolumeHostCIDRDenylist是一个CIDRs的列表，这些CIDRs不应该被控制器从插件中到达。
	VolumeHostCIDRDenylist []string
	// VolumeHostAllowLocalLoopback indicates if local loopback hosts (127.0.0.1, etc)
	// should be allowed from plugins.
	// VolumeHostAllowLocalLoopback表示是否应该允许来自插件的本地回环主机（127.0.0.1，等等）。
	VolumeHostAllowLocalLoopback bool
}

// VolumeConfiguration contains *all* enumerated flags meant to configure all volume
// plugins. From this config, the controller-manager binary will create many instances of
// volume.VolumeConfig, each containing only the configuration needed for that plugin which
// are then passed to the appropriate plugin. The ControllerManager binary is the only part
// of the code which knows what plugins are supported and which flags correspond to each plugin.
// VolumeConfiguration包含了所有用于配置所有卷插件的枚举标志。
// 从这个配置中，Controller-manager二进制文件将创建许多volume.VolumeConfig的实例，
// 每个实例只包含该插件所需的配置，然后传递给相应的插件。
// ControllerManager二进制文件是代码中唯一知道支持哪些插件以及每个插件对应哪些标志的部分。
type VolumeConfiguration struct {
	// enableHostPathProvisioning enables HostPath PV provisioning when running without a
	// cloud provider. This allows testing and development of provisioning features. HostPath
	// provisioning is not supported in any way, won't work in a multi-node cluster, and
	// should not be used for anything other than testing or development.
	// enableHostPathProvisioning在没有云提供商的情况下运行时，启用HostPath PV配置。
	// 这允许对配置功能进行测试和开发。HostPath配置不受任何支持，不会在多节点集群中工作，并且不应该用于测试或开发以外的其他用途。
	EnableHostPathProvisioning bool
	// enableDynamicProvisioning enables the provisioning of volumes when running within an environment
	// that supports dynamic provisioning. Defaults to true.
	// enableDynamicProvisioning在支持动态配置的环境中运行时，启用卷的配置。默认为true。
	EnableDynamicProvisioning bool
	// persistentVolumeRecyclerConfiguration holds configuration for persistent volume plugins.
	// persistentVolumeRecyclerConfiguration持有持久化卷插件的配置。
	PersistentVolumeRecyclerConfiguration PersistentVolumeRecyclerConfiguration
	// volumePluginDir is the full path of the directory in which the flex
	// volume plugin should search for additional third party volume plugins
	// volumePluginDir是目录的完整路径，flex volume插件应该在该目录中搜索额外的第三方体积插件。
	FlexVolumePluginDir string
}

// PersistentVolumeRecyclerConfiguration contains elements describing persistent volume plugins.
type PersistentVolumeRecyclerConfiguration struct {
	// maximumRetry is number of retries the PV recycler will execute on failure to recycle
	// PV.
	// maximumRetry是PV回收器在回收PV失败时将执行的重试次数。
	MaximumRetry int32
	// minimumTimeoutNFS is the minimum ActiveDeadlineSeconds to use for an NFS Recycler
	// pod.
	// minimumTimeoutNFS是用于NFS回收器pod的最小ActiveDeadlineSeconds。
	MinimumTimeoutNFS int32
	// podTemplateFilePathNFS is the file path to a pod definition used as a template for
	// NFS persistent volume recycling
	// podTemplateFilePathNFS是作为NFS持久卷回收模板的pod定义的文件路径。
	PodTemplateFilePathNFS string
	// incrementTimeoutNFS is the increment of time added per Gi to ActiveDeadlineSeconds
	// for an NFS scrubber pod.
	// incrementTimeoutNFS是指NFS洗涤器的每个Gi添加到ActiveDeadlineSeconds的时间增量。
	IncrementTimeoutNFS int32
	// podTemplateFilePathHostPath is the file path to a pod definition used as a template for
	// HostPath persistent volume recycling. This is for development and testing only and
	// will not work in a multi-node cluster.
	// podTemplateFilePathHostPath是一个pod定义的文件路径，用作HostPath持久化卷回收的模板。这仅用于开发和测试，在多节点集群中不会工作。
	PodTemplateFilePathHostPath string
	// minimumTimeoutHostPath is the minimum ActiveDeadlineSeconds to use for a HostPath
	// Recycler pod.  This is for development and testing only and will not work in a multi-node
	// cluster.
	// minimumTimeoutHostPath是用于HostPath回收器pod的最小ActiveDeadlineSeconds。这仅用于开发和测试，在多节点集群中不会工作。
	MinimumTimeoutHostPath int32
	// incrementTimeoutHostPath is the increment of time added per Gi to ActiveDeadlineSeconds
	// for a HostPath scrubber pod.  This is for development and testing only and will not work
	// in a multi-node cluster.
	// incrementTimeoutHostPath是指HostPath刷子的每Gi向ActiveDeadlineSeconds增加的时间增量。这仅用于开发和测试，在多节点集群中不会起作用。
	IncrementTimeoutHostPath int32
}
