package store

import (
	"time"
)

type GeneralStore map[string]interface{}
type GeneralStoreInfo struct {
	UpdatedAt int64
	Other     GeneralStore
}

var general_store []GeneralStore = []GeneralStore{}
var general_store_info GeneralStoreInfo = GeneralStoreInfo{}
var general_store_q = make(chan GeneralStore, 1000000)
var watching_general_store_q = false

func GeneralStore_AppendFromQ() {
	watching_general_store_q = true
	go func() {
		for lc := range general_store_q {
			general_store = append(general_store, lc)
		}
	}()
}

func GeneralStore_Append(lc GeneralStore, updated_at int64) {
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

func GeneralStore_Reset(lc GeneralStore) {
	general_store = []GeneralStore{}
}

func GeneralStore_GetInfo() GeneralStoreInfo {
	return general_store_info
}

func GeneralStore_Get(index int64) GeneralStore {
	return general_store[index]
}
func GeneralStore_GetAll() *[]GeneralStore {
	return &general_store
}

func GeneralStore_GetAllWithInfo() (*[]GeneralStore, *GeneralStoreInfo) {
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
