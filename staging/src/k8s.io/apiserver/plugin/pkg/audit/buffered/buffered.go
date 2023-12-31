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

package buffered

import (
	"fmt"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	auditinternal "k8s.io/apiserver/pkg/apis/audit"
	"k8s.io/apiserver/pkg/audit"
	"k8s.io/client-go/util/flowcontrol"
)

// PluginName is the name reported in error metrics.
const PluginName = "buffered"

// BatchConfig represents batching delegate audit backend configuration.
// BatchConfig 表示批处理委托审计后端配置。
type BatchConfig struct {
	// BufferSize defines a size of the buffering queue.
	// BufferSize 定义缓冲队列的大小。
	BufferSize int
	// MaxBatchSize defines maximum size of a batch.
	// MaxBatchSize 定义批处理的最大大小。
	MaxBatchSize int
	// MaxBatchWait indicates the maximum interval between two batches.
	// MaxBatchWait 指示两个批处理之间的最大间隔。
	MaxBatchWait time.Duration

	// ThrottleEnable defines whether throttling will be applied to the batching process.
	// ThrottleEnable 定义是否将限流应用于批处理过程。
	ThrottleEnable bool
	// ThrottleQPS defines the allowed rate of batches per second sent to the delegate backend.
	// ThrottleQPS 定义发送到委托后端的每秒允许的批处理速率。
	ThrottleQPS float32
	// ThrottleBurst defines the maximum number of requests sent to the delegate backend at the same moment in case
	// the capacity defined by ThrottleQPS was not utilized.
	// ThrottleBurst 定义在未使用由 ThrottleQPS 定义的容量的情况下，同时发送到委托后端的最大请求数。
	ThrottleBurst int

	// Whether the delegate backend should be called asynchronously.
	// 是否应异步调用委托后端.
	AsyncDelegate bool
}

type bufferedBackend struct {
	// The delegate backend that actually exports events.
	// 实际导出事件的委托后端。
	delegateBackend audit.Backend

	// Channel to buffer events before sending to the delegate backend.
	// 用于在发送到委托后端之前缓冲事件的通道。
	buffer chan *auditinternal.Event
	// Maximum number of events in a batch sent to the delegate backend.
	// 发送到委托后端的批处理中的最大事件数。
	maxBatchSize int
	// Amount of time to wait after sending a batch to the delegate backend before sending another one.
	// 发送批处理到委托后端后等待的时间，然后再发送另一个批处理。
	//
	// Receiving maxBatchSize events will always trigger sending a batch, regardless of the amount of time passed.
	// 接收 maxBatchSize 事件将始终触发发送批处理，而不管经过的时间量。
	maxBatchWait time.Duration

	// Whether the delegate backend should be called asynchronously.
	// 是否应异步调用委托后端。
	asyncDelegate bool

	// Channel to signal that the batching routine has processed all remaining events and exited.
	// Once `shutdownCh` is closed no new events will be sent to the delegate backend.
	// 用于表示批处理例程已处理所有剩余事件并退出的通道。
	shutdownCh chan struct{}

	// WaitGroup to control the concurrency of sending batches to the delegate backend.
	// Worker routine calls Add before sending a batch and
	// then spawns a routine that calls Done after batch was processed by the delegate backend.
	// This WaitGroup is used to wait for all sending routines to finish before shutting down audit backend.
	// 用于控制将批处理发送到委托后端的并发性的 WaitGroup。
	wg sync.WaitGroup

	// Limits the number of batches sent to the delegate backend per second.
	// 限制每秒发送到委托后端的批处理数。
	throttle flowcontrol.RateLimiter
}

var _ audit.Backend = &bufferedBackend{}

// NewBackend returns a buffered audit backend that wraps delegate backend.
// Buffered backend automatically runs and shuts down the delegate backend.
// NewBackend 返回一个缓冲的审计后端，该后端包装委托后端。
// 缓冲后端会自动运行和关闭委托后端。
func NewBackend(delegate audit.Backend, config BatchConfig) audit.Backend {
	var throttle flowcontrol.RateLimiter
	if config.ThrottleEnable {
		throttle = flowcontrol.NewTokenBucketRateLimiter(config.ThrottleQPS, config.ThrottleBurst)
	}
	return &bufferedBackend{
		delegateBackend: delegate,
		buffer:          make(chan *auditinternal.Event, config.BufferSize),
		maxBatchSize:    config.MaxBatchSize,
		maxBatchWait:    config.MaxBatchWait,
		asyncDelegate:   config.AsyncDelegate,
		shutdownCh:      make(chan struct{}),
		wg:              sync.WaitGroup{},
		throttle:        throttle,
	}
}

func (b *bufferedBackend) Run(stopCh <-chan struct{}) error {
	go func() {
		// Signal that the working routine has exited.
		defer close(b.shutdownCh)

		b.processIncomingEvents(stopCh)

		// Handle the events that were received after the last buffer
		// scraping and before this line. Since the buffer is closed, no new
		// events will come through.
		allEventsProcessed := false
		timer := make(chan time.Time)
		for !allEventsProcessed {
			allEventsProcessed = func() bool {
				// Recover from any panic in order to try to process all remaining events.
				// Note, that in case of a panic, the return value will be false and
				// the loop execution will continue.
				defer runtime.HandleCrash()

				events := b.collectEvents(timer, wait.NeverStop)
				b.processEvents(events)
				return len(events) == 0
			}()
		}
	}()
	return b.delegateBackend.Run(stopCh)
}

