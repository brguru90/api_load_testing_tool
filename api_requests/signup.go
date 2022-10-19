package api_requests

import (
	"apis_load_test/benchmark/my_modules"
	"apis_load_test/benchmark/store"
	"fmt"
	// "encoding/json"
)

var signup_credentials store.LoginCredential

func SignUp(total_req int64, concurrent_req int64) interface{} {

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
	signup_credentials=store.NewLoginCredential(concurrent_req)
	iteration_data, all_data := my_modules.BenchmarkAPIAsMultiUser(total_req, concurrent_req, _url, "post", headers, nil, func(uid int64) map[string]interface{} {
		signup_payload := map[string]interface{}{
			"email":       my_modules.RandomString(100) + "@gmail.com",
			"name":        my_modules.RandomString(20),
			"description": my_modules.RandomString(100),
		}
		signup_credentials.LoginCredential_Append(store.LoginCredentialStruct{
			Name:  signup_payload["name"].(string),
			Email: signup_payload["email"].(string),
		})
		return signup_payload
	}, nil, nil)

	// iteration_data,all_data := my_modules.BenchmarkAPI(total_req,concurrent_req,_url, "post", headers, nil,func() map[string]interface{} {
	// 	return map[string]interface{}{
	// 		"email":       my_modules.RandomString(100) + "@gmail.com",
	// 		"name":        my_modules.RandomString(20),
	// 		"description": my_modules.RandomString(100),
	// 	}
	// })

	fmt.Println("bench mark on api finished")
	signup_credentials.LoginCredential_WaitForAppend()

	result := make(map[string]interface{})
	result[_url] = map[string]interface{}{
		"iteration_data": iteration_data,
		"all_data":       all_data,
	}
	return result
}
