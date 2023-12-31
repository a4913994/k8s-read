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

// CSRSigningControllerConfiguration contains elements describing CSRSigningController.
type CSRSigningControllerConfiguration struct {
	// clusterSigningCertFile is the filename containing a PEM-encoded
	// X509 CA certificate used to issue cluster-scoped certificates
	// clusterSigningCertFile 是包含用于颁发集群范围证书的 PEM 编码 X509 CA 证书的文件名
	ClusterSigningCertFile string
	// clusterSigningCertFile is the filename containing a PEM-encoded
	// RSA or ECDSA private key used to issue cluster-scoped certificates
	// clusterSigningCertFile 是包含用于颁发集群范围证书的 PEM 编码 RSA 或 ECDSA 私钥的文件名
	ClusterSigningKeyFile string

	// kubeletServingSignerConfiguration holds the certificate and key used to issue certificates for the kubernetes.io/kubelet-serving signer
	// kubeletServingSignerConfiguration 持有用于为 kubernetes.iokubelet-serving signer 颁发证书的证书和密钥
	KubeletServingSignerConfiguration CSRSigningConfiguration
	// kubeletClientSignerConfiguration holds the certificate and key used to issue certificates for the kubernetes.io/kube-apiserver-client-kubelet
	// kubeletClientSignerConfiguration 持有用于为 kubernetes.iokube-apiserver-client-kubelet 颁发证书的证书和密钥
	KubeletClientSignerConfiguration CSRSigningConfiguration
	// kubeAPIServerClientSignerConfiguration holds the certificate and key used to issue certificates for the kubernetes.io/kube-apiserver-client
	// kubeAPIServerClientSignerConfiguration 持有用于为 kubernetes.iokube-apiserver-client 颁发证书的证书和密钥
	KubeAPIServerClientSignerConfiguration CSRSigningConfiguration
	// legacyUnknownSignerConfiguration holds the certificate and key used to issue certificates for the kubernetes.io/legacy-unknown
	// legacyUnknownSignerConfiguration 持有用于为 kubernetes.iolegacy-unknown 颁发证书的证书和密钥
	LegacyUnknownSignerConfiguration CSRSigningConfiguration

	// clusterSigningDuration is the max length of duration signed certificates will be given.
	// Individual CSRs may request shorter certs by setting spec.expirationSeconds.
	// clusterSigningDuration 是将给出签名证书的最大持续时间。个别 CSR 可以通过设置 spec.expirationSeconds 来请求更短的证书
	ClusterSigningDuration metav1.Duration
}

// CSRSigningConfiguration holds information about a particular CSR signer
type CSRSigningConfiguration struct {
	// certFile is the filename containing a PEM-encoded
	// X509 CA certificate used to issue certificates
	CertFile string
	// keyFile is the filename containing a PEM-encoded
	// RSA or ECDSA private key used to issue certificates
	KeyFile string
}
