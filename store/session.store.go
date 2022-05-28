package store

import (
	"net/http"
	"time"
)

type RequestSideSession struct {
	Cookies    []*http.Cookie
	CSRF_token string
}

var req_sessions []RequestSideSession = []RequestSideSession{}
var req_sessions_q = make(chan RequestSideSession)
var watching_req_sessions_q = false

func RequestSideSession_AppendFromQ() {
	watching_req_sessions_q = true
	go func() {
		for rs := range req_sessions_q {
			req_sessions = append(req_sessions, rs)
		}
	}()
}

func GetAllSessions() []RequestSideSession {
	return req_sessions
}

func GetSessionsRefs() *[]RequestSideSession {
	return &req_sessions
}

func ResetSessions() {
	req_sessions = []RequestSideSession{}
}

func AppendCSession(req_session RequestSideSession) {
	if !watching_req_sessions_q {
		RequestSideSession_AppendFromQ()
	}
	req_sessions_q <- req_session
	req_sessions = append(req_sessions, req_session)
}

func PopSession() RequestSideSession {
	req_sessions = req_sessions[:len(req_sessions)-1]
	return req_sessions[len(req_sessions)-1]
}

func RequestSideSession_WaitForAppend() {
	for {
		if len(req_sessions_q) == 0 {
			break
		}
		time.Sleep(time.Second * 1)
	}
}
