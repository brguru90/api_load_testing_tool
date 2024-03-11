package store

import (
	"sync/atomic"
)

type CredentialStore[T any] struct {
	store_data   []T
	store_data_q chan T
	queue_open   atomic.Bool
}

func (CredentialStore[T]) NewCredentialStore(buffer_size int64) CredentialStore[T] {
	t := CredentialStore[T]{
		store_data:   []T{},
		store_data_q: make(chan T, buffer_size),
	}
	t.queue_open.Store(true)
	return t
}

func (e *CredentialStore[T]) CredentialStore_AppendFromQ() {
	for lc := range e.store_data_q {
		e.store_data = append(e.store_data, lc)
		if len(e.store_data_q) == 0 {
			break
		}
	}
}

func (e *CredentialStore[T]) CredentialStore_Append(lc T) {
	if !e.queue_open.Load() {
		return
	}
	e.store_data_q <- lc
	if len(e.store_data_q) == cap(e.store_data_q) {
		e.queue_open.Store(false)
	}
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

func (e *CredentialStore[T]) Len() int {
	return len(e.store_data)
}

func (e *CredentialStore[T]) CredentialStore_GetCount() int {
	return len(e.store_data)
}

func (e *CredentialStore[T]) CredentialStore_Pop() []T {
	e.store_data = e.store_data[:len(e.store_data)-1]
	return e.store_data
}

func (e *CredentialStore[T]) CredentialStore_WaitForAppend() {
	if len(e.store_data_q) > 0 {
		e.CredentialStore_AppendFromQ()
	}
	e.CloseQ()
}

func (e *CredentialStore[T]) CloseQ() {
	if e.store_data_q != nil {
		close(e.store_data_q)
	}
	e.store_data_q = nil
	e.queue_open.Store(false)
}

func (e *CredentialStore[T]) Dispose() {
	e.CloseQ()
	e.store_data = []T{}
}
