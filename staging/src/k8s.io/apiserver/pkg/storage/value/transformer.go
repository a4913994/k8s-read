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

// Package value contains methods for assisting with transformation of values in storage.
package value

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/util/errors"
)

func init() {
	RegisterMetrics()
}

// Context is additional information that a storage transformation may need to verify the data at rest.
// Context 是存储转换可能需要验证静态数据的附加信息。
type Context interface {
	// AuthenticatedData should return an array of bytes that describes the current value. If the value changes,
	// the transformer may report the value as unreadable or tampered. This may be nil if no such description exists
	// or is needed. For additional verification, set this to data that strongly identifies the value, such as
	// the key and creation version of the stored data.
	// AuthenticatedData 应返回描述当前值的字节数组。如果值发生变化，转换器可能会报告该值不可读或被篡改。
	// 如果不存在或不需要这样的描述，这可能是零。对于额外的验证，将此设置为强烈标识值的数据，例如存储数据的密钥和创建版本。
	AuthenticatedData() []byte
}

// Transformer allows a value to be transformed before being read from or written to the underlying store. The methods
// must be able to undo the transformation caused by the other.
// Transformer 允许在从底层存储读取或写入底层存储之前转换值。这些方法必须能够撤消由其他方法引起的转换。
type Transformer interface {
	// TransformFromStorage may transform the provided data from its underlying storage representation or return an error.
	// Stale is true if the object on disk is stale and a write to etcd should be issued, even if the contents of the object
	// have not changed.
	// TransformFromStorage 可能会从其底层存储表示转换提供的数据或返回错误。
	// 如果磁盘上的对象陈旧并且应该向 etcd 发出写入，则 Stale 为真，即使对象的内容没有更改。
	TransformFromStorage(ctx context.Context, data []byte, dataCtx Context) (out []byte, stale bool, err error)
	// TransformToStorage may transform the provided data into the appropriate form in storage or return an error.
	// TransformToStorage 可能会将提供的数据转换为存储中的适当形式或返回错误。
	TransformToStorage(ctx context.Context, data []byte, dataCtx Context) (out []byte, err error)
}

// DefaultContext is a simple implementation of Context for a slice of bytes.
// DefaultContext 是 Context 的一个字节切片的简单实现。
type DefaultContext []byte

// AuthenticatedData returns itself.
func (c DefaultContext) AuthenticatedData() []byte { return c }

// PrefixTransformer holds a transformer interface and the prefix that the transformation is located under.
// PrefixTransformer 持有转换器接口和转换所在的前缀。
type PrefixTransformer struct {
	Prefix      []byte
	Transformer Transformer
}

type prefixTransformers struct {
	transformers []PrefixTransformer
	err          error
}

var _ Transformer = &prefixTransformers{}

// NewPrefixTransformers supports the Transformer interface by checking the incoming data against the provided
// prefixes in order. The first matching prefix will be used to transform the value (the prefix is stripped
// before the Transformer interface is invoked). The first provided transformer will be used when writing to
// the store.
// NewPrefixTransformers 通过按顺序检查传入数据与提供的前缀来支持 Transformer 接口。
// 第一个匹配的前缀将用于转换值（前缀在调用 Transformer 接口之前被剥离）。写入存储时将使用第一个提供的转换器
func NewPrefixTransformers(err error, transformers ...PrefixTransformer) Transformer {
	if err == nil {
		err = fmt.Errorf("the provided value does not match any of the supported transformers")
	}
	return &prefixTransformers{
		transformers: transformers,
		err:          err,
	}
}

// TransformFromStorage finds the first transformer with a prefix matching the provided data and returns
// the result of transforming the value. It will always mark any transformation as stale that is not using
// the first transformer.
// TransformFromStorage 找到前缀与提供的数据匹配的第一个转换器，并返回转换值的结果。它将始终将任何未使用第一个转换器的转换标记为陈旧
func (t *prefixTransformers) TransformFromStorage(ctx context.Context, data []byte, dataCtx Context) ([]byte, bool, error) {
	start := time.Now()
	var errs []error
	for i, transformer := range t.transformers {
		if bytes.HasPrefix(data, transformer.Prefix) {
			result, stale, err := transformer.Transformer.TransformFromStorage(ctx, data[len(transformer.Prefix):], dataCtx)
			// To migrate away from encryption, user can specify an identity transformer higher up
			// (in the config file) than the encryption transformer. In that scenario, the identity transformer needs to
			// identify (during reads from disk) whether the data being read is encrypted or not. If the data is encrypted,
			// it shall throw an error, but that error should not prevent the next subsequent transformer from being tried.
			if len(transformer.Prefix) == 0 && err != nil {
				continue
			}
			if len(transformer.Prefix) == 0 {
				RecordTransformation("from_storage", "identity", start, err)
			} else {
				RecordTransformation("from_storage", string(transformer.Prefix), start, err)
			}

			// It is valid to have overlapping prefixes when the same encryption provider
			// is specified multiple times but with different keys (the first provider is
			// being rotated to and some later provider is being rotated away from).
			//
			// Example:
			//
			//  {
			//    "aescbc": {
			//      "keys": [
			//        {
			//          "name": "2",
			//          "secret": "some key 2"
			//        }
			//      ]
			//    }
			//  },
			//  {
			//    "aescbc": {
			//      "keys": [
			//        {
			//          "name": "1",
			//          "secret": "some key 1"
			//        }
			//      ]
			//    }
			//  },
			//
			// The transformers for both aescbc configs share the prefix k8s:enc:aescbc:v1:
			// but a failure in the first one should not prevent a later match from being attempted.
			// Thus we never short-circuit on a prefix match that results in an error.
			if err != nil {
				errs = append(errs, err)
				continue
			}

			return result, stale || i != 0, err
		}
	}
	if err := errors.Reduce(errors.NewAggregate(errs)); err != nil {
		return nil, false, err
	}
	RecordTransformation("from_storage", "unknown", start, t.err)
	return nil, false, t.err
}

// TransformToStorage uses the first transformer and adds its prefix to the data.
// TransformToStorage 使用第一个转换器并将其前缀添加到数据中。
func (t *prefixTransformers) TransformToStorage(ctx context.Context, data []byte, dataCtx Context) ([]byte, error) {
	start := time.Now()
	transformer := t.transformers[0]
	result, err := transformer.Transformer.TransformToStorage(ctx, data, dataCtx)
	RecordTransformation("to_storage", string(transformer.Prefix), start, err)
	if err != nil {
		return nil, err
	}
	prefixedData := make([]byte, len(transformer.Prefix), len(result)+len(transformer.Prefix))
	copy(prefixedData, transformer.Prefix)
	prefixedData = append(prefixedData, result...)
	return prefixedData, nil
}
