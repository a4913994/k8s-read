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
	"fmt"
	"sync"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// FullChannelBehavior controls how the Broadcaster reacts if a watcher's watch
// channel is full.
// FullChannelBehavior 描述了 Broadcaster 如何处理观察者的通道已满的情况。
type FullChannelBehavior int

const (
	WaitIfChannelFull FullChannelBehavior = iota
	DropIfChannelFull
)

// Buffer the incoming queue a little bit even though it should rarely ever accumulate
// anything, just in case a few events are received in such a short window that
// Broadcaster can't move them onto the watchers' queues fast enough.
// 即使 incomingQueueLength 应该很少积累任何内容，也要缓冲 incoming 队列，以防在这样短的时间窗口内收到几个事件，
// 并且 Broadcaster 无法足够快地将它们移动到观察者的队列上。
const incomingQueueLength = 25

// Broadcaster distributes event notifications among any number of watchers. Every event
// is delivered to every watcher.
// Broadcaster 描述了如何将事件通知分发给任意数量的观察者。每个事件都会被分发给每个观察者。
type Broadcaster struct {
	watchers     map[int64]*broadcasterWatcher
	nextWatcher  int64
	distributing sync.WaitGroup

	// incomingBlock allows us to ensure we don't race and end up sending events
	// to a closed channel following a broadcaster shutdown.
	// incomingBlock 用于确保我们不会发生竞争，并且在广播器关闭后不会向关闭的通道发送事件。
	incomingBlock sync.Mutex
	incoming      chan Event
	stopped       chan struct{}

	// How large to make watcher's channel.
	watchQueueLength int
	// If one of the watch channels is full, don't wait for it to become empty.
	// Instead just deliver it to the watchers that do have space in their
	// channels and move on to the next event.
	// It's more fair to do this on a per-watcher basis than to do it on the
	// "incoming" channel, which would allow one slow watcher to prevent all
	// other watchers from getting new events.
	// 如果一个观察者的通道已满，不要等待它变为空。而是只将事件分发给有空间的观察者，并继续处理下一个事件。
	// 在每个观察者的基础上这样做比在“incoming”通道上这样做更公平，这允许一个缓慢的观察者阻止所有其他观察者获取新事件。
	fullChannelBehavior FullChannelBehavior
}

// NewBroadcaster creates a new Broadcaster. queueLength is the maximum number of events to queue per watcher.
// It is guaranteed that events will be distributed in the order in which they occur,
// but the order in which a single event is distributed among all of the watchers is unspecified.
// NewBroadcaster 创建一个新的 Broadcaster。queueLength 是每个观察者的最大事件数。保证事件将按发生的顺序分发，
// 但是在所有观察者之间分发单个事件的顺序是未指定的。
func NewBroadcaster(queueLength int, fullChannelBehavior FullChannelBehavior) *Broadcaster {
	m := &Broadcaster{
		watchers:            map[int64]*broadcasterWatcher{},
		incoming:            make(chan Event, incomingQueueLength),
		stopped:             make(chan struct{}),
		watchQueueLength:    queueLength,
		fullChannelBehavior: fullChannelBehavior,
	}
	m.distributing.Add(1)
	go m.loop()
	return m
}

// NewLongQueueBroadcaster functions nearly identically to NewBroadcaster,
// except that the incoming queue is the same size as the outgoing queues
// (specified by queueLength).
func NewLongQueueBroadcaster(queueLength int, fullChannelBehavior FullChannelBehavior) *Broadcaster {
	m := &Broadcaster{
		watchers:            map[int64]*broadcasterWatcher{},
		incoming:            make(chan Event, queueLength),
		stopped:             make(chan struct{}),
		watchQueueLength:    queueLength,
		fullChannelBehavior: fullChannelBehavior,
	}
	m.distributing.Add(1)
	go m.loop()
	return m
}

const internalRunFunctionMarker = "internal-do-function"

// a function type we can shoehorn into the queue.
// 一个我们可以强行塞进队列的函数类型。
type functionFakeRuntimeObject func()

