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

package rest

import (
	"context"
	"io"
	"net/http"
	"net/url"

	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"sigs.k8s.io/structured-merge-diff/v4/fieldpath"
)

//TODO:
// Storage interfaces need to be separated into two groups; those that operate
// on collections and those that operate on individually named items.
// Collection interfaces:
// (Method: Current -> Proposed)
//    GET: Lister -> CollectionGetter
//    WATCH: Watcher -> CollectionWatcher
//    CREATE: Creater -> CollectionCreater
//    DELETE: (n/a) -> CollectionDeleter
//    UPDATE: (n/a) -> CollectionUpdater
//
// Single item interfaces:
// (Method: Current -> Proposed)
//    GET: Getter -> NamedGetter
//    WATCH: (n/a) -> NamedWatcher
//    CREATE: (n/a) -> NamedCreater
//    DELETE: Deleter -> NamedDeleter
//    UPDATE: Update -> NamedUpdater

// Storage is a generic interface for RESTful storage services.
// Resources which are exported to the RESTful API of apiserver need to implement this interface. It is expected
// that objects may implement any of the below interfaces.
// Storage 是一个通用的接口，用于 RESTful 存储服务。
// 需要实现此接口的资源将导出到 apiserver 的 RESTful API。 期望对象可能实现以下任何接口。
type Storage interface {
	// New returns an empty object that can be used with Create and Update after request data has been put into it.
	// This object must be a pointer type for use with Codec.DecodeInto([]byte, runtime.Object)
	// New 返回一个空对象，可以在将请求数据放入其中后使用 Create 和 Update。
	// 此对象必须是指针类型，以便与 Codec.DecodeInto([]byte, runtime.Object) 一起使用。
	New() runtime.Object

	// Destroy cleans up its resources on shutdown.
	// Destroy has to be implemented in thread-safe way and be prepared
	// for being called more than once.
	// Destroy 在关闭时清理其资源。
	// Destroy 必须以线程安全的方式实现，并且必须准备好多次调用。
	Destroy()
}

// Scoper indicates what scope the resource is at. It must be specified.
// It is usually provided automatically based on your strategy.
// Scoper 表示资源所在的范围。 必须指定它。
// 它通常是根据您的策略自动提供的。
type Scoper interface {
	// NamespaceScoped returns true if the storage is namespaced
	// NamespaceScoped 如果存储是有命名空间的，则返回 true
	NamespaceScoped() bool
}

// KindProvider specifies a different kind for its API than for its internal storage.  This is necessary for external
// objects that are not compiled into the api server.  For such objects, there is no in-memory representation for
// the object, so they must be represented as generic objects (e.g. runtime.Unknown), but when we present the object as part of
// API discovery we want to present the specific kind, not the generic internal representation.
// KindProvider 指定其 API 的类型与其内部存储的类型不同。 这对于未编译到 API 服务器中的外部对象是必需的。
// 对于这样的对象，没有用于对象的内存表示，因此它们必须表示为通用对象（例如 runtime.Unknown），但是当我们将对象作为 API 发现的一部分呈现时，我们希望呈现特定的类型，而不是通用的内部表示。
type KindProvider interface {
	Kind() string
}

// ShortNamesProvider is an interface for RESTful storage services. Delivers a list of short names for a resource. The list is used by kubectl to have short names representation of resources.
// ShortNamesProvider 是 RESTful 存储服务的接口。 为资源提供短名称列表。 列表由 kubectl 使用，以便具有资源的短名称表示形式。
type ShortNamesProvider interface {
	ShortNames() []string
}

// CategoriesProvider allows a resource to specify which groups of resources (categories) it's part of. Categories can
// be used by API clients to refer to a batch of resources by using a single name (e.g. "all" could translate to "pod,rc,svc,...").
// CategoriesProvider 允许资源指定它所属的哪些资源组（类别）。 类别可以由 API 客户端使用，以便通过使用单个名称（例如，“all”可以转换为“pod，rc，svc，...”）来引用一批资源。
type CategoriesProvider interface {
	Categories() []string
}

