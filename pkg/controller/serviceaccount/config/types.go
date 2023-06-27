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

// SAControllerConfiguration contains elements describing ServiceAccountController.
type SAControllerConfiguration struct {
	// serviceAccountKeyFile is the filename containing a PEM-encoded private RSA key
	// used to sign service account tokens.
	// serviceAccountKeyFile是包含用于签署服务账户令牌的PEM编码的私人RSA密钥的文件名。
	ServiceAccountKeyFile string
	// concurrentSATokenSyncs is the number of service account token syncing operations
	// that will be done concurrently.
	// concurrentSATokenSyncs是将并发进行的服务账户令牌同步操作的数量。
	ConcurrentSATokenSyncs int32
	// rootCAFile is the root certificate authority will be included in service
	// account's token secret. This must be a valid PEM-encoded CA bundle.
	// rootCAFile是根证书授权将被包括在服务帐户的令牌秘密中。这必须是一个有效的PEM编码的CA捆绑文件。
	RootCAFile string
}
