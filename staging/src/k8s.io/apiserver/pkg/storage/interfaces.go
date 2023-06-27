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

package storage

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
)

// Versioner abstracts setting and retrieving metadata fields from database response
// onto the object ot list. It is required to maintain storage invariants - updating an
// object twice with the same data except for the ResourceVersion and SelfLink must be
// a no-op. A resourceVersion of type uint64 is a 'raw' resourceVersion,
// intended to be sent directly to or from the backend. A resourceVersion of
// type string is a 'safe' resourceVersion, intended for consumption by users.
// Versioner 是一个接口，用于将数据库响应中的元数据字段设置到对象或列表中。它是维护存储不变量的必要条件
// - 使用相同的数据（除了ResourceVersion和SelfLink之外）更新对象两次必须是无操作的
// 。uint64类型的resourceVersion是“原始”resourceVersion，用于直接发送到或从后端。
// string类型的resourceVersion是“安全”的resourceVersion，用于用户消费。
type Versioner interface {
	// UpdateObject sets storage metadata into an API object. Returns an error if the object
	// cannot be updated correctly. May return nil if the requested object does not need metadata
	// from database.
	// UpdateObject 将存储元数据设置到 API 对象中。如果对象无法正确更新，则返回错误。如果请求的对象不需要来自数据库的元数据，则可能返回 nil。
	UpdateObject(obj runtime.Object, resourceVersion uint64) error
	// UpdateList sets the resource version into an API list object. Returns an error if the object
	// cannot be updated correctly. May return nil if the requested object does not need metadata from
	// database. continueValue is optional and indicates that more results are available if the client
	// passes that value to the server in a subsequent call. remainingItemCount indicates the number
	// of remaining objects if the list is partial. The remainingItemCount field is omitted during
	// serialization if it is set to nil.
	// UpdateList 将资源版本设置到 API 列表对象中。如果对象无法正确更新，则返回错误。如果请求的对象不需要来自数据库的元数据，则可能返回 nil。
	// continueValue 是可选的，如果客户端在后续调用中将该值传递给服务器，则表示如果有更多结果可用。
	// remainingItemCount 表示列表是否部分的剩余对象数。如果 remainingItemCount 字段设置为 nil，则在序列化期间将省略该字段。
	UpdateList(obj runtime.Object, resourceVersion uint64, continueValue string, remainingItemCount *int64) error
	// PrepareObjectForStorage should set SelfLink and ResourceVersion to the empty value. Should
	// return an error if the specified object cannot be updated.
	// PrepareObjectForStorage 应将 SelfLink 和 ResourceVersion 设置为空值。如果指定的对象无法更新，则应返回错误。
	PrepareObjectForStorage(obj runtime.Object) error
	// ObjectResourceVersion returns the resource version (for persistence) of the specified object.
	// Should return an error if the specified object does not have a persistable version.
	// ObjectResourceVersion 返回指定对象的资源版本（用于持久性）。如果指定的对象没有可持久化的版本，则应返回错误。
	ObjectResourceVersion(obj runtime.Object) (uint64, error)

	// ParseResourceVersion takes a resource version argument and
	// converts it to the storage backend. For watch we should pass to helper.Watch().
	// Because resourceVersion is an opaque value, the default watch
	// behavior for non-zero watch is to watch the next value (if you pass
	// "1", you will see updates from "2" onwards).
	// ParseResourceVersion 将资源版本参数转换为存储后端。对于 watch，我们应该传递给 helper.Watch()。
	// 因为 resourceVersion 是一个不透明的值，所以对于非零 watch，默认的 watch 行为是观察下一个值（如果您传递“1”，则将从“2”开始看到更新）。
	ParseResourceVersion(resourceVersion string) (uint64, error)
}

// ResponseMeta contains information about the database metadata that is associated with
// an object. It abstracts the actual underlying objects to prevent coupling with concrete
// database and to improve testability.
// ResponseMeta 包含与对象关联的数据库元数据的信息。它抽象了实际的底层对象，以防止与具体的数据库耦合，并提高可测试性。
type ResponseMeta struct {
	// TTL is the time to live of the node that contained the returned object. It may be
	// zero or negative in some cases (objects may be expired after the requested
	// expiration time due to server lag).
	// TTL 是包含返回对象的节点的生存时间。在某些情况下，它可能为零或负数（由于服务器滞后，对象可能在请求的过期时间之后过期）。
	TTL int64
	// The resource version of the node that contained the returned object.
	// 包含返回对象的节点的资源版本。
	ResourceVersion uint64
}

