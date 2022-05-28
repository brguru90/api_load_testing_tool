package tests

import (
	"apis_load_test/api_requests"
	"apis_load_test/api_requests/user"
	"apis_load_test/my_modules"
	"fmt"
)

func TestAsMultiUser() {
	my_modules.LogToJSON(api_requests.LoginAsMultiUser(), nil)
	fmt.Println("--> LoginAsMultiUser finished")
	my_modules.LogToJSON(user.GetUserDetailAsMultiUser("http://localhost:8000/api/user/"), nil)
	go func() {
		my_modules.LogToJSON(api_requests.SignUp(2000, 1000), "./log2.json")
	}()
	my_modules.LogToJSON(user.GetUserDetailAsMultiUser("http://localhost:8000/api/user/?page=1000&limit=20"), nil)
}
