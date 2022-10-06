package benchmark

import (
	"apis_load_test/benchmark/my_modules"
	"apis_load_test/benchmark/server"
)

func RunBenchmark(callback func()) {
	go func() {
		my_modules.InitBeforeBenchMarkStart()
		callback()
		my_modules.OnBenchmarkEnd()
	}()
	server.RunServer()
}
