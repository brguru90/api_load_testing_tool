package tests

import (
	"fmt"
	"github.com/brguru90/api_load_testing_tool/api_requests"
	"github.com/brguru90/api_load_testing_tool/api_requests/user"
	"github.com/brguru90/api_load_testing_tool/benchmark/my_modules"
	// "sync"
)

func TestAsMultiUser() {
	// var test_wg sync.WaitGroup
	my_modules.LogToJSON(api_requests.LoginAsMultiUser(10, 2), nil)
	fmt.Println("--> LoginAsMultiUser finished")
	my_modules.LogToJSON(user.GetUserDetailAsMultiUser("http://localhost:8000/api/user/", 10, 2, false), nil)
	// test_wg.Add(1)
	// go func() {
	// 	my_modules.LogToJSON(api_requests.SignUp(40, 10), "./log2.json")
	// 	test_wg.Done()
	// }()
	my_modules.LogToJSON(user.GetUserDetailAsMultiUser("http://localhost:8000/api/user/", 10, 2, true), nil)
	// test_wg.Wait()
}
