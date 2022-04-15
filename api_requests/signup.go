package api_requests

import (
	"apis_load_test/my_modules"
	"fmt"
	// "encoding/json"
)

func SignUp() interface{} {
	var total_req int64=100
	var concurrent_req int64=10

	_url := "http://localhost:8000/api/sign_up/"

	// payload_obj := map[string]interface{}{
	// 	"email":       my_modules.RandomString(100) + "@gmail.com",
	// 	"name":        my_modules.RandomString(20),
	// 	"description": my_modules.RandomString(100),
	// }
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	// iteration_data,all_data := my_modules.BenchmarkAPI(10,2,_url, "post", headers, payload_obj,nil)
	iteration_data,all_data := my_modules.BenchmarkAPIAsMultiUser(total_req,concurrent_req,_url, "post", headers, nil,func(uid int64) map[string]interface{} {
		return map[string]interface{}{
			"email":       my_modules.RandomString(100) + "@gmail.com",
			"name":        my_modules.RandomString(20),
			"description": my_modules.RandomString(100),
		}
	},nil,nil)

	fmt.Println("bench mark on api finished")

	return map[string]interface{}{
		"iteration_data":iteration_data,
		"all_data":all_data,
	}
}
