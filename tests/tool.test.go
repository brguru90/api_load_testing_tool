package tests

import (
	"fmt"
	"net/http"

	"github.com/brguru90/api_load_testing_tool/benchmark/my_modules"
)

func TestTool() {
	// var total_req int64 = 500000
	// var concurrent_req int64 = 100000

	var total_req int64 = 200000
	var concurrent_req int64 = 100000

	// _url := "http://localhost:8000/api/test"
	_url := "http://localhost"
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	request_interceptor := func(req *http.Request, uid int64) {
		// fmt.Printf("request interceptor uid--> %v\n", uid)
		// req.Header.Add("secret", "1234")
		// req.URL.RawQuery = "/" + strconv.FormatInt(uid, 10)
	}

	response_interceptor := func(resp *http.Response, uid int64) {
		// fmt.Printf("response interceptor uid--> %v\n", uid)
	}

	my_modules.BenchmarkAPIAsMultiUser(total_req, concurrent_req, _url, "get", headers, nil, nil, request_interceptor, response_interceptor)

	fmt.Println("bench mark on api finished")
}
