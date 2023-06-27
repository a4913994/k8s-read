/*
Copyright 2015 The Kubernetes Authors.

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

//go:generate mockgen -destination=testing/mock_manager.go -package=testing -build_flags=-mod=mod . Manager
package pod

import (
	"sync"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/kubelet/metrics"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
)

// Manager stores and manages access to pods, maintaining the mappings
// between static pods and mirror pods.
//
// Manager 存储并管理对 pod 的访问，维护静态 pod 和镜像 pod 之间的映射。
//
// The kubelet discovers pod updates from 3 sources: file, http, and
// apiserver. Pods from non-apiserver sources are called static pods, and API
// server is not aware of the existence of static pods. In order to monitor
// the status of such pods, the kubelet creates a mirror pod for each static
// pod via the API server.
//
// kubelet 从 3 个源发现 pod 更新： 文件、http 和 apiserver。
// 来自非 apiserver 源的 pod 称为静态 pod，而 API 服务器不知道静态 pod 的存在。
// 为了监控这些 pod 的状态，kubelet 通过 API 服务器为每个静态 pod 创建一个镜像 pod。
//
// A mirror pod has the same pod full name (name and namespace) as its static
// counterpart (albeit different metadata such as UID, etc). By leveraging the
// fact that the kubelet reports the pod status using the pod full name, the
// status of the mirror pod always reflects the actual status of the static
// pod. When a static pod gets deleted, the associated orphaned mirror pod
// will also be removed.
// 一个镜像 pod 的 pod 全名（名称和命名空间）与其静态副本相同（尽管元数据不同，例如 UID 等）。
// 通过利用 kubelet 使用 pod 全名报告 pod 状态的事实，镜像 pod 的状态始终反映静态 pod 的实际状态。
// 当删除静态 pod 时，相关的孤立镜像 pod 也将被删除。
type Manager interface {
	// GetPods returns the regular pods bound to the kubelet and their spec.
	// GetPods 返回绑定到 kubelet 的常规 pod 和它们的 spec。
	GetPods() []*v1.Pod
	// GetPodByFullName returns the (non-mirror) pod that matches full name, as well as
	// whether the pod was found.
	// GetPodByFullName 返回与全名匹配的（非镜像）pod，以及是否找到了 pod。
	GetPodByFullName(podFullName string) (*v1.Pod, bool)
	// GetPodByName provides the (non-mirror) pod that matches namespace and
	// name, as well as whether the pod was found.
	// GetPodByName 提供与命名空间和名称匹配的（非镜像）pod，以及是否找到了 pod。
	GetPodByName(namespace, name string) (*v1.Pod, bool)
	// GetPodByUID provides the (non-mirror) pod that matches pod UID, as well as
	// whether the pod is found.
	// GetPodByUID 提供与 pod UID 匹配的（非镜像）pod，以及是否找到了 pod。
	GetPodByUID(types.UID) (*v1.Pod, bool)
	// GetPodByMirrorPod returns the static pod for the given mirror pod and
	// whether it was known to the pod manager.
	// GetPodByMirrorPod 返回给定镜像 pod 的静态 pod，以及 pod 管理器是否知道它。
	GetPodByMirrorPod(*v1.Pod) (*v1.Pod, bool)
	// GetMirrorPodByPod returns the mirror pod for the given static pod and
	// whether it was known to the pod manager.
	// GetMirrorPodByPod 返回给定静态 pod 的镜像 pod，以及 pod 管理器是否知道它。
	GetMirrorPodByPod(*v1.Pod) (*v1.Pod, bool)
	// GetPodsAndMirrorPods returns the both regular and mirror pods.
	// GetPodsAndMirrorPods 返回常规和镜像 pod。
	GetPodsAndMirrorPods() ([]*v1.Pod, []*v1.Pod)
	// SetPods replaces the internal pods with the new pods.
	// It is currently only used for testing.
	// SetPods 用新的 pod 替换内部 pod。 目前仅用于测试。
	SetPods(pods []*v1.Pod)
	// AddPod adds the given pod to the manager.
	// AddPod 将给定的 pod 添加到管理器中。
	AddPod(pod *v1.Pod)
	// UpdatePod updates the given pod in the manager.
	// UpdatePod 在管理器中更新给定的 pod。
	UpdatePod(pod *v1.Pod)
	// DeletePod deletes the given pod from the manager.  For mirror pods,
	// this means deleting the mappings related to mirror pods.  For non-
	// mirror pods, this means deleting from indexes for all non-mirror pods.
	// DeletePod 从管理器中删除给定的 pod。 对于镜像 pod，这意味着删除与镜像 pod 相关的映射。
	// 对于非镜像 pod，这意味着从所有非镜像 pod 的索引中删除。
	DeletePod(pod *v1.Pod)
	// GetOrphanedMirrorPodNames returns names of orphaned mirror pods
	// GetOrphanedMirrorPodNames 返回孤立的镜像 pod 的名称
	GetOrphanedMirrorPodNames() []string
	// TranslatePodUID returns the actual UID of a pod. If the UID belongs to
	// a mirror pod, returns the UID of its static pod. Otherwise, returns the
	// original UID.
	//
	// All public-facing functions should perform this translation for UIDs
	// because user may provide a mirror pod UID, which is not recognized by
	// internal Kubelet functions.
	// TranslatePodUID 返回 pod 的实际 UID。 如果 UID 属于镜像 pod，则返回其静态 pod 的 UID。
	// 否则，返回原始 UID。
	//
	// 所有面向公众的函数都应该为 UID 执行此转换，因为用户可能会提供镜像 pod 的 UID，
	// 该 UID 不被内部 Kubelet 函数识别。
	TranslatePodUID(uid types.UID) kubetypes.ResolvedPodUID
	// GetUIDTranslations returns the mappings of static pod UIDs to mirror pod
	// UIDs and mirror pod UIDs to static pod UIDs.
	// GetUIDTranslations 返回静态 pod UID 到镜像 pod UID 的映射和镜像 pod UID 到静态 pod UID 的映射。
	GetUIDTranslations() (podToMirror map[kubetypes.ResolvedPodUID]kubetypes.MirrorPodUID, mirrorToPod map[kubetypes.MirrorPodUID]kubetypes.ResolvedPodUID)
	// IsMirrorPodOf returns true if mirrorPod is a correct representation of
	// pod; false otherwise.
	// IsMirrorPodOf 如果 mirrorPod 是 pod 的正确表示，则返回 true；否则返回 false。
	IsMirrorPodOf(mirrorPod, pod *v1.Pod) bool

	MirrorClient
}

// basicManager is a functional Manager.
//
// All fields in basicManager are read-only and are updated calling SetPods,
// AddPod, UpdatePod, or DeletePod.
// basicManager 是一个功能性的 Manager。
//
// basicManager 中的所有字段都是只读的，并且通过调用 SetPods、AddPod、UpdatePod 或 DeletePod 进行更新。
type basicManager struct {
	// Protects all internal maps.
	// 保护所有内部映射。
	lock sync.RWMutex

	// Regular pods indexed by UID.
	podByUID map[kubetypes.ResolvedPodUID]*v1.Pod
	// Mirror pods indexed by UID.
	mirrorPodByUID map[kubetypes.MirrorPodUID]*v1.Pod

	// Pods indexed by full name for easy access.
	// 通过全名索引的 pod，以便轻松访问。
	podByFullName       map[string]*v1.Pod
	mirrorPodByFullName map[string]*v1.Pod

	// Mirror pod UID to pod UID map.
	// 镜像 pod UID 到 pod UID 的映射。
	translationByUID map[kubetypes.MirrorPodUID]kubetypes.ResolvedPodUID

	// A mirror pod client to create/delete mirror pods.
	// 一个镜像 pod 客户端来创建/删除镜像 pod。
	MirrorClient
}

// NewBasicPodManager returns a functional Manager.
// NewBasicPodManager 返回一个功能性的 Manager。
func NewBasicPodManager(client MirrorClient) Manager {
	pm := &basicManager{}
	pm.MirrorClient = client
	pm.SetPods(nil)
	return pm
}

// Set the internal pods based on the new pods.
// 根据新的 pod 设置内部 pod。
func (pm *basicManager) SetPods(newPods []*v1.Pod) {
	pm.lock.Lock()
	defer pm.lock.Unlock()

	pm.podByUID = make(map[kubetypes.ResolvedPodUID]*v1.Pod)
	pm.podByFullName = make(map[string]*v1.Pod)
	pm.mirrorPodByUID = make(map[kubetypes.MirrorPodUID]*v1.Pod)
	pm.mirrorPodByFullName = make(map[string]*v1.Pod)
	pm.translationByUID = make(map[kubetypes.MirrorPodUID]kubetypes.ResolvedPodUID)

	pm.updatePodsInternal(newPods...)
}

func (pm *basicManager) AddPod(pod *v1.Pod) {
	pm.UpdatePod(pod)
}

func (pm *basicManager) UpdatePod(pod *v1.Pod) {
	pm.lock.Lock()
	defer pm.lock.Unlock()
	pm.updatePodsInternal(pod)
}

// updateMetrics updates the metrics surfaced by the pod manager.
// oldPod or newPod may be nil to signify creation or deletion.
// updateMetrics 更新 pod 管理器公开的指标。 oldPod 或 newPod 可能为 nil，以表示创建或删除。
func updateMetrics(oldPod, newPod *v1.Pod) {
	var numEC int
	if oldPod != nil {
		numEC -= len(oldPod.Spec.EphemeralContainers)
	}
	if newPod != nil {
		numEC += len(newPod.Spec.EphemeralContainers)
	}
	if numEC != 0 {
		metrics.ManagedEphemeralContainers.Add(float64(numEC))
	}
}

// updatePodsInternal replaces the given pods in the current state of the
// manager, updating the various indices. The caller is assumed to hold the
// lock.
// updatePodsInternal 在管理器的当前状态中替换给定的 pod，并更新各种索引。 假定调用者持有锁。
func (pm *basicManager) updatePodsInternal(pods ...*v1.Pod) {
	for _, pod := range pods {
		podFullName := kubecontainer.GetPodFullName(pod)
		// This logic relies on a static pod and its mirror to have the same name.
		// It is safe to type convert here due to the IsMirrorPod guard.
		if kubetypes.IsMirrorPod(pod) {
			mirrorPodUID := kubetypes.MirrorPodUID(pod.UID)
			pm.mirrorPodByUID[mirrorPodUID] = pod
			pm.mirrorPodByFullName[podFullName] = pod
			if p, ok := pm.podByFullName[podFullName]; ok {
				pm.translationByUID[mirrorPodUID] = kubetypes.ResolvedPodUID(p.UID)
			}
		} else {
			resolvedPodUID := kubetypes.ResolvedPodUID(pod.UID)
			updateMetrics(pm.podByUID[resolvedPodUID], pod)
			pm.podByUID[resolvedPodUID] = pod
			pm.podByFullName[podFullName] = pod
			if mirror, ok := pm.mirrorPodByFullName[podFullName]; ok {
				pm.translationByUID[kubetypes.MirrorPodUID(mirror.UID)] = resolvedPodUID
			}
		}
	}
}

func (pm *basicManager) DeletePod(pod *v1.Pod) {
	updateMetrics(pod, nil)
	pm.lock.Lock()
	defer pm.lock.Unlock()
	podFullName := kubecontainer.GetPodFullName(pod)
	// It is safe to type convert here due to the IsMirrorPod guard.
	if kubetypes.IsMirrorPod(pod) {
		mirrorPodUID := kubetypes.MirrorPodUID(pod.UID)
		delete(pm.mirrorPodByUID, mirrorPodUID)
		delete(pm.mirrorPodByFullName, podFullName)
		delete(pm.translationByUID, mirrorPodUID)
	} else {
		delete(pm.podByUID, kubetypes.ResolvedPodUID(pod.UID))
		delete(pm.podByFullName, podFullName)
	}
}

func (pm *basicManager) GetPods() []*v1.Pod {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	return podsMapToPods(pm.podByUID)
}

func (pm *basicManager) GetPodsAndMirrorPods() ([]*v1.Pod, []*v1.Pod) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	pods := podsMapToPods(pm.podByUID)
	mirrorPods := mirrorPodsMapToMirrorPods(pm.mirrorPodByUID)
	return pods, mirrorPods
}

func (pm *basicManager) GetPodByUID(uid types.UID) (*v1.Pod, bool) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	pod, ok := pm.podByUID[kubetypes.ResolvedPodUID(uid)] // Safe conversion, map only holds non-mirrors.
	return pod, ok
}

func (pm *basicManager) GetPodByName(namespace, name string) (*v1.Pod, bool) {
	podFullName := kubecontainer.BuildPodFullName(name, namespace)
	return pm.GetPodByFullName(podFullName)
}

func (pm *basicManager) GetPodByFullName(podFullName string) (*v1.Pod, bool) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	pod, ok := pm.podByFullName[podFullName]
	return pod, ok
}

func (pm *basicManager) TranslatePodUID(uid types.UID) kubetypes.ResolvedPodUID {
	// It is safe to type convert to a resolved UID because type conversion is idempotent.
	if uid == "" {
		return kubetypes.ResolvedPodUID(uid)
	}

	pm.lock.RLock()
	defer pm.lock.RUnlock()
	if translated, ok := pm.translationByUID[kubetypes.MirrorPodUID(uid)]; ok {
		return translated
	}
	return kubetypes.ResolvedPodUID(uid)
}

func (pm *basicManager) GetUIDTranslations() (podToMirror map[kubetypes.ResolvedPodUID]kubetypes.MirrorPodUID,
	mirrorToPod map[kubetypes.MirrorPodUID]kubetypes.ResolvedPodUID) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()

	podToMirror = make(map[kubetypes.ResolvedPodUID]kubetypes.MirrorPodUID, len(pm.translationByUID))
	mirrorToPod = make(map[kubetypes.MirrorPodUID]kubetypes.ResolvedPodUID, len(pm.translationByUID))
	// Insert empty translation mapping for all static pods.
	for uid, pod := range pm.podByUID {
		if !kubetypes.IsStaticPod(pod) {
			continue
		}
		podToMirror[uid] = ""
	}
	// Fill in translations. Notice that if there is no mirror pod for a
	// static pod, its uid will be translated into empty string "". This
	// is WAI, from the caller side we can know that the static pod doesn't
	// have a corresponding mirror pod instead of using static pod uid directly.
	for k, v := range pm.translationByUID {
		mirrorToPod[k] = v
		podToMirror[v] = k
	}
	return podToMirror, mirrorToPod
}

func (pm *basicManager) GetOrphanedMirrorPodNames() []string {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	var podFullNames []string
	for podFullName := range pm.mirrorPodByFullName {
		if _, ok := pm.podByFullName[podFullName]; !ok {
			podFullNames = append(podFullNames, podFullName)
		}
	}
	return podFullNames
}

func (pm *basicManager) IsMirrorPodOf(mirrorPod, pod *v1.Pod) bool {
	// Check name and namespace first.
	if pod.Name != mirrorPod.Name || pod.Namespace != mirrorPod.Namespace {
		return false
	}
	hash, ok := getHashFromMirrorPod(mirrorPod)
	if !ok {
		return false
	}
	return hash == getPodHash(pod)
}

func podsMapToPods(UIDMap map[kubetypes.ResolvedPodUID]*v1.Pod) []*v1.Pod {
	pods := make([]*v1.Pod, 0, len(UIDMap))
	for _, pod := range UIDMap {
		pods = append(pods, pod)
	}
	return pods
}

func mirrorPodsMapToMirrorPods(UIDMap map[kubetypes.MirrorPodUID]*v1.Pod) []*v1.Pod {
	pods := make([]*v1.Pod, 0, len(UIDMap))
	for _, pod := range UIDMap {
		pods = append(pods, pod)
	}
	return pods
}

func (pm *basicManager) GetMirrorPodByPod(pod *v1.Pod) (*v1.Pod, bool) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	mirrorPod, ok := pm.mirrorPodByFullName[kubecontainer.GetPodFullName(pod)]
	return mirrorPod, ok
}

func (pm *basicManager) GetPodByMirrorPod(mirrorPod *v1.Pod) (*v1.Pod, bool) {
	pm.lock.RLock()
	defer pm.lock.RUnlock()
	pod, ok := pm.podByFullName[kubecontainer.GetPodFullName(mirrorPod)]
	return pod, ok
}