// IndexerFunc is a function that for a given object computes
// `<value of an index>` for a particular `<index>`.
// IndexerFunc 是一个函数，用于计算给定对象的特定索引的索引值。
type IndexerFunc func(obj runtime.Object) string

// IndexerFuncs is a mapping from `<index name>` to function that
// for a given object computes `<value for that index>`.
// IndexerFuncs 是从索引名称到函数的映射，用于计算给定对象的特定索引的索引值。
type IndexerFuncs map[string]IndexerFunc

// Everything accepts all objects.
var Everything = SelectionPredicate{
	Label: labels.Everything(),
	Field: fields.Everything(),
}

// MatchValue defines a pair (`<index name>`, `<value for that index>`).
type MatchValue struct {
	IndexName string
	Value     string
}

// Pass an UpdateFunc to Interface.GuaranteedUpdate to make an update
// that is guaranteed to succeed.
// See the comment for GuaranteedUpdate for more details.
// 将 UpdateFunc 传递给 Interface.GuaranteedUpdate 以进行更新，该更新保证成功。
// 有关更多详细信息，请参阅 GuaranteedUpdate 的注释。
type UpdateFunc func(input runtime.Object, res ResponseMeta) (output runtime.Object, ttl *uint64, err error)

// ValidateObjectFunc is a function to act on a given object. An error may be returned
// if the hook cannot be completed. The function may NOT transform the provided
// object.
// ValidateObjectFunc 是一个函数，用于对给定对象执行操作。如果无法完成钩子，则可能返回错误。该函数不得转换提供的对象。
type ValidateObjectFunc func(ctx context.Context, obj runtime.Object) error

// ValidateAllObjectFunc is a "admit everything" instance of ValidateObjectFunc.
// ValidateAllObjectFunc 是 ValidateObjectFunc 的“允许所有”实例。
func ValidateAllObjectFunc(ctx context.Context, obj runtime.Object) error {
	return nil
}

// Preconditions must be fulfilled before an operation (update, delete, etc.) is carried out.
// Preconditions 必须在执行操作（更新、删除等）之前得到满足。
type Preconditions struct {
	// Specifies the target UID.
	// +optional
	UID *types.UID `json:"uid,omitempty"`
	// Specifies the target ResourceVersion
	// +optional
	ResourceVersion *string `json:"resourceVersion,omitempty"`
}

// NewUIDPreconditions returns a Preconditions with UID set.
// NewUIDPreconditions 返回一个设置了 UID 的 Preconditions。
func NewUIDPreconditions(uid string) *Preconditions {
	u := types.UID(uid)
	return &Preconditions{UID: &u}
}

func (p *Preconditions) Check(key string, obj runtime.Object) error {
	if p == nil {
		return nil
	}
	objMeta, err := meta.Accessor(obj)
	if err != nil {
		return NewInternalErrorf(
			"can't enforce preconditions %v on un-introspectable object %v, got error: %v",
			*p,
			obj,
			err)
	}
	if p.UID != nil && *p.UID != objMeta.GetUID() {
		err := fmt.Sprintf(
			"Precondition failed: UID in precondition: %v, UID in object meta: %v",
			*p.UID,
			objMeta.GetUID())
		return NewInvalidObjError(key, err)
	}
	if p.ResourceVersion != nil && *p.ResourceVersion != objMeta.GetResourceVersion() {
		err := fmt.Sprintf(
			"Precondition failed: ResourceVersion in precondition: %v, ResourceVersion in object meta: %v",
			*p.ResourceVersion,
			objMeta.GetResourceVersion())
		return NewInvalidObjError(key, err)
	}
	return nil
}

