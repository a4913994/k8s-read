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

package pod

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog/v2"
	kubecontainer "k8s.io/kubernetes/pkg/kubelet/container"
	kubetypes "k8s.io/kubernetes/pkg/kubelet/types"
)

// MirrorClient knows how to create/delete a mirror pod in the API server.
// MirrorClient 知道如何在 API 服务器中创建/删除镜像 Pod。
type MirrorClient interface {
	// CreateMirrorPod creates a mirror pod in the API server for the given
	// pod or returns an error.  The mirror pod will have the same annotations
	// as the given pod as well as an extra annotation containing the hash of
	// the static pod.
	// CreateMirrorPod 函数为给定的 Pod 在 API 服务器中创建一个镜像 Pod，或者返回一个错误。
	// 镜像 Pod 将具有与给定 Pod 相同的注解，以及一个包含静态 Pod 的哈希值的额外注解。
	CreateMirrorPod(pod *v1.Pod) error
	// DeleteMirrorPod deletes the mirror pod with the given full name from
	// the API server or returns an error.
	// DeleteMirrorPod 函数从 API 服务器中删除具有给定全名的镜像 Pod，或者返回一个错误。
	DeleteMirrorPod(podFullName string, uid *types.UID) (bool, error)
}

// nodeGetter is a subset of NodeLister, simplified for testing.
// nodeGetter 是 NodeLister 的子集，用于测试。
type nodeGetter interface {
	// Get retrieves the Node for a given name.
	Get(name string) (*v1.Node, error)
}

// basicMirrorClient is a functional MirrorClient.  Mirror pods are stored in
// the kubelet directly because they need to be in sync with the internal
// pods.
// basicMirrorClient 是一个功能完备的 MirrorClient。镜像 Pod 存储在 kubelet 中，因为它们需要与内部 Pod 同步。
type basicMirrorClient struct {
	apiserverClient clientset.Interface
	nodeGetter      nodeGetter
	nodeName        string
}

// NewBasicMirrorClient returns a new MirrorClient.
// NewBasicMirrorClient 函数返回一个新的 MirrorClient。
func NewBasicMirrorClient(apiserverClient clientset.Interface, nodeName string, nodeGetter nodeGetter) MirrorClient {
	return &basicMirrorClient{
		apiserverClient: apiserverClient,
		nodeName:        nodeName,
		nodeGetter:      nodeGetter,
	}
}

func (mc *basicMirrorClient) CreateMirrorPod(pod *v1.Pod) error {
	if mc.apiserverClient == nil {
		return nil
	}
	// Make a copy of the pod.
	copyPod := *pod
	copyPod.Annotations = make(map[string]string)

	for k, v := range pod.Annotations {
		copyPod.Annotations[k] = v
	}
	hash := getPodHash(pod)
	copyPod.Annotations[kubetypes.ConfigMirrorAnnotationKey] = hash

	// With the MirrorPodNodeRestriction feature, mirror pods are required to have an owner reference
	// to the owning node.
	// See https://git.k8s.io/enhancements/keps/sig-auth/1314-node-restriction-pods/README.md
	nodeUID, err := mc.getNodeUID()
	if err != nil {
		return fmt.Errorf("failed to get node UID: %v", err)
	}
	controller := true
	copyPod.OwnerReferences = []metav1.OwnerReference{{
		APIVersion: v1.SchemeGroupVersion.String(),
		Kind:       "Node",
		Name:       mc.nodeName,
		UID:        nodeUID,
		Controller: &controller,
	}}

	apiPod, err := mc.apiserverClient.CoreV1().Pods(copyPod.Namespace).Create(context.TODO(), &copyPod, metav1.CreateOptions{})
	if err != nil && apierrors.IsAlreadyExists(err) {
		// Check if the existing pod is the same as the pod we want to create.
		if h, ok := apiPod.Annotations[kubetypes.ConfigMirrorAnnotationKey]; ok && h == hash {
			return nil
		}
	}
	return err
}

// DeleteMirrorPod deletes a mirror pod.
// It takes the full name of the pod and optionally a UID.  If the UID
// is non-nil, the pod is deleted only if its UID matches the supplied UID.
// It returns whether the pod was actually deleted, and any error returned
// while parsing the name of the pod.
// Non-existence of the pod or UID mismatch is not treated as an error; the
// routine simply returns false in that case.
// DeleteMirrorPod 函数删除一个镜像 Pod。 它需要 Pod 的全名，并且可以选择一个 UID。 如果 UID 不为 nil，则只有当 Pod 的 UID 与提供的 UID 匹配时，
// 才会删除 Pod。 它返回 Pod 是否实际被删除以及在解析 Pod 名称时返回的任何错误。 Pod 不存在或 UID 不匹配不被视为错误；在这种情况下，该例程只返回 false。
func (mc *basicMirrorClient) DeleteMirrorPod(podFullName string, uid *types.UID) (bool, error) {
	if mc.apiserverClient == nil {
		return false, nil
	}
	name, namespace, err := kubecontainer.ParsePodFullName(podFullName)
	if err != nil {
		klog.ErrorS(err, "Failed to parse a pod full name", "podFullName", podFullName)
		return false, err
	}

	var uidValue types.UID
	if uid != nil {
		uidValue = *uid
	}
	klog.V(2).InfoS("Deleting a mirror pod", "pod", klog.KRef(namespace, name), "podUID", uidValue)

	var GracePeriodSeconds int64
	if err := mc.apiserverClient.CoreV1().Pods(namespace).Delete(context.TODO(), name, metav1.DeleteOptions{GracePeriodSeconds: &GracePeriodSeconds, Preconditions: &metav1.Preconditions{UID: uid}}); err != nil {
		// Unfortunately, there's no generic error for failing a precondition
		if !(apierrors.IsNotFound(err) || apierrors.IsConflict(err)) {
			// We should return the error here, but historically this routine does
			// not return an error unless it can't parse the pod name
			klog.ErrorS(err, "Failed deleting a mirror pod", "pod", klog.KRef(namespace, name))
		}
		return false, nil
	}
	return true, nil
}

func (mc *basicMirrorClient) getNodeUID() (types.UID, error) {
	node, err := mc.nodeGetter.Get(mc.nodeName)
	if err != nil {
		return "", err
	}
	if node.UID == "" {
		return "", fmt.Errorf("UID unset for node %s", mc.nodeName)
	}
	return node.UID, nil
}

func getHashFromMirrorPod(pod *v1.Pod) (string, bool) {
	hash, ok := pod.Annotations[kubetypes.ConfigMirrorAnnotationKey]
	return hash, ok
}

func getPodHash(pod *v1.Pod) string {
	// The annotation exists for all static pods.
	return pod.Annotations[kubetypes.ConfigHashAnnotationKey]
}
