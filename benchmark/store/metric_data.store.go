package store

import (
	"time"
)

type BenchmarkDataStoreInfo struct {
	UpdatedAt int64
	Other     interface{}
}

var benchmark_data_store []interface{} = []interface{}{}
var benchmark_data_store_info BenchmarkDataStoreInfo = BenchmarkDataStoreInfo{}
var benchmark_data_store_q = make(chan interface{}, 1000000)
var watching_benchmark_data_store_q = false
var benchmark_data_store_callback *func([]interface{},interface{}) []interface{}=nil

func BenchmarkDataStore_ManualAppendFromQ(callback *func([]interface{},interface{}) []interface{}) {
	benchmark_data_store_callback=callback
}

func BenchmarkDataStore_AppendFromQ() {
	watching_benchmark_data_store_q = true
	go func() {
		for lc := range benchmark_data_store_q {
			if benchmark_data_store_callback!=nil{
				benchmark_data_store=(*benchmark_data_store_callback)(benchmark_data_store,lc)
			} else{
				benchmark_data_store = append(benchmark_data_store, lc)
			}
		}
	}()
}

func BenchmarkDataStore_Append(lc interface{}, updated_at int64) {
	if !watching_benchmark_data_store_q {
		BenchmarkDataStore_AppendFromQ()
	}
	benchmark_data_store_q <- lc
	if updated_at > 0 {
		benchmark_data_store_info = BenchmarkDataStoreInfo{
			UpdatedAt: updated_at,
		}
	}
}

func BenchmarkDataStore_Reset(lc interface{}) {
	benchmark_data_store = []interface{}{}
}

func BenchmarkDataStore_GetInfo() BenchmarkDataStoreInfo {
	return benchmark_data_store_info
}

func BenchmarkDataStore_Get(index int64) interface{} {
	return benchmark_data_store[index]
}
func BenchmarkDataStore_GetAll() *[]interface{} {
	return &benchmark_data_store
}

func BenchmarkDataStore_GetAllWithInfo() (*[]interface{}, *BenchmarkDataStoreInfo) {
	return &benchmark_data_store, &benchmark_data_store_info
}

func BenchmarkDataStore_WaitForAppend() {
	for {
		if len(benchmark_data_store_q) == 0 {
			break
		}
		time.Sleep(time.Second * 1)
	}
}

func BenchmarkDataStore_CloseQ(){
	close(benchmark_data_store_q)
	benchmark_data_store_q=nil
}

func BenchmarkDataStore_Dispose(){
	close(benchmark_data_store_q)
	benchmark_data_store_q=nil
	benchmark_data_store = []interface{}{}
	benchmark_data_store_info = BenchmarkDataStoreInfo{}
}