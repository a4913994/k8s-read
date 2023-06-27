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

package kuberuntime

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	cadvisorapi "github.com/google/cadvisor/info/v1"
	crierror "k8s.io/cri-api/pkg/errors"
	"k8s.io/klog/v2"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubetypes "k8s.io/apimachinery/pkg/types"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	utilversion "k8s.io/apimachinery/pkg/util/version"
	utilfeature "k8s.io/apiserver/pkg/util/feature"
	"k8s.io/client-go/tools/record"
	ref "k8s.io/client-go/tools/reference"
	"k8s.io/client-go/util/flowcontrol"
	"k8s.io/component-base/logs/logreduction"
	internalapi "k8s.io/cri-api/pkg/apis"
	runtimeapi "k8s.io/cri-api/pkg/apis/runtime/v1"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	"k8s.io/kubernetes/pkg/credentialprovider"
	"k8s.io/kubernetes/pkg/credentialprovider/plugin"
	"k8s.io/kubernetes/pkg/features"
	"k8s.io/kubernetes/pkg/kubelet/cm"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/events"
	"k8s.io/kubernetes/pkg/kubelet/images"
	runtimeutil "k8s.io/kubernetes/pkg/kubelet/kuberuntime/util"
	"k8s.io/kubernetes/pkg/kubelet/lifecycle"
	"k8s.io/kubernetes/pkg/kubelet/logs"
	"k8s.io/kubernetes/pkg/kubelet/metrics"
	proberesults "k8s.io/kubernetes/pkg/kubelet/prober/results"
	"k8s.io/kubernetes/pkg/kubelet/runtimeclass"
	"k8s.io/kubernetes/pkg/kubelet/sysctl"
	"k8s.io/kubernetes/pkg/kubelet/types"
	"k8s.io/kubernetes/pkg/kubelet/util/cache"
	"k8s.io/kubernetes/pkg/kubelet/util/format"
	sc "k8s.io/kubernetes/pkg/securitycontext"
)

const (
	// The api version of kubelet runtime api
	// kubelet运行时api的api版本
	kubeRuntimeAPIVersion = "0.1.0"
	// The root directory for pod logs
	// pod日志的根目录
	podLogsRootDirectory = "/var/log/pods"
	// A minimal shutdown window for avoiding unnecessary SIGKILLs
	// 避免不必要的SIGKILL的最小关闭窗口
	minimumGracePeriodInSeconds = 2

	// The expiration time of version cache.
	// 版本缓存的过期时间
	versionCacheTTL = 60 * time.Second
	// How frequently to report identical errors
	// 报告相同错误的频率
	identicalErrorDelay = 1 * time.Minute
)

var (
	// ErrVersionNotSupported is returned when the api version of runtime interface is not supported
	// 当运行时接口的api版本不受支持时返回ErrVersionNotSupported
	ErrVersionNotSupported = errors.New("runtime api version is not supported")
)

// podStateProvider can determine if none of the elements are necessary to retain (pod content)
// or if none of the runtime elements are necessary to retain (containers)
// podStateProvider可以确定是否不需要保留任何元素（pod内容）或者是否不需要保留任何运行时元素（容器）
type podStateProvider interface {
	IsPodTerminationRequested(kubetypes.UID) bool
	ShouldPodContentBeRemoved(kubetypes.UID) bool
	ShouldPodRuntimeBeRemoved(kubetypes.UID) bool
}

type kubeGenericRuntimeManager struct {
	runtimeName string
	recorder    record.EventRecorder
	osInterface kubecontainer.OSInterface

	// machineInfo contains the machine information.
	// machineInfo包含机器信息
	machineInfo *cadvisorapi.MachineInfo

	// Container GC manager
	// 容器GC管理器
	containerGC *containerGC

	// Keyring for pulling images
	// 拉取镜像的
	keyring credentialprovider.DockerKeyring

	// Runner of lifecycle events.
	// 生命周期事件的运行者
	runner kubecontainer.HandlerRunner

	// RuntimeHelper that wraps kubelet to generate runtime container options.
	// RuntimeHelper包装kubelet以生成运行时容器选项
	runtimeHelper kubecontainer.RuntimeHelper

	// Health check results.
	// 健康检查结果
	livenessManager  proberesults.Manager
	readinessManager proberesults.Manager
	startupManager   proberesults.Manager

	// If true, enforce container cpu limits with CFS quota support
	// 如果为true，则使用CFS配额支持强制执行容器cpu限制
	cpuCFSQuota bool

	// CPUCFSQuotaPeriod sets the CPU CFS quota period value, cpu.cfs_period_us, defaults to 100ms
	// CPUCFSQuotaPeriod设置CPU CFS配额期值，cpu.cfs_period_us，默认为100ms
	cpuCFSQuotaPeriod metav1.Duration

	// wrapped image puller.
	// 包装的镜像拉取器
	imagePuller images.ImageManager

	// gRPC service clients
	// gRPC服务客户端
	runtimeService internalapi.RuntimeService
	imageService   internalapi.ImageManagerService

	// The version cache of runtime daemon.
	// 运行时守护进程的版本缓存
	versionCache *cache.ObjectCache

	// The directory path for seccomp profiles.
	// seccomp配置文件的目录路径
	seccompProfileRoot string

	// Internal lifecycle event handlers for container resource management.
	// 用于容器资源管理的内部生命周期事件处理程序
	internalLifecycle cm.InternalContainerLifecycle

	// Manage container logs.
	// 管理容器日志
	logManager logs.ContainerLogManager

	// Manage RuntimeClass resources.
	// 管理RuntimeClass资源
	runtimeClassManager *runtimeclass.Manager

	// Cache last per-container error message to reduce log spam
	// 缓存上一个容器的错误消息以减少日志垃圾
	logReduction *logreduction.LogReduction

	// PodState provider instance
	// PodState提供程序实例
	podStateProvider podStateProvider

	// Use RuntimeDefault as the default seccomp profile for all workloads.
	// 将RuntimeDefault用作所有工作负载的默认seccomp配置文件
	seccompDefault bool

	// MemorySwapBehavior defines how swap is used
	// MemorySwapBehavior定义了如何使用swap
	memorySwapBehavior string

	//Function to get node allocatable resources
	// 用于获取节点可分配资源的函数
	getNodeAllocatable func() v1.ResourceList

	// Memory throttling factor for MemoryQoS
	// MemoryQoS的内存节流因子
	memoryThrottlingFactor float64
}

// KubeGenericRuntime is a interface contains interfaces for container runtime and command.
// KubeGenericRuntime是一个包含容器运行时和命令的接口
type KubeGenericRuntime interface {
	kubecontainer.Runtime
	kubecontainer.StreamingRuntime
	kubecontainer.CommandRunner
}

