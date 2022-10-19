package store

import (
	"time"
)

type LoginCredentialStruct struct {
	Name  string
	Email string
}

type LoginCredential struct{
	login_credential []LoginCredentialStruct
	login_credential_q chan LoginCredentialStruct
	watching_login_credential_q bool
}

// var login_credential []LoginCredentialStruct = []LoginCredentialStruct{}
// var login_credential_q = make(chan LoginCredentialStruct,1000000)
// var watching_login_credential_q = false

func NewLoginCredential(buffer_size int64) LoginCredential{
	return LoginCredential{
		login_credential:[]LoginCredentialStruct{},
		login_credential_q:make(chan LoginCredentialStruct,buffer_size),
		watching_login_credential_q:false,
	}
}

func (e *LoginCredential) LoginCredential_AppendFromQ() {
	e.watching_login_credential_q = true
	go func() {
		for lc := range e.login_credential_q {
			e.login_credential = append(e.login_credential, lc)
		}
	}()
}

func (e *LoginCredential) LoginCredential_Append(lc LoginCredentialStruct) {
	if !e.watching_login_credential_q {
		e.LoginCredential_AppendFromQ()
	}
	e.login_credential_q <- lc
}

func (e *LoginCredential) LoginCredential_Reset(lc LoginCredentialStruct) {
	e.login_credential = []LoginCredentialStruct{}
}

func (e *LoginCredential) LoginCredential_Get(index int64) LoginCredentialStruct {
	return e.login_credential[index]
}
func (e *LoginCredential) LoginCredential_GetAll() *[]LoginCredentialStruct {
	return &(e.login_credential)
}

func (e *LoginCredential) LoginCredential_WaitForAppend() {
	for {
		if len(e.login_credential_q) == 0 {
			break
		}
		time.Sleep(time.Second * 1)
	}
}

func  (e *LoginCredential)  Dispose(){
	close(e.login_credential_q)
	e.login_credential_q=nil
	e.login_credential=[]LoginCredentialStruct{}
}