// Interface offers a common interface for object marshaling/unmarshaling operations and
// hides all the storage-related operations behind it.
// Interface 提供了对象编组/解组操作的通用接口，并将所有与存储相关的操作隐藏在后面。
type Interface interface {
	// Returns Versioner associated with this interface.
	// 返回与此接口关联的 Versioner。
	Versioner() Versioner

	// Create adds a new object at a key unless it already exists. 'ttl' is time-to-live
	// in seconds (0 means forever). If no error is returned and out is not nil, out will be
	// set to the read value from database.
	// Create 在键处添加一个新对象，除非它已经存在。ttl 是秒为单位的生存时间（0 表示永远）。如果没有返回错误并且 out 不为 nil，则 out 将被设置为从数据库中读取的值。
	Create(ctx context.Context, key string, obj, out runtime.Object, ttl uint64) error

	// Delete removes the specified key and returns the value that existed at that spot.
	// If key didn't exist, it will return NotFound storage error.
	// If 'cachedExistingObject' is non-nil, it can be used as a suggestion about the
	// current version of the object to avoid read operation from storage to get it.
	// However, the implementations have to retry in case suggestion is stale.
	// Delete 删除指定的键并返回该位置存在的值。如果键不存在，则将返回 NotFound 存储错误。如果 cachedExistingObject 非空，则可以将其用作关于对象当前版本的建议，以避免从存储中读取操作来获取它。但是，实现必须重试，以防建议过时。
	Delete(
		ctx context.Context, key string, out runtime.Object, preconditions *Preconditions,
		validateDeletion ValidateObjectFunc, cachedExistingObject runtime.Object) error

	// Watch begins watching the specified key. Events are decoded into API objects,
	// and any items selected by 'p' are sent down to returned watch.Interface.
	// resourceVersion may be used to specify what version to begin watching,
	// which should be the current resourceVersion, and no longer rv+1
	// (e.g. reconnecting without missing any updates).
	// If resource version is "0", this interface will get current object at given key
	// and send it in an "ADDED" event, before watch starts.
	// Watch 开始监视指定的键。事件被解码为 API 对象，并且任何由 p 选择的项目都将发送到返回的 watch.Interface。resourceVersion 可用于指定要开始监视的版本，该版本应为当前 resourceVersion，并且不再是 rv+1（例如，重新连接而不丢失任何更新）。如果资源版本为“0”，则此接口将在给定键处获取当前对象，并在监视开始之前将其发送为“ADDED”事件。
	Watch(ctx context.Context, key string, opts ListOptions) (watch.Interface, error)

	// Get unmarshals object found at key into objPtr. On a not found error, will either
	// return a zero object of the requested type, or an error, depending on 'opts.ignoreNotFound'.
	// Treats empty responses and nil response nodes exactly like a not found error.
	// The returned contents may be delayed, but it is guaranteed that they will
	// match 'opts.ResourceVersion' according 'opts.ResourceVersionMatch'.
	// Get 将在键处找到的对象解组到 objPtr。在找不到错误的情况下，将返回请求类型的零对象，或者根据 opts.ignoreNotFound 返回错误。将空响应和空响应节点视为找不到错误。返回的内容可能会延迟，但是保证它们将根据 opts.ResourceVersionMatch 匹配 opts.ResourceVersion。
	Get(ctx context.Context, key string, opts GetOptions, objPtr runtime.Object) error

	// GetList unmarshalls objects found at key into a *List api object (an object
	// that satisfies runtime.IsList definition).
	// If 'opts.Recursive' is false, 'key' is used as an exact match. If `opts.Recursive'
	// is true, 'key' is used as a prefix.
	// The returned contents may be delayed, but it is guaranteed that they will
	// match 'opts.ResourceVersion' according 'opts.ResourceVersionMatch'.
	// GetList 将在键处找到的对象解组到 *List api 对象（满足 runtime.IsList 定义的对象）中。如果 opts.Recursive 为 false，则 key 用作精确匹配。如果 opts.Recursive 为 true，则 key 用作前缀。返回的内容可能会延迟，但是保证它们将根据 opts.ResourceVersionMatch 匹配 opts.ResourceVersion。
	GetList(ctx context.Context, key string, opts ListOptions, listObj runtime.Object) error

	// GuaranteedUpdate keeps calling 'tryUpdate()' to update key 'key' (of type 'destination')
	// retrying the update until success if there is index conflict.
	// Note that object passed to tryUpdate may change across invocations of tryUpdate() if
	// other writers are simultaneously updating it, so tryUpdate() needs to take into account
	// the current contents of the object when deciding how the update object should look.
	// If the key doesn't exist, it will return NotFound storage error if ignoreNotFound=false
	// else `destination` will be set to the zero value of it's type.
	// If the eventual successful invocation of `tryUpdate` returns an output with the same serialized
	// contents as the input, it won't perform any update, but instead set `destination` to an object with those
	// contents.
	// If 'cachedExistingObject' is non-nil, it can be used as a suggestion about the
	// current version of the object to avoid read operation from storage to get it.
	// However, the implementations have to retry in case suggestion is stale.
	//
	// Example:
	//
	// s := /* implementation of Interface */
	// err := s.GuaranteedUpdate(
	//     "myKey", &MyType{}, true, preconditions,
	//     func(input runtime.Object, res ResponseMeta) (runtime.Object, *uint64, error) {
	//       // Before each invocation of the user defined function, "input" is reset to
	//       // current contents for "myKey" in database.
	//       curr := input.(*MyType)  // Guaranteed to succeed.
	//
	//       // Make the modification
	//       curr.Counter++
	//
	//       // Return the modified object - return an error to stop iterating. Return
	//       // a uint64 to alter the TTL on the object, or nil to keep it the same value.
	//       return cur, nil, nil
	//    }, cachedExistingObject
	// )
	// GuaranteedUpdate 一直调用 tryUpdate() 来更新键 key（类型为 destination）直到成功，如果有索引冲突则重试更新。请注意，传递给 tryUpdate() 的对象可能会在调用 tryUpdate() 时发生更改，如果其他写入器同时更新它，则 tryUpdate() 需要考虑对象的当前内容，以确定更新对象应该如何查看。如果键不存在，则如果 ignoreNotFound=false 则会返回 NotFound 存储错误，否则 destination 将设置为其类型的零值。如果最终成功调用 tryUpdate 返回与输入具有相同序列化内容的输出，则不会执行任何更新，而是将 destination 设置为具有这些内容的对象。如果 cachedExistingObject 非空，则可以将其用作关于避免从存储中读取以获取它的对象的当前版本的建议。但是，实现必须重试，因为建议可能已过时。
	GuaranteedUpdate(
		ctx context.Context, key string, destination runtime.Object, ignoreNotFound bool,
		preconditions *Preconditions, tryUpdate UpdateFunc, cachedExistingObject runtime.Object) error

	// Count returns number of different entries under the key (generally being path prefix).
	// Count 返回键下的不同条目数（通常是路径前缀）。
	Count(key string) (int64, error)
}

