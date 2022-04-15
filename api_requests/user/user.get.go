package user

import (
	"apis_load_test/my_modules"
	"apis_load_test/store"
	"fmt"
	"net/http"
)



func GetUserDetailAsMultiUser() interface{} {
	var total_req int64 = 10
	var concurrent_req int64 = 2

	_url := "http://localhost:8000/api/user/"
	headers := map[string]string{
		"Content-Type": "application/json",
	}


	request_interceptor := func(req *http.Request, uid int64) {
		fmt.Printf("request interceptor uid--> %v\n", uid)

		req.Header.Add("csrf_token",(*store.GetSessionsRefs())[uid].CSRF_token)
		
		for _,cookie := range (*store.GetSessionsRefs())[uid].Cookies{
			req.AddCookie(cookie)
		}
	}

	response_interceptor := func(resp *http.Response, uid int64) {
		fmt.Printf("response interceptor uid--> %v\n", uid)
	}

	iteration_data, all_data := my_modules.BenchmarkAPIAsMultiUser(total_req, concurrent_req, _url, "get", headers, nil, nil, request_interceptor, response_interceptor)

	fmt.Println("bench mark on api finished")

	return map[string]interface{}{
		"iteration_data": iteration_data,
		"all_data":       all_data,
	}
}
