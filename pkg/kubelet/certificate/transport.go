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

package certificate

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"k8s.io/klog/v2"

	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/wait"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/util/certificate"
	"k8s.io/client-go/util/connrotation"
)

// UpdateTransport 用一个动态使用管理器提供的 TLS 客户端认证的运输工具来记录一个 restconfig。
// 的证书用于TLS客户端认证。
//
// 该配置必须没有提供一个明确的传输。
//
// 返回的函数允许强行关闭所有活动连接。
//
// 返回的运输工具会定期检查管理器以确定
// 证书是否有变化。如果发生了变化，传输会关闭所有现有的客户
// 连接，迫使客户端与服务器重新握手并使用
// 新的证书。
//
// 如果设置了 exitAfter 持续时间，则将终止当前进程，如果证书
// 储存器中没有可用的证书（因为它已在磁盘上被删除或已损坏），则 exitAfter 持续时间将终止当前进程。
// 或者，如果证书已经过期，而服务器是有反应的。这允许
// 进程的父代或启动凭证有机会检索到一个新的初始
// 证书。
//
// stopCh 应该被用来表示当运输工具未被使用并且不需要
// 继续检查管理器。
func UpdateTransport(stopCh <-chan struct{}, clientConfig *restclient.Config, clientCertificateManager certificate.Manager, exitAfter time.Duration) (func(), error) {
	return updateTransport(stopCh, 10*time.Second, clientConfig, clientCertificateManager, exitAfter)
}

// updateTransport is an internal method that exposes how often this method checks that the
// client cert has changed.
func updateTransport(stopCh <-chan struct{}, period time.Duration, clientConfig *restclient.Config, clientCertificateManager certificate.Manager, exitAfter time.Duration) (func(), error) {
	if clientConfig.Transport != nil || clientConfig.Dial != nil {
		return nil, fmt.Errorf("there is already a transport or dialer configured")
	}

	d := connrotation.NewDialer((&net.Dialer{Timeout: 30 * time.Second, KeepAlive: 30 * time.Second}).DialContext)

	if clientCertificateManager != nil {
		if err := addCertRotation(stopCh, period, clientConfig, clientCertificateManager, exitAfter, d); err != nil {
			return nil, err
		}
	} else {
		clientConfig.Dial = d.DialContext
	}

	return d.CloseAll, nil
}

func addCertRotation(stopCh <-chan struct{}, period time.Duration, clientConfig *restclient.Config, clientCertificateManager certificate.Manager, exitAfter time.Duration, d *connrotation.Dialer) error {
	tlsConfig, err := restclient.TLSConfigFor(clientConfig)
	if err != nil {
		return fmt.Errorf("unable to configure TLS for the rest client: %v", err)
	}
	if tlsConfig == nil {
		tlsConfig = &tls.Config{}
	}

	tlsConfig.Certificates = nil
	tlsConfig.GetClientCertificate = func(requestInfo *tls.CertificateRequestInfo) (*tls.Certificate, error) {
		cert := clientCertificateManager.Current()
		if cert == nil {
			return &tls.Certificate{Certificate: nil}, nil
		}
		return cert, nil
	}

	lastCertAvailable := time.Now()
	lastCert := clientCertificateManager.Current()
	go wait.Until(func() {
		curr := clientCertificateManager.Current()

		if exitAfter > 0 {
			now := time.Now()
			if curr == nil {
				// the certificate has been deleted from disk or is otherwise corrupt
				if now.After(lastCertAvailable.Add(exitAfter)) {
					if clientCertificateManager.ServerHealthy() {
						klog.ErrorS(nil, "No valid client certificate is found and the server is responsive, exiting.", "lastCertificateAvailabilityTime", lastCertAvailable, "shutdownThreshold", exitAfter)
						os.Exit(1)
					} else {
						klog.ErrorS(nil, "No valid client certificate is found but the server is not responsive. A restart may be necessary to retrieve new initial credentials.", "lastCertificateAvailabilityTime", lastCertAvailable, "shutdownThreshold", exitAfter)
					}
				}
			} else {
				// the certificate is expired
				if now.After(curr.Leaf.NotAfter) {
					if clientCertificateManager.ServerHealthy() {
						klog.ErrorS(nil, "The currently active client certificate has expired and the server is responsive, exiting.")
						os.Exit(1)
					} else {
						klog.ErrorS(nil, "The currently active client certificate has expired, but the server is not responsive. A restart may be necessary to retrieve new initial credentials.")
					}
				}
				lastCertAvailable = now
			}
		}

		if curr == nil || lastCert == curr {
			// Cert hasn't been rotated.
			return
		}
		lastCert = curr

		klog.InfoS("Certificate rotation detected, shutting down client connections to start using new credentials")
		// The cert has been rotated. Close all existing connections to force the client
		// to reperform its TLS handshake with new cert.
		//
		// See: https://github.com/kubernetes-incubator/bootkube/pull/663#issuecomment-318506493
		d.CloseAll()
	}, period, stopCh)

	clientConfig.Transport = utilnet.SetTransportDefaults(&http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     tlsConfig,
		MaxIdleConnsPerHost: 25,
		DialContext:         d.DialContext,
	})

	// Zero out all existing TLS options since our new transport enforces them.
	clientConfig.CertData = nil
	clientConfig.KeyData = nil
	clientConfig.CertFile = ""
	clientConfig.KeyFile = ""
	clientConfig.CAData = nil
	clientConfig.CAFile = ""
	clientConfig.Insecure = false
	clientConfig.NextProtos = nil

	return nil
}
