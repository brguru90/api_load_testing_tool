package main

import (
	"apis_load_test/my_modules"

	"apis_load_test/tests"
)


func main() {
	my_modules.LogPath = "./log.json"
	// tests.TestAsSingleUser()
	tests.TestAsMultiUser()
}
