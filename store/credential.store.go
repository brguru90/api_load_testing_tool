package store

import (
	"time"
)

type LoginCredential struct {
	Name  string
	Email string
}

var login_credential []LoginCredential = []LoginCredential{}
var login_credential_q = make(chan LoginCredential)
var watching_login_credential_q = false

func LoginCredential_AppendFromQ() {
	watching_login_credential_q = true
	go func() {
		for lc := range login_credential_q {
			login_credential = append(login_credential, lc)
		}
	}()
}

func LoginCredential_Append(lc LoginCredential) {
	if !watching_login_credential_q {
		LoginCredential_AppendFromQ()
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
