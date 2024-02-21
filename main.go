package main

import (
	"time"

	"github.com/brguru90/api_load_testing_tool/benchmark"
	"github.com/brguru90/api_load_testing_tool/benchmark/my_modules"
	"github.com/brguru90/api_load_testing_tool/tests"
)

func main() {
	my_modules.HTTPTimeout = time.Minute * 1
	my_modules.LogPath = "./log.json"
	my_modules.DisableLogging = true
	benchmark.RunBenchmark(func() {
		// tests.TestTool()
		tests.TestAsSingleUser()
		tests.TestAsMultiUser()
	})
}