// GroupVersionKindProvider is used to specify a particular GroupVersionKind to discovery.  This is used for polymorphic endpoints
// which generally point to foreign versions.  Scale refers to Scale.v1beta1.extensions for instance.
// This trumps KindProvider since it is capable of providing the information required.
// TODO KindProvider (only used by federation) should be removed and replaced with this, but that presents greater risk late in 1.8.
// GroupVersionKindProvider 用于指定特定的 GroupVersionKind 到发现。 这用于多态端点，通常指向外部版本。 例如，Scale 指向 Scale.v1beta1.extensions。
// 这个优先于 KindProvider，因为它能够提供所需的信息。
type GroupVersionKindProvider interface {
	GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind
}

// GroupVersionAcceptor is used to determine if a particular GroupVersion is acceptable to send to an endpoint.
// This is used for endpoints which accept multiple versions (which is extremely rare).
// The only known instance is pods/evictions which accepts policy/v1, but also policy/v1beta1 for backwards compatibility.
// GroupVersionAcceptor 用于确定特定 GroupVersion 是否适合发送到端点。 这用于接受多个版本的端点（这是非常罕见的）。
// 已知的唯一实例是 pods/evictions，它接受 policy/v1，但也接受 policy/v1beta1 以保持向后兼容性。
type GroupVersionAcceptor interface {
	AcceptsGroupVersion(gv schema.GroupVersion) bool
}

// Lister is an object that can retrieve resources that match the provided field and label criteria.
// Lister 是一个可以检索与提供的字段和标签标准匹配的资源的对象。
type Lister interface {
	// NewList returns an empty object that can be used with the List call.
	// This object must be a pointer type for use with Codec.DecodeInto([]byte, runtime.Object)
	// NewList 返回一个空对象，可以在将请求数据放入其中后使用 List。
	// 此对象必须是指针类型，以便与 Codec.DecodeInto([]byte, runtime.Object) 一起使用。
	NewList() runtime.Object
	// List selects resources in the storage which match to the selector. 'options' can be nil.
	// List 选择与选择器匹配的存储中的资源。 'options' 可以为 nil。
	List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error)
	// TableConvertor ensures all list implementers also implement table conversion
	// TableConvertor 确保所有列表实现者也实现了表格转换
	TableConvertor
}

// Getter is an object that can retrieve a named RESTful resource.
// Getter 是一个可以检索命名 RESTful 资源的对象。
type Getter interface {
	// Get finds a resource in the storage by name and returns it.
	// Although it can return an arbitrary error value, IsNotFound(err) is true for the
	// returned error value err when the specified resource is not found.
	// Get 按名称在存储中查找资源并将其返回。
	// 虽然它可以返回任意错误值，但当指定的资源未找到时，IsNotFound(err) 对于返回的错误值 err 是 true。
	Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error)
}

// GetterWithOptions is an object that retrieve a named RESTful resource and takes
// additional options on the get request. It allows a caller to also receive the
// subpath of the GET request.
// GetterWithOptions 是一个可以检索命名 RESTful 资源的对象，并在 get 请求上获取其他选项。 它允许调用者还接收 GET 请求的子路径。
type GetterWithOptions interface {
	// Get finds a resource in the storage by name and returns it.
	// Although it can return an arbitrary error value, IsNotFound(err) is true for the
	// returned error value err when the specified resource is not found.
	// The options object passed to it is of the same type returned by the NewGetOptions
	// method.
	// TODO: Pass metav1.GetOptions.
	// Get 按名称在存储中查找资源并将其返回。
	// 虽然它可以返回任意错误值，但当指定的资源未找到时，IsNotFound(err) 对于返回的错误值 err 是 true。
	// 传递给它的选项对象与 NewGetOptions 方法返回的类型相同。
	Get(ctx context.Context, name string, options runtime.Object) (runtime.Object, error)

	// NewGetOptions returns an empty options object that will be used to pass
	// options to the Get method. It may return a bool and a string, if true, the
	// value of the request path below the object will be included as the named
	// string in the serialization of the runtime object. E.g., returning "path"
	// will convert the trailing request scheme value to "path" in the map[string][]string
	// passed to the converter.
	// NewGetOptions 返回一个空选项对象，该对象将用于将选项传递给 Get 方法。 如果为 true，请求路径下的值将被包含在运行时对象的序列化中，作为命名的字符串。
	// 例如，返回 "path" 将将请求方案值转换为 "path"，并将其作为 map[string][]string 传递给转换器。
	NewGetOptions() (runtime.Object, bool, string)
}

type TableConvertor interface {
	ConvertToTable(ctx context.Context, object runtime.Object, tableOptions runtime.Object) (*metav1.Table, error)
}

