/*
Copyright 2022 The Kubernetes Authors.

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

package v2beta1

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:prerelease-lifecycle-gen:introduced=1.26
// +k8s:prerelease-lifecycle-gen:deprecated=1.32
// +k8s:prerelease-lifecycle-gen:removed=1.35
// The deprecate and remove versions stated above are rough estimates and may be subject to change. We are estimating v2 types will be available in 1.28 and will support 4 versions where both v2beta1 and v2 are supported before deprecation.

// APIGroupDiscoveryList is a resource containing a list of APIGroupDiscovery.
// This is one of the types able to be returned from the /api and /apis endpoint and contains an aggregated
// list of API resources (built-ins, Custom Resource Definitions, resources from aggregated servers)
// that a cluster supports.
// APIGroupDiscoveryList 是一个包含 APIGroupDiscovery 列表的资源。
// 这是可以从 /api 和 /apis 端点返回的类型之一，并包含集群支持的聚合的 API 资源列表（内置的、自定义资源定义、来自聚合服务器的资源）。
type APIGroupDiscoveryList struct {
	v1.TypeMeta `json:",inline"`
	// ResourceVersion will not be set, because this does not have a replayable ordering among multiple apiservers.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	v1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// items is the list of groups for discovery. The groups are listed in priority order.
	Items []APIGroupDiscovery `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:prerelease-lifecycle-gen:introduced=1.26
// +k8s:prerelease-lifecycle-gen:deprecated=1.32
// +k8s:prerelease-lifecycle-gen:removed=1.35
// The deprecate and remove versions stated above are rough estimates and may be subject to change. We are estimating v2 types will be available in 1.28 and will support 4 versions where both v2beta1 and v2 are supported before deprecation.

// APIGroupDiscovery holds information about which resources are being served for all version of the API Group.
// It contains a list of APIVersionDiscovery that holds a list of APIResourceDiscovery types served for a version.
// Versions are in descending order of preference, with the first version being the preferred entry.
// APIGroupDiscovery 包含有关正在为 API 组的所有版本提供服务的资源的信息。 它包含一个 APIVersionDiscovery 的列表，该列表包含为版本提供服务的 APIResourceDiscovery 类型。 版本按优先顺序降序排列，第一个版本是首选项。
type APIGroupDiscovery struct {
	v1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// The only field completed will be name. For instance, resourceVersion will be empty.
	// name is the name of the API group whose discovery information is presented here.
	// name is allowed to be "" to represent the legacy, ungroupified resources.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	v1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// versions are the versions supported in this group. They are sorted in descending order of preference,
	// with the preferred version being the first entry.
	// +listType=map
	// +listMapKey=version
	Versions []APIVersionDiscovery `json:"versions,omitempty" protobuf:"bytes,2,rep,name=versions"`
}

// APIVersionDiscovery holds a list of APIResourceDiscovery types that are served for a particular version within an API Group.
// APIVersionDiscovery 包含为 API 组中特定版本提供服务的 APIResourceDiscovery 类型的列表。
type APIVersionDiscovery struct {
	// version is the name of the version within a group version.
	Version string `json:"version" protobuf:"bytes,1,opt,name=version"`
	// resources is a list of APIResourceDiscovery objects for the corresponding group version.
	// +listType=map
	// +listMapKey=resource
	// 资源是相应组版本的 APIResourceDiscovery 对象列表。
	Resources []APIResourceDiscovery `json:"resources,omitempty" protobuf:"bytes,2,rep,name=resources"`
	// freshness marks whether a group version's discovery document is up to date.
	// "Current" indicates the discovery document was recently
	// refreshed. "Stale" indicates the discovery document could not
	// be retrieved and the returned discovery document may be
	// significantly out of date. Clients that require the latest
	// version of the discovery information be retrieved before
	// performing an operation should not use the aggregated document
	// freshness标志着组版本的发现文档是否是最新的。
	// “当前”表示最近刷新了发现文档。 “陈旧”表示无法检索发现文档，并且返回的发现文档可能已经过时。
	// 需要在执行操作之前检索最新版本的发现信息的客户端不应使用聚合文档的新鲜度。
	Freshness DiscoveryFreshness `json:"freshness,omitempty" protobuf:"bytes,3,opt,name=freshness"`
}

// APIResourceDiscovery provides information about an API resource for discovery.
// APIResourceDiscovery 提供有关 API 资源的发现信息。
type APIResourceDiscovery struct {
	// resource is the plural name of the resource.  This is used in the URL path and is the unique identifier
	// for this resource across all versions in the API group.
	// Resources with non-empty groups are located at /apis/<APIGroupDiscovery.objectMeta.name>/<APIVersionDiscovery.version>/<APIResourceDiscovery.Resource>
	// Resources with empty groups are located at /api/v1/<APIResourceDiscovery.Resource>
	// resource 是资源的复数名称。 这用于 URL 路径，并且在 API 组中所有版本之间，此资源是唯一标识符。 具有非空组的资源位于 /apis/<APIGroupDiscovery.objectMeta.name>/<APIVersionDiscovery.version>/<APIResourceDiscovery.Resource> 具有空组的资源位于 /api/v1/<APIResourceDiscovery.Resource>
	Resource string `json:"resource" protobuf:"bytes,1,opt,name=resource"`
	// responseKind describes the group, version, and kind of the serialization schema for the object type this endpoint typically returns.
	// APIs may return other objects types at their discretion, such as error conditions, requests for alternate representations, or other operation specific behavior.
	// This value will be null if an APIService reports subresources but supports no operations on the parent resource
	// responseKind 描述了此端点通常返回的对象类型的序列化模式的组、版本和类型。 API 可以根据自己的意愿返回其他对象类型，例如错误条件、请求替代表示或其他操作特定行为。 如果 APIService 报告子资源但不支持对父资源的任何操作，则此值将为 null
	ResponseKind *v1.GroupVersionKind `json:"responseKind,omitempty" protobuf:"bytes,2,opt,name=responseKind"`
	// scope indicates the scope of a resource, either Cluster or Namespaced
	// scope 表示资源的范围，可以是集群或命名空间
	Scope ResourceScope `json:"scope" protobuf:"bytes,3,opt,name=scope"`
	// singularResource is the singular name of the resource.  This allows clients to handle plural and singular opaquely.
	// For many clients the singular form of the resource will be more understandable to users reading messages and should be used when integrating the name of the resource into a sentence.
	// The command line tool kubectl, for example, allows use of the singular resource name in place of plurals.
	// The singular forAm of a resource should always be an optional element - when in doubt use the canonical resource name.
	// singularResource 是资源的单数名称。 这允许客户端无差别地处理复数和单数。 对于许多客户端，资源的单数形式对于阅读消息的用户更容易理解，并且在将资源名称集成到句子中时应该使用它。 例如，命令行工具 kubectl 允许在复数形式的资源名称中使用单数形式的资源名称。 资源的单数形式应始终是可选元素 - 当有疑问时，请使用规范的资源名称。
	SingularResource string `json:"singularResource" protobuf:"bytes,4,opt,name=singularResource"`
	// verbs is a list of supported API operation types (this includes
	// but is not limited to get, list, watch, create, update, patch,
	// delete, deletecollection, and proxy).
	// +listType=set
	// verbs 是支持的 API 操作类型的列表（这包括但不限于 get、list、watch、create、update、patch、delete、deletecollection 和 proxy）。
	Verbs []string `json:"verbs" protobuf:"bytes,5,opt,name=verbs"`
	// shortNames is a list of suggested short names of the resource.
	// +listType=set
	// shortNames 是资源的建议缩写名称的列表。
	ShortNames []string `json:"shortNames,omitempty" protobuf:"bytes,6,rep,name=shortNames"`
	// categories is a list of the grouped resources this resource belongs to (e.g. 'all').
	// Clients may use this to simplify acting on multiple resource types at once.
	// +listType=set
	// categories 是此资源所属的分组资源的列表（例如“all”）。 客户端可以使用它来简化同时操作多种资源类型的操作。
	Categories []string `json:"categories,omitempty" protobuf:"bytes,7,rep,name=categories"`
	// subresources is a list of subresources provided by this resource. Subresources are located at /apis/<APIGroupDiscovery.objectMeta.name>/<APIVersionDiscovery.version>/<APIResourceDiscovery.Resource>/name-of-instance/<APIResourceDiscovery.subresources[i].subresource>
	// +listType=map
	// +listMapKey=subresource
	// subresources 是此资源提供的子资源的列表。 子资源位于 /apis/<APIGroupDiscovery.objectMeta.name>/<APIVersionDiscovery.version>/<APIResourceDiscovery.Resource>/name-of-instance/<APIResourceDiscovery.subresources[i].subresource>
	Subresources []APISubresourceDiscovery `json:"subresources,omitempty" protobuf:"bytes,8,rep,name=subresources"`
}

// ResourceScope is an enum defining the different scopes available to a resource.
// ResourceScope 是定义资源可用范围的枚举。
type ResourceScope string

const (
	ScopeCluster   ResourceScope = "Cluster"
	ScopeNamespace ResourceScope = "Namespaced"
)

// DiscoveryFreshness is an enum defining whether the Discovery document published by an apiservice is up to date (fresh).
// DiscoveryFreshness 是定义 apiservice 发布的 Discovery 文档是否是最新的（fresh）的枚举。
type DiscoveryFreshness string

const (
	DiscoveryFreshnessCurrent DiscoveryFreshness = "Current"
	DiscoveryFreshnessStale   DiscoveryFreshness = "Stale"
)

// APISubresourceDiscovery provides information about an API subresource for discovery.
// APISubresourceDiscovery 提供有关 API 子资源的发现信息。
type APISubresourceDiscovery struct {
	// subresource is the name of the subresource.  This is used in the URL path and is the unique identifier
	// for this resource across all versions.
	// subresource 是子资源的名称。 这用于 URL 路径，并且在所有版本中，此资源都是唯一标识符。
	Subresource string `json:"subresource" protobuf:"bytes,1,opt,name=subresource"`
	// responseKind describes the group, version, and kind of the serialization schema for the object type this endpoint typically returns.
	// Some subresources do not return normal resources, these will have null return types.
	// responseKind 描述了此端点通常返回的对象类型的序列化模式的组、版本和类型。 一些子资源不返回正常资源，这些资源将具有空返回类型。
	ResponseKind *v1.GroupVersionKind `json:"responseKind,omitempty" protobuf:"bytes,2,opt,name=responseKind"`
	// acceptedTypes describes the kinds that this endpoint accepts.
	// Subresources may accept the standard content types or define
	// custom negotiation schemes. The list may not be exhaustive for
	// all operations.
	// +listType=map
	// +listMapKey=group
	// +listMapKey=version
	// +listMapKey=kind
	// acceptedTypes 描述了此端点接受的类型。 子资源可以接受标准内容类型或定义自定义协商方案。 列表可能不完整，适用于所有操作。
	AcceptedTypes []v1.GroupVersionKind `json:"acceptedTypes,omitempty" protobuf:"bytes,3,rep,name=acceptedTypes"`
	// verbs is a list of supported API operation types (this includes
	// but is not limited to get, list, watch, create, update, patch,
	// delete, deletecollection, and proxy). Subresources may define
	// custom verbs outside the standard Kubernetes verb set. Clients
	// should expect the behavior of standard verbs to align with
	// Kubernetes interaction conventions.
	// +listType=set
	// verbs 是支持的 API 操作类型的列表（这包括但不限于 get、list、watch、create、update、patch、delete、deletecollection 和 proxy）。 子资源可以在标准 Kubernetes 动词集之外定义自定义动词。 客户端应该期望标准动词的行为与 Kubernetes 交互约定保持一致。
	Verbs []string `json:"verbs" protobuf:"bytes,4,opt,name=verbs"`
}
