package api_requests

import (
	"apis_load_test/my_modules"
	"fmt"
	"net/http"
)



func TestInvalidateCache() interface{} {
	var total_req int64 = 10
	var concurrent_req int64 = 2

	_url := "http://localhost:8000/api/del_user_cache/600"
	headers := map[string]string{
		"Content-Type": "application/json",
	}


	request_interceptor := func(req *http.Request, uid int64) {
		fmt.Printf("request interceptor uid--> %v\n", uid)
		req.Header.Add("secret","1234")
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
