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

package kubelet

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"

	cadvisorapiv1 "github.com/google/cadvisor/info/v1"
	cadvisorv2 "github.com/google/cadvisor/info/v2"
	"k8s.io/klog/v2"
	"k8s.io/mount-utils"
	utilpath "k8s.io/utils/path"
	utilstrings "k8s.io/utils/strings"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/kubelet/cm"
	"k8s.io/kubernetes/pkg/kubelet/config"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	kubelettypes "k8s.io/kubernetes/pkg/kubelet/types"
	utilnode "k8s.io/kubernetes/pkg/util/node"
	"k8s.io/kubernetes/pkg/volume/csi"
)

// getRootDir returns the full path to the directory under which kubelet can
// store data.  These functions are useful to pass interfaces to other modules
// that may need to know where to write data without getting a whole kubelet
// instance.
// getRootDir 返回 kubelet 可以存储数据的目录的完整路径。这些函数对于将接口传递给其他模块非常有用，
// 这些模块可能需要知道在哪里写入数据而不必获取整个 kubelet 实例。
func (kl *Kubelet) getRootDir() string {
	return kl.rootDirectory
}

// getPodsDir returns the full path to the directory under which pod
// directories are created.
// getPodsDir 返回 pod 目录创建的目录的完整路径。
func (kl *Kubelet) getPodsDir() string {
	return filepath.Join(kl.getRootDir(), config.DefaultKubeletPodsDirName)
}

// getPluginsDir returns the full path to the directory under which plugin
// directories are created.  Plugins can use these directories for data that
// they need to persist.  Plugins should create subdirectories under this named
// after their own names.
// getPluginsDir 返回插件目录创建的目录的完整路径。插件可以使用这些目录来存储它们需要持久化的数据。
func (kl *Kubelet) getPluginsDir() string {
	return filepath.Join(kl.getRootDir(), config.DefaultKubeletPluginsDirName)
}

// getPluginsRegistrationDir returns the full path to the directory under which
// plugins socket should be placed to be registered.
// More information is available about plugin registration in the pluginwatcher
// module
// getPluginsRegistrationDir 返回插件套接字应放置以进行注册的目录的完整路径。
func (kl *Kubelet) getPluginsRegistrationDir() string {
	return filepath.Join(kl.getRootDir(), config.DefaultKubeletPluginsRegistrationDirName)
}

// getPluginDir returns a data directory name for a given plugin name.
// Plugins can use these directories to store data that they need to persist.
// For per-pod plugin data, see getPodPluginDir.
// getPluginDir 返回给定插件名称的数据目录名称。插件可以使用这些目录来存储它们需要持久化的数据。
func (kl *Kubelet) getPluginDir(pluginName string) string {
	return filepath.Join(kl.getPluginsDir(), pluginName)
}

// getCheckpointsDir returns a data directory name for checkpoints.
// Checkpoints can be stored in this directory for further use.
// getCheckpointsDir 返回检查点的数据目录名称。检查点可以存储在此目录中以供进一步使用。
func (kl *Kubelet) getCheckpointsDir() string {
	return filepath.Join(kl.getRootDir(), config.DefaultKubeletCheckpointsDirName)
}

// getVolumeDevicePluginsDir returns the full path to the directory under which plugin
// directories are created.  Plugins can use these directories for data that
// they need to persist.  Plugins should create subdirectories under this named
// after their own names.
// getVolumeDevicePluginsDir 返回插件目录创建的目录的完整路径。插件可以使用这些目录来存储它们需要持久化的数据。
func (kl *Kubelet) getVolumeDevicePluginsDir() string {
	return filepath.Join(kl.getRootDir(), config.DefaultKubeletPluginsDirName)
}

// getVolumeDevicePluginDir returns a data directory name for a given plugin name.
// Plugins can use these directories to store data that they need to persist.
// For per-pod plugin data, see getVolumeDevicePluginsDir.
// getVolumeDevicePluginDir 返回给定插件名称的数据目录名称。插件可以使用这些目录来存储它们需要持久化的数据。
func (kl *Kubelet) getVolumeDevicePluginDir(pluginName string) string {
	return filepath.Join(kl.getVolumeDevicePluginsDir(), pluginName, config.DefaultKubeletVolumeDevicesDirName)
}

