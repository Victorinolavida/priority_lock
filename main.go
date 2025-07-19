package main

import (
	"fmt"
	"sync"
	prioritylock "testgg/priorityLock"
	"time"
)

func main() {

	lock := prioritylock.NewPriorityPreferenceLock()
	var wg sync.WaitGroup
	routines := 8
	wg.Add(2 * routines)

	go func() {
		for range routines/2 + routines {
			lock.LowPriorityLock()
			wg.Done()
			lock.LowPriorityUnlock()

		}

	}()

	go func() {
		for range routines / 2 {
			lock.HighPriorityLock()
			wg.Done()
			lock.HighPriorityUnlock()

		}
	}()

	wg.Wait()
}

func ProcessOrder(isin, process string) string {
	fmt.Printf("start processing (%s) order for %s\n", process, isin)
	time.Sleep(5 * time.Microsecond)
	return fmt.Sprintf("%s,%s", isin, process)
}
