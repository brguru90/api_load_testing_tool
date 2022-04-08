package api_requests

import (
	"apis_load_test/my_modules"
	"fmt"
)

func SignUp() map[string]interface{} {
	_url := "http://localhost:8000/api/sign_up/"

	payload_obj := map[string]interface{}{
		"email":       my_modules.RandomString(100) + "@gmail.com",
		"name":        my_modules.RandomString(20),
		"description": my_modules.RandomString(100),
	}
	headers := map[string]string{
		"Content-Type": "application/json",
	}
	data, _, _, err := my_modules.APIReq(_url, "post", headers, payload_obj)
	if err == nil {
		fmt.Println(data)
		return data
	}
	panic(err)
}
