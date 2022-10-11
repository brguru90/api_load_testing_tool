package main

import (
	"apis_load_test/benchmark"
	"apis_load_test/benchmark/my_modules"
	"apis_load_test/tests"
	"time"
)

func main() {
	my_modules.HTTPTimeout = time.Minute * 1
	my_modules.LogPath = "./log.json"
	my_modules.DisableLogging=true
	benchmark.RunBenchmark(func() {
		tests.TestAsSingleUser()
		tests.TestAsMultiUser()
	})
}
