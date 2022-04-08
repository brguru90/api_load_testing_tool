package main

import (
	"apis_load_test/api_requests"
	"apis_load_test/my_modules"
)

func main() {
	my_modules.LogPath = "./log.json"

	// fmt.Printf("%v", api_requests.SignUp())

	my_modules.LogToJSON(api_requests.SignUp())

	// _map:=map[string]interface{}{
	// 	"fname": "guru",
	// 	"lname": "prasad",
	// }

	// fmt.Println(my_modules.LogToJSON(_map))
	// fmt.Println(my_modules.LogToJSON(_map))
	// my_modules.LogToJSON(_map)
	// my_modules.LogToJSON(_map)
}
