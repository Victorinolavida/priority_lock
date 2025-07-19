package main

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testgg/lock"
	prioritylock "testgg/priorityLock"
	"testing"
)

var (
	createKey = "create"
	cancelKey = "cancel"
)

func TestUniqueLocker(t *testing.T) {
	numIsins := 10
	numOrders := 10
	lock := prioritylock.NewPriorityPreferenceLock()

	var isins []string
	for i := range numIsins {
		isin := fmt.Sprintf("MX%d", i)
		isins = append(isins, isin)
	}
	totalProcessedOrderfs := len(isins) * numOrders * 2

	var wg sync.WaitGroup
	wg.Add(totalProcessedOrderfs)
	chanOrders := make(chan string, totalProcessedOrderfs)

	for i, isin := range isins {

		for range numOrders {
			if i%2 == 0 {
				go func() {
					lock.HighPriorityLock()
					results := ProcessOrder(isin, cancelKey)
					_, operation := splitResult(results)
					chanOrders <- operation

					wg.Done()
					lock.HighPriorityUnlock()

				}()
			}
		}
		for range numOrders {
			if i%2 != 0 {
				go func() {
					lock.HighPriorityLock()
					results := ProcessOrder(isin, cancelKey)
					_, operation := splitResult(results)
					chanOrders <- operation
					wg.Done()
					lock.HighPriorityUnlock()

				}()
			}
		}

		for range numOrders {
			go func() {
				lock.LowPriorityLock()
				results := ProcessOrder(isin, createKey)
				_, operation := splitResult(results)
				chanOrders <- operation
				wg.Done()
				lock.LowPriorityUnlock()

			}()
		}
	}

	wg.Wait()
	close(chanOrders)
	assertLen(t, totalProcessedOrderfs, chanOrders)
	var totalCreated []string
	var totalCancel []string
	for result := range chanOrders {
		switch result {
		case cancelKey:
			totalCancel = append(totalCancel, result)
		case createKey:
			totalCreated = append(totalCreated, result)
		}
	}

	assertLen(t, totalProcessedOrderfs/2, totalCreated)
	assertLen(t, totalProcessedOrderfs/2, totalCancel)
}

func TestPriorityLockMap(t *testing.T) {

	numIsins := 10
	numOrders := 10
	// lock := prioritylock.NewPriorityPreferenceLock()
	globalLock := lock.GetGlobalLock()

	var isins []string
	for i := range numIsins {
		isin := fmt.Sprintf("MX%d", i)
		isins = append(isins, isin)
	}
	totalProcessedOrderfs := len(isins) * numOrders * 2

	var wg sync.WaitGroup
	wg.Add(totalProcessedOrderfs)
	chanOrders := make(chan string, totalProcessedOrderfs)

	for i, isin := range isins {

		for range numOrders {
			if i%2 == 0 {
				go func() {
					lock := globalLock.Get(isin)
					lock.HighPriorityLock()
					results := ProcessOrder(isin, cancelKey)
					_, operation := splitResult(results)
					chanOrders <- operation

					wg.Done()
					lock.HighPriorityUnlock()

				}()
			}
		}
		for range numOrders {
			if i%2 != 0 {
				go func() {
					lock := globalLock.Get(isin)
					lock.HighPriorityLock()
					results := ProcessOrder(isin, cancelKey)
					_, operation := splitResult(results)
					chanOrders <- operation
					wg.Done()
					lock.HighPriorityUnlock()

				}()
			}
		}

		for range numOrders {
			go func() {
				lock := globalLock.Get(isin)
				lock.LowPriorityLock()
				results := ProcessOrder(isin, createKey)
				_, operation := splitResult(results)
				chanOrders <- operation
				wg.Done()
				lock.LowPriorityUnlock()

			}()
		}
	}

	wg.Wait()
	close(chanOrders)
	assertLen(t, totalProcessedOrderfs, chanOrders)
	var totalCreated []string
	var totalCancel []string
	for result := range chanOrders {
		switch result {
		case cancelKey:
			totalCancel = append(totalCancel, result)
		case createKey:
			totalCreated = append(totalCreated, result)
		}
	}

	assertLen(t, totalProcessedOrderfs/2, totalCreated)
	assertLen(t, totalProcessedOrderfs/2, totalCancel)

}

func splitResult(result string) (string, string) {
	resultedSplit := strings.Split(result, ",")
	if len(resultedSplit) != 2 {
		return "", ""
	}
	return resultedSplit[0], resultedSplit[1]
}

func assertLen(t *testing.T, expected int, got any) {
	t.Helper()

	rv := reflect.ValueOf(got)

	switch rv.Kind() {
	case reflect.Array, reflect.Slice, reflect.Map, reflect.String, reflect.Chan:

		gotLen := rv.Len()
		if gotLen != expected {
			t.Errorf("got %d, expected %d", gotLen, expected)
		}

	default:
		t.Errorf("invalid data type %s", rv.Kind())
	}

}
