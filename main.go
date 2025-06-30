// package main
//
// import (
// 	"fmt"
// 	"sync"
// )
//
// type mapLock struct {
// 	mux   sync.Mutex
// 	locks map[string]*sync.Mutex
// }
//
// func (m *mapLock) Lock(key string) {
//
// 	fmt.Printf("lock for %s", key)
// 	m.mux.Lock()
// 	lock, exist := m.locks[key]
// 	if !exist {
// 		fmt.Printf("lock for %s not exist", key)
// 		lock = &sync.Mutex{}
// 		m.locks[key] = lock
// 	}
// 	m.mux.Unlock()
// 	lock.Lock()
//
// }
//
// func (m *mapLock) Unlock(key string) {
// 	m.mux.Lock()
//
// 	fmt.Printf("unlock for %s", key)
// 	lock, exist := m.locks[key]
// 	if !exist {
// 		return
// 	}
//
// 	lock.Unlock()
// 	m.mux.Unlock()
// }
//
// var (
// 	locks *mapLock
// 	once  sync.Once
// )
//
// func newLockMap() *mapLock {
// 	return &mapLock{
// 		mux:   sync.Mutex{},
// 		locks: make(map[string]*sync.Mutex),
// 	}
// }
//
// func GetLocks() *mapLock {
// 	once.Do(func() {
// 		locks = newLockMap()
// 	})
// 	return locks
// }
//
// func main() {}
//
// // ========================================
// // ========================================
// // ========================================
//
// // type LockManager struct {
// // 	locks map[string]*sync.Mutex
// // 	mu    sync.RWMutex
// // }
// //
// // func NewLockManager() *LockManager {
// // 	return &LockManager{
// // 		locks: make(map[string]*sync.Mutex),
// // 	}
// // }
// //
// // func (lm *LockManager) Lock(identifier string) {
// // 	lm.mu.Lock()
// // 	if _, exists := lm.locks[identifier]; !exists {
// // 		lm.locks[identifier] = &sync.Mutex{}
// // 	}
// // 	mutex := lm.locks[identifier]
// // 	lm.mu.Unlock()
// //
// // 	mutex.Lock()
// // }
// //
// // func (lm *LockManager) Unlock(identifier string) {
// // 	lm.mu.RLock()
// // 	if mutex, exists := lm.locks[identifier]; exists {
// // 		mutex.Unlock()
// // 	}
// // 	lm.mu.RUnlock()
// // }
//
// type PriorityLock interface {
// 	Lock()
// 	Unlock()
// 	HighPriorityLock()
// 	HighPriorityUnlock()
// }
//
// // PriorityPreferenceLock implements a simple triple-mutex priority lock
// // patterns are like:
// //
// //	Low Priority would do: lock lowPriorityMutex, wait for high priority groups, lock nextToAccess, lock dataMutex, unlock nextToAccess, do stuff, unlock dataMutex, unlock lowPriorityMutex
// //	High Priority would do: increment high priority waiting, lock nextToAccess, lock dataMutex, unlock nextToAccess, do stuff, unlock dataMutex, decrement high priority waiting
// type PriorityPreferenceLock struct {
// 	dataMutex           sync.Mutex
// 	nextToAccess        sync.Mutex
// 	lowPriorityMutex    sync.Mutex
// 	highPriorityWaiting sync.WaitGroup
// }
//
// func NewPriorityPreferenceLock() *PriorityPreferenceLock {
// 	lock := PriorityPreferenceLock{
// 		highPriorityWaiting: sync.WaitGroup{},
// 	}
// 	return &lock
// }
//
// // Lock will acquire a low-priority lock
// // it must wait until both low priority and all high priority lock holders are released.
// func (lock *PriorityPreferenceLock) Lock() {
// 	lock.lowPriorityMutex.Lock()
// 	lock.highPriorityWaiting.Wait()
// 	lock.nextToAccess.Lock()
// 	lock.dataMutex.Lock()
// 	lock.nextToAccess.Unlock()
// }
//
// // Unlock will unlock the low-priority lock
// func (lock *PriorityPreferenceLock) Unlock() {
// 	lock.dataMutex.Unlock()
// 	lock.lowPriorityMutex.Unlock()
// }
//
// // HighPriorityLock will acquire a high-priority lock
// // it must still wait until a low-priority lock has been released and then potentially other high priority lock contenders.
// func (lock *PriorityPreferenceLock) HighPriorityLock() {
// 	lock.highPriorityWaiting.Add(1)
// 	lock.nextToAccess.Lock()
// 	lock.dataMutex.Lock()
// 	lock.nextToAccess.Unlock()
// }
//
// // HighPriorityUnlock will unlock the high-priority lock
// func (lock *PriorityPreferenceLock) HighPriorityUnlock() {
// 	lock.dataMutex.Unlock()
// 	lock.highPriorityWaiting.Done()
// }
//
// type LockManager struct {
// 	locks map[string]*PriorityPreferenceLock
// 	mu    sync.RWMutex
// }
//
// func NewLockManager() *LockManager {
// 	return &LockManager{
// 		locks: make(map[string]*PriorityPreferenceLock),
// 	}
// }
//
// func (lm *LockManager) set(identifier string) {
// 	lm.mu.Lock()
// 	lm.locks[identifier] = NewPriorityPreferenceLock()
// 	lm.mu.Unlock()
// }
//
// func (lm *LockManager) get(identifier string) *PriorityPreferenceLock {
// 	lm.mu.RLock()
// 	mutex, exits := lm.locks[identifier]
// 	if !exits {
// 		return nil
// 	}
// 	return mutex
// }
//
// func (lm *LockManager) loadOrGet(identifier int) *PriorityPreferenceLock {
// 	mutex := lm.ge(ideke int
// }
//
// var (
// 	locksM    *LockManager
// 	onceLocks sync.Once
// )
//
// func GetLocksM() *LockManager {
// 	onceLocks.Do(func() {
// 		if locksM == nil {
// 			locksM = NewLockManager()
// 		}
// 	})
// 	return locksM
// }
//
// func (l *LockManager) LowPriorityLockByID(id string) {
// 	l.locks.get(id)
// }
//
// // Usage
// // func main() {
// // 	lm := NewLockManager()
// //
// // 	var wg sync.WaitGroup
// //
// // 	for i := 0; i < 5; i++ {
// // 		wg.Add(1)
// // 		go func(id int) {
// // 			defer wg.Done()
// // 			identifier := fmt.Sprintf("resource-%d", id%2)
// //
// // 			lm.Lock(identifier)
// // 			fmt.Printf("Goroutine %d acquired lock for %s\n", id, identifier)
// // 			time.Sleep(time.Second)
// // 			fmt.Printf("Goroutine %d releasing lock for %s\n", id, identifier)
// // 			lm.Unlock(identifier)
// // 		}(i)
// // 	}
// //
// // 	wg.Wait()
// // }

