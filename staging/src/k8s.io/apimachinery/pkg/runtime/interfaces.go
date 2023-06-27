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

package runtime

import (
	"io"
	"net/url"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	// APIVersionInternal may be used if you are registering a type that should not
	// be considered stable or serialized - it is a convention only and has no
	// special behavior in this package.
	// 如果您正在注册一个不应该被认为是稳定或序列化的类型，则可以使用APIVersionInternal—它只是一个约定，在这个包中没有特殊的行为。
	APIVersionInternal = "__internal"
)

// GroupVersioner refines a set of possible conversion targets into a single option.
// GroupVersioner将一组可能的转换目标细化为单个选项。
type GroupVersioner interface {
	// KindForGroupVersionKinds returns a desired target group version kind for the given input, or returns ok false if no
	// target is known. In general, if the return target is not in the input list, the caller is expected to invoke
	// Scheme.New(target) and then perform a conversion between the current Go type and the destination Go type.
	// Sophisticated implementations may use additional information about the input kinds to pick a destination kind.
	// kindforgroupversiontypes返回给定输入所需的目标组版本类型，如果不知道目标，则返回ok false。通常，如果返回目标不在输入列表中，则调用方将调用Scheme.New(目标)，然后在当前Go类型和目标Go类型之间执行转换。
	KindForGroupVersionKinds(kinds []schema.GroupVersionKind) (target schema.GroupVersionKind, ok bool)
	// Identifier returns string representation of the object.
	// Identifiers of two different encoders should be equal only if for every input
	// kinds they return the same result.
	// 标识符返回对象的字符串表示形式。两个不同编码器的标识符只有在每个输入类型返回相同结果时才应该相等。
	Identifier() string
}

// Identifier represents an identifier.
// Identitier of two different objects should be equal if and only if for every
// input the output they produce is exactly the same.
// 两个不同对象的标识符应该相等当且仅当它们产生的每个输入的输出都完全相同。
type Identifier string

// Encoder writes objects to a serialized form
type Encoder interface {
	// Encode writes an object to a stream. Implementations may return errors if the versions are
	// incompatible, or if no conversion is defined.
	Encode(obj Object, w io.Writer) error
	// Identifier returns an identifier of the encoder.
	// Identifiers of two different encoders should be equal if and only if for every input
	// object it will be encoded to the same representation by both of them.
	//
	// Identifier is intended for use with CacheableObject#CacheEncode method. In order to
	// correctly handle CacheableObject, Encode() method should look similar to below, where
	// doEncode() is the encoding logic of implemented encoder:
	//   func (e *MyEncoder) Encode(obj Object, w io.Writer) error {
	//     if co, ok := obj.(CacheableObject); ok {
	//       return co.CacheEncode(e.Identifier(), e.doEncode, w)
	//     }
	//     return e.doEncode(obj, w)
	//   }
	Identifier() Identifier
}

// MemoryAllocator is responsible for allocating memory.
// By encapsulating memory allocation into its own interface, we can reuse the memory
// across many operations in places we know it can significantly improve the performance.
// 通过将内存分配封装到它自己的接口中，我们可以在许多操作中重用内存，我们知道这样做可以显著提高性能。
type MemoryAllocator interface {
	// Allocate reserves memory for n bytes.
	// Note that implementations of this method are not required to zero the returned array.
	// It is the caller's responsibility to clean the memory if needed.
	Allocate(n uint64) []byte
}

// EncoderWithAllocator  serializes objects in a way that allows callers to manage any additional memory allocations.
// EncoderWithAllocator以一种允许调用者管理任何额外内存分配的方式序列化对象。
type EncoderWithAllocator interface {
	Encoder
	// EncodeWithAllocator writes an object to a stream as Encode does.
	// In addition, it allows for providing a memory allocator for efficient memory usage during object serialization
	// EncodeWithAllocator像Encode一样将对象写入流。此外，它还允许提供内存分配器，以便在对象序列化期间有效地使用内存
	EncodeWithAllocator(obj Object, w io.Writer, memAlloc MemoryAllocator) error
}

