package tests

import (
	"apis_load_test/api_requests"
	"apis_load_test/my_modules"
)

func TestAsSingleUser() {
	my_modules.LogToJSON(api_requests.SignUp(10000, 1000), nil)
}
