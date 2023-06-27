/*
Copyright 2014 The Kubernetes Authors.

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

package names

import (
	"fmt"

	utilrand "k8s.io/apimachinery/pkg/util/rand"
)

// NameGenerator generates names for objects. Some backends may have more information
// available to guide selection of new names and this interface hides those details.
// NameGenerator 为对象生成名称。 一些后端可能有更多的信息可用于指导新名称的选择，而此接口隐藏了这些细节。
type NameGenerator interface {
	// GenerateName generates a valid name from the base name, adding a random suffix to
	// the base. If base is valid, the returned name must also be valid. The generator is
	// responsible for knowing the maximum valid name length.
	// GenerateName 从基本名称生成有效名称，将随机后缀添加到基本名称。 如果基本名称有效，则返回的名称也必须有效。 生成器负责知道最大有效名称长度。
	GenerateName(base string) string
}

// simpleNameGenerator generates random names.
// simpleNameGenerator 生成随机名称。
type simpleNameGenerator struct{}

// SimpleNameGenerator is a generator that returns the name plus a random suffix of five alphanumerics
// when a name is requested. The string is guaranteed to not exceed the length of a standard Kubernetes
// name (63 characters)
// SimpleNameGenerator 是一个生成器，当请求名称时，它将名称加上五个字母数字的随机后缀。 该字符串保证不超过标准 Kubernetes 名称的长度（63 个字符）
var SimpleNameGenerator NameGenerator = simpleNameGenerator{}

const (
	// TODO: make this flexible for non-core resources with alternate naming rules.
	maxNameLength          = 63
	randomLength           = 5
	MaxGeneratedNameLength = maxNameLength - randomLength
)

func (simpleNameGenerator) GenerateName(base string) string {
	if len(base) > MaxGeneratedNameLength {
		base = base[:MaxGeneratedNameLength]
	}
	return fmt.Sprintf("%s%s", base, utilrand.String(randomLength))
}
