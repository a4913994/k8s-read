/*
Copyright 2021 The Kubernetes Authors.

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

package server

import (
	"sync"
)

/*
We make an attempt here to identify the events that take place during
lifecycle of the apiserver.
我们在这里尝试认识在 apiserver 的生命周期中发生的事件。

We also identify each event with a name so we can refer to it.
我们还用名称标识每个事件，以便我们可以引用它。

Events:
- ShutdownInitiated: KILL signal received
- ShutdownInitiated：收到 KILL 信号

- AfterShutdownDelayDuration: shutdown delay duration has passed // 关闭延迟时间
- InFlightRequestsDrained: all in flight request(s) have been drained // 所有处理中的请求都已完成
- HasBeenReady is signaled when the readyz endpoint succeeds for the first time // 当 readyz 端点第一次成功时发出信号

The following is a sequence of shutdown events that we expect to see with
  'ShutdownSendRetryAfter' = false:

T0: ShutdownInitiated: KILL signal received
	- /readyz starts returning red
    - run pre shutdown hooks

T0+70s: AfterShutdownDelayDuration: shutdown delay duration has passed
	- the default value of 'ShutdownDelayDuration' is '70s'
	- it's time to initiate shutdown of the HTTP Server, server.Shutdown is invoked
	- as a consequene, the Close function has is called for all listeners
 	- the HTTP Server stops listening immediately
	- any new request arriving on a new TCP socket is denied with
      a network error similar to 'connection refused'
    - the HTTP Server waits gracefully for existing requests to complete
      up to '60s' (dictated by ShutdownTimeout)
	- active long running requests will receive a GOAWAY.

T0+70s: HTTPServerStoppedListening:
	- this event is signaled when the HTTP Server has stopped listening
      which is immediately after server.Shutdown has been invoked
	- 在 server.Shutdown 被调用后 HTTP 服务器停止侦听时会发出此事件信号

T0 + 70s + up-to 60s: InFlightRequestsDrained: existing in flight requests have been drained
	- long running requests are outside of this scope
	- up-to 60s: the default value of 'ShutdownTimeout' is 60s, this means that
      any request in flight has a hard timeout of 60s.
	- it's time to call 'Shutdown' on the audit events since all
	  in flight request(s) have drained.


The following is a sequence of shutdown events that we expect to see with
  'ShutdownSendRetryAfter' = true:

T0: ShutdownInitiated: KILL signal received
	- /readyz starts returning red
    - run pre shutdown hooks

T0+70s: AfterShutdownDelayDuration: shutdown delay duration has passed
	- the default value of 'ShutdownDelayDuration' is '70s'
	- the HTTP Server will continue to listen
	- the apiserver is not accepting new request(s)
		- it includes new request(s) on a new or an existing TCP connection
		- new request(s) arriving after this point are replied with a 429
      	  and the  response headers: 'Retry-After: 1` and 'Connection: close'
	- note: these new request(s) will not show up in audit logs

T0 + 70s + up to 60s: InFlightRequestsDrained: existing in flight requests have been drained
	- long running requests are outside of this scope
	- up to 60s: the default value of 'ShutdownTimeout' is 60s, this means that
      any request in flight has a hard timeout of 60s.
	- server.Shutdown is called, the HTTP Server stops listening immediately
    - the HTTP Server waits gracefully for existing requests to complete
      up to '2s' (it's hard coded right now)
*/

// lifecycleSignal encapsulates a named apiserver event
// lifecycleSignal 封装了一个命名的 apiserver 事件
type lifecycleSignal interface {
	// Signal signals the event, indicating that the event has occurred.
	// Signal is idempotent, once signaled the event stays signaled and
	// it immediately unblocks any goroutine waiting for this event.
	// Signal 向事件发出信号，表明事件已经发生。信号是幂等的，一旦发出信号，事件就会保持信号状态，它会立即解除阻塞等待该事件的任何 goroutine
	Signal()

	// Signaled returns a channel that is closed when the underlying event
	// has been signaled. Successive calls to Signaled return the same value.
	// Signaled 返回一个通道，该通道在发出基础事件信号时关闭。对 Signaled 的连续调用返回相同的值。
	Signaled() <-chan struct{}

	// Name returns the name of the signal, useful for logging.
	// Name 返回信号的名称，对记录很有用。
	Name() string
}