// GetPodDir returns the full path to the per-pod data directory for the
// specified pod. This directory may not exist if the pod does not exist.
// GetPodDir 返回指定 pod 的 per-pod 数据目录的完整路径。如果 pod 不存在，则此目录可能不存在。
func (kl *Kubelet) GetPodDir(podUID types.UID) string {
	return kl.getPodDir(podUID)
}

// getPodDir returns the full path to the per-pod directory for the pod with
// the given UID.
// getPodDir 返回具有给定 UID 的 pod 的 per-pod 目录的完整路径。
func (kl *Kubelet) getPodDir(podUID types.UID) string {
	return filepath.Join(kl.getPodsDir(), string(podUID))
}

// getPodVolumesSubpathsDir returns the full path to the per-pod subpaths directory under
// which subpath volumes are created for the specified pod.  This directory may not
// exist if the pod does not exist or subpaths are not specified.
// getPodVolumesSubpathsDir 返回 per-pod subpaths 目录的完整路径，该目录为指定的 pod 创建 subpath 卷。
func (kl *Kubelet) getPodVolumeSubpathsDir(podUID types.UID) string {
	return filepath.Join(kl.getPodDir(podUID), config.DefaultKubeletVolumeSubpathsDirName)
}

// getPodVolumesDir returns the full path to the per-pod data directory under
// which volumes are created for the specified pod.  This directory may not
// exist if the pod does not exist.
// getPodVolumesDir 返回 per-pod 数据目录的完整路径，该目录为指定的 pod 创建卷。
func (kl *Kubelet) getPodVolumesDir(podUID types.UID) string {
	return filepath.Join(kl.getPodDir(podUID), config.DefaultKubeletVolumesDirName)
}

// getPodVolumeDir returns the full path to the directory which represents the
// named volume under the named plugin for specified pod.  This directory may not
// exist if the pod does not exist.
// getPodVolumeDir 返回表示指定 pod 下的命名插件下的命名卷的目录的完整路径。
func (kl *Kubelet) getPodVolumeDir(podUID types.UID, pluginName string, volumeName string) string {
	return filepath.Join(kl.getPodVolumesDir(podUID), pluginName, volumeName)
}

// getPodVolumeDevicesDir returns the full path to the per-pod data directory under
// which volumes are created for the specified pod. This directory may not
// exist if the pod does not exist.
// getPodVolumeDevicesDir 返回 per-pod 数据目录的完整路径，该目录为指定的 pod 创建卷。
func (kl *Kubelet) getPodVolumeDevicesDir(podUID types.UID) string {
	return filepath.Join(kl.getPodDir(podUID), config.DefaultKubeletVolumeDevicesDirName)
}

// getPodVolumeDeviceDir returns the full path to the directory which represents the
// named plugin for specified pod. This directory may not exist if the pod does not exist.
// getPodVolumeDeviceDir 返回表示指定 pod 下的命名插件的目录的完整路径。
func (kl *Kubelet) getPodVolumeDeviceDir(podUID types.UID, pluginName string) string {
	return filepath.Join(kl.getPodVolumeDevicesDir(podUID), pluginName)
}

// getPodPluginsDir returns the full path to the per-pod data directory under
// which plugins may store data for the specified pod.  This directory may not
// exist if the pod does not exist.
// getPodPluginsDir 返回 per-pod 数据目录的完整路径，该目录下的插件可以为指定的 pod 存储数据。
func (kl *Kubelet) getPodPluginsDir(podUID types.UID) string {
	return filepath.Join(kl.getPodDir(podUID), config.DefaultKubeletPluginsDirName)
}

// getPodPluginDir returns a data directory name for a given plugin name for a
// given pod UID.  Plugins can use these directories to store data that they
// need to persist.  For non-per-pod plugin data, see getPluginDir.
// getPodPluginDir 返回给定 pod UID 的给定插件名称的数据目录名称。
func (kl *Kubelet) getPodPluginDir(podUID types.UID, pluginName string) string {
	return filepath.Join(kl.getPodPluginsDir(podUID), pluginName)
}

// getPodContainerDir returns the full path to the per-pod data directory under
// which container data is held for the specified pod.  This directory may not
// exist if the pod or container does not exist.
// getPodContainerDir 返回 per-pod 数据目录的完整路径，该目录下的容器数据为指定的 pod 所保留。
func (kl *Kubelet) getPodContainerDir(podUID types.UID, ctrName string) string {
	return filepath.Join(kl.getPodDir(podUID), config.DefaultKubeletContainersDirName, ctrName)
}