// Decoder attempts to load an object from data.
// Decoder 尝试从数据中加载对象。
type Decoder interface {
	// Decode attempts to deserialize the provided data using either the innate typing of the scheme or the
	// default kind, group, and version provided. It returns a decoded object as well as the kind, group, and
	// version from the serialized data, or an error. If into is non-nil, it will be used as the target type
	// and implementations may choose to use it rather than reallocating an object. However, the object is not
	// guaranteed to be populated. The returned object is not guaranteed to match into. If defaults are
	// provided, they are applied to the data by default. If no defaults or partial defaults are provided, the
	// type of the into may be used to guide conversion decisions.
	// Decode尝试使用方案的固有类型或提供的默认类型、组和版本来反序列化所提供的数据。它返回一个已解码的对象以及序列化数据的类型、组和版本，或者返回一个错误。如果into非nil，它将被用作目标类型，实现可以选择使用它而不是重新分配对象。但是，不保证填充对象。返回的对象不能保证匹配到。如果提供了默认值，则默认将它们应用于数据。如果没有提供默认值或部分默认值，则可以使用into的类型来指导转换决策。
	Decode(data []byte, defaults *schema.GroupVersionKind, into Object) (Object, *schema.GroupVersionKind, error)
}

// Serializer is the core interface for transforming objects into a serialized format and back.
// Implementations may choose to perform conversion of the object, but no assumptions should be made.
// Serializer是将对象转换为序列化格式和反向转换的核心接口。实现可以选择执行对象的转换，但不应做任何假设。
type Serializer interface {
	Encoder
	Decoder
}

// Codec is a Serializer that deals with the details of versioning objects. It offers the same
// interface as Serializer, so this is a marker to consumers that care about the version of the objects
// they receive.
// Codec是处理对象版本控制细节的序列化器。它提供了与Serializer相同的接口，因此对于关心所接收对象的版本的消费者来说，这是一个标记。
type Codec Serializer

// ParameterCodec defines methods for serializing and deserializing API objects to url.Values and
// performing any necessary conversion. Unlike the normal Codec, query parameters are not self describing
// and the desired version must be specified.
// ParameterCodec定义了将API对象序列化和反序列化为url的方法。值并执行任何必要的转换。与普通的编解码器不同，查询参数不是自描述的，必须指定所需的版本。
type ParameterCodec interface {
	// DecodeParameters takes the given url.Values in the specified group version and decodes them
	// into the provided object, or returns an error.
	DecodeParameters(parameters url.Values, from schema.GroupVersion, into Object) error
	// EncodeParameters encodes the provided object as query parameters or returns an error.
	EncodeParameters(obj Object, to schema.GroupVersion) (url.Values, error)
}

// Framer is a factory for creating readers and writers that obey a particular framing pattern.
// Framer是一个工厂，用于创建遵循特定框架模式的读取器和写入器。
type Framer interface {
	NewFrameReader(r io.ReadCloser) io.ReadCloser
	NewFrameWriter(w io.Writer) io.Writer
}

// SerializerInfo contains information about a specific serialization format
// SerializerInfo 包含有关特定序列化格式的信息
type SerializerInfo struct {
	// MediaType is the value that represents this serializer over the wire.
	// MediaType 是通过线路表示此序列化程序的值。
	MediaType string
	// MediaTypeType is the first part of the MediaType ("application" in "application/json").
	// MediaTypeType 是 MediaType 的第一部分（“applicationjson”中的“application”）。
	MediaTypeType string
	// MediaTypeSubType is the second part of the MediaType ("json" in "application/json").
	// MediaTypeSubType 是 MediaType 的第二部分（“applicationjson”中的“json”）。
	MediaTypeSubType string
	// EncodesAsText indicates this serializer can be encoded to UTF-8 safely.
	// EncodesAsText 指示此序列化程序可以安全地编码为 UTF-8。
	EncodesAsText bool
	// Serializer is the individual object serializer for this media type.
	// Serializer 是此媒体类型的单独对象序列化程序。
	Serializer Serializer
	// PrettySerializer, if set, can serialize this object in a form biased towards
	// readability.
	// PrettySerializer，如果设置，可以以偏向于可读性的形式序列化此对象。
	PrettySerializer Serializer
	// StrictSerializer, if set, deserializes this object strictly,
	// erring on unknown fields.
	// StrictSerializer，如果设置，严格反序列化此对象，在未知字段上出错。
	StrictSerializer Serializer
	// StreamSerializer, if set, describes the streaming serialization format
	// for this media type.
	// StreamSerializer（如果已设置）描述此媒体类型的流式序列化格式。
	StreamSerializer *StreamSerializerInfo
}

