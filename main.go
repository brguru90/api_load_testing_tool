package main

import (
	"github.com/brguru90/api_load_testing_tool/benchmark"
	"github.com/brguru90/api_load_testing_tool/benchmark/my_modules"
	"github.com/brguru90/api_load_testing_tool/tests"
	"time"
)

func main() {
	my_modules.HTTPTimeout = time.Minute * 1
	my_modules.LogPath = "./log.json"
	my_modules.DisableLogging=true
	benchmark.RunBenchmark(func() {
		tests.TestTool()
		tests.TestAsSingleUser()
		tests.TestAsMultiUser()
	})
}