// GracefulDeleter knows how to pass deletion options to allow delayed deletion of a
// RESTful object.
// GracefulDeleter 知道如何传递删除选项，以允许延迟删除 RESTful 对象。
type GracefulDeleter interface {
	// Delete finds a resource in the storage and deletes it.
	// The delete attempt is validated by the deleteValidation first.
	// If options are provided, the resource will attempt to honor them or return an invalid
	// request error.
	// Although it can return an arbitrary error value, IsNotFound(err) is true for the
	// returned error value err when the specified resource is not found.
	// Delete *may* return the object that was deleted, or a status object indicating additional
	// information about deletion.
	// It also returns a boolean which is set to true if the resource was instantly
	// deleted or false if it will be deleted asynchronously.
	// Delete 发现存储中的资源并删除它。
	// 首先通过 deleteValidation 验证删除尝试。
	// 如果提供了选项，则资源将尝试遵守它们或返回无效的请求错误。
	// 虽然它可以返回任意错误值，但当指定的资源未找到时，IsNotFound(err) 对于返回的错误值 err 是 true。
	// Delete *may* 返回已删除的对象，或者是一个状态对象，指示有关删除的其他信息。
	// 它还返回一个布尔值，如果资源被立即删除，则为 true；如果它将被异步删除，则为 false。
	Delete(ctx context.Context, name string, deleteValidation ValidateObjectFunc, options *metav1.DeleteOptions) (runtime.Object, bool, error)
}

// MayReturnFullObjectDeleter may return deleted object (instead of a simple status) on deletion.
// MayReturnFullObjectDeleter 可能会在删除时返回已删除的对象（而不是简单的状态）。
type MayReturnFullObjectDeleter interface {
	DeleteReturnsDeletedObject() bool
}

// CollectionDeleter is an object that can delete a collection
// of RESTful resources.
// CollectionDeleter 是一个可以删除 RESTful 资源集合的对象。
type CollectionDeleter interface {
	// DeleteCollection selects all resources in the storage matching given 'listOptions'
	// and deletes them. The delete attempt is validated by the deleteValidation first.
	// If 'options' are provided, the resource will attempt to honor them or return an
	// invalid request error.
	// DeleteCollection may not be atomic - i.e. it may delete some objects and still
	// return an error after it. On success, returns a list of deleted objects.
	// DeleteCollection 选择与给定 'listOptions' 匹配的存储中的所有资源，并将其删除。 首先通过 deleteValidation 验证删除尝试。
	// 如果提供了 'options'，则资源将尝试遵守它们或返回无效的请求错误。
	// DeleteCollection 可能不是原子的 - 即它可能删除一些对象，但仍然在删除后返回错误。 成功时，返回已删除对象的列表。
	DeleteCollection(ctx context.Context, deleteValidation ValidateObjectFunc, options *metav1.DeleteOptions, listOptions *metainternalversion.ListOptions) (runtime.Object, error)
}

// Creater is an object that can create an instance of a RESTful object.
// Creater 是一个可以创建 RESTful 对象实例的对象。
type Creater interface {
	// New returns an empty object that can be used with Create after request data has been put into it.
	// This object must be a pointer type for use with Codec.DecodeInto([]byte, runtime.Object)
	// New 返回一个空对象，可以在将请求数据放入其中后使用 Create。
	// 此对象必须是指针类型，以便与 Codec.DecodeInto([]byte, runtime.Object) 一起使用。
	New() runtime.Object

	// Create creates a new version of a resource.
	// Create 创建资源的新版本。
	Create(ctx context.Context, obj runtime.Object, createValidation ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error)
}

// NamedCreater is an object that can create an instance of a RESTful object using a name parameter.
// NamedCreater 是一个可以使用名称参数创建 RESTful 对象实例的对象。
type NamedCreater interface {
	// New returns an empty object that can be used with Create after request data has been put into it.
	// This object must be a pointer type for use with Codec.DecodeInto([]byte, runtime.Object)
	// New 返回一个空对象，可以在将请求数据放入其中后使用 Create。
	// 此对象必须是指针类型，以便与 Codec.DecodeInto([]byte, runtime.Object) 一起使用。
	New() runtime.Object

	// Create creates a new version of a resource. It expects a name parameter from the path.
	// This is needed for create operations on subresources which include the name of the parent
	// resource in the path.
	// Create 创建资源的新版本。 它期望从路径中获取名称参数。
	// 这对于在路径中包含父资源名称的子资源上的创建操作是必需的。
	Create(ctx context.Context, name string, obj runtime.Object, createValidation ValidateObjectFunc, options *metav1.CreateOptions) (runtime.Object, error)
}