// NewKubeGenericRuntimeManager creates a new kubeGenericRuntimeManager
// NewKubeGenericRuntimeManager创建一个新的kubeGenericRuntimeManager
func NewKubeGenericRuntimeManager(
	recorder record.EventRecorder,
	livenessManager proberesults.Manager,
	readinessManager proberesults.Manager,
	startupManager proberesults.Manager,
	rootDirectory string,
	machineInfo *cadvisorapi.MachineInfo,
	podStateProvider podStateProvider,
	osInterface kubecontainer.OSInterface,
	runtimeHelper kubecontainer.RuntimeHelper,
	insecureContainerLifecycleHTTPClient types.HTTPDoer,
	imageBackOff *flowcontrol.Backoff,
	serializeImagePulls bool,
	imagePullQPS float32,
	imagePullBurst int,
	imageCredentialProviderConfigFile string,
	imageCredentialProviderBinDir string,
	cpuCFSQuota bool,
	cpuCFSQuotaPeriod metav1.Duration,
	runtimeService internalapi.RuntimeService,
	imageService internalapi.ImageManagerService,
	internalLifecycle cm.InternalContainerLifecycle,
	logManager logs.ContainerLogManager,
	runtimeClassManager *runtimeclass.Manager,
	seccompDefault bool,
	memorySwapBehavior string,
	getNodeAllocatable func() v1.ResourceList,
	memoryThrottlingFactor float64,
	podPullingTimeRecorder images.ImagePodPullingTimeRecorder,
) (KubeGenericRuntime, error) {
	ctx := context.Background()
	runtimeService = newInstrumentedRuntimeService(runtimeService)
	imageService = newInstrumentedImageManagerService(imageService)
	kubeRuntimeManager := &kubeGenericRuntimeManager{
		recorder:               recorder,
		cpuCFSQuota:            cpuCFSQuota,
		cpuCFSQuotaPeriod:      cpuCFSQuotaPeriod,
		seccompProfileRoot:     filepath.Join(rootDirectory, "seccomp"),
		livenessManager:        livenessManager,
		readinessManager:       readinessManager,
		startupManager:         startupManager,
		machineInfo:            machineInfo,
		osInterface:            osInterface,
		runtimeHelper:          runtimeHelper,
		runtimeService:         runtimeService,
		imageService:           imageService,
		internalLifecycle:      internalLifecycle,
		logManager:             logManager,
		runtimeClassManager:    runtimeClassManager,
		logReduction:           logreduction.NewLogReduction(identicalErrorDelay),
		seccompDefault:         seccompDefault,
		memorySwapBehavior:     memorySwapBehavior,
		getNodeAllocatable:     getNodeAllocatable,
		memoryThrottlingFactor: memoryThrottlingFactor,
	}

	typedVersion, err := kubeRuntimeManager.getTypedVersion(ctx)
	if err != nil {
		klog.ErrorS(err, "Get runtime version failed")
		return nil, err
	}

	// Only matching kubeRuntimeAPIVersion is supported now
	// TODO: Runtime API machinery is under discussion at https://github.com/kubernetes/kubernetes/issues/28642
	if typedVersion.Version != kubeRuntimeAPIVersion {
		klog.ErrorS(err, "This runtime api version is not supported",
			"apiVersion", typedVersion.Version,
			"supportedAPIVersion", kubeRuntimeAPIVersion)
		return nil, ErrVersionNotSupported
	}

	kubeRuntimeManager.runtimeName = typedVersion.RuntimeName
	klog.InfoS("Container runtime initialized",
		"containerRuntime", typedVersion.RuntimeName,
		"version", typedVersion.RuntimeVersion,
		"apiVersion", typedVersion.RuntimeApiVersion)

	// If the container logs directory does not exist, create it.
	// TODO: create podLogsRootDirectory at kubelet.go when kubelet is refactored to
	// new runtime interface
	if _, err := osInterface.Stat(podLogsRootDirectory); os.IsNotExist(err) {
		if err := osInterface.MkdirAll(podLogsRootDirectory, 0755); err != nil {
			klog.ErrorS(err, "Failed to create pod log directory", "path", podLogsRootDirectory)
		}
	}

	if imageCredentialProviderConfigFile != "" || imageCredentialProviderBinDir != "" {
		if err := plugin.RegisterCredentialProviderPlugins(imageCredentialProviderConfigFile, imageCredentialProviderBinDir); err != nil {
			klog.ErrorS(err, "Failed to register CRI auth plugins")
			os.Exit(1)
		}
	}
	kubeRuntimeManager.keyring = credentialprovider.NewDockerKeyring()

	kubeRuntimeManager.imagePuller = images.NewImageManager(
		kubecontainer.FilterEventRecorder(recorder),
		kubeRuntimeManager,
		imageBackOff,
		serializeImagePulls,
		imagePullQPS,
		imagePullBurst,
		podPullingTimeRecorder)
	kubeRuntimeManager.runner = lifecycle.NewHandlerRunner(insecureContainerLifecycleHTTPClient, kubeRuntimeManager, kubeRuntimeManager, recorder)
	kubeRuntimeManager.containerGC = newContainerGC(runtimeService, podStateProvider, kubeRuntimeManager)
	kubeRuntimeManager.podStateProvider = podStateProvider

	kubeRuntimeManager.versionCache = cache.NewObjectCache(
		func() (interface{}, error) {
			return kubeRuntimeManager.getTypedVersion(ctx)
		},
		versionCacheTTL,
	)

	return kubeRuntimeManager, nil
}

// Type returns the type of the container runtime.
// Type返回容器运行时的类型。
func (m *kubeGenericRuntimeManager) Type() string {
	return m.runtimeName
}

func newRuntimeVersion(version string) (*utilversion.Version, error) {
	if ver, err := utilversion.ParseSemantic(version); err == nil {
		return ver, err
	}
	return utilversion.ParseGeneric(version)
}

func (m *kubeGenericRuntimeManager) getTypedVersion(ctx context.Context) (*runtimeapi.VersionResponse, error) {
	typedVersion, err := m.runtimeService.Version(ctx, kubeRuntimeAPIVersion)
	if err != nil {
		return nil, fmt.Errorf("get remote runtime typed version failed: %v", err)
	}
	return typedVersion, nil
}

// Version returns the version information of the container runtime.
func (m *kubeGenericRuntimeManager) Version(ctx context.Context) (kubecontainer.Version, error) {
	typedVersion, err := m.getTypedVersion(ctx)
	if err != nil {
		return nil, err
	}

	return newRuntimeVersion(typedVersion.RuntimeVersion)
}

