package benchmark

import (
	"apis_load_test/benchmark/my_modules"
	"apis_load_test/benchmark/server"
	"net/http"
	_ "net/http/pprof"
	"os"
)

func RunBenchmark(callback func()) {

	var gin_mode string = os.Getenv("GIN_MODE")

	if gin_mode != "release" {
		go func() {
			http.ListenAndServe("localhost:7777", nil)
			// http://localhost:7777/debug/pprof/
		}()
	}

	go func() {
		my_modules.InitBeforeBenchMarkStart()
		callback()
		my_modules.OnBenchmarkEnd()
	}()
	server.RunServer(
		os.Getenv("DISABLE_COLOR") == "true",
		gin_mode,
	)
}