func (obj functionFakeRuntimeObject) GetObjectKind() schema.ObjectKind {
	return schema.EmptyObjectKind
}
func (obj functionFakeRuntimeObject) DeepCopyObject() runtime.Object {
	if obj == nil {
		return nil
	}
	// funcs are immutable. Hence, just return the original func.
	// funcs 是不可变的。因此，只需返回原始 func。
	return obj
}

// Execute f, blocking the incoming queue (and waiting for it to drain first).
// The purpose of this terrible hack is so that watchers added after an event
// won't ever see that event, and will always see any event after they are
// added.
// 执行 f，阻塞 incoming 队列（并等待它先被清空）。这个糟糕的 hack 的目的是，以便在事件之后添加的观察者永远不会看到该事件，
// 并且将始终看到在添加之后的任何事件。
func (m *Broadcaster) blockQueue(f func()) {
	m.incomingBlock.Lock()
	defer m.incomingBlock.Unlock()
	select {
	case <-m.stopped:
		return
	default:
	}
	var wg sync.WaitGroup
	wg.Add(1)
	m.incoming <- Event{
		Type: internalRunFunctionMarker,
		Object: functionFakeRuntimeObject(func() {
			defer wg.Done()
			f()
		}),
	}
	wg.Wait()
}

// Watch adds a new watcher to the list and returns an Interface for it.
// Note: new watchers will only receive new events. They won't get an entire history
// of previous events. It will block until the watcher is actually added to the
// broadcaster.
// Watch 添加一个新的观察者到列表中，并返回一个 Interface。注意：新的观察者只会接收新的事件。他们不会得到之前事件的整个历史记录。
func (m *Broadcaster) Watch() (Interface, error) {
	var w *broadcasterWatcher
	m.blockQueue(func() {
		id := m.nextWatcher
		m.nextWatcher++
		w = &broadcasterWatcher{
			result:  make(chan Event, m.watchQueueLength),
			stopped: make(chan struct{}),
			id:      id,
			m:       m,
		}
		m.watchers[id] = w
	})
	if w == nil {
		return nil, fmt.Errorf("broadcaster already stopped")
	}
	return w, nil
}

// WatchWithPrefix adds a new watcher to the list and returns an Interface for it. It sends
// queuedEvents down the new watch before beginning to send ordinary events from Broadcaster.
// The returned watch will have a queue length that is at least large enough to accommodate
// all of the items in queuedEvents. It will block until the watcher is actually added to
// the broadcaster.
// WatchWithPrefix 添加一个新的观察者到列表中，并返回一个 Interface。它在开始从 Broadcaster 发送普通事件之前，
// 将 queuedEvents 通过新的 watch 发送。返回的 watch 将具有一个足够大的队列长度，以容纳 queuedEvents 中的所有项目。
// 它将阻塞，直到观察者实际上被添加到广播器。
func (m *Broadcaster) WatchWithPrefix(queuedEvents []Event) (Interface, error) {
	var w *broadcasterWatcher
	m.blockQueue(func() {
		id := m.nextWatcher
		m.nextWatcher++
		length := m.watchQueueLength
		if n := len(queuedEvents) + 1; n > length {
			length = n
		}
		w = &broadcasterWatcher{
			result:  make(chan Event, length),
			stopped: make(chan struct{}),
			id:      id,
			m:       m,
		}
		m.watchers[id] = w
		for _, e := range queuedEvents {
			w.result <- e
		}
	})
	if w == nil {
		return nil, fmt.Errorf("broadcaster already stopped")
	}
	return w, nil
}

// stopWatching stops the given watcher and removes it from the list.
// stopWatching 停止给定的观察者并将其从列表中删除。
func (m *Broadcaster) stopWatching(id int64) {
	m.blockQueue(func() {
		w, ok := m.watchers[id]
		if !ok {
			// No need to do anything, it's already been removed from the list.
			return
		}
		delete(m.watchers, id)
		close(w.result)
	})
}

// closeAll disconnects all watchers (presumably in response to a Shutdown call).
// closeAll 断开所有观察者的连接（可能是对 Shutdown 调用的响应）。
func (m *Broadcaster) closeAll() {
	for _, w := range m.watchers {
		close(w.result)
	}
	// Delete everything from the map, since presence/absence in the map is used
	// by stopWatching to avoid double-closing the channel.
	// 从映射中删除所有内容，因为映射中的存在/不存在由 stopWatching 用于避免双重关闭通道。
	m.watchers = map[int64]*broadcasterWatcher{}
}