// GetOptions provides the options that may be provided for storage get operations.
// GetOptions 提供可能用于存储 get 操作的选项。
type GetOptions struct {
	// IgnoreNotFound determines what is returned if the requested object is not found. If
	// true, a zero object is returned. If false, an error is returned.
	// IgnoreNotFound 确定如果找不到请求的对象，则返回什么。如果为 true，则返回零对象。如果为 false，则返回错误。
	IgnoreNotFound bool
	// ResourceVersion provides a resource version constraint to apply to the get operation
	// as a "not older than" constraint: the result contains data at least as new as the provided
	// ResourceVersion. The newest available data is preferred, but any data not older than this
	// ResourceVersion may be served.
	// ResourceVersion 提供一个资源版本约束，用于将 get 操作应用为“不旧于”约束：结果包含至少与提供的 ResourceVersion 一样新的数据。首选最新可用的数据，但是可以提供比此版本更旧的数据。
	ResourceVersion string
}

// ListOptions provides the options that may be provided for storage list operations.
// ListOptions 提供可能用于存储列表操作的选项。
type ListOptions struct {
	// ResourceVersion provides a resource version constraint to apply to the list operation
	// as a "not older than" constraint: the result contains data at least as new as the provided
	// ResourceVersion. The newest available data is preferred, but any data not older than this
	// ResourceVersion may be served.
	// ResourceVersion 提供一个资源版本约束，用于将列表操作应用为“不旧于”约束：结果包含至少与提供的 ResourceVersion 一样新的数据。首选最新可用的数据，但是可以提供比此版本更旧的数据。
	ResourceVersion string
	// ResourceVersionMatch provides the rule for how the resource version constraint applies. If set
	// to the default value "" the legacy resource version semantic apply.
	// ResourceVersionMatch 提供了如何应用资源版本约束的规则。如果设置为默认值“”，则将应用旧版资源版本语义。
	ResourceVersionMatch metav1.ResourceVersionMatch
	// Predicate provides the selection rules for the list operation.
	// Predicate 为列表操作提供选择规则。
	Predicate SelectionPredicate
	// Recursive determines whether the list or watch is defined for a single object located at the
	// given key, or for the whole set of objects with the given key as a prefix.
	// Recursive 确定列表或监视是否为位于给定键处的单个对象定义，还是为具有给定键作为前缀的整个对象集定义。
	Recursive bool
	// ProgressNotify determines whether storage-originated bookmark (progress notify) events should
	// be delivered to the users. The option is ignored for non-watch requests.
	// ProgressNotify 确定是否应将存储源生成的书签（进度通知）事件传递给用户。选项对非监视请求被忽略。
	ProgressNotify bool
}