// APIVersion returns the cached API version information of the container
// runtime. Implementation is expected to update this cache periodically.
// This may be different from the runtime engine's version.
// APIVersion返回容器运行时的缓存API版本信息。实现应该定期更新此缓存。
// 这可能与运行时引擎的版本不同。
func (m *kubeGenericRuntimeManager) APIVersion() (kubecontainer.Version, error) {
	versionObject, err := m.versionCache.Get(m.machineInfo.MachineID)
	if err != nil {
		return nil, err
	}
	typedVersion := versionObject.(*runtimeapi.VersionResponse)

	return newRuntimeVersion(typedVersion.RuntimeApiVersion)
}

// Status returns the status of the runtime. An error is returned if the Status
// function itself fails, nil otherwise.
// Status返回运行时的状态。如果Status函数本身失败，则返回错误，否则返回nil。
func (m *kubeGenericRuntimeManager) Status(ctx context.Context) (*kubecontainer.RuntimeStatus, error) {
	resp, err := m.runtimeService.Status(ctx, false)
	if err != nil {
		return nil, err
	}
	if resp.GetStatus() == nil {
		return nil, errors.New("runtime status is nil")
	}
	return toKubeRuntimeStatus(resp.GetStatus()), nil
}

// GetPods returns a list of containers grouped by pods. The boolean parameter
// specifies whether the runtime returns all containers including those already
// exited and dead containers (used for garbage collection).
// GetPods返回按pod分组的容器列表。布尔参数指定运行时是否返回所有容器，包括已退出和死亡的容器（用于垃圾回收）。
func (m *kubeGenericRuntimeManager) GetPods(ctx context.Context, all bool) ([]*kubecontainer.Pod, error) {
	pods := make(map[kubetypes.UID]*kubecontainer.Pod)
	sandboxes, err := m.getKubeletSandboxes(ctx, all)
	if err != nil {
		return nil, err
	}
	for i := range sandboxes {
		s := sandboxes[i]
		if s.Metadata == nil {
			klog.V(4).InfoS("Sandbox does not have metadata", "sandbox", s)
			continue
		}
		podUID := kubetypes.UID(s.Metadata.Uid)
		if _, ok := pods[podUID]; !ok {
			pods[podUID] = &kubecontainer.Pod{
				ID:        podUID,
				Name:      s.Metadata.Name,
				Namespace: s.Metadata.Namespace,
			}
		}
		p := pods[podUID]
		converted, err := m.sandboxToKubeContainer(s)
		if err != nil {
			klog.V(4).InfoS("Convert sandbox of pod failed", "runtimeName", m.runtimeName, "sandbox", s, "podUID", podUID, "err", err)
			continue
		}
		p.Sandboxes = append(p.Sandboxes, converted)
		p.CreatedAt = uint64(s.GetCreatedAt())
	}

	containers, err := m.getKubeletContainers(ctx, all)
	if err != nil {
		return nil, err
	}
	for i := range containers {
		c := containers[i]
		if c.Metadata == nil {
			klog.V(4).InfoS("Container does not have metadata", "container", c)
			continue
		}

		labelledInfo := getContainerInfoFromLabels(c.Labels)
		pod, found := pods[labelledInfo.PodUID]
		if !found {
			pod = &kubecontainer.Pod{
				ID:        labelledInfo.PodUID,
				Name:      labelledInfo.PodName,
				Namespace: labelledInfo.PodNamespace,
			}
			pods[labelledInfo.PodUID] = pod
		}

		converted, err := m.toKubeContainer(c)
		if err != nil {
			klog.V(4).InfoS("Convert container of pod failed", "runtimeName", m.runtimeName, "container", c, "podUID", labelledInfo.PodUID, "err", err)
			continue
		}

		pod.Containers = append(pod.Containers, converted)
	}

	// Convert map to list.
	var result []*kubecontainer.Pod
	for _, pod := range pods {
		result = append(result, pod)
	}

	// There are scenarios where multiple pods are running in parallel having
	// the same name, because one of them have not been fully terminated yet.
	// To avoid unexpected behavior on container name based search (for example
	// by calling *Kubelet.findContainer() without specifying a pod ID), we now
	// return the list of pods ordered by their creation time.
	sort.SliceStable(result, func(i, j int) bool {
		return result[i].CreatedAt > result[j].CreatedAt
	})

	return result, nil
}

// containerKillReason explains what killed a given container
// containerKillReason解释了为什么杀死了给定的容器
type containerKillReason string

const (
	reasonStartupProbe        containerKillReason = "StartupProbe"
	reasonLivenessProbe       containerKillReason = "LivenessProbe"
	reasonFailedPostStartHook containerKillReason = "FailedPostStartHook"
	reasonUnknown             containerKillReason = "Unknown"
)

// containerToKillInfo contains necessary information to kill a container.
// containerToKillInfo包含杀死容器所需的必要信息。
type containerToKillInfo struct {
	// The spec of the container.
	container *v1.Container
	// The name of the container.
	name string
	// The message indicates why the container will be killed.
	message string
	// The reason is a clearer source of info on why a container will be killed
	// TODO: replace message with reason?
	reason containerKillReason
}

// podActions keeps information what to do for a pod.
// podActions保存关于pod的信息。
type podActions struct {
	// Stop all running (regular, init and ephemeral) containers and the sandbox for the pod.
	// 停止所有运行的(常规的、初始化的和临时的)容器和pod的沙箱。
	KillPod bool
	// Whether need to create a new sandbox. If needed to kill pod and create
	// a new pod sandbox, all init containers need to be purged (i.e., removed).
	// 是否需要创建新的沙盒。如果需要杀死pod并创建一个新的pod沙箱，所有的init容器都需要被清除(即删除)。
	CreateSandbox bool
	// The id of existing sandbox. It is used for starting containers in ContainersToStart.
	// 已存在沙箱的id。它用于在ContainersToStart中启动容器。
	SandboxID string
	// The attempt number of creating sandboxes for the pod.
	// 为pod创建沙箱的尝试次数。
	Attempt uint32

	// The next init container to start.
	// 要启动的下一个init容器。
	NextInitContainerToStart *v1.Container
	// ContainersToStart keeps a list of indexes for the containers to start,
	// where the index is the index of the specific container in the pod spec (
	// pod.Spec.Containers.
	// ContainersToStart保留要启动的容器的索引列表，其中索引是pod spec中特定容器的索引(pod.Spec.Containers。
	ContainersToStart []int
	// ContainersToKill keeps a map of containers that need to be killed, note that
	// the key is the container ID of the container, while
	// the value contains necessary information to kill a container.
	// ContainersToKill保存了一个需要杀死的容器的映射，注意键是容器的容器ID，而值包含了杀死容器的必要信息。
	ContainersToKill map[kubecontainer.ContainerID]containerToKillInfo
	// EphemeralContainersToStart is a list of indexes for the ephemeral containers to start,
	// where the index is the index of the specific container in pod.Spec.EphemeralContainers.
	// EphemeralContainersToStart是要启动的临时容器的索引列表，其中索引是pod.Spec.EphemeralContainers中特定容器的索引。
	EphemeralContainersToStart []int
}