// UpdatedObjectInfo provides information about an updated object to an Updater.
// It requires access to the old object in order to return the newly updated object.
// UpdatedObjectInfo 提供有关更新对象的信息以更新 Updater。
// 它需要访问旧对象才能返回新更新的对象。
type UpdatedObjectInfo interface {
	// Returns preconditions built from the updated object, if applicable.
	// May return nil, or a preconditions object containing nil fields,
	// if no preconditions can be determined from the updated object.
	// 返回从更新对象构建的前提条件（如果适用）。
	// 如果无法从更新的对象确定前提条件，则可以返回 nil，或者包含 nil 字段的前提条件对象。
	Preconditions() *metav1.Preconditions

	// UpdatedObject returns the updated object, given a context and old object.
	// The only time an empty oldObj should be passed in is if a "create on update" is occurring (there is no oldObj).
	// UpdatedObject 返回更新的对象，给定上下文和旧对象。
	// 仅当发生“创建更新”时才应传入空 oldObj（没有 oldObj）。
	UpdatedObject(ctx context.Context, oldObj runtime.Object) (newObj runtime.Object, err error)
}

// ValidateObjectFunc is a function to act on a given object. An error may be returned
// if the hook cannot be completed. A ValidateObjectFunc may NOT transform the provided
// object.
// ValidateObjectFunc 是一个函数，用于对给定对象执行操作。 如果无法完成钩子，则可能返回错误。
// ValidateObjectFunc 不得转换提供的对象。
type ValidateObjectFunc func(ctx context.Context, obj runtime.Object) error

// ValidateAllObjectFunc is a "admit everything" instance of ValidateObjectFunc.
// ValidateAllObjectFunc 是 ValidateObjectFunc 的“允许所有”实例。
func ValidateAllObjectFunc(ctx context.Context, obj runtime.Object) error {
	return nil
}

// ValidateObjectUpdateFunc is a function to act on a given object and its predecessor.
// An error may be returned if the hook cannot be completed. An UpdateObjectFunc
// may NOT transform the provided object.
// ValidateObjectUpdateFunc 是一个函数，用于对给定对象和其前身执行操作。
// 如果无法完成钩子，则可能返回错误。 UpdateObjectFunc 不得转换提供的对象。
type ValidateObjectUpdateFunc func(ctx context.Context, obj, old runtime.Object) error

// ValidateAllObjectUpdateFunc is a "admit everything" instance of ValidateObjectUpdateFunc.
// ValidateAllObjectUpdateFunc 是 ValidateObjectUpdateFunc 的“允许所有”实例。
func ValidateAllObjectUpdateFunc(ctx context.Context, obj, old runtime.Object) error {
	return nil
}

// Updater is an object that can update an instance of a RESTful object.
// Updater 是一个可以更新 RESTful 对象实例的对象。
type Updater interface {
	// New returns an empty object that can be used with Update after request data has been put into it.
	// This object must be a pointer type for use with Codec.DecodeInto([]byte, runtime.Object)
	// New 返回一个空对象，可以在将请求数据放入其中后使用 Update。
	// 此对象必须是指针类型，以便与 Codec.DecodeInto([]byte, runtime.Object) 一起使用。
	New() runtime.Object

	// Update finds a resource in the storage and updates it. Some implementations
	// may allow updates creates the object - they should set the created boolean
	// to true.
	// Update 在存储中查找资源并更新它。 一些实现可能允许更新创建对象 - 它们应将 created 布尔值设置为 true。
	Update(ctx context.Context, name string, objInfo UpdatedObjectInfo, createValidation ValidateObjectFunc, updateValidation ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error)
}

// CreaterUpdater is a storage object that must support both create and update.
// Go prevents embedded interfaces that implement the same method.
// CreaterUpdater 是一个必须支持创建和更新的存储对象。
type CreaterUpdater interface {
	Creater
	Update(ctx context.Context, name string, objInfo UpdatedObjectInfo, createValidation ValidateObjectFunc, updateValidation ValidateObjectUpdateFunc, forceAllowCreate bool, options *metav1.UpdateOptions) (runtime.Object, bool, error)
}

