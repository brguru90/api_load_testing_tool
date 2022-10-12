package tests

import (
	"apis_load_test/api_requests"
	"apis_load_test/api_requests/user"
	"apis_load_test/benchmark/my_modules"
	"fmt"
)

func TestAsMultiUser() {
	my_modules.LogToJSON(api_requests.LoginAsMultiUser(10000, 500), nil)
	fmt.Println("--> LoginAsMultiUser finished")
	my_modules.LogToJSON(user.GetUserDetailAsMultiUser("http://localhost:8000/api/user/", 10000, 1000), nil)
	// go func() {
	// 	my_modules.LogToJSON(api_requests.SignUp(400, 100), "./log2.json")
	// }()
	my_modules.LogToJSON(user.GetUserDetailAsMultiUser("http://localhost:8000/api/user/?page=1000&limit=20", 10000, 1000), nil)
}