func containerChanged(container *v1.Container, containerStatus *kubecontainer.Status) (uint64, uint64, bool) {
	expectedHash := kubecontainer.HashContainer(container)
	return expectedHash, containerStatus.Hash, containerStatus.Hash != expectedHash
}

func shouldRestartOnFailure(pod *v1.Pod) bool {
	return pod.Spec.RestartPolicy != v1.RestartPolicyNever
}

func containerSucceeded(c *v1.Container, podStatus *kubecontainer.PodStatus) bool {
	cStatus := podStatus.FindContainerStatusByName(c.Name)
	if cStatus == nil || cStatus.State == kubecontainer.ContainerStateRunning {
		return false
	}
	return cStatus.ExitCode == 0
}

// computePodActions checks whether the pod spec has changed and returns the changes if true.
// computePodActions检查pod spec是否已更改，并在true时返回更改。
func (m *kubeGenericRuntimeManager) computePodActions(pod *v1.Pod, podStatus *kubecontainer.PodStatus) podActions {
	klog.V(5).InfoS("Syncing Pod", "pod", klog.KObj(pod))

	createPodSandbox, attempt, sandboxID := runtimeutil.PodSandboxChanged(pod, podStatus)
	changes := podActions{
		KillPod:           createPodSandbox,
		CreateSandbox:     createPodSandbox,
		SandboxID:         sandboxID,
		Attempt:           attempt,
		ContainersToStart: []int{},
		ContainersToKill:  make(map[kubecontainer.ContainerID]containerToKillInfo),
	}

	// If we need to (re-)create the pod sandbox, everything will need to be
	// killed and recreated, and init containers should be purged.
	if createPodSandbox {
		if !shouldRestartOnFailure(pod) && attempt != 0 && len(podStatus.ContainerStatuses) != 0 {
			// Should not restart the pod, just return.
			// we should not create a sandbox for a pod if it is already done.
			// if all containers are done and should not be started, there is no need to create a new sandbox.
			// this stops confusing logs on pods whose containers all have exit codes, but we recreate a sandbox before terminating it.
			//
			// If ContainerStatuses is empty, we assume that we've never
			// successfully created any containers. In this case, we should
			// retry creating the sandbox.
			changes.CreateSandbox = false
			return changes
		}

		// Get the containers to start, excluding the ones that succeeded if RestartPolicy is OnFailure.
		var containersToStart []int
		for idx, c := range pod.Spec.Containers {
			if pod.Spec.RestartPolicy == v1.RestartPolicyOnFailure && containerSucceeded(&c, podStatus) {
				continue
			}
			containersToStart = append(containersToStart, idx)
		}
		// We should not create a sandbox for a Pod if initialization is done and there is no container to start.
		if len(containersToStart) == 0 {
			_, _, done := findNextInitContainerToRun(pod, podStatus)
			if done {
				changes.CreateSandbox = false
				return changes
			}
		}

		if len(pod.Spec.InitContainers) != 0 {
			// Pod has init containers, return the first one.
			changes.NextInitContainerToStart = &pod.Spec.InitContainers[0]
			return changes
		}
		changes.ContainersToStart = containersToStart
		return changes
	}

	// Ephemeral containers may be started even if initialization is not yet complete.
	for i := range pod.Spec.EphemeralContainers {
		c := (*v1.Container)(&pod.Spec.EphemeralContainers[i].EphemeralContainerCommon)

		// Ephemeral Containers are never restarted
		if podStatus.FindContainerStatusByName(c.Name) == nil {
			changes.EphemeralContainersToStart = append(changes.EphemeralContainersToStart, i)
		}
	}

	// Check initialization progress.
	initLastStatus, next, done := findNextInitContainerToRun(pod, podStatus)
	if !done {
		if next != nil {
			initFailed := initLastStatus != nil && isInitContainerFailed(initLastStatus)
			if initFailed && !shouldRestartOnFailure(pod) {
				changes.KillPod = true
			} else {
				// Always try to stop containers in unknown state first.
				if initLastStatus != nil && initLastStatus.State == kubecontainer.ContainerStateUnknown {
					changes.ContainersToKill[initLastStatus.ID] = containerToKillInfo{
						name:      next.Name,
						container: next,
						message: fmt.Sprintf("Init container is in %q state, try killing it before restart",
							initLastStatus.State),
						reason: reasonUnknown,
					}
				}
				changes.NextInitContainerToStart = next
			}
		}
		// Initialization failed or still in progress. Skip inspecting non-init
		// containers.
		return changes
	}

	// Number of running containers to keep.
	keepCount := 0
	// check the status of containers.
	for idx, container := range pod.Spec.Containers {
		containerStatus := podStatus.FindContainerStatusByName(container.Name)

		// Call internal container post-stop lifecycle hook for any non-running container so that any
		// allocated cpus are released immediately. If the container is restarted, cpus will be re-allocated
		// to it.
		if containerStatus != nil && containerStatus.State != kubecontainer.ContainerStateRunning {
			if err := m.internalLifecycle.PostStopContainer(containerStatus.ID.ID); err != nil {
				klog.ErrorS(err, "Internal container post-stop lifecycle hook failed for container in pod with error",
					"containerName", container.Name, "pod", klog.KObj(pod))
			}
		}

		// If container does not exist, or is not running, check whether we
		// need to restart it.
		if containerStatus == nil || containerStatus.State != kubecontainer.ContainerStateRunning {
			if kubecontainer.ShouldContainerBeRestarted(&container, pod, podStatus) {
				klog.V(3).InfoS("Container of pod is not in the desired state and shall be started", "containerName", container.Name, "pod", klog.KObj(pod))
				changes.ContainersToStart = append(changes.ContainersToStart, idx)
				if containerStatus != nil && containerStatus.State == kubecontainer.ContainerStateUnknown {
					// If container is in unknown state, we don't know whether it
					// is actually running or not, always try killing it before
					// restart to avoid having 2 running instances of the same container.
					changes.ContainersToKill[containerStatus.ID] = containerToKillInfo{
						name:      containerStatus.Name,
						container: &pod.Spec.Containers[idx],
						message: fmt.Sprintf("Container is in %q state, try killing it before restart",
							containerStatus.State),
						reason: reasonUnknown,
					}
				}
			}
			continue
		}
		// The container is running, but kill the container if any of the following condition is met.
		var message string
		var reason containerKillReason
		restart := shouldRestartOnFailure(pod)
		if _, _, changed := containerChanged(&container, containerStatus); changed {
			message = fmt.Sprintf("Container %s definition changed", container.Name)
			// Restart regardless of the restart policy because the container
			// spec changed.
			restart = true
		} else if liveness, found := m.livenessManager.Get(containerStatus.ID); found && liveness == proberesults.Failure {
			// If the container failed the liveness probe, we should kill it.
			message = fmt.Sprintf("Container %s failed liveness probe", container.Name)
			reason = reasonLivenessProbe
		} else if startup, found := m.startupManager.Get(containerStatus.ID); found && startup == proberesults.Failure {
			// If the container failed the startup probe, we should kill it.
			message = fmt.Sprintf("Container %s failed startup probe", container.Name)
			reason = reasonStartupProbe
		} else {
			// Keep the container.
			keepCount++
			continue
		}

		// We need to kill the container, but if we also want to restart the
		// container afterwards, make the intent clear in the message. Also do
		// not kill the entire pod since we expect container to be running eventually.
		if restart {
			message = fmt.Sprintf("%s, will be restarted", message)
			changes.ContainersToStart = append(changes.ContainersToStart, idx)
		}

		changes.ContainersToKill[containerStatus.ID] = containerToKillInfo{
			name:      containerStatus.Name,
			container: &pod.Spec.Containers[idx],
			message:   message,
			reason:    reason,
		}
		klog.V(2).InfoS("Message for Container of pod", "containerName", container.Name, "containerStatusID", containerStatus.ID, "pod", klog.KObj(pod), "containerMessage", message)
	}

	if keepCount == 0 && len(changes.ContainersToStart) == 0 {
		changes.KillPod = true
	}

	return changes
}

