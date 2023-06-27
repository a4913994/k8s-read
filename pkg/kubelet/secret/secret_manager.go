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

package secret

import (
	"context"
	"fmt"
	"time"

	v1 "k8s.io/api/core/v1"
	clientset "k8s.io/client-go/kubernetes"
	podutil "k8s.io/kubernetes/pkg/api/v1/pod"
	corev1 "k8s.io/kubernetes/pkg/apis/core/v1"
	"k8s.io/kubernetes/pkg/kubelet/util/manager"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/utils/clock"
)

// Manager manages Kubernetes secrets. This includes retrieving
// secrets or registering/unregistering them via Pods.
// Manager 管理 Kubernetes 的 Secret。这包括获取 Secret 或通过 Pod 注册/注销它们。
type Manager interface {
	// GetSecret 函数通过 Secret 的命名空间和名称获取 Secret。
	// Get secret by secret namespace and name.
	GetSecret(namespace, name string) (*v1.Secret, error)

	// WARNING: Register/UnregisterPod functions should be efficient,
	// i.e. should not block on network operations.
	// 警告：Register/UnregisterPod 函数应该是高效的，即不应该阻塞网络操作。

	// RegisterPod registers all secrets from a given pod.
	// RegisterPod 函数从给定的 Pod 中注册所有的 Secret。
	RegisterPod(pod *v1.Pod)

	// UnregisterPod unregisters secrets from a given pod that are not
	// used by any other registered pod.
	// UnregisterPod 函数从给定的 Pod 中注销所有未被其他已注册 Pod 使用的 Secret。
	UnregisterPod(pod *v1.Pod)
}

// simpleSecretManager implements SecretManager interfaces with
// simple operations to apiserver.
// simpleSecretManager 实现了 SecretManager 接口，使用简单的操作来访问 apiserver。
type simpleSecretManager struct {
	kubeClient clientset.Interface
}

// NewSimpleSecretManager creates a new SecretManager instance.
// NewSimpleSecretManager 函数创建一个新的 SecretManager 实例。
func NewSimpleSecretManager(kubeClient clientset.Interface) Manager {
	return &simpleSecretManager{kubeClient: kubeClient}
}

func (s *simpleSecretManager) GetSecret(namespace, name string) (*v1.Secret, error) {
	return s.kubeClient.CoreV1().Secrets(namespace).Get(context.TODO(), name, metav1.GetOptions{})
}

func (s *simpleSecretManager) RegisterPod(pod *v1.Pod) {
}

func (s *simpleSecretManager) UnregisterPod(pod *v1.Pod) {
}

// secretManager keeps a store with secrets necessary
// for registered pods. Different implementations of the store
// may result in different semantics for freshness of secrets
// (e.g. ttl-based implementation vs watch-based implementation).
// secretManager 保存了已注册 Pod 所需的 Secret。存储的不同实现可能会导致 Secret 的不同语义（例如基于 ttl 的实现与基于 watch 的实现）。
type secretManager struct {
	manager manager.Manager
}

func (s *secretManager) GetSecret(namespace, name string) (*v1.Secret, error) {
	object, err := s.manager.GetObject(namespace, name)
	if err != nil {
		return nil, err
	}
	if secret, ok := object.(*v1.Secret); ok {
		return secret, nil
	}
	return nil, fmt.Errorf("unexpected object type: %v", object)
}

func (s *secretManager) RegisterPod(pod *v1.Pod) {
	s.manager.RegisterPod(pod)
}

func (s *secretManager) UnregisterPod(pod *v1.Pod) {
	s.manager.UnregisterPod(pod)
}

func getSecretNames(pod *v1.Pod) sets.String {
	result := sets.NewString()
	podutil.VisitPodSecretNames(pod, func(name string) bool {
		result.Insert(name)
		return true
	})
	return result
}

const (
	defaultTTL = time.Minute
)

// NewCachingSecretManager creates a manager that keeps a cache of all secrets
// necessary for registered pods.
// It implements the following logic:
//   - whenever a pod is created or updated, the cached versions of all secrets
//     are invalidated
//   - every GetObject() call tries to fetch the value from local cache; if it is
//     not there, invalidated or too old, we fetch it from apiserver and refresh the
//     value in cache; otherwise it is just fetched from cache
//
// NewCachingSecretManager 函数创建一个管理器，该管理器会缓存已注册 Pod 所需的所有 Secret。
// 它实现了以下逻辑：
//   - 每当创建或更新 Pod 时，所有 Secret 的缓存版本都会失效
//   - 每次调用 GetObject() 函数时，都会尝试从本地缓存中获取值；如果它不存在、失效或过期，我们会从 apiserver 中获取它并刷新缓存中的值；否则，它将从缓存中获取
func NewCachingSecretManager(kubeClient clientset.Interface, getTTL manager.GetObjectTTLFunc) Manager {
	getSecret := func(namespace, name string, opts metav1.GetOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().Secrets(namespace).Get(context.TODO(), name, opts)
	}
	secretStore := manager.NewObjectStore(getSecret, clock.RealClock{}, getTTL, defaultTTL)
	return &secretManager{
		manager: manager.NewCacheBasedManager(secretStore, getSecretNames),
	}
}

// NewWatchingSecretManager creates a manager that keeps a cache of all secrets
// necessary for registered pods.
// It implements the following logic:
//   - whenever a pod is created or updated, we start individual watches for all
//     referenced objects that aren't referenced from other registered pods
//   - every GetObject() returns a value from local cache propagated via watches
//
// NewWatchingSecretManager 函数创建一个管理器，该管理器会缓存已注册 Pod 所需的所有 Secret。
// 它实现了以下逻辑：
//   - 每当创建或更新 Pod 时，我们会为所有未被其他已注册 Pod 引用的对象启动单独的 watch
//   - 每次调用 GetObject() 函数时，都会从通过 watch 传播的本地缓存中返回值
func NewWatchingSecretManager(kubeClient clientset.Interface, resyncInterval time.Duration) Manager {
	listSecret := func(namespace string, opts metav1.ListOptions) (runtime.Object, error) {
		return kubeClient.CoreV1().Secrets(namespace).List(context.TODO(), opts)
	}
	watchSecret := func(namespace string, opts metav1.ListOptions) (watch.Interface, error) {
		return kubeClient.CoreV1().Secrets(namespace).Watch(context.TODO(), opts)
	}
	newSecret := func() runtime.Object {
		return &v1.Secret{}
	}
	isImmutable := func(object runtime.Object) bool {
		if secret, ok := object.(*v1.Secret); ok {
			return secret.Immutable != nil && *secret.Immutable
		}
		return false
	}
	gr := corev1.Resource("secret")
	return &secretManager{
		manager: manager.NewWatchBasedManager(listSecret, watchSecret, newSecret, isImmutable, gr, resyncInterval, getSecretNames),
	}
}
