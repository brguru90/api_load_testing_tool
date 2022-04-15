package store

import "net/http"

type RequestSideSession struct{
	Cookies []*http.Cookie
	CSRF_token string
}

var req_sessions []RequestSideSession=[]RequestSideSession{}


func GetAllSessions() []RequestSideSession {
	return req_sessions
}

func GetSessionsRefs() *[]RequestSideSession{
	return &req_sessions
}

func ResetSessions(){
	req_sessions=[]RequestSideSession{}
}

func AppendCSession(req_session RequestSideSession)  {
	req_sessions = append(req_sessions, req_session)
}

func PopSession() RequestSideSession {
	req_sessions=req_sessions[:len(req_sessions)-1]
	return req_sessions[len(req_sessions)-1]
}

