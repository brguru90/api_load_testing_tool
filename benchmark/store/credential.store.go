package store

import (
	"sync"
	"time"
)

type CredentialStore[T any] struct {
	store_data   []T
	store_data_q chan T
	// Pointers so the struct itself stays copy-safe; vet flags by-value mutex
	// in a value-receiver constructor, and a copied lock would violate its
	// invariants. All copies share the same heap-allocated mutex/once.
	mu        *sync.RWMutex
	startOnce *sync.Once
	closed    bool
}

func (CredentialStore[T]) NewCredentialStore(buffer_size int64) CredentialStore[T] {
	return CredentialStore[T]{
		store_data:   []T{},
		store_data_q: make(chan T, buffer_size),
		mu:           &sync.RWMutex{},
		startOnce:    &sync.Once{},
	}
}

// Spawned exactly once via startOnce. Drain goroutine is the sole writer of
// store_data; it holds e.mu.Lock() per append so concurrent readers (Get,
// GetCount) see a consistent slice header.
func (e *CredentialStore[T]) CredentialStore_AppendFromQ() {
	e.startOnce.Do(func() {
		go func() {
			for lc := range e.store_data_q {
				e.mu.Lock()
				e.store_data = append(e.store_data, lc)
				e.mu.Unlock()
			}
		}()
	})
}

func (e *CredentialStore[T]) CredentialStore_Append(lc T) {
	// Idempotent: only first caller spawns the drain goroutine.
	e.CredentialStore_AppendFromQ()
	e.store_data_q <- lc
}

func (e *CredentialStore[T]) CredentialStore_Reset(lc T) {
	e.mu.Lock()
	e.store_data = []T{}
	e.mu.Unlock()
}

func (e *CredentialStore[T]) CredentialStore_Get(index int64) T {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.store_data[index]
}

// Returns a defensive copy. Use this when you need a stable snapshot.
func (e *CredentialStore[T]) CredentialStore_GetAll() []T {
	e.mu.RLock()
	defer e.mu.RUnlock()
	out := make([]T, len(e.store_data))
	copy(out, e.store_data)
	return out
}

// Returns a pointer to the live slice. Callers MUST hold no concurrent
// writers; safe only after CredentialStore_WaitForAppend + CloseQ, or in
// phases where the producer has stopped.
func (e *CredentialStore[T]) CredentialStore_GetAllRefs() *[]T {
	return &(e.store_data)
}

func (e *CredentialStore[T]) CredentialStore_GetQRefs() *chan T {
	return &e.store_data_q
}

func (e *CredentialStore[T]) CredentialStore_GetCount() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return len(e.store_data) + len(e.store_data_q)
}

func (e *CredentialStore[T]) CredentialStore_Pop() []T {
	for len(e.store_data_q) <= 0 {
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.store_data = e.store_data[:len(e.store_data)-1]
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
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.store_data_q != nil && !e.closed {
		close(e.store_data_q)
		e.closed = true
	}
}

func (e *CredentialStore[T]) Dispose() {
	e.mu.Lock()
	defer e.mu.Unlock()
	if e.store_data_q != nil && !e.closed {
		close(e.store_data_q)
		e.closed = true
	}
	e.store_data = []T{}
}
