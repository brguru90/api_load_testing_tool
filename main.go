package main

import (
	"apis_load_test/my_modules"
	"apis_load_test/server"
	"time"

	"apis_load_test/tests"
)

func main() {
	my_modules.HTTPTimeout = time.Minute * 1
	my_modules.LogPath = "./log.json"
	go func() {
		tests.TestAsSingleUser()
		tests.TestAsMultiUser()
		close(my_modules.BenchmarkMetricStream)
	}()
	server.RunServer()
}
