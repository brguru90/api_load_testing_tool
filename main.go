package main

import (
	"apis_load_test/my_modules"
	"apis_load_test/server"
	"apis_load_test/tests"
	"time"
)

func main() {
	my_modules.HTTPTimeout = time.Minute * 1
	my_modules.LogPath = "./log.json"
	go func() {
		my_modules.InitBeforeBenchMarkStart()
		tests.TestAsSingleUser()
		tests.TestAsMultiUser()
		my_modules.OnBenchmarkEnd()
	}()
	server.RunServer()
}