// SyncPod syncs the running pod into the desired pod by executing following steps:
//
//  1. Compute sandbox and container changes.
//  2. Kill pod sandbox if necessary.
//  3. Kill any containers that should not be running.
//  4. Create sandbox if necessary.
//  5. Create ephemeral containers.
//  6. Create init containers.
//  7. Create normal containers.
//
// SyncPod同步运行中的pod到期望的pod，通过执行以下步骤：
// 1. 计算沙箱和容器的变化。
// 2. 如果需要，杀死pod沙箱。
// 3. 杀死任何不应该运行的容器。
// 4. 如果需要，创建沙箱。
// 5. 创建临时容器。
// 6. 创建初始化容器。
// 7. 创建正常容器。
func (m *kubeGenericRuntimeManager) SyncPod(ctx context.Context, pod *v1.Pod, podStatus *kubecontainer.PodStatus, pullSecrets []v1.Secret, backOff *flowcontrol.Backoff) (result kubecontainer.PodSyncResult) {
	// Step 1: Compute sandbox and container changes.
	podContainerChanges := m.computePodActions(pod, podStatus)
	klog.V(3).InfoS("computePodActions got for pod", "podActions", podContainerChanges, "pod", klog.KObj(pod))
	if podContainerChanges.CreateSandbox {
		ref, err := ref.GetReference(legacyscheme.Scheme, pod)
		if err != nil {
			klog.ErrorS(err, "Couldn't make a ref to pod", "pod", klog.KObj(pod))
		}
		if podContainerChanges.SandboxID != "" {
			m.recorder.Eventf(ref, v1.EventTypeNormal, events.SandboxChanged, "Pod sandbox changed, it will be killed and re-created.")
		} else {
			klog.V(4).InfoS("SyncPod received new pod, will create a sandbox for it", "pod", klog.KObj(pod))
		}
	}

	// Step 2: Kill the pod if the sandbox has changed.
	if podContainerChanges.KillPod {
		if podContainerChanges.CreateSandbox {
			klog.V(4).InfoS("Stopping PodSandbox for pod, will start new one", "pod", klog.KObj(pod))
		} else {
			klog.V(4).InfoS("Stopping PodSandbox for pod, because all other containers are dead", "pod", klog.KObj(pod))
		}

		killResult := m.killPodWithSyncResult(ctx, pod, kubecontainer.ConvertPodStatusToRunningPod(m.runtimeName, podStatus), nil)
		result.AddPodSyncResult(killResult)
		if killResult.Error() != nil {
			klog.ErrorS(killResult.Error(), "killPodWithSyncResult failed")
			return
		}

		if podContainerChanges.CreateSandbox {
			m.purgeInitContainers(ctx, pod, podStatus)
		}
	} else {
		// Step 3: kill any running containers in this pod which are not to keep.
		for containerID, containerInfo := range podContainerChanges.ContainersToKill {
			klog.V(3).InfoS("Killing unwanted container for pod", "containerName", containerInfo.name, "containerID", containerID, "pod", klog.KObj(pod))
			killContainerResult := kubecontainer.NewSyncResult(kubecontainer.KillContainer, containerInfo.name)
			result.AddSyncResult(killContainerResult)
			if err := m.killContainer(ctx, pod, containerID, containerInfo.name, containerInfo.message, containerInfo.reason, nil); err != nil {
				killContainerResult.Fail(kubecontainer.ErrKillContainer, err.Error())
				klog.ErrorS(err, "killContainer for pod failed", "containerName", containerInfo.name, "containerID", containerID, "pod", klog.KObj(pod))
				return
			}
		}
	}

	// Keep terminated init containers fairly aggressively controlled
	// This is an optimization because container removals are typically handled
	// by container garbage collector.
	m.pruneInitContainersBeforeStart(ctx, pod, podStatus)

	// We pass the value of the PRIMARY podIP and list of podIPs down to
	// generatePodSandboxConfig and generateContainerConfig, which in turn
	// passes it to various other functions, in order to facilitate functionality
	// that requires this value (hosts file and downward API) and avoid races determining
	// the pod IP in cases where a container requires restart but the
	// podIP isn't in the status manager yet. The list of podIPs is used to
	// generate the hosts file.
	//
	// We default to the IPs in the passed-in pod status, and overwrite them if the
	// sandbox needs to be (re)started.
	var podIPs []string
	if podStatus != nil {
		podIPs = podStatus.IPs
	}

	// Step 4: Create a sandbox for the pod if necessary.
	podSandboxID := podContainerChanges.SandboxID
	if podContainerChanges.CreateSandbox {
		var msg string
		var err error

		klog.V(4).InfoS("Creating PodSandbox for pod", "pod", klog.KObj(pod))
		metrics.StartedPodsTotal.Inc()
		createSandboxResult := kubecontainer.NewSyncResult(kubecontainer.CreatePodSandbox, format.Pod(pod))
		result.AddSyncResult(createSandboxResult)

		// ConvertPodSysctlsVariableToDotsSeparator converts sysctl variable
		// in the Pod.Spec.SecurityContext.Sysctls slice into a dot as a separator.
		// runc uses the dot as the separator to verify whether the sysctl variable
		// is correct in a separate namespace, so when using the slash as the sysctl
		// variable separator, runc returns an error: "sysctl is not in a separate kernel namespace"
		// and the podSandBox cannot be successfully created. Therefore, before calling runc,
		// we need to convert the sysctl variable, the dot is used as a separator to separate the kernel namespace.
		// When runc supports slash as sysctl separator, this function can no longer be used.
		sysctl.ConvertPodSysctlsVariableToDotsSeparator(pod.Spec.SecurityContext)

		podSandboxID, msg, err = m.createPodSandbox(ctx, pod, podContainerChanges.Attempt)
		if err != nil {
			// createPodSandbox can return an error from CNI, CSI,
			// or CRI if the Pod has been deleted while the POD is
			// being created. If the pod has been deleted then it's
			// not a real error.
			//
			// SyncPod can still be running when we get here, which
			// means the PodWorker has not acked the deletion.
			if m.podStateProvider.IsPodTerminationRequested(pod.UID) {
				klog.V(4).InfoS("Pod was deleted and sandbox failed to be created", "pod", klog.KObj(pod), "podUID", pod.UID)
				return
			}
			metrics.StartedPodsErrorsTotal.Inc()
			createSandboxResult.Fail(kubecontainer.ErrCreatePodSandbox, msg)
			klog.ErrorS(err, "CreatePodSandbox for pod failed", "pod", klog.KObj(pod))
			ref, referr := ref.GetReference(legacyscheme.Scheme, pod)
			if referr != nil {
				klog.ErrorS(referr, "Couldn't make a ref to pod", "pod", klog.KObj(pod))
			}
			m.recorder.Eventf(ref, v1.EventTypeWarning, events.FailedCreatePodSandBox, "Failed to create pod sandbox: %v", err)
			return
		}
		klog.V(4).InfoS("Created PodSandbox for pod", "podSandboxID", podSandboxID, "pod", klog.KObj(pod))

		resp, err := m.runtimeService.PodSandboxStatus(ctx, podSandboxID, false)
		if err != nil {
			ref, referr := ref.GetReference(legacyscheme.Scheme, pod)
			if referr != nil {
				klog.ErrorS(referr, "Couldn't make a ref to pod", "pod", klog.KObj(pod))
			}
			m.recorder.Eventf(ref, v1.EventTypeWarning, events.FailedStatusPodSandBox, "Unable to get pod sandbox status: %v", err)
			klog.ErrorS(err, "Failed to get pod sandbox status; Skipping pod", "pod", klog.KObj(pod))
			result.Fail(err)
			return
		}
		if resp.GetStatus() == nil {
			result.Fail(errors.New("pod sandbox status is nil"))
			return
		}

		// If we ever allow updating a pod from non-host-network to
		// host-network, we may use a stale IP.
		if !kubecontainer.IsHostNetworkPod(pod) {
			// Overwrite the podIPs passed in the pod status, since we just started the pod sandbox.
			podIPs = m.determinePodSandboxIPs(pod.Namespace, pod.Name, resp.GetStatus())
			klog.V(4).InfoS("Determined the ip for pod after sandbox changed", "IPs", podIPs, "pod", klog.KObj(pod))
		}
	}

	// the start containers routines depend on pod ip(as in primary pod ip)
	// instead of trying to figure out if we have 0 < len(podIPs)
	// everytime, we short circuit it here
	podIP := ""
	if len(podIPs) != 0 {
		podIP = podIPs[0]
	}

	// Get podSandboxConfig for containers to start.
	configPodSandboxResult := kubecontainer.NewSyncResult(kubecontainer.ConfigPodSandbox, podSandboxID)
	result.AddSyncResult(configPodSandboxResult)
	podSandboxConfig, err := m.generatePodSandboxConfig(pod, podContainerChanges.Attempt)
	if err != nil {
		message := fmt.Sprintf("GeneratePodSandboxConfig for pod %q failed: %v", format.Pod(pod), err)
		klog.ErrorS(err, "GeneratePodSandboxConfig for pod failed", "pod", klog.KObj(pod))
		configPodSandboxResult.Fail(kubecontainer.ErrConfigPodSandbox, message)
		return
	}

	// Helper containing boilerplate common to starting all types of containers.
	// typeName is a description used to describe this type of container in log messages,
	// currently: "container", "init container" or "ephemeral container"
	// metricLabel is the label used to describe this type of container in monitoring metrics.
	// currently: "container", "init_container" or "ephemeral_container"
	start := func(ctx context.Context, typeName, metricLabel string, spec *startSpec) error {
		startContainerResult := kubecontainer.NewSyncResult(kubecontainer.StartContainer, spec.container.Name)
		result.AddSyncResult(startContainerResult)

		isInBackOff, msg, err := m.doBackOff(pod, spec.container, podStatus, backOff)
		if isInBackOff {
			startContainerResult.Fail(err, msg)
			klog.V(4).InfoS("Backing Off restarting container in pod", "containerType", typeName, "container", spec.container, "pod", klog.KObj(pod))
			return err
		}

		metrics.StartedContainersTotal.WithLabelValues(metricLabel).Inc()
		if sc.HasWindowsHostProcessRequest(pod, spec.container) {
			metrics.StartedHostProcessContainersTotal.WithLabelValues(metricLabel).Inc()
		}
		klog.V(4).InfoS("Creating container in pod", "containerType", typeName, "container", spec.container, "pod", klog.KObj(pod))
		// NOTE (aramase) podIPs are populated for single stack and dual stack clusters. Send only podIPs.
		if msg, err := m.startContainer(ctx, podSandboxID, podSandboxConfig, spec, pod, podStatus, pullSecrets, podIP, podIPs); err != nil {
			// startContainer() returns well-defined error codes that have reasonable cardinality for metrics and are
			// useful to cluster administrators to distinguish "server errors" from "user errors".
			metrics.StartedContainersErrorsTotal.WithLabelValues(metricLabel, err.Error()).Inc()
			if sc.HasWindowsHostProcessRequest(pod, spec.container) {
				metrics.StartedHostProcessContainersErrorsTotal.WithLabelValues(metricLabel, err.Error()).Inc()
			}
			startContainerResult.Fail(err, msg)
			// known errors that are logged in other places are logged at higher levels here to avoid
			// repetitive log spam
			switch {
			case err == images.ErrImagePullBackOff:
				klog.V(3).InfoS("Container start failed in pod", "containerType", typeName, "container", spec.container, "pod", klog.KObj(pod), "containerMessage", msg, "err", err)
			default:
				utilruntime.HandleError(fmt.Errorf("%v %+v start failed in pod %v: %v: %s", typeName, spec.container, format.Pod(pod), err, msg))
			}
			return err
		}

		return nil
	}

	// Step 5: start ephemeral containers
	// These are started "prior" to init containers to allow running ephemeral containers even when there
	// are errors starting an init container. In practice init containers will start first since ephemeral
	// containers cannot be specified on pod creation.
	for _, idx := range podContainerChanges.EphemeralContainersToStart {
		start(ctx, "ephemeral container", metrics.EphemeralContainer, ephemeralContainerStartSpec(&pod.Spec.EphemeralContainers[idx]))
	}

	// Step 6: start the init container.
	if container := podContainerChanges.NextInitContainerToStart; container != nil {
		// Start the next init container.
		if err := start(ctx, "init container", metrics.InitContainer, containerStartSpec(container)); err != nil {
			return
		}

		// Successfully started the container; clear the entry in the failure
		klog.V(4).InfoS("Completed init container for pod", "containerName", container.Name, "pod", klog.KObj(pod))
	}

	// Step 7: start containers in podContainerChanges.ContainersToStart.
	for _, idx := range podContainerChanges.ContainersToStart {
		start(ctx, "container", metrics.Container, containerStartSpec(&pod.Spec.Containers[idx]))
	}

	return
}