package main

import (
	"fmt"
	"strings"
	"sync"
	"testgg/lock"
	"time"
)

func main() {

	locks := lock.GetGlobalLock()

	numIsins := 30
	numOrders := 5

	var isins []string

	for i := range numIsins {
		isin := fmt.Sprintf("MX%d", i)
		isins = append(isins, isin)
	}

	results := make(chan string, numIsins*numOrders*2)
	var wg sync.WaitGroup
	wg.Add(len(isins) * numOrders * 2)

	for _, isin := range isins {
		for range numOrders {
			go func() {
				lock := locks.Get(isin)
				lock.Lock()
				r := processOrder(isin, "create")
				results <- r
				lock.Unlock()
				wg.Done()

			}()
		}
	}

	for _, isin := range isins {
		for i := range numOrders {
			if i%2 != 0 {
				continue
			}
			go func() {
				lock := locks.Get(isin)
				lock.HighPriorityLock()
				r := processOrder(isin, "cancel")
				results <- r
				lock.HighPriorityUnlock()
				wg.Done()

			}()
		}

		for i := range numOrders {
			if i%2 == 0 {
				continue
			}
			go func() {
				lock := locks.Get(isin)
				lock.HighPriorityLock()
				r := processOrder(isin, "cancel")
				results <- r
				lock.HighPriorityUnlock()
				wg.Done()

			}()
		}
	}

	wg.Wait()
	close(results)

	var create []string
	var cancel []string
	for r := range results {
		fmt.Printf("%s\n", r)
		if strings.Index(r, "create") > 0 {
			create = append(create, r)
		}

		if strings.Index(r, "cancel") > 0 {
			cancel = append(cancel, r)
		}
	}

	fmt.Printf("cancel %d, created %d", len(cancel), len(create))

}

func processOrder(isin, process string) string {
	time.Sleep(20 * time.Microsecond)
	return fmt.Sprintf("order %s %s", isin, process)
}
