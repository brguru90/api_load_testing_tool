package api_requests

import (
	"fmt"

	"github.com/brguru90/api_load_testing_tool/benchmark/my_modules"
	"github.com/brguru90/api_load_testing_tool/benchmark/store"
	// "encoding/json"
)

type LoginCredentialStruct struct {
	Name  string
	Email string
}

var signup_credentials store.CredentialStore[LoginCredentialStruct]

func SignUp(total_req int64, concurrent_req int64) interface{} {

	_url := "http://localhost:8000/api/sign_up/"

	headers := map[string]string{
		"Content-Type": "application/json",
	}

	payload_generator_callback := func(uid int64) map[string]interface{} {
		signup_payload := map[string]interface{}{
			"email":       my_modules.RandomString(100) + "@gmail.com",
			"name":        my_modules.RandomString(20),
			"description": my_modules.RandomString(100),
		}
		signup_credentials.CredentialStore_Append(LoginCredentialStruct{
			Name:  signup_payload["name"].(string),
			Email: signup_payload["email"].(string),
		})
		return signup_payload
	}

	signup_credentials = signup_credentials.NewCredentialStore(concurrent_req)
	iteration_data, all_data := my_modules.BenchmarkAPIAsMultiUser(total_req, concurrent_req, _url, "post", headers, nil, payload_generator_callback, nil, nil)

	fmt.Println("bench mark on api finished")
	signup_credentials.CredentialStore_WaitForAppend()

	result := make(map[string]interface{})
	result[_url] = map[string]interface{}{
		"iteration_data": iteration_data,
		"all_data":       all_data,
	}
	return result
}