// getPodResourcesSocket returns the full path to the directory containing the pod resources socket
// getPodResourcesSocket返回包含 pod 资源套接字的目录的完整路径
func (kl *Kubelet) getPodResourcesDir() string {
	return filepath.Join(kl.getRootDir(), config.DefaultKubeletPodResourcesDirName)
}

// GetPods returns all pods bound to the kubelet and their spec, and the mirror
// pods.
// GetPods 返回绑定到 kubelet 和其规范的所有 pod，以及镜像 pod。
func (kl *Kubelet) GetPods() []*v1.Pod {
	pods := kl.podManager.GetPods()
	// a kubelet running without apiserver requires an additional
	// update of the static pod status. See #57106
	for _, p := range pods {
		if kubelettypes.IsStaticPod(p) {
			if status, ok := kl.statusManager.GetPodStatus(p.UID); ok {
				klog.V(2).InfoS("Pod status updated", "pod", klog.KObj(p), "status", status.Phase)
				p.Status = status
			}
		}
	}
	return pods
}

// GetRunningPods returns all pods running on kubelet from looking at the
// container runtime cache. This function converts kubecontainer.Pod to
// v1.Pod, so only the fields that exist in both kubecontainer.Pod and
// v1.Pod are considered meaningful.
// GetRunningPods 从容器运行时缓存中查看 kubelet 上运行的所有 pod。 此函数将 kubecontainer.Pod 转换为 v1.Pod，
// 因此只有在 kubecontainer.Pod 和 v1.Pod 中都存在的字段才被认为是有意义的。
func (kl *Kubelet) GetRunningPods(ctx context.Context) ([]*v1.Pod, error) {
	pods, err := kl.runtimeCache.GetPods(ctx)
	if err != nil {
		return nil, err
	}

	apiPods := make([]*v1.Pod, 0, len(pods))
	for _, pod := range pods {
		apiPods = append(apiPods, pod.ToAPIPod())
	}
	return apiPods, nil
}

// GetPodByFullName gets the pod with the given 'full' name, which
// incorporates the namespace as well as whether the pod was found.
// GetPodByFullName 获取具有给定“完整”名称的 pod，该名称包括命名空间以及是否找到了 pod。
func (kl *Kubelet) GetPodByFullName(podFullName string) (*v1.Pod, bool) {
	return kl.podManager.GetPodByFullName(podFullName)
}

// GetPodByName provides the first pod that matches namespace and name, as well
// as whether the pod was found.
// GetPodByName 提供与命名空间和名称匹配的第一个 pod，以及是否找到了 pod。
func (kl *Kubelet) GetPodByName(namespace, name string) (*v1.Pod, bool) {
	return kl.podManager.GetPodByName(namespace, name)
}

// GetPodByCgroupfs provides the pod that maps to the specified cgroup, as well
// as whether the pod was found.
// GetPodByCgroupfs 提供映射到指定 cgroup 的 pod，以及是否找到了 pod。
func (kl *Kubelet) GetPodByCgroupfs(cgroupfs string) (*v1.Pod, bool) {
	pcm := kl.containerManager.NewPodContainerManager()
	if result, podUID := pcm.IsPodCgroup(cgroupfs); result {
		return kl.podManager.GetPodByUID(podUID)
	}
	return nil, false
}

// GetHostname Returns the hostname as the kubelet sees it.
// GetHostname 返回 kubelet 所看到的主机名。
func (kl *Kubelet) GetHostname() string {
	return kl.hostname
}

// getRuntime returns the current Runtime implementation in use by the kubelet.
// getRuntime 返回 kubelet 正在使用的当前 Runtime 实现。
func (kl *Kubelet) getRuntime() kubecontainer.Runtime {
	return kl.containerRuntime
}

// GetNode returns the node info for the configured node name of this Kubelet.
// GetNode 返回此 Kubelet 配置的节点名称的节点信息。
func (kl *Kubelet) GetNode() (*v1.Node, error) {
	if kl.kubeClient == nil {
		return kl.initialNode(context.TODO())
	}
	return kl.nodeLister.Get(string(kl.nodeName))
}

