package main

import (
	"apis_load_test/tests"
	"apis_load_test/my_modules"	
)

func main() {
	my_modules.LogPath = "./log.json"

	// tests.TestAsSingleUser()
	tests.TestAsMultiUser()
}
