package main

import (
	"apis_load_test/api_requests"
	"apis_load_test/my_modules"
	"fmt"
)

func main()  {
	fmt.Println(api_requests.SignUp())

	my_modules.LogPath="./log.json"


	// _map:=map[string]interface{}{
	// 	"fname": "guru",
	// 	"lname": "prasad",
	// }

	// fmt.Println(my_modules.LogToJSON(_map))
	// fmt.Println(my_modules.LogToJSON(_map))
	// my_modules.LogToJSON(_map)
	// my_modules.LogToJSON(_map)
}