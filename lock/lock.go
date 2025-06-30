package lock

import (
	"sync"
	maplock "testgg/mapLock"
	prioritylock "testgg/priorityLock"
)

type Lock struct {
	locks *maplock.MapLock[prioritylock.PriorityLock]
}

var (
	globalLock *Lock
	o          sync.Once
)

func newLock() *Lock {
	l := &Lock{
		locks: maplock.NewMapLock[prioritylock.PriorityLock](),
	}
	return l
}

func GetGlobalLock() *Lock {
	o.Do(func() {
		if globalLock == nil {
			globalLock = newLock()
		}
	})
	return globalLock
}

// ================================
// ================================
// ================================

func (l *Lock) Get(key string) prioritylock.PriorityLock {
	priorityLock := prioritylock.NewPriorityPreferenceLock()
	pl := l.locks.LoadOrSet(key, priorityLock)
	return *pl
}
