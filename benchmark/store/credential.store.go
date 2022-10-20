package store

import (
	"time"
)

type CredentialStore[T any] struct{
	store_data []T
	store_data_q chan T
	watching_store_q bool
}

func (CredentialStore[T]) NewCredentialStore(buffer_size int64) CredentialStore[T]{
	return CredentialStore[T]{
		store_data:[]T{},
		store_data_q:make(chan T,buffer_size),
		watching_store_q:false,
	}
}

func (e *CredentialStore[T]) CredentialStore_AppendFromQ() {
	e.watching_store_q = true
	go func() {
		for lc := range e.store_data_q {
			e.store_data = append(e.store_data, lc)
		}
	}()
}

func (e *CredentialStore[T]) CredentialStore_Append(lc T) {
	if !e.watching_store_q {
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

func (e *CredentialStore[T]) CredentialStore_GetQRefs() *chan T{
	return &e.store_data_q
}

func (e *CredentialStore[T]) CredentialStore_GetCount() int{
	return len(e.store_data)+len(e.store_data_q)
}

func (e *CredentialStore[T]) CredentialStore_Pop() []T{
	for ;len(e.store_data_q)<=0;{}
	e.store_data=e.store_data[:len(e.store_data)-1]
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

func  (e *CredentialStore[T])  CloseQ(){
	close(e.store_data_q)
	e.store_data_q=nil
}

func  (e *CredentialStore[T])  Dispose(){
	close(e.store_data_q)
	e.store_data_q=nil
	e.store_data=[]T{}
}


