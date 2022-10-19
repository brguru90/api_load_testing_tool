package user

import (
	"apis_load_test/benchmark/my_modules"
	"apis_load_test/benchmark/store"
	"fmt"
	"math/rand"
	"net/http"
)

func GetUserDetailAsMultiUser(_url string, total_req int64, concurrent_req int64) interface{} {

	// _url := "http://localhost:8000/api/user/"
	// _url :="http://localhost:8000/api/user/?page=1000&limit=20"
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	request_interceptor := func(req *http.Request, uid int64) {
		// fmt.Printf("request interceptor uid--> %v\n", uid)

		if len(*store.GetSessionsRefs()) <= int(uid%concurrent_req){
			return 
		}

		req.Header.Add("csrf_token", (*store.GetSessionsRefs())[uid%concurrent_req].CSRF_token)

		for _, cookie := range (*store.GetSessionsRefs())[uid%concurrent_req].Cookies {
			req.AddCookie(cookie)
		}

		q := req.URL.Query()
		q.Add("limit","20")
		q.Add("page",fmt.Sprintf("%d",rand.Intn(1000)))
		req.URL.RawQuery = q.Encode()
	}

	iteration_data, all_data := my_modules.BenchmarkAPIAsMultiUser(total_req, concurrent_req, _url, "get", headers, nil, nil, request_interceptor, nil)

	fmt.Println("bench mark on api finished")

	result := make(map[string]interface{})
	result[_url] = map[string]interface{}{
		"iteration_data": iteration_data,
		"all_data":       all_data,
	}
	return result
}
