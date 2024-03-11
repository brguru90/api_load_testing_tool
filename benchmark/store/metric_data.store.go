package store

import (
	"sync"
	"sync/atomic"
	"time"
)

type BenchmarkDataStoreInfo struct {
	UpdatedAt int64
	Other     interface{}
}

var benchmark_data_store []interface{} = []interface{}{}
var benchmark_data_store_info BenchmarkDataStoreInfo = BenchmarkDataStoreInfo{}
var benchmark_data_store_lock sync.Mutex
var benchmark_data_store_info_lock sync.Mutex
var benchmark_data_store_q = make(chan interface{}, 1000000)
var watching_benchmark_data_store_q atomic.Bool
var benchmark_data_store_callback *func([]interface{}, interface{}) []interface{} = nil

func init() {
	watching_benchmark_data_store_q.Store(false)
}

func BenchmarkDataStore_ManualAppendFromQ(callback *func([]interface{}, interface{}) []interface{}) {
	benchmark_data_store_callback = callback
}

func BenchmarkDataStore_AppendFromQ() {
	watching_benchmark_data_store_q.Store(true)
	go func() {
		defer watching_benchmark_data_store_q.Store(false)
		for lc := range benchmark_data_store_q {
			benchmark_data_store_lock.Lock()
			if benchmark_data_store_callback != nil {
				benchmark_data_store = (*benchmark_data_store_callback)(benchmark_data_store, lc)
			} else {
				benchmark_data_store = append(benchmark_data_store, lc)
			}
			benchmark_data_store_lock.Unlock()
			// fmt.Println("Number of goroutines:", runtime.NumGoroutine())
			if len(benchmark_data_store_q) == 0 {
				break
			}
		}
	}()
}

func BenchmarkDataStore_Append(lc interface{}, updated_at int64) {
	if !watching_benchmark_data_store_q.Load() {
		BenchmarkDataStore_AppendFromQ()
	}
	benchmark_data_store_q <- lc
	defer benchmark_data_store_info_lock.Unlock()
	benchmark_data_store_info_lock.Lock()
	if updated_at > 0 {
		benchmark_data_store_info = BenchmarkDataStoreInfo{
			UpdatedAt: updated_at,
		}
	}
}

func BenchmarkDataStore_Reset(lc interface{}) {
	defer benchmark_data_store_lock.Unlock()
	benchmark_data_store_lock.Lock()
	benchmark_data_store = []interface{}{}
}

func BenchmarkDataStore_GetInfo() BenchmarkDataStoreInfo {
	return benchmark_data_store_info
}

func BenchmarkDataStore_Get(index int64) interface{} {
	defer benchmark_data_store_lock.Unlock()
	benchmark_data_store_lock.Lock()
	return benchmark_data_store[index]
}
func BenchmarkDataStore_GetAll() *[]interface{} {
	defer benchmark_data_store_lock.Unlock()
	benchmark_data_store_lock.Lock()
	return &benchmark_data_store
}

func BenchmarkDataStore_GetAllWithInfo() ([]interface{}, BenchmarkDataStoreInfo) {
	defer (func() {
		benchmark_data_store_lock.Unlock()
		benchmark_data_store_info_lock.Unlock()
	})()
	benchmark_data_store_lock.Lock()
	benchmark_data_store_info_lock.Lock()
	return benchmark_data_store, benchmark_data_store_info
}

func BenchmarkDataStore_WaitForAppend() {
	if !watching_benchmark_data_store_q.Load() && len(benchmark_data_store_q) > 0 {
		BenchmarkDataStore_AppendFromQ()
	}
	for {
		if len(benchmark_data_store_q) == 0 {
			break
		}
		time.Sleep(time.Second * 1)
	}
}

func BenchmarkDataStore_CloseQ() {
	if benchmark_data_store_q != nil {
		close(benchmark_data_store_q)
	}
	benchmark_data_store_q = nil
}

func BenchmarkDataStore_Dispose() {
	defer benchmark_data_store_lock.Unlock()
	benchmark_data_store_lock.Lock()
	if benchmark_data_store_q != nil {
		close(benchmark_data_store_q)
	}
	benchmark_data_store_q = nil
	benchmark_data_store = []interface{}{}
	benchmark_data_store_info = BenchmarkDataStoreInfo{}
}