// CreaterUpdater must satisfy the Updater interface.
// CreaterUpdater 必须满足 Updater 接口。
var _ Updater = CreaterUpdater(nil)

// Patcher is a storage object that supports both get and update.
// Patcher 是一个支持 get 和 update 的存储对象。
type Patcher interface {
	Getter
	Updater
}

// Watcher should be implemented by all Storage objects that
// want to offer the ability to watch for changes through the watch api.
// Watcher 应由所有希望通过 watch api 提供监视更改的存储对象实现。
type Watcher interface {
	// 'label' selects on labels; 'field' selects on the object's fields. Not all fields
	// are supported; an error should be returned if 'field' tries to select on a field that
	// isn't supported. 'resourceVersion' allows for continuing/starting a watch at a
	// particular version.
	// 'label' 选择标签；'field' 选择对象的字段。 不支持所有字段； 如果 'field' 尝试选择不受支持的字段，则应返回错误。
	// 'resourceVersion' 允许在特定版本上继续/开始监视。
	Watch(ctx context.Context, options *metainternalversion.ListOptions) (watch.Interface, error)
}

// StandardStorage is an interface covering the common verbs. Provided for testing whether a
// resource satisfies the normal storage methods. Use Storage when passing opaque storage objects.
// StandardStorage 是涵盖常见动词的接口。 用于测试资源是否满足正常存储方法。 使用 Storage 传递不透明的存储对象。
type StandardStorage interface {
	Getter
	Lister
	CreaterUpdater
	GracefulDeleter
	CollectionDeleter
	Watcher

	// Destroy cleans up its resources on shutdown.
	// Destroy has to be implemented in thread-safe way and be prepared
	// for being called more than once.
	Destroy()
}

// Redirector know how to return a remote resource's location.
// Redirector 知道如何返回远程资源的位置。
type Redirector interface {
	// ResourceLocation should return the remote location of the given resource, and an optional transport to use to request it, or an error.
	// ResourceLocation 应返回给定资源的远程位置以及用于请求它的可选传输，或者返回错误。
	ResourceLocation(ctx context.Context, id string) (remoteLocation *url.URL, transport http.RoundTripper, err error)
}

// Responder abstracts the normal response behavior for a REST method and is passed to callers that
// may wish to handle the response directly in some cases, but delegate to the normal error or object
// behavior in other cases.
// Responder 抽象了 REST 方法的正常响应行为，并将其传递给可能希望在某些情况下直接处理响应，但在其他情况下委托给正常错误或对象行为的调用者。
type Responder interface {
	// Object writes the provided object to the response. Invoking this method multiple times is undefined.
	// Object 将提供的对象写入响应。 调用此方法多次是未定义的。
	Object(statusCode int, obj runtime.Object)
	// Error writes the provided error to the response. This method may only be invoked once.
	// Error 将提供的错误写入响应。 只能调用此方法一次。
	Error(err error)
}

// Connecter is a storage object that responds to a connection request.
// Connecter 是一个响应连接请求的存储对象。
type Connecter interface {
	// Connect returns an http.Handler that will handle the request/response for a given API invocation.
	// The provided responder may be used for common API responses. The responder will write both status
	// code and body, so the ServeHTTP method should exit after invoking the responder. The Handler will
	// be used for a single API request and then discarded. The Responder is guaranteed to write to the
	// same http.ResponseWriter passed to ServeHTTP.
	// Connect 返回一个 http.Handler，该 http.Handler 将处理给定 API 调用的请求/响应。
	// 提供的响应者可用于常见的 API 响应。 响应者将写入状态代码和正文，因此在调用响应者后，ServeHTTP 方法应退出。
	// 该处理程序将用于单个 API 请求，然后丢弃。 保证响应者将写入传递给 ServeHTTP 的相同 http.ResponseWriter。
	Connect(ctx context.Context, id string, options runtime.Object, r Responder) (http.Handler, error)

	// NewConnectOptions returns an empty options object that will be used to pass
	// options to the Connect method. If nil, then a nil options object is passed to
	// Connect. It may return a bool and a string. If true, the value of the request
	// path below the object will be included as the named string in the serialization
	// of the runtime object.
	// NewConnectOptions 返回一个空选项对象，该选项对象将用于将选项传递给 Connect 方法。
	// 如果为 nil，则将传递一个空选项对象给 Connect。 它可以返回一个布尔值和一个字符串。
	// 如果为 true，则请求路径下的对象的值将作为运行时对象的序列化中的命名字符串包含在内。
	NewConnectOptions() (runtime.Object, bool, string)

	// ConnectMethods returns the list of HTTP methods handled by Connect
	// ConnectMethods 返回由 Connect 处理的 HTTP 方法列表
	ConnectMethods() []string
}

