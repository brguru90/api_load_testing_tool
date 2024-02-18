package tests

import (
	"fmt"
	"sync"

	"github.com/brguru90/api_load_testing_tool/api_requests"
	"github.com/brguru90/api_load_testing_tool/api_requests/user"
	"github.com/brguru90/api_load_testing_tool/benchmark/my_modules"
)

func TestAsMultiUser() {
	var test_wg sync.WaitGroup
	my_modules.LogToJSON(api_requests.LoginAsMultiUser(1000, 100), nil)
	fmt.Println("--> LoginAsMultiUser finished")
	my_modules.LogToJSON(user.GetUserDetailAsMultiUser("http://localhost:8000/api/user/", 1000, 100, false), nil)
	test_wg.Add(1)
	go func() {
		my_modules.LogToJSON(api_requests.SignUp(400, 100), "./log2.json")
		test_wg.Done()
	}()
	my_modules.LogToJSON(user.GetUserDetailAsMultiUser("http://localhost:8000/api/user/", 1000, 100, true), nil)
	test_wg.Wait()
}
