package store

import (
	"sync"
	"sync/atomic"
	"time"
)

type CredentialStore[T any] struct {
	store_data       []T
	store_data_q     chan T
	watching_store_q atomic.Bool
	mutex            sync.Mutex
}

func (CredentialStore[T]) NewCredentialStore(buffer_size int64) CredentialStore[T] {
	t := CredentialStore[T]{
		store_data:   []T{},
		store_data_q: make(chan T, buffer_size),
	}
	t.watching_store_q.Store(false)
	return t
}

func (e *CredentialStore[T]) CredentialStore_AppendFromQ() {
	e.watching_store_q.Store(true)
	go func() {
		for lc := range e.store_data_q {
			e.mutex.Lock()
			e.store_data = append(e.store_data, lc)
			e.mutex.Unlock()

		}
	}()
}

func (e *CredentialStore[T]) CredentialStore_Append(lc T) {
	if !e.watching_store_q.Load() {
		e.CredentialStore_AppendFromQ()
	}
	e.store_data_q <- lc
}

func (e *CredentialStore[T]) CredentialStore_Reset(lc T) {
	e.store_data = []T{}
}

func (e *CredentialStore[T]) CredentialStore_Get(index int64) T {
	return e.store_data[index]
}
func (e *CredentialStore[T]) CredentialStore_GetAll() []T {
	return (e.store_data)
}

func (e *CredentialStore[T]) CredentialStore_GetAllRefs() *[]T {
	return &(e.store_data)
}

func (e *CredentialStore[T]) CredentialStore_GetQRefs() *chan T {
	return &e.store_data_q
}

func (e *CredentialStore[T]) Len() int {
	defer e.mutex.Unlock()
	e.mutex.Lock()
	return len(e.store_data)
}

func (e *CredentialStore[T]) CredentialStore_GetCount() int {
	defer e.mutex.Unlock()
	e.mutex.Lock()
	return len(e.store_data) + len(e.store_data_q)
}

func (e *CredentialStore[T]) CredentialStore_Pop() []T {
	for len(e.store_data_q) <= 0 {
	}
	e.mutex.Lock()
	e.store_data = e.store_data[:len(e.store_data)-1]
	e.mutex.Unlock()
	return e.store_data
}

func (e *CredentialStore[T]) CredentialStore_WaitForAppend() {
	for {
		if len(e.store_data_q) == 0 {
			break
		}
		time.Sleep(time.Second * 1)
	}
}

func (e *CredentialStore[T]) CloseQ() {
	if e.store_data_q == nil {
		close(e.store_data_q)
	}
	e.store_data_q = nil
}

func (e *CredentialStore[T]) Dispose() {
	if e.store_data_q == nil {
		close(e.store_data_q)
	}
	e.store_data_q = nil
	e.store_data = []T{}
}
