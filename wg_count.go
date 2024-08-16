// This code wraps a waitgroup to expose a "count," which can be used to get an (approximate, and immediately stale) idea of how many goroutines using this WaitGroup are in flight
// From https://stackoverflow.com/questions/68995144/how-to-get-the-number-of-goroutines-associated-with-a-waitgroup
package main

import (
	"sync"
	"sync/atomic"
)

type WaitGroupCount struct {
	sync.WaitGroup
	count int64
}

func (wg *WaitGroupCount) Add(delta int) {
	atomic.AddInt64(&wg.count, int64(delta))
	wg.WaitGroup.Add(delta)
}

func (wg *WaitGroupCount) Done() {
	atomic.AddInt64(&wg.count, -1)
	wg.WaitGroup.Done()
}

func (wg *WaitGroupCount) GetCount() int {
	return int(atomic.LoadInt64(&wg.count))
}