// StreamSerializerInfo contains information about a specific stream serialization format
type StreamSerializerInfo struct {
	// EncodesAsText indicates this serializer can be encoded to UTF-8 safely.
	EncodesAsText bool
	// Serializer is the top level object serializer for this type when streaming
	Serializer
	// Framer is the factory for retrieving streams that separate objects on the wire
	Framer
}

// NegotiatedSerializer is an interface used for obtaining encoders, decoders, and serializers
// for multiple supported media types. This would commonly be accepted by a server component
// that performs HTTP content negotiation to accept multiple formats.
// NegotiatedSerializer是一个接口，用于获取多种支持媒体类型的编码器、解码器和序列化器。这通常会被执行HTTP内容协商以接受多种格式的服务器组件所接受
type NegotiatedSerializer interface {
	// SupportedMediaTypes is the media types supported for reading and writing single objects.
	SupportedMediaTypes() []SerializerInfo

	// EncoderForVersion returns an encoder that ensures objects being written to the provided
	// serializer are in the provided group version.
	EncoderForVersion(serializer Encoder, gv GroupVersioner) Encoder
	// DecoderToVersion returns a decoder that ensures objects being read by the provided
	// serializer are in the provided group version by default.
	DecoderToVersion(serializer Decoder, gv GroupVersioner) Decoder
}

// ClientNegotiator handles turning an HTTP content type into the appropriate encoder.
// Use NewClientNegotiator or NewVersionedClientNegotiator to create this interface from
// a NegotiatedSerializer.
// ClientNegotiator负责将HTTP内容类型转换为适当的编码器。使用NewClientNegotiator或NewVersionedClientNegotiator从NegotiatedSerializer创建此接口。
type ClientNegotiator interface {
	// Encoder returns the appropriate encoder for the provided contentType (e.g. application/json)
	// and any optional mediaType parameters (e.g. pretty=1), or an error. If no serializer is found
	// a NegotiateError will be returned. The current client implementations consider params to be
	// optional modifiers to the contentType and will ignore unrecognized parameters.
	Encoder(contentType string, params map[string]string) (Encoder, error)
	// Decoder returns the appropriate decoder for the provided contentType (e.g. application/json)
	// and any optional mediaType parameters (e.g. pretty=1), or an error. If no serializer is found
	// a NegotiateError will be returned. The current client implementations consider params to be
	// optional modifiers to the contentType and will ignore unrecognized parameters.
	Decoder(contentType string, params map[string]string) (Decoder, error)
	// StreamDecoder returns the appropriate stream decoder for the provided contentType (e.g.
	// application/json) and any optional mediaType parameters (e.g. pretty=1), or an error. If no
	// serializer is found a NegotiateError will be returned. The Serializer and Framer will always
	// be returned if a Decoder is returned. The current client implementations consider params to be
	// optional modifiers to the contentType and will ignore unrecognized parameters.
	StreamDecoder(contentType string, params map[string]string) (Decoder, Serializer, Framer, error)
}

// StorageSerializer is an interface used for obtaining encoders, decoders, and serializers
// that can read and write data at rest. This would commonly be used by client tools that must
// read files, or server side storage interfaces that persist restful objects.
// StorageSerializer是一个接口，用于获取可静态读写数据的编码器、解码器和序列化器。这通常由必须读取文件的客户端工具或持久化restful对象的服务器端存储接口使用。
type StorageSerializer interface {
	// SupportedMediaTypes are the media types supported for reading and writing objects.
	SupportedMediaTypes() []SerializerInfo

	// UniversalDeserializer returns a Serializer that can read objects in multiple supported formats
	// by introspecting the data at rest.
	UniversalDeserializer() Decoder

	// EncoderForVersion returns an encoder that ensures objects being written to the provided
	// serializer are in the provided group version.
	EncoderForVersion(serializer Encoder, gv GroupVersioner) Encoder
	// DecoderForVersion returns a decoder that ensures objects being read by the provided
	// serializer are in the provided group version by default.
	DecoderToVersion(serializer Decoder, gv GroupVersioner) Decoder
}

// NestedObjectEncoder is an optional interface that objects may implement to be given
// an opportunity to encode any nested Objects / RawExtensions during serialization.
type NestedObjectEncoder interface {
	EncodeNestedObjects(e Encoder) error
}