// getNodeAnyWay() must return a *v1.Node which is required by RunGeneralPredicates().
// The *v1.Node is obtained as follows:
// Return kubelet's nodeInfo for this node, except on error or if in standalone mode,
// in which case return a manufactured nodeInfo representing a node with no pods,
// zero capacity, and the default labels.
// getNodeAnyWay() 必须返回 RunGeneralPredicates() 所需的 *v1.Node。 *v1.Node 通过以下方式获取：
// 返回此节点的 kubelet 的 nodeInfo，除非出错或处于独立模式，否则返回表示没有 pod、零容量和默认标签的节点的 nodeInfo。
func (kl *Kubelet) getNodeAnyWay() (*v1.Node, error) {
	if kl.kubeClient != nil {
		if n, err := kl.nodeLister.Get(string(kl.nodeName)); err == nil {
			return n, nil
		}
	}
	return kl.initialNode(context.TODO())
}

// GetNodeConfig returns the container manager node config.
// GetNodeConfig 返回容器管理器节点配置。
func (kl *Kubelet) GetNodeConfig() cm.NodeConfig {
	return kl.containerManager.GetNodeConfig()
}

// GetPodCgroupRoot returns the listeral cgroupfs value for the cgroup containing all pods
// GetPodCgroupRoot 返回包含所有 pod 的 cgroup 的 listeral cgroupfs 值
func (kl *Kubelet) GetPodCgroupRoot() string {
	return kl.containerManager.GetPodCgroupRoot()
}

// GetHostIPs returns host IPs or nil in case of error.
// GetHostIPs 返回主机 IP 或错误的情况下返回 nil。
func (kl *Kubelet) GetHostIPs() ([]net.IP, error) {
	node, err := kl.GetNode()
	if err != nil {
		return nil, fmt.Errorf("cannot get node: %v", err)
	}
	return utilnode.GetNodeHostIPs(node)
}

// getHostIPsAnyWay attempts to return the host IPs from kubelet's nodeInfo, or
// the initialNode.
// getHostIPsAnyWay 尝试从 kubelet 的 nodeInfo 或 initialNode 返回主机 IP。
func (kl *Kubelet) getHostIPsAnyWay() ([]net.IP, error) {
	node, err := kl.getNodeAnyWay()
	if err != nil {
		return nil, err
	}
	return utilnode.GetNodeHostIPs(node)
}

// GetExtraSupplementalGroupsForPod returns a list of the extra
// supplemental groups for the Pod. These extra supplemental groups come
// from annotations on persistent volumes that the pod depends on.
// GetExtraSupplementalGroupsForPod 返回 pod 的额外补充组的列表。这些额外的补充组来自 pod 依赖的持久卷上的注释。
func (kl *Kubelet) GetExtraSupplementalGroupsForPod(pod *v1.Pod) []int64 {
	return kl.volumeManager.GetExtraSupplementalGroupsForPod(pod)
}

// getPodVolumePathListFromDisk returns a list of the volume paths by reading the
// volume directories for the given pod from the disk.
// getPodVolumePathListFromDisk 通过从磁盘读取给定 pod 的卷目录来返回卷路径的列表。
func (kl *Kubelet) getPodVolumePathListFromDisk(podUID types.UID) ([]string, error) {
	volumes := []string{}
	podVolDir := kl.getPodVolumesDir(podUID)

	if pathExists, pathErr := mount.PathExists(podVolDir); pathErr != nil {
		return volumes, fmt.Errorf("error checking if path %q exists: %v", podVolDir, pathErr)
	} else if !pathExists {
		klog.V(6).InfoS("Path does not exist", "path", podVolDir)
		return volumes, nil
	}

	volumePluginDirs, err := os.ReadDir(podVolDir)
	if err != nil {
		klog.ErrorS(err, "Could not read directory", "path", podVolDir)
		return volumes, err
	}
	for _, volumePluginDir := range volumePluginDirs {
		volumePluginName := volumePluginDir.Name()
		volumePluginPath := filepath.Join(podVolDir, volumePluginName)
		volumeDirs, err := utilpath.ReadDirNoStat(volumePluginPath)
		if err != nil {
			return volumes, fmt.Errorf("could not read directory %s: %v", volumePluginPath, err)
		}
		unescapePluginName := utilstrings.UnescapeQualifiedName(volumePluginName)

		if unescapePluginName != csi.CSIPluginName {
			for _, volumeDir := range volumeDirs {
				volumes = append(volumes, filepath.Join(volumePluginPath, volumeDir))
			}
		} else {
			// For CSI volumes, the mounted volume path has an extra sub path "/mount", so also add it
			// to the list if the mounted path exists.
			for _, volumeDir := range volumeDirs {
				path := filepath.Join(volumePluginPath, volumeDir)
				csimountpath := csi.GetCSIMounterPath(path)
				if pathExists, _ := mount.PathExists(csimountpath); pathExists {
					volumes = append(volumes, csimountpath)
				}
			}
		}
	}
	return volumes, nil
}

