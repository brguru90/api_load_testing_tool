package store

import (
	"sync"
	"time"
)

type BenchmarkDataStoreInfo struct {
	UpdatedAt int64
	Other     interface{}
}

var (
	benchmark_data_store          = []interface{}{}
	benchmark_data_store_info     = BenchmarkDataStoreInfo{}
	benchmark_data_store_q        = make(chan interface{}, 1000000)
	benchmark_data_store_callback *func([]interface{}, interface{}) []interface{}
	benchmark_data_store_mu       sync.RWMutex
	benchmark_data_store_once     sync.Once
	benchmark_data_store_closed   bool
)

func BenchmarkDataStore_ManualAppendFromQ(callback *func([]interface{}, interface{}) []interface{}) {
	benchmark_data_store_callback = callback
}

// Spawned exactly once. Sole writer of benchmark_data_store; readers must
// hold benchmark_data_store_mu.RLock().
func BenchmarkDataStore_AppendFromQ() {
	benchmark_data_store_once.Do(func() {
		go func() {
			for lc := range benchmark_data_store_q {
				benchmark_data_store_mu.Lock()
				if benchmark_data_store_callback != nil {
					benchmark_data_store = (*benchmark_data_store_callback)(benchmark_data_store, lc)
				} else {
					benchmark_data_store = append(benchmark_data_store, lc)
				}
				benchmark_data_store_mu.Unlock()
			}
		}()
	})
}

func BenchmarkDataStore_Append(lc interface{}, updated_at int64) {
	BenchmarkDataStore_AppendFromQ()
	benchmark_data_store_q <- lc
	if updated_at > 0 {
		benchmark_data_store_mu.Lock()
		benchmark_data_store_info = BenchmarkDataStoreInfo{
			UpdatedAt: updated_at,
		}
		benchmark_data_store_mu.Unlock()
	}
}

func BenchmarkDataStore_Reset(lc interface{}) {
	benchmark_data_store_mu.Lock()
	benchmark_data_store = []interface{}{}
	benchmark_data_store_mu.Unlock()
}

func BenchmarkDataStore_GetInfo() BenchmarkDataStoreInfo {
	benchmark_data_store_mu.RLock()
	defer benchmark_data_store_mu.RUnlock()
	return benchmark_data_store_info
}

func BenchmarkDataStore_Get(index int64) interface{} {
	benchmark_data_store_mu.RLock()
	defer benchmark_data_store_mu.RUnlock()
	return benchmark_data_store[index]
}

// Returns a defensive copy so the caller can iterate without holding the lock.
func BenchmarkDataStore_GetAll() *[]interface{} {
	benchmark_data_store_mu.RLock()
	defer benchmark_data_store_mu.RUnlock()
	out := make([]interface{}, len(benchmark_data_store))
	copy(out, benchmark_data_store)
	return &out
}

// Snapshot copy of slice + info — safe to use from another goroutine.
func BenchmarkDataStore_GetAllWithInfo() (*[]interface{}, *BenchmarkDataStoreInfo) {
	benchmark_data_store_mu.RLock()
	defer benchmark_data_store_mu.RUnlock()
	out := make([]interface{}, len(benchmark_data_store))
	copy(out, benchmark_data_store)
	info := benchmark_data_store_info
	return &out, &info
}

func BenchmarkDataStore_WaitForAppend() {
	for {
		if len(benchmark_data_store_q) == 0 {
			break
		}
		time.Sleep(time.Second * 1)
	}
}

func BenchmarkDataStore_CloseQ() {
	benchmark_data_store_mu.Lock()
	defer benchmark_data_store_mu.Unlock()
	if !benchmark_data_store_closed {
		close(benchmark_data_store_q)
		benchmark_data_store_closed = true
	}
}

func BenchmarkDataStore_Dispose() {
	benchmark_data_store_mu.Lock()
	defer benchmark_data_store_mu.Unlock()
	if !benchmark_data_store_closed {
		close(benchmark_data_store_q)
		benchmark_data_store_closed = true
	}
	benchmark_data_store = []interface{}{}
	benchmark_data_store_info = BenchmarkDataStoreInfo{}
}
