package tests

import (
	"apis_load_test/api_requests"
	"apis_load_test/benchmark/my_modules"
)

func TestAsSingleUser() {
	my_modules.LogToJSON(api_requests.SignUp(400, 200), nil)
}