func (kl *Kubelet) getMountedVolumePathListFromDisk(podUID types.UID) ([]string, error) {
	mountedVolumes := []string{}
	volumePaths, err := kl.getPodVolumePathListFromDisk(podUID)
	if err != nil {
		return mountedVolumes, err
	}
	// Only use IsLikelyNotMountPoint to check might not cover all cases. For CSI volumes that
	// either: 1) don't mount or 2) bind mount in the rootfs, the mount check will not work as expected.
	// We plan to remove this mountpoint check as a condition before deleting pods since it is
	// not reliable and the condition might be different for different types of volumes. But it requires
	// a reliable way to clean up unused volume dir to avoid problems during pod deletion. See discussion in issue #74650
	for _, volumePath := range volumePaths {
		isNotMount, err := kl.mounter.IsLikelyNotMountPoint(volumePath)
		if err != nil {
			return mountedVolumes, fmt.Errorf("fail to check mount point %q: %v", volumePath, err)
		}
		if !isNotMount {
			mountedVolumes = append(mountedVolumes, volumePath)
		}
	}
	return mountedVolumes, nil
}

// getPodVolumeSubpathListFromDisk returns a list of the volume-subpath paths by reading the
// subpath directories for the given pod from the disk.
func (kl *Kubelet) getPodVolumeSubpathListFromDisk(podUID types.UID) ([]string, error) {
	volumes := []string{}
	podSubpathsDir := kl.getPodVolumeSubpathsDir(podUID)

	if pathExists, pathErr := mount.PathExists(podSubpathsDir); pathErr != nil {
		return nil, fmt.Errorf("error checking if path %q exists: %v", podSubpathsDir, pathErr)
	} else if !pathExists {
		return volumes, nil
	}

	// Explicitly walks /<volume>/<container name>/<subPathIndex>
	volumePluginDirs, err := os.ReadDir(podSubpathsDir)
	if err != nil {
		klog.ErrorS(err, "Could not read directory", "path", podSubpathsDir)
		return volumes, err
	}
	for _, volumePluginDir := range volumePluginDirs {
		volumePluginName := volumePluginDir.Name()
		volumePluginPath := filepath.Join(podSubpathsDir, volumePluginName)
		containerDirs, err := os.ReadDir(volumePluginPath)
		if err != nil {
			return volumes, fmt.Errorf("could not read directory %s: %v", volumePluginPath, err)
		}
		for _, containerDir := range containerDirs {
			containerName := containerDir.Name()
			containerPath := filepath.Join(volumePluginPath, containerName)
			// Switch to ReadDirNoStat at the subPathIndex level to prevent issues with stat'ing
			// mount points that may not be responsive
			subPaths, err := utilpath.ReadDirNoStat(containerPath)
			if err != nil {
				return volumes, fmt.Errorf("could not read directory %s: %v", containerPath, err)
			}
			for _, subPathDir := range subPaths {
				volumes = append(volumes, filepath.Join(containerPath, subPathDir))
			}
		}
	}
	return volumes, nil
}

// GetRequestedContainersInfo returns container info.
func (kl *Kubelet) GetRequestedContainersInfo(containerName string, options cadvisorv2.RequestOptions) (map[string]*cadvisorapiv1.ContainerInfo, error) {
	return kl.cadvisor.GetRequestedContainersInfo(containerName, options)
}

// GetVersionInfo returns information about the version of cAdvisor in use.
func (kl *Kubelet) GetVersionInfo() (*cadvisorapiv1.VersionInfo, error) {
	return kl.cadvisor.VersionInfo()
}

// GetCachedMachineInfo assumes that the machine info can't change without a reboot
func (kl *Kubelet) GetCachedMachineInfo() (*cadvisorapiv1.MachineInfo, error) {
	kl.machineInfoLock.RLock()
	defer kl.machineInfoLock.RUnlock()
	return kl.machineInfo, nil
}

func (kl *Kubelet) setCachedMachineInfo(info *cadvisorapiv1.MachineInfo) {
	kl.machineInfoLock.Lock()
	defer kl.machineInfoLock.Unlock()
	kl.machineInfo = info
}
