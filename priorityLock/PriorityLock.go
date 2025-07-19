package prioritylock

import (
	"sync"
	"sync/atomic"
)

type PriorityLock interface {
	LowPriorityUnlock()
	LowPriorityLock()
	HighPriorityLock()
	HighPriorityUnlock()
}

type PriorityPreferenceLock struct {
	dataMutex         sync.Mutex
	lowPriorityMutex  sync.Mutex
	cond              *sync.Cond
	highPriorityCount int32
}

var _ PriorityLock = (*PriorityPreferenceLock)(nil)

func NewPriorityPreferenceLock() *PriorityPreferenceLock {
	lock := PriorityPreferenceLock{
		highPriorityCount: 0,
	}
	lock.cond = sync.NewCond(&lock.dataMutex)
	return &lock
}

func (pl *PriorityPreferenceLock) HighPriorityLock() {
	atomic.AddInt32(&pl.highPriorityCount, 1)
	pl.dataMutex.Lock()
}

func (pl *PriorityPreferenceLock) HighPriorityUnlock() {
	atomic.AddInt32(&pl.highPriorityCount, -1)

	if pl.highPriorityCount <= 0 {
		pl.cond.Broadcast()
	}
	pl.dataMutex.Unlock()

}

func (pl *PriorityPreferenceLock) LowPriorityLock() {
	pl.lowPriorityMutex.Lock()
	pl.dataMutex.Lock()
	for pl.highPriorityCount > 0 {
		pl.cond.Wait()
	}

}

func (pl *PriorityPreferenceLock) LowPriorityUnlock() {
	pl.dataMutex.Unlock()

	pl.lowPriorityMutex.Unlock()

}
