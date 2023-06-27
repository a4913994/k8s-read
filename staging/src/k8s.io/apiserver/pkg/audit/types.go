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

package audit

import (
	auditinternal "k8s.io/apiserver/pkg/apis/audit"
)

type Sink interface {
	// ProcessEvents handles events. Per audit ID it might be that ProcessEvents is called up to three times.
	// Errors might be logged by the sink itself. If an error should be fatal, leading to an internal
	// error, ProcessEvents is supposed to panic. The event must not be mutated and is reused by the caller
	// after the call returns, i.e. the sink has to make a deepcopy to keep a copy around if necessary.
	// Returns true on success, may return false on error.
	// ProcessEvents 处理事件。每个审计 ID 可能最多调用 ProcessEvents 三次。接收器本身可能会记录错误。如果错误应该是致命的，
	// 导致内部错误，ProcessEvents 应该会恐慌。该事件不得发生变化，并在调用返回后由调用者重用，即接收器必须进行深层复制以在必要时保留副本。
	// 成功返回 true，错误可能返回 false。
	ProcessEvents(events ...*auditinternal.Event) bool
}

type Backend interface {
	Sink

	// Run will initialize the backend. It must not block, but may run go routines in the background. If
	// stopCh is closed, it is supposed to stop them. Run will be called before the first call to ProcessEvents.
	// Run 将初始化后端。它不得阻塞，但可以在后台运行 go 程序。如果 stopCh 被关闭，它应该停止它们。在第一次调用 ProcessEvents 之前将调用 Run。
	Run(stopCh <-chan struct{}) error

	// Shutdown will synchronously shut down the backend while making sure that all pending
	// events are delivered. It can be assumed that this method is called after
	// the stopCh channel passed to the Run method has been closed.
	// Shutdown 将同步关闭后端，同时确保所有待处理事件都被传递。可以假定在调用此方法之前已将传递给 Run 方法的 stopCh 通道关闭。
	Shutdown()

	// Returns the backend PluginName.
	// 返回后端插件名称。
	String() string
}