// NestedObjectDecoder is an optional interface that objects may implement to be given
// an opportunity to decode any nested Objects / RawExtensions during serialization.
// It is possible for DecodeNestedObjects to return a non-nil error but for the decoding
// to have succeeded in the case of strict decoding errors (e.g. unknown/duplicate fields).
// As such it is important for callers of DecodeNestedObjects to check to confirm whether
// an error is a runtime.StrictDecodingError before short circuiting.
// Similarly, implementations of DecodeNestedObjects should ensure that a runtime.StrictDecodingError
// is only returned when the rest of decoding has succeeded.
type NestedObjectDecoder interface {
	DecodeNestedObjects(d Decoder) error
}

///////////////////////////////////////////////////////////////////////////////
// Non-codec interfaces

type ObjectDefaulter interface {
	// Default takes an object (must be a pointer) and applies any default values.
	// Defaulters may not error.
	// Default 接受一个对象（必须是一个指针）并应用任何默认值。违约者不得错误。
	Default(in Object)
}

type ObjectVersioner interface {
	ConvertToVersion(in Object, gv GroupVersioner) (out Object, err error)
}

// ObjectConvertor converts an object to a different version.
// ObjectConvertor将一个对象转换为不同的版本。
type ObjectConvertor interface {
	// Convert attempts to convert one object into another, or returns an error. This
	// method does not mutate the in object, but the in and out object might share data structures,
	// i.e. the out object cannot be mutated without mutating the in object as well.
	// The context argument will be passed to all nested conversions.
	// Convert尝试将一个对象转换为另一个对象，或返回错误。这个方法不会改变in对象，但是in和out对象可能共享数据结构，也就是说，如果不改变in对象，就不能改变out对象。context参数将被传递给所有嵌套转换。
	Convert(in, out, context interface{}) error
	// ConvertToVersion takes the provided object and converts it the provided version. This
	// method does not mutate the in object, but the in and out object might share data structures,
	// i.e. the out object cannot be mutated without mutating the in object as well.
	// This method is similar to Convert() but handles specific details of choosing the correct
	// output version.
	// ConvertToVersion接受所提供的对象并将其转换为所提供的版本。这个方法不会改变in对象，但是in和out对象可能共享数据结构，也就是说，如果不改变in对象，就不能改变out对象。此方法类似于Convert()，但处理选择正确输出版本的具体细节。
	ConvertToVersion(in Object, gv GroupVersioner) (out Object, err error)
	ConvertFieldLabel(gvk schema.GroupVersionKind, label, value string) (string, string, error)
}

// ObjectTyper contains methods for extracting the APIVersion and Kind
// of objects.
// objecttype包含提取APIVersion和Kind对象的方法。
type ObjectTyper interface {
	// ObjectKinds returns the all possible group,version,kind of the provided object, true if
	// the object is unversioned, or an error if the object is not recognized
	// (IsNotRegisteredError will return true).
	// objecttypes返回所提供对象的所有可能的组、版本和类型，如果对象未被版本控制，则返回true;如果对象无法识别，则返回错误(IsNotRegisteredError将返回true)。
	ObjectKinds(Object) ([]schema.GroupVersionKind, bool, error)
	// Recognizes returns true if the scheme is able to handle the provided version and kind,
	// or more precisely that the provided version is a possible conversion or decoding
	// target.
	// 如果方案能够处理所提供的版本和种类，或者更准确地说，所提供的版本是可能的转换或解码目标，则recognized返回true。
	Recognizes(gvk schema.GroupVersionKind) bool
}

// ObjectCreater contains methods for instantiating an object by kind and version.
// objectcreator包含按类型和版本实例化对象的方法。
type ObjectCreater interface {
	New(kind schema.GroupVersionKind) (out Object, err error)
}

// EquivalentResourceMapper provides information about resources that address the same underlying data as a specified resource
// EquivalentResourceMapper 提供有关处理与指定资源相同的基础数据的资源的信息
type EquivalentResourceMapper interface {
	// EquivalentResourcesFor returns a list of resources that address the same underlying data as resource.
	// If subresource is specified, only equivalent resources which also have the same subresource are included.
	// The specified resource can be included in the returned list.
	// EquivalentResourcesFor 返回一个资源列表，这些资源处理与资源相同的基础数据。
	// 如果指定了 subresource，则仅包含也具有相同子资源的等效资源。指定的资源可以包含在返回的列表中。
	EquivalentResourcesFor(resource schema.GroupVersionResource, subresource string) []schema.GroupVersionResource
	// KindFor returns the kind expected by the specified resource[/subresource].
	// A zero value is returned if the kind is unknown.
	// KindFor 返回指定资源 [子资源] 所期望的种类。如果种类未知，则返回零值。
	KindFor(resource schema.GroupVersionResource, subresource string) schema.GroupVersionKind
}

