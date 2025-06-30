package maplock

import (
	"sync"
)

// type MapLock[T any] interface {
// 	set(identifier string)
// 	load(identifier string)
// 	LoadOrSet(identifier string) T
// }

type MapLock[T any] struct {
	mux       *sync.RWMutex
	mapValues map[string]*T
}

func NewMapLock[T any]() *MapLock[T] {
	return &MapLock[T]{
		mux:       &sync.RWMutex{},
		mapValues: make(map[string]*T),
	}
}

func (m *MapLock[T]) LoadOrSet(identifier string, newVal T) *T {
	m.mux.Lock()
	defer m.mux.Unlock()

	val, exist := m.mapValues[identifier]
	if exist {
		return val
	}

	m.mapValues[identifier] = &newVal
	return &newVal
}