// Action distributes the given event among all watchers.
// Action 将给定的事件分发给所有观察者。
func (m *Broadcaster) Action(action EventType, obj runtime.Object) error {
	m.incomingBlock.Lock()
	defer m.incomingBlock.Unlock()
	select {
	case <-m.stopped:
		return fmt.Errorf("broadcaster already stopped")
	default:
	}

	m.incoming <- Event{action, obj}
	return nil
}

// ActionOrDrop distributes the given event among all watchers, or drops it on the floor
// if too many incoming actions are queued up.  Returns true if the action was sent,
// false if dropped.
// ActionOrDrop 将给定的事件分发给所有观察者，或者丢弃它，如果有太多的传入动作排队。如果动作被发送，则返回 true，如果丢弃，则返回 false。
func (m *Broadcaster) ActionOrDrop(action EventType, obj runtime.Object) (bool, error) {
	m.incomingBlock.Lock()
	defer m.incomingBlock.Unlock()

	// Ensure that if the broadcaster is stopped we do not send events to it.
	select {
	case <-m.stopped:
		return false, fmt.Errorf("broadcaster already stopped")
	default:
	}

	select {
	case m.incoming <- Event{action, obj}:
		return true, nil
	default:
		return false, nil
	}
}

// Shutdown disconnects all watchers (but any queued events will still be distributed).
// You must not call Action or Watch* after calling Shutdown. This call blocks
// until all events have been distributed through the outbound channels. Note
// that since they can be buffered, this means that the watchers might not
// have received the data yet as it can remain sitting in the buffered
// channel. It will block until the broadcaster stop request is actually executed
// Shutdown 断开所有观察者的连接（但是任何排队的事件仍将被分发）。在调用 Shutdown 之后，您不得调用 Action 或 Watch*。
// 此调用将阻塞，直到所有事件都通过出站通道分发。请注意，由于它们可以缓冲，这意味着观察者可能尚未接收到数据，因为它可以保留在缓冲的通道中。
// 它将阻塞，直到广播器停止请求实际执行
func (m *Broadcaster) Shutdown() {
	m.blockQueue(func() {
		close(m.stopped)
		close(m.incoming)
	})
	m.distributing.Wait()
}

// loop receives from m.incoming and distributes to all watchers.
// loop接收m.incoming并分发给所有的watchers
func (m *Broadcaster) loop() {
	// Deliberately not catching crashes here. Yes, bring down the process if there's a
	// bug in watch.Broadcaster.
	for event := range m.incoming {
		if event.Type == internalRunFunctionMarker {
			event.Object.(functionFakeRuntimeObject)()
			continue
		}
		m.distribute(event)
	}
	m.closeAll()
	m.distributing.Done()
}

// distribute sends event to all watchers. Blocking.
// distribute将事件发送给所有的观察者。阻塞
func (m *Broadcaster) distribute(event Event) {
	if m.fullChannelBehavior == DropIfChannelFull {
		for _, w := range m.watchers {
			select {
			case w.result <- event:
			case <-w.stopped:
			default: // Don't block if the event can't be queued.
			}
		}
	} else {
		for _, w := range m.watchers {
			select {
			case w.result <- event:
			case <-w.stopped:
			}
		}
	}
}

// broadcasterWatcher handles a single watcher of a broadcaster
// broadcasterWatcher 处理广播器的单个观察者
type broadcasterWatcher struct {
	result  chan Event
	stopped chan struct{}
	stop    sync.Once
	id      int64
	m       *Broadcaster
}

// ResultChan returns a channel to use for waiting on events.
// ResultChan 返回用于等待事件的通道。
func (mw *broadcasterWatcher) ResultChan() <-chan Event {
	return mw.result
}

// Stop stops watching and removes mw from its list.
// It will block until the watcher stop request is actually executed
// Stop 停止观看并从其列表中删除 mw。它将阻塞，直到观察者停止请求实际执行
func (mw *broadcasterWatcher) Stop() {
	mw.stop.Do(func() {
		close(mw.stopped)
		mw.m.stopWatching(mw.id)
	})
}