// If a container is still in backoff, the function will return a brief backoff error and
// a detailed error message.
// 如果容器仍在回退，则该函数将返回一个简短的回退错误和详细的错误消息。
func (m *kubeGenericRuntimeManager) doBackOff(pod *v1.Pod, container *v1.Container, podStatus *kubecontainer.PodStatus, backOff *flowcontrol.Backoff) (bool, string, error) {
	var cStatus *kubecontainer.Status
	for _, c := range podStatus.ContainerStatuses {
		if c.Name == container.Name && c.State == kubecontainer.ContainerStateExited {
			cStatus = c
			break
		}
	}

	if cStatus == nil {
		return false, "", nil
	}

	klog.V(3).InfoS("Checking backoff for container in pod", "containerName", container.Name, "pod", klog.KObj(pod))
	// Use the finished time of the latest exited container as the start point to calculate whether to do back-off.
	ts := cStatus.FinishedAt
	// backOff requires a unique key to identify the container.
	key := getStableKey(pod, container)
	if backOff.IsInBackOffSince(key, ts) {
		if containerRef, err := kubecontainer.GenerateContainerRef(pod, container); err == nil {
			m.recorder.Eventf(containerRef, v1.EventTypeWarning, events.BackOffStartContainer,
				fmt.Sprintf("Back-off restarting failed container %s in pod %s", container.Name, format.Pod(pod)))
		}
		err := fmt.Errorf("back-off %s restarting failed container=%s pod=%s", backOff.Get(key), container.Name, format.Pod(pod))
		klog.V(3).InfoS("Back-off restarting failed container", "err", err.Error())
		return true, err.Error(), kubecontainer.ErrCrashLoopBackOff
	}

	backOff.Next(key, ts)
	return false, "", nil
}

