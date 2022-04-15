package store

import "net/http"

type ListOfCookies []*http.Cookie

var cookies []ListOfCookies=[]ListOfCookies{}


func GetAllCookies() []ListOfCookies {
	return cookies
}

func GetCookiesRefs() *[]ListOfCookies{
	return &cookies
}

func ResetCookies(){
	cookies=[]ListOfCookies{}
}

func AppendCookie(cookie []*http.Cookie)  {
	cookies = append(cookies, cookie)
}

func PopCookie() ListOfCookies {
	cookies=cookies[:len(cookies)-1]
	return cookies[len(cookies)-1]
}

