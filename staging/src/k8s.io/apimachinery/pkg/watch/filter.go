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

package watch

import (
	"sync"
)

// FilterFunc should take an event, possibly modify it in some way, and return
// the modified event. If the event should be ignored, then return keep=false.
// FilterFunc应该接收一个事件，可能会以某种方式修改它，并返回修改后的事件。如果应该忽略事件，则返回keep=false。
type FilterFunc func(in Event) (out Event, keep bool)

// Filter passes all events through f before allowing them to pass on.
// Putting a filter on a watch, as an unavoidable side-effect due to the way
// go channels work, effectively causes the watch's event channel to have its
// queue length increased by one.
//
// Filter传递所有事件通过f，然后允许它们通过。将过滤器放在监视上，由于go通道的工作方式，不可避免地导致监视的事件通道的队列长度增加了一个。
// WARNING: filter has a fatal flaw, in that it can't properly update the
// Type field (Add/Modified/Deleted) to reflect items beginning to pass the
// filter when they previously didn't.
// WARNING: filter有一个致命的缺陷，即它无法正确更新Type字段（Add/Modified/Deleted）以反映项目开始通过过滤器时，它们以前没有通过过滤器。
func Filter(w Interface, f FilterFunc) Interface {
	fw := &filteredWatch{
		incoming: w,
		result:   make(chan Event),
		f:        f,
	}
	go fw.loop()
	return fw
}

type filteredWatch struct {
	incoming Interface
	result   chan Event
	f        FilterFunc
}

// ResultChan returns a channel which will receive filtered events.
// ResultChan 返回一个通道，该通道将接收过滤后的事件。
func (fw *filteredWatch) ResultChan() <-chan Event {
	return fw.result
}

// Stop stops the upstream watch, which will eventually stop this watch.
// Stop 停止上游监视，这将最终停止此监视。
func (fw *filteredWatch) Stop() {
	fw.incoming.Stop()
}

// loop waits for new values, filters them, and resends them.
// loop 等待新值，过滤它们，并重新发送它们。
func (fw *filteredWatch) loop() {
	defer close(fw.result)
	for event := range fw.incoming.ResultChan() {
		filtered, keep := fw.f(event)
		if keep {
			fw.result <- filtered
		}
	}
}

// Recorder records all events that are sent from the watch until it is closed.
// Recorder 记录从监视发送的所有事件，直到它关闭为止。
type Recorder struct {
	Interface

	lock   sync.Mutex
	events []Event
}

var _ Interface = &Recorder{}

// NewRecorder wraps an Interface and records any changes sent across it.
func NewRecorder(w Interface) *Recorder {
	r := &Recorder{}
	r.Interface = Filter(w, r.record)
	return r
}

// record is a FilterFunc and tracks each received event.
// record 是一个FilterFunc，并跟踪每个接收到的事件。
func (r *Recorder) record(in Event) (Event, bool) {
	r.lock.Lock()
	defer r.lock.Unlock()
	r.events = append(r.events, in)
	return in, true
}

// Events returns a copy of the events sent across this recorder.
// Events 返回通过此记录器发送的事件的副本。
func (r *Recorder) Events() []Event {
	r.lock.Lock()
	defer r.lock.Unlock()
	copied := make([]Event, len(r.events))
	copy(copied, r.events)
	return copied
}