// EquivalentResourceRegistry provides an EquivalentResourceMapper interface,
// and allows registering known resource[/subresource] -> kind
// EquivalentResourceRegistry 提供了一个 EquivalentResourceMapper 接口，并允许注册已知资源[subresource] -> kind
type EquivalentResourceRegistry interface {
	EquivalentResourceMapper
	// RegisterKindFor registers the existence of the specified resource[/subresource] along with its expected kind.
	// RegisterKindFor 注册指定资源 [subresource] 的存在及其预期种类。
	RegisterKindFor(resource schema.GroupVersionResource, subresource string, kind schema.GroupVersionKind)
}

// ResourceVersioner provides methods for setting and retrieving
// the resource version from an API object.
// ResourceVersioner 提供了从 API 对象设置和检索资源版本的方法。
type ResourceVersioner interface {
	SetResourceVersion(obj Object, version string) error
	ResourceVersion(obj Object) (string, error)
}

// Namer provides methods for retrieving name and namespace of an API object.
// Namer 提供用于检索 API 对象的名称和命名空间的方法。
type Namer interface {
	// Name returns the name of a given object.
	Name(obj Object) (string, error)
	// Namespace returns the name of a given object.
	Namespace(obj Object) (string, error)
}

// Object interface must be supported by all API types registered with Scheme. Since objects in a scheme are
// expected to be serialized to the wire, the interface an Object must provide to the Scheme allows
// serializers to set the kind, version, and group the object is represented as. An Object may choose
// to return a no-op ObjectKindAccessor in cases where it is not expected to be serialized.
// 对象接口必须为注册到Scheme的所有API类型所支持。由于方案中的对象被序列化到线路，因此对象必须提供给方案的接口允许序列化器设置对象表示的类型、版本和组。在不希望被序列化的情况下，对象可以选择返回一个无操作ObjectKindAccessor。
type Object interface {
	GetObjectKind() schema.ObjectKind
	DeepCopyObject() Object
}

// CacheableObject allows an object to cache its different serializations
// to avoid performing the same serialization multiple times.
// CacheableObject允许一个对象缓存其不同的序列化，以避免多次执行相同的序列化。
type CacheableObject interface {
	// CacheEncode writes an object to a stream. The <encode> function will
	// be used in case of cache miss. The <encode> function takes ownership
	// of the object.
	// If CacheableObject is a wrapper, then deep-copy of the wrapped object
	// should be passed to <encode> function.
	// CacheEncode assumes that for two different calls with the same <id>,
	// <encode> function will also be the same.
	CacheEncode(id Identifier, encode func(Object, io.Writer) error, w io.Writer) error
	// GetObject returns a deep-copy of an object to be encoded - the caller of
	// GetObject() is the owner of returned object. The reason for making a copy
	// is to avoid bugs, where caller modifies the object and forgets to copy it,
	// thus modifying the object for everyone.
	// The object returned by GetObject should be the same as the one that is supposed
	// to be passed to <encode> function in CacheEncode method.
	// If CacheableObject is a wrapper, the copy of wrapped object should be returned.
	GetObject() Object
}

// Unstructured objects store values as map[string]interface{}, with only values that can be serialized
// to JSON allowed.
// 非结构化对象将值存储为map[string]interface{}，只允许将值序列化为JSON。
type Unstructured interface {
	Object
	// NewEmptyInstance returns a new instance of the concrete type containing only kind/apiVersion and no other data.
	// This should be called instead of reflect.New() for unstructured types because the go type alone does not preserve kind/apiVersion info.
	NewEmptyInstance() Unstructured
	// UnstructuredContent returns a non-nil map with this object's contents. Values may be
	// []interface{}, map[string]interface{}, or any primitive type. Contents are typically serialized to
	// and from JSON. SetUnstructuredContent should be used to mutate the contents.
	UnstructuredContent() map[string]interface{}
	// SetUnstructuredContent updates the object content to match the provided map.
	SetUnstructuredContent(map[string]interface{})
	// IsList returns true if this type is a list or matches the list convention - has an array called "items".
	IsList() bool
	// EachListItem should pass a single item out of the list as an Object to the provided function. Any
	// error should terminate the iteration. If IsList() returns false, this method should return an error
	// instead of calling the provided function.
	EachListItem(func(Object) error) error
}
