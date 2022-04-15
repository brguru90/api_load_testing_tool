package tests

import (
	"apis_load_test/my_modules"	
	"apis_load_test/store"	
	"fmt"
	"net/http"
)

func TestAsMultiUser()  {
	my_modules.LogPath = "./log.json"
	signUpAPI()
}

func signUpAPI()  {

	_url := "http://localhost:8000/api/login/"
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	// api_payload:=map[string]interface{}{
	// 	"email":       my_modules.RandomString(100) + "@gmail.com",
	// 	"name":        my_modules.RandomString(20),
	// 	"description": my_modules.RandomString(100),
	// }

	api_payload:=map[string]interface{}{
		"email":       "lU2zXfRFVocs1p7rkO/6/P+890zFv9QUYqn/D9ihM1NnOO9/JrieMT/sH1DS8Jz/yiNjpDfFL5wp0NDzGEIrnqI6oTwP3Tv5jP8x@gmail.com",
	}	

	request_interceptor:=func (req *http.Request,uid int64)  {
		// fmt.Printf("request interceptor --> %v\n",req)
	}

	response_interceptor:=func (resp *http.Response,uid int64)  {
		// fmt.Printf("response interceptor --> %v\n",resp)	
		
		fmt.Printf("Got cookie --> %v\n\n",len(resp.Cookies()))
		
		for _,cookie:= range resp.Cookies(){
			if(cookie.Name=="access_token"){
				fmt.Printf("access_token--> %v\n\n",cookie)
				store.AppendCookie(cookie)
			}
		}
	}



	fmt.Println(my_modules.APIReq(_url, "post", headers, api_payload,-1,request_interceptor,response_interceptor))

	fmt.Println(store.GetAllCookies())

}