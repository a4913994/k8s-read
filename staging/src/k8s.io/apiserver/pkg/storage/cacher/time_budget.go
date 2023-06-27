/*
Copyright 2016 The Kubernetes Authors.

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

package cacher

import (
	"sync"
	"time"

	"k8s.io/utils/clock"
)

const (
	refreshPerSecond = 50 * time.Millisecond
	maxBudget        = 100 * time.Millisecond
)

// timeBudget implements a budget of time that you can use and is
// periodically being refreshed. The pattern to use it is:
//
//	budget := newTimeBudget(...)
//	...
//	timeout := budget.takeAvailable()
//	// Now you can spend at most timeout on doing stuff
//	...
//	// If you didn't use all timeout, return what you didn't use
//	budget.returnUnused(<unused part of timeout>)
//
// NOTE: It's not recommended to be used concurrently from multiple threads -
// if first user takes the whole timeout, the second one will get 0 timeout
// even though the first one may return something later.
// timeBudget 实现了一个你可以使用的时间预算，并且会定期刷新。
// 使用它的模式是：
// budget := newTimeBudget(...)
// ...
// timeout := budget.takeAvailable()
// // 现在你可以在做事情上花费最多 timeout
// ...
// // 如果你没有使用所有的 timeout，返回你没有使用的部分
// budget.returnUnused(<unused part of timeout>)
// 注意：不建议从多个线程并发使用 - 如果第一个用户使用了整个超时，第二个用户将获得 0 超时，即使第一个用户稍后可能会返回一些内容。
type timeBudget interface {
	takeAvailable() time.Duration
	returnUnused(unused time.Duration)
}

type timeBudgetImpl struct {
	sync.Mutex
	clock     clock.Clock
	budget    time.Duration
	maxBudget time.Duration
	refresh   time.Duration
	// last store last access time
	last time.Time
}

func newTimeBudget() timeBudget {
	result := &timeBudgetImpl{
		clock:     clock.RealClock{},
		budget:    time.Duration(0),
		refresh:   refreshPerSecond,
		maxBudget: maxBudget,
	}
	result.last = result.clock.Now()
	return result
}

func (t *timeBudgetImpl) takeAvailable() time.Duration {
	t.Lock()
	defer t.Unlock()
	// budget accumulated since last access
	now := t.clock.Now()
	acc := now.Sub(t.last).Seconds() * t.refresh.Seconds()
	if acc < 0 {
		acc = 0
	}
	// update current budget and store the current time
	if t.budget = t.budget + time.Duration(acc*1e9); t.budget > t.maxBudget {
		t.budget = t.maxBudget
	}
	t.last = now
	result := t.budget
	t.budget = time.Duration(0)
	return result
}

func (t *timeBudgetImpl) returnUnused(unused time.Duration) {
	t.Lock()
	defer t.Unlock()
	if unused < 0 {
		// We used more than allowed.
		return
	}
	// add the unused time directly to the budget
	// takeAvailable() will take into account the elapsed time
	if t.budget = t.budget + unused; t.budget > t.maxBudget {
		t.budget = t.maxBudget
	}
}
