package store

import (
	"fmt"
	"time"
)

type LoginCredential struct {
	Name  string
	Email string
}

var login_credential_q = make(chan LoginCredential)
var watching_q = false

func AppendFromQ() {
	watching_q = true
	fmt.Println("LoginCredential: watching_q")
	go func() {
		for lc := range login_credential_q {
			login_credential = append(login_credential, lc)
		}
	}()
}

var login_credential []LoginCredential = []LoginCredential{}

func LoginCredential_Append(lc LoginCredential) {
	if !watching_q {
		AppendFromQ()
	}
	login_credential_q <- lc
}

func LoginCredential_Reset(lc LoginCredential) {
	login_credential = []LoginCredential{}
}

func LoginCredential_Get(index int64) LoginCredential {
	return login_credential[index]
}
func LoginCredential_GetAll() *[]LoginCredential {
	return &login_credential
}

func LoginCredential_WaitForAppend() {
	for {
		if len(login_credential_q) == 0 {
			break
		}
		time.Sleep(time.Second * 1)
	}
}
