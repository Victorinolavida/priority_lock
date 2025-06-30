package main

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	createKey = "create"
	cancelKey = "cancel"
)

func TestUniqueLocker_SerialExecutionSameKey(t *testing.T) {
	locker := GetLocks()

	// var mu sync.Mutex
	totalProcess := 100

	results := make(chan string, totalProcess)

	wg := sync.WaitGroup{}
	wg.Add(totalProcess)

	go func() {
		for i := range totalProcess {
			if i%2 == 0 {
				continue
			}
			order, result := proccessOrder(100, cancelKey)
			locker.Lock(order)
			time.Sleep(10 * time.Millisecond) // aseguramos que el primero bloquee
			locker.Unlock(order)

			results <- result

			wg.Done()

		}
	}()

	go func() {

		for i := range totalProcess {
			if i%2 != 0 {
				continue
			}
			order, result := proccessOrder(100, createKey)
			locker.Lock(order)
			time.Sleep(10 * time.Millisecond) // aseguramos que el primero bloquee
			locker.Unlock(order)

			results <- result

			wg.Done()
		}

	}()

	wg.Wait()
	close(results)
	for r := range results {
		fmt.Printf("%v\n", r)
	}

	// La goroutine que hizo Lock primero debe haber ejecutado primero
}
func TestUniqueLocker_ParallelExecutionDifferentKeys(t *testing.T) {
	locker := GetLocks()

	var result []int
	var mu sync.Mutex

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		locker.Lock("order-A")
		defer locker.Unlock("order-A")
		defer wg.Done()

		time.Sleep(50 * time.Millisecond)
		mu.Lock()
		result = append(result, 1)
		mu.Unlock()
	}()

	go func() {
		locker.Lock("order-B")
		defer locker.Unlock("order-B")
		defer wg.Done()

		time.Sleep(10 * time.Millisecond)
		mu.Lock()
		result = append(result, 2)
		mu.Unlock()
	}()

	wg.Wait()

	// Como las claves son distintas, el orden puede variar
	assert.ElementsMatch(t, []int{1, 2}, result)
}

func proccessOrder(i int, process string) (string, string) {
	id := fmt.Sprintf("order-%d\n", i)
	fullProcess := fmt.Sprintf("%s ---> %s", process, id)
	return id, fullProcess
}