// KillPod kills all the containers of a pod. Pod may be nil, running pod must not be.
// gracePeriodOverride if specified allows the caller to override the pod default grace period.
// only hard kill paths are allowed to specify a gracePeriodOverride in the kubelet in order to not corrupt user data.
// it is useful when doing SIGKILL for hard eviction scenarios, or max grace period during soft eviction scenarios.
// KillPod 杀死一个Pod的所有容器。Pod可能为零，运行Pod一定不是。
// 如果指定了gracePeriodOverride，则允许调用者覆盖Pod默认的宽限期。
// 只有硬杀路径允许在kubelet中指定gracePeriodOverride，以避免破坏用户数据。
// 在硬逐出场景中进行SIGKILL时很有用，或者在软逐出场景中进行最大宽限期。
func (m *kubeGenericRuntimeManager) KillPod(ctx context.Context, pod *v1.Pod, runningPod kubecontainer.Pod, gracePeriodOverride *int64) error {
	err := m.killPodWithSyncResult(ctx, pod, runningPod, gracePeriodOverride)
	return err.Error()
}

// killPodWithSyncResult kills a runningPod and returns SyncResult.
// Note: The pod passed in could be *nil* when kubelet restarted.
// killPodWithSyncResult 杀死一个运行的Pod并返回SyncResult。 注意：当kubelet重启时，传入的Pod可能为*nil*。
func (m *kubeGenericRuntimeManager) killPodWithSyncResult(ctx context.Context, pod *v1.Pod, runningPod kubecontainer.Pod, gracePeriodOverride *int64) (result kubecontainer.PodSyncResult) {
	killContainerResults := m.killContainersWithSyncResult(ctx, pod, runningPod, gracePeriodOverride)
	for _, containerResult := range killContainerResults {
		result.AddSyncResult(containerResult)
	}

	// stop sandbox, the sandbox will be removed in GarbageCollect
	killSandboxResult := kubecontainer.NewSyncResult(kubecontainer.KillPodSandbox, runningPod.ID)
	result.AddSyncResult(killSandboxResult)
	// Stop all sandboxes belongs to same pod
	for _, podSandbox := range runningPod.Sandboxes {
		if err := m.runtimeService.StopPodSandbox(ctx, podSandbox.ID.ID); err != nil && !crierror.IsNotFound(err) {
			killSandboxResult.Fail(kubecontainer.ErrKillPodSandbox, err.Error())
			klog.ErrorS(nil, "Failed to stop sandbox", "podSandboxID", podSandbox.ID)
		}
	}

	return
}

func (m *kubeGenericRuntimeManager) GeneratePodStatus(event *runtimeapi.ContainerEventResponse) (*kubecontainer.PodStatus, error) {
	podIPs := m.determinePodSandboxIPs(event.PodSandboxStatus.Metadata.Namespace, event.PodSandboxStatus.Metadata.Name, event.PodSandboxStatus)

	kubeContainerStatuses := []*kubecontainer.Status{}
	for _, status := range event.ContainersStatuses {
		kubeContainerStatuses = append(kubeContainerStatuses, m.convertToKubeContainerStatus(status))
	}

	sort.Sort(containerStatusByCreated(kubeContainerStatuses))

	return &kubecontainer.PodStatus{
		ID:                kubetypes.UID(event.PodSandboxStatus.Metadata.Uid),
		Name:              event.PodSandboxStatus.Metadata.Name,
		Namespace:         event.PodSandboxStatus.Metadata.Namespace,
		IPs:               podIPs,
		SandboxStatuses:   []*runtimeapi.PodSandboxStatus{event.PodSandboxStatus},
		ContainerStatuses: kubeContainerStatuses,
	}, nil
}

