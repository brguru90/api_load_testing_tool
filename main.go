package main

import (
	"apis_load_test/my_modules"
	"time"

	"apis_load_test/tests"
)


func main() {
	my_modules.HTTPTimeout=time.Minute*10
	my_modules.LogPath = "./log.json"
	tests.TestAsSingleUser()
	tests.TestAsMultiUser()	
}
