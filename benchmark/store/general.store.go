package store

import (
	"time"
)

type GeneralStoreInfo struct {
	UpdatedAt int64
	Other     interface{}
}

var general_store []interface{} = []interface{}{}
var general_store_info GeneralStoreInfo = GeneralStoreInfo{}
var general_store_q = make(chan interface{}, 1000000)
var watching_general_store_q = false
var general_store_callback *func([]interface{},interface{}) []interface{}=nil

func GeneralStore_ManualAppendFromQ(callback *func([]interface{},interface{}) []interface{}) {
	general_store_callback=callback
}

func GeneralStore_AppendFromQ() {
	watching_general_store_q = true
	go func() {
		for lc := range general_store_q {
			if general_store_callback!=nil{
				general_store=(*general_store_callback)(general_store,lc)
			} else{
				general_store = append(general_store, lc)
			}
		}
	}()
}

func GeneralStore_Append(lc interface{}, updated_at int64) {
	if !watching_general_store_q {
		GeneralStore_AppendFromQ()
	}
	general_store_q <- lc
	if updated_at > 0 {
		general_store_info = GeneralStoreInfo{
			UpdatedAt: updated_at,
		}
	}
}

func GeneralStore_Reset(lc interface{}) {
	general_store = []interface{}{}
}

func GeneralStore_GetInfo() GeneralStoreInfo {
	return general_store_info
}

func GeneralStore_Get(index int64) interface{} {
	return general_store[index]
}
func GeneralStore_GetAll() *[]interface{} {
	return &general_store
}

func GeneralStore_GetAllWithInfo() (*[]interface{}, *GeneralStoreInfo) {
	return &general_store, &general_store_info
}

func GeneralStore_WaitForAppend() {
	for {
		if len(general_store_q) == 0 {
			break
		}
		time.Sleep(time.Second * 1)
	}
}