// GetPodStatus retrieves the status of the pod, including the
// information of all containers in the pod that are visible in Runtime.
// GetPodStatus检索Pod的状态，包括Pod中所有可见的容器的信息。
func (m *kubeGenericRuntimeManager) GetPodStatus(ctx context.Context, uid kubetypes.UID, name, namespace string) (*kubecontainer.PodStatus, error) {
	// Now we retain restart count of container as a container label. Each time a container
	// restarts, pod will read the restart count from the registered dead container, increment
	// it to get the new restart count, and then add a label with the new restart count on
	// the newly started container.
	// However, there are some limitations of this method:
	//	1. When all dead containers were garbage collected, the container status could
	//	not get the historical value and would be *inaccurate*. Fortunately, the chance
	//	is really slim.
	//	2. When working with old version containers which have no restart count label,
	//	we can only assume their restart count is 0.
	// Anyhow, we only promised "best-effort" restart count reporting, we can just ignore
	// these limitations now.
	// TODO: move this comment to SyncPod.
	podSandboxIDs, err := m.getSandboxIDByPodUID(ctx, uid, nil)
	if err != nil {
		return nil, err
	}

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
			UID:       uid,
		},
	}

	podFullName := format.Pod(pod)

	klog.V(4).InfoS("getSandboxIDByPodUID got sandbox IDs for pod", "podSandboxID", podSandboxIDs, "pod", klog.KObj(pod))

	sandboxStatuses := []*runtimeapi.PodSandboxStatus{}
	containerStatuses := []*kubecontainer.Status{}
	var timestamp time.Time

	podIPs := []string{}
	for idx, podSandboxID := range podSandboxIDs {
		resp, err := m.runtimeService.PodSandboxStatus(ctx, podSandboxID, false)
		// Between List (getSandboxIDByPodUID) and check (PodSandboxStatus) another thread might remove a container, and that is normal.
		// The previous call (getSandboxIDByPodUID) never fails due to a pod sandbox not existing.
		// Therefore, this method should not either, but instead act as if the previous call failed,
		// which means the error should be ignored.
		if crierror.IsNotFound(err) {
			continue
		}
		if err != nil {
			klog.ErrorS(err, "PodSandboxStatus of sandbox for pod", "podSandboxID", podSandboxID, "pod", klog.KObj(pod))
			return nil, err
		}
		if resp.GetStatus() == nil {
			return nil, errors.New("pod sandbox status is nil")

		}
		sandboxStatuses = append(sandboxStatuses, resp.Status)
		// Only get pod IP from latest sandbox
		if idx == 0 && resp.Status.State == runtimeapi.PodSandboxState_SANDBOX_READY {
			podIPs = m.determinePodSandboxIPs(namespace, name, resp.Status)
		}

		if idx == 0 && utilfeature.DefaultFeatureGate.Enabled(features.EventedPLEG) {
			if resp.Timestamp == 0 {
				// If the Evented PLEG is enabled in the kubelet, but not in the runtime
				// then the pod status we get will not have the timestamp set.
				// e.g. CI job 'pull-kubernetes-e2e-gce-alpha-features' will runs with
				// features gate enabled, which includes Evented PLEG, but uses the
				// runtime without Evented PLEG support.
				klog.V(4).InfoS("Runtime does not set pod status timestamp", "pod", klog.KObj(pod))
				containerStatuses, err = m.getPodContainerStatuses(ctx, uid, name, namespace)
				if err != nil {
					if m.logReduction.ShouldMessageBePrinted(err.Error(), podFullName) {
						klog.ErrorS(err, "getPodContainerStatuses for pod failed", "pod", klog.KObj(pod))
					}
					return nil, err
				}
			} else {
				// Get the statuses of all containers visible to the pod and
				// timestamp from sandboxStatus.
				timestamp = time.Unix(resp.Timestamp, 0)
				for _, cs := range resp.ContainersStatuses {
					cStatus := m.convertToKubeContainerStatus(cs)
					containerStatuses = append(containerStatuses, cStatus)
				}
			}
		}
	}

	if !utilfeature.DefaultFeatureGate.Enabled(features.EventedPLEG) {
		// Get statuses of all containers visible in the pod.
		containerStatuses, err = m.getPodContainerStatuses(ctx, uid, name, namespace)
		if err != nil {
			if m.logReduction.ShouldMessageBePrinted(err.Error(), podFullName) {
				klog.ErrorS(err, "getPodContainerStatuses for pod failed", "pod", klog.KObj(pod))
			}
			return nil, err
		}
	}

	m.logReduction.ClearID(podFullName)

	return &kubecontainer.PodStatus{
		ID:                uid,
		Name:              name,
		Namespace:         namespace,
		IPs:               podIPs,
		SandboxStatuses:   sandboxStatuses,
		ContainerStatuses: containerStatuses,
		TimeStamp:         timestamp,
	}, nil
}

// GarbageCollect removes dead containers using the specified container gc policy.
// GarbageCollect 删除死亡容器使用指定的容器 gc 策略。
func (m *kubeGenericRuntimeManager) GarbageCollect(ctx context.Context, gcPolicy kubecontainer.GCPolicy, allSourcesReady bool, evictNonDeletedPods bool) error {
	return m.containerGC.GarbageCollect(ctx, gcPolicy, allSourcesReady, evictNonDeletedPods)
}

// UpdatePodCIDR is just a passthrough method to update the runtimeConfig of the shim
// with the podCIDR supplied by the kubelet.
// UpdatePodCIDR 只是一个传递方法，用于更新带有 kubelet 提供的 podCIDR 的 shim 的 runtimeConfig。
func (m *kubeGenericRuntimeManager) UpdatePodCIDR(ctx context.Context, podCIDR string) error {
	// TODO(#35531): do we really want to write a method on this manager for each
	// field of the config?
	klog.InfoS("Updating runtime config through cri with podcidr", "CIDR", podCIDR)
	return m.runtimeService.UpdateRuntimeConfig(ctx,
		&runtimeapi.RuntimeConfig{
			NetworkConfig: &runtimeapi.NetworkConfig{
				PodCidr: podCIDR,
			},
		})
}

func (m *kubeGenericRuntimeManager) CheckpointContainer(ctx context.Context, options *runtimeapi.CheckpointContainerRequest) error {
	return m.runtimeService.CheckpointContainer(ctx, options)
}

func (m *kubeGenericRuntimeManager) ListMetricDescriptors(ctx context.Context) ([]*runtimeapi.MetricDescriptor, error) {
	return m.runtimeService.ListMetricDescriptors(ctx)
}

func (m *kubeGenericRuntimeManager) ListPodSandboxMetrics(ctx context.Context) ([]*runtimeapi.PodSandboxMetrics, error) {
	return m.runtimeService.ListPodSandboxMetrics(ctx)
}