// lifecycleSignals provides an abstraction of the events that
// transpire during the lifecycle of the apiserver. This abstraction makes it easy
// for us to write unit tests that can verify expected graceful termination behavior.
// lifecycleSignals 提供了在 apiserver 的生命周期中发生的事件的抽象。这种抽象使我们很容易编写单元测试来验证预期的优雅终止行为。
//
// GenericAPIServer can use these to either:
//   - signal that a particular termination event has transpired
//   - wait for a designated termination event to transpire and do some action.
//
// GenericAPIServer 可以使用这些来：
//   - 发出特定终止事件已经发生的信号
//   - 等待指定的终止事件发生并执行一些操作。
type lifecycleSignals struct {
	// ShutdownInitiated event is signaled when an apiserver shutdown has been initiated.
	// It is signaled when the `stopCh` provided by the main goroutine
	// receives a KILL signal and is closed as a consequence.
	// ShutdownInitiated 事件在启动 apiserver 关闭时发出信号。当主 goroutine 提供的 stopCh 接收到 KILL 信号并因此关闭时，它会发出信号。
	ShutdownInitiated lifecycleSignal

	// AfterShutdownDelayDuration event is signaled as soon as ShutdownDelayDuration
	// has elapsed since the ShutdownInitiated event.
	// ShutdownDelayDuration allows the apiserver to delay shutdown for some time.
	// 自 ShutdownInitiated 事件以来，一旦 ShutdownDelayDuration 过去，AfterShutdownDelayDuration 事件就会发出信号。 ShutdownDelayDuration 允许 apiserver 延迟关闭一段时间。
	AfterShutdownDelayDuration lifecycleSignal

	// PreShutdownHooksStopped event is signaled when all registered
	// preshutdown hook(s) have finished running.
	// PreShutdownHooksStopped 事件在所有已注册的预关闭挂钩完成运行时发出信号。
	PreShutdownHooksStopped lifecycleSignal

	// NotAcceptingNewRequest event is signaled when the server is no
	// longer accepting any new request, from this point on any new
	// request will receive an error.
	// NotAcceptingNewRequest 事件在服务器不再接受任何新请求时发出信号，从此时开始任何新请求都将收到错误。
	NotAcceptingNewRequest lifecycleSignal

	// InFlightRequestsDrained event is signaled when the existing requests
	// in flight have completed. This is used as signal to shut down the audit backends
	// InFlightRequestsDrained 事件在现有的正在处理的请求完成时发出信号。这用作关闭审计后端的信号
	InFlightRequestsDrained lifecycleSignal

	// HTTPServerStoppedListening termination event is signaled when the
	// HTTP Server has stopped listening to the underlying socket.
	// HTTPServerStoppedListening 终止事件在 HTTP 服务器停止侦听底层套接字时发出信号。
	HTTPServerStoppedListening lifecycleSignal

	// HasBeenReady is signaled when the readyz endpoint succeeds for the first time.
	// HasBeenReady 在 readyz 端点第一次成功时发出信号。
	HasBeenReady lifecycleSignal

	// MuxAndDiscoveryComplete is signaled when all known HTTP paths have been installed.
	// It exists primarily to avoid returning a 404 response when a resource actually exists but we haven't installed the path to a handler.
	// The actual logic is implemented by an APIServer using the generic server library.
	// MuxAndDiscoveryComplete 在安装了所有已知的 HTTP 路径后发出信号。它的存在主要是为了避免在资源实际存在但我们尚未安装处理程序路径时返回 404 响应。实际逻辑由 APIServer 使用通用服务器库实现。
	MuxAndDiscoveryComplete lifecycleSignal
}

// newLifecycleSignals returns an instance of lifecycleSignals interface to be used
// to coordinate lifecycle of the apiserver
func newLifecycleSignals() lifecycleSignals {
	return lifecycleSignals{
		ShutdownInitiated:          newNamedChannelWrapper("ShutdownInitiated"),
		AfterShutdownDelayDuration: newNamedChannelWrapper("AfterShutdownDelayDuration"),
		PreShutdownHooksStopped:    newNamedChannelWrapper("PreShutdownHooksStopped"),
		NotAcceptingNewRequest:     newNamedChannelWrapper("NotAcceptingNewRequest"),
		InFlightRequestsDrained:    newNamedChannelWrapper("InFlightRequestsDrained"),
		HTTPServerStoppedListening: newNamedChannelWrapper("HTTPServerStoppedListening"),
		HasBeenReady:               newNamedChannelWrapper("HasBeenReady"),
		MuxAndDiscoveryComplete:    newNamedChannelWrapper("MuxAndDiscoveryComplete"),
	}
}

func newNamedChannelWrapper(name string) lifecycleSignal {
	return &namedChannelWrapper{
		name: name,
		once: sync.Once{},
		ch:   make(chan struct{}),
	}
}

type namedChannelWrapper struct {
	name string
	once sync.Once
	ch   chan struct{}
}

func (e *namedChannelWrapper) Signal() {
	e.once.Do(func() {
		close(e.ch)
	})
}

func (e *namedChannelWrapper) Signaled() <-chan struct{} {
	return e.ch
}

func (e *namedChannelWrapper) Name() string {
	return e.name
}
