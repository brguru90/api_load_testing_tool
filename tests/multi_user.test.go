package tests

import (
	"apis_load_test/api_requests"
	// "apis_load_test/api_requests/user"
	"apis_load_test/my_modules"
	"fmt"
)

func TestAsMultiUser() {
	my_modules.LogToJSON(api_requests.LoginAsMultiUser(1000, 100), nil)
	fmt.Println("--> LoginAsMultiUser finished")
	// my_modules.LogToJSON(user.GetUserDetailAsMultiUser("http://localhost:8000/api/user/", 1000, 100), nil)
	// go func() {
	// 	my_modules.LogToJSON(api_requests.SignUp(20000, 100), "./log2.json")
	// }()
	// my_modules.LogToJSON(user.GetUserDetailAsMultiUser("http://localhost:8000/api/user/?page=1000&limit=20", 1000, 100), nil)
}