// Shutdown blocks until stopCh passed to the Run method is closed and all
// events added prior to that moment are batched and sent to the delegate backend.
func (b *bufferedBackend) Shutdown() {
	// Wait until the routine spawned in Run method exits.
	<-b.shutdownCh

	// Wait until all sending routines exit.
	//
	// - When b.shutdownCh is closed, we know that the goroutine in Run has terminated.
	// - This means that processIncomingEvents has terminated.
	// - Which means that b.buffer is closed and cannot accept any new events anymore.
	// - Because processEvents is called synchronously from the Run goroutine, the waitgroup has its final value.
	// Hence wg.Wait will not miss any more outgoing batches.
	b.wg.Wait()

	b.delegateBackend.Shutdown()
}

// processIncomingEvents runs a loop that collects events from the buffer. When
// b.stopCh is closed, processIncomingEvents stops and closes the buffer.
func (b *bufferedBackend) processIncomingEvents(stopCh <-chan struct{}) {
	defer close(b.buffer)

	var (
		maxWaitChan  <-chan time.Time
		maxWaitTimer *time.Timer
	)
	// Only use max wait batching if batching is enabled.
	if b.maxBatchSize > 1 {
		maxWaitTimer = time.NewTimer(b.maxBatchWait)
		maxWaitChan = maxWaitTimer.C
		defer maxWaitTimer.Stop()
	}

	for {
		func() {
			// Recover from any panics caused by this function so a panic in the
			// goroutine can't bring down the main routine.
			defer runtime.HandleCrash()

			if b.maxBatchSize > 1 {
				maxWaitTimer.Reset(b.maxBatchWait)
			}
			b.processEvents(b.collectEvents(maxWaitChan, stopCh))
		}()

		select {
		case <-stopCh:
			return
		default:
		}
	}
}

// collectEvents attempts to collect some number of events in a batch.
//
// The following things can cause collectEvents to stop and return the list
// of events:
//
//   - Maximum number of events for a batch.
//   - Timer has passed.
//   - Buffer channel is closed and empty.
//   - stopCh is closed.
func (b *bufferedBackend) collectEvents(timer <-chan time.Time, stopCh <-chan struct{}) []*auditinternal.Event {
	var events []*auditinternal.Event

L:
	for i := 0; i < b.maxBatchSize; i++ {
		select {
		case ev, ok := <-b.buffer:
			// Buffer channel was closed and no new events will follow.
			if !ok {
				break L
			}
			events = append(events, ev)
		case <-timer:
			// Timer has expired. Send currently accumulated batch.
			break L
		case <-stopCh:
			// Backend has been stopped. Send currently accumulated batch.
			break L
		}
	}

	return events
}

// processEvents process the batch events in a goroutine using delegateBackend's ProcessEvents.
func (b *bufferedBackend) processEvents(events []*auditinternal.Event) {
	if len(events) == 0 {
		return
	}

	// TODO(audit): Should control the number of active goroutines
	// if one goroutine takes 5 seconds to finish, the number of goroutines can be 5 * defaultBatchThrottleQPS
	if b.throttle != nil {
		b.throttle.Accept()
	}

	if b.asyncDelegate {
		b.wg.Add(1)
		go func() {
			defer b.wg.Done()
			defer runtime.HandleCrash()

			// Execute the real processing in a goroutine to keep it from blocking.
			// This lets the batching routine continue draining the queue immediately.
			b.delegateBackend.ProcessEvents(events...)
		}()
	} else {
		func() {
			defer runtime.HandleCrash()

			// Execute the real processing in a goroutine to keep it from blocking.
			// This lets the batching routine continue draining the queue immediately.
			b.delegateBackend.ProcessEvents(events...)
		}()
	}
}

func (b *bufferedBackend) ProcessEvents(ev ...*auditinternal.Event) bool {
	// The following mechanism is in place to support the situation when audit
	// events are still coming after the backend was stopped.
	var sendErr error
	var evIndex int

	// If the delegateBackend was shutdown and the buffer channel was closed, an
	// attempt to add an event to it will result in panic that we should
	// recover from.
	defer func() {
		if err := recover(); err != nil {
			sendErr = fmt.Errorf("audit backend shut down")
		}
		if sendErr != nil {
			audit.HandlePluginError(PluginName, sendErr, ev[evIndex:]...)
		}
	}()

	for i, e := range ev {
		evIndex = i
		// Per the audit.Backend interface these events are reused after being
		// sent to the Sink. Deep copy and send the copy to the queue.
		event := e.DeepCopy()

		select {
		case b.buffer <- event:
		default:
			sendErr = fmt.Errorf("audit buffer queue blocked")
			return true
		}
	}
	return true
}

func (b *bufferedBackend) String() string {
	return fmt.Sprintf("%s<%s>", PluginName, b.delegateBackend)
}