// ResourceStreamer is an interface implemented by objects that prefer to be streamed from the server
// instead of decoded directly.
// ResourceStreamer 是由希望从服务器流式传输而不是直接解码的对象实现的接口。
type ResourceStreamer interface {
	// InputStream should return an io.ReadCloser if the provided object supports streaming. The desired
	// api version and an accept header (may be empty) are passed to the call. If no error occurs,
	// the caller may return a flag indicating whether the result should be flushed as writes occur
	// and a content type string that indicates the type of the stream.
	// If a null stream is returned, a StatusNoContent response wil be generated.
	// InputStream 应返回一个 io.ReadCloser，如果提供的对象支持流式传输。 传递所需的 api 版本和接受标头（可能为空）。
	// 如果没有发生错误，则调用者可以返回一个标志，指示是否应在写入发生时刷新结果以及指示流类型的内容类型字符串。
	// 如果返回空流，则将生成 StatusNoContent 响应。
	InputStream(ctx context.Context, apiVersion, acceptHeader string) (stream io.ReadCloser, flush bool, mimeType string, err error)
}

// StorageMetadata is an optional interface that callers can implement to provide additional
// information about their Storage objects.
// StorageMetadata 是调用者可以实现的可选接口，用于提供有关其 Storage 对象的其他信息。
type StorageMetadata interface {
	// ProducesMIMETypes returns a list of the MIME types the specified HTTP verb (GET, POST, DELETE,
	// PATCH) can respond with.
	// ProducesMIMETypes 返回指定 HTTP 动词（GET、POST、DELETE、PATCH）可以响应的 MIME 类型列表。
	ProducesMIMETypes(verb string) []string

	// ProducesObject returns an object the specified HTTP verb respond with. It will overwrite storage object if
	// it is not nil. Only the type of the return object matters, the value will be ignored.
	// ProducesObject 返回指定 HTTP 动词响应的对象。 如果不为 nil，它将覆盖存储对象。 仅返回对象的类型有关，值将被忽略。
	ProducesObject(verb string) interface{}
}

// StorageVersionProvider is an optional interface that a storage object can
// implement if it wishes to disclose its storage version.
// StorageVersionProvider 是一个可选接口，如果存储对象希望公开其存储版本，则可以实现它。
type StorageVersionProvider interface {
	// StorageVersion returns a group versioner, which will outputs the gvk
	// an object will be converted to before persisted in etcd, given a
	// list of kinds the object might belong to.
	// StorageVersion 返回一个 group versioner，它将输出对象在存储到 etcd 之前将转换为的 gvk，
	// 给定对象可能属于的一组类型。
	StorageVersion() runtime.GroupVersioner
}

// ResetFieldsStrategy is an optional interface that a storage object can
// implement if it wishes to provide the fields reset by its strategies.
// ResetFieldsStrategy 是一个可选接口，如果存储对象希望提供其策略重置的字段，则可以实现它。
type ResetFieldsStrategy interface {
	GetResetFields() map[fieldpath.APIVersion]*fieldpath.Set
}

// CreateUpdateResetFieldsStrategy is a union of RESTCreateUpdateStrategy
// and ResetFieldsStrategy.
// CreateUpdateResetFieldsStrategy 是 RESTCreateUpdateStrategy 和 ResetFieldsStrategy 的联合。
type CreateUpdateResetFieldsStrategy interface {
	RESTCreateUpdateStrategy
	ResetFieldsStrategy
}

// UpdateResetFieldsStrategy is a union of RESTUpdateStrategy
// and ResetFieldsStrategy.
// UpdateResetFieldsStrategy 是 RESTUpdateStrategy 和 ResetFieldsStrategy 的联合。
type UpdateResetFieldsStrategy interface {
	RESTUpdateStrategy
	ResetFieldsStrategy
}
