/*
Copyright 2018 The Kubernetes Authors.

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

package truncate

import (
	"fmt"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	auditinternal "k8s.io/apiserver/pkg/apis/audit"
	"k8s.io/apiserver/pkg/audit"
)

const (
	// PluginName is the name reported in error metrics.
	// PluginName 是报告错误指标的名称。
	PluginName = "truncate"

	// annotationKey defines the name of the annotation used to indicate truncation.
	// annotationKey 定义用于指示截断的注释的名称。
	annotationKey = "audit.k8s.io/truncated"
	// annotationValue defines the value of the annotation used to indicate truncation.
	// annotationValue 定义用于指示截断的注释的值。
	annotationValue = "true"
)

// Config represents truncating backend configuration.
// Config 表示截断后端配置。
type Config struct {
	// MaxEventSize defines max allowed size of the event. If the event is larger,
	// truncating will be performed.
	// MaxEventSize 定义事件的最大允许大小。如果事件更大，则将执行截断。
	MaxEventSize int64

	// MaxBatchSize defined max allowed size of the batch of events, passed to the backend.
	// If the total size of the batch is larger than this number, batch will be split. Actual
	// size of the serialized request might be slightly higher, on the order of hundreds of bytes.
	// MaxBatchSize 定义传递给后端的事件批处理的最大允许大小。如果批处理的总大小大于此数字，则批处理将被拆分。序列化请求的实际大小可能略高，约为数百字节。
	MaxBatchSize int64
}

type backend struct {
	// The delegate backend that actually exports events.
	// 实际导出事件的委托后端。
	delegateBackend audit.Backend

	// Configuration used for truncation.
	c Config

	// Encoder used to calculate audit event sizes.
	e runtime.Encoder
}

var _ audit.Backend = &backend{}

// NewBackend returns a new truncating backend, using configuration passed in the parameters.
// Truncate backend automatically runs and shut downs the delegate backend.
// NewBackend 返回一个新的截断后端，使用参数中传递的配置。截断后端会自动运行和关闭委托后端。
func NewBackend(delegateBackend audit.Backend, config Config, groupVersion schema.GroupVersion) audit.Backend {
	return &backend{
		delegateBackend: delegateBackend,
		c:               config,
		e:               audit.Codecs.LegacyCodec(groupVersion),
	}
}

func (b *backend) ProcessEvents(events ...*auditinternal.Event) bool {
	var errors []error
	var impacted []*auditinternal.Event
	var batch []*auditinternal.Event
	var batchSize int64
	success := true
	for _, event := range events {
		size, err := b.calcSize(event)
		// If event was correctly serialized, but the size is more than allowed
		// and it makes sense to do trimming, i.e. there's a request and/or
		// response present, try to strip away request and response.
		if err == nil && size > b.c.MaxEventSize && event.Level.GreaterOrEqual(auditinternal.LevelRequest) {
			event = truncate(event)
			size, err = b.calcSize(event)
		}
		if err != nil {
			errors = append(errors, err)
			impacted = append(impacted, event)
			continue
		}
		if size > b.c.MaxEventSize {
			errors = append(errors, fmt.Errorf("event is too large even after truncating"))
			impacted = append(impacted, event)
			continue
		}

		if len(batch) > 0 && batchSize+size > b.c.MaxBatchSize {
			success = b.delegateBackend.ProcessEvents(batch...) && success
			batch = []*auditinternal.Event{}
			batchSize = 0
		}

		batchSize += size
		batch = append(batch, event)
	}

	if len(batch) > 0 {
		success = b.delegateBackend.ProcessEvents(batch...) && success
	}

	if len(impacted) > 0 {
		audit.HandlePluginError(PluginName, utilerrors.NewAggregate(errors), impacted...)
	}
	return success
}

// truncate removed request and response objects from the audit events,
// to try and keep at least metadata.
// truncate 从审计事件中删除请求和响应对象，以尝试保留至少元数据。
func truncate(e *auditinternal.Event) *auditinternal.Event {
	// Make a shallow copy to avoid copying response/request objects.
	newEvent := &auditinternal.Event{}
	*newEvent = *e

	newEvent.RequestObject = nil
	newEvent.ResponseObject = nil

	if newEvent.Annotations == nil {
		newEvent.Annotations = make(map[string]string)
	}
	newEvent.Annotations[annotationKey] = annotationValue

	return newEvent
}

func (b *backend) Run(stopCh <-chan struct{}) error {
	return b.delegateBackend.Run(stopCh)
}

func (b *backend) Shutdown() {
	b.delegateBackend.Shutdown()
}

func (b *backend) calcSize(e *auditinternal.Event) (int64, error) {
	s := &sizer{}
	if err := b.e.Encode(e, s); err != nil {
		return 0, err
	}
	return s.Size, nil
}

func (b *backend) String() string {
	return fmt.Sprintf("%s<%s>", PluginName, b.delegateBackend)
}

type sizer struct {
	Size int64
}

func (s *sizer) Write(p []byte) (n int, err error) {
	s.Size += int64(len(p))
	return len(p), nil
}
