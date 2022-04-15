package tests

import (
	"apis_load_test/api_requests"
	"apis_load_test/api_requests/user"
	"apis_load_test/my_modules"
	"fmt"
)

func TestAsMultiUser() {
	my_modules.LogToJSON(api_requests.LoginAsMultiUser())
	fmt.Println("--> LoginAsMultiUser finished")
	my_modules.LogToJSON(user.GetUserDetailAsMultiUser("http://localhost:8000/api/user/"))
	my_modules.LogToJSON(user.GetUserDetailAsMultiUser("http://localhost:8000/api/user/?page=1000&limit=20"))
}

