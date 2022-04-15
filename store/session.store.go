package store

import "net/http"

var cookies []*http.Cookie=[]*http.Cookie{}


func GetAllCookies() []*http.Cookie {
	return cookies
}

func GetCookiesRefs() *[]*http.Cookie {
	return &cookies
}

func ResetCookies(){
	cookies=[]*http.Cookie{}
}

func AppendCookie(cookie *http.Cookie)  {
	cookies = append(cookies, cookie)
}

func PopCookie() *http.Cookie{
	cookies=cookies[:len(cookies)-1]
	return cookies[len(cookies)-1]
}

