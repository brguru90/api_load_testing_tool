package tests

import (
	"github.com/brguru90/api_load_testing_tool/api_requests"
	"github.com/brguru90/api_load_testing_tool/benchmark/my_modules"
)

func TestAsSingleUser() {
	my_modules.LogToJSON(api_requests.SignUp(100, 10), nil)
}
