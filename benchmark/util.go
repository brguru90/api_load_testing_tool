package benchmark

import (
	"apis_load_test/benchmark/my_modules"
	"apis_load_test/benchmark/server"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime/pprof"
)

func RunBenchmark(callback func()) {

	var gin_mode string = os.Getenv("GIN_MODE")

	if gin_mode != "release" {
		go func() {
			http.ListenAndServe("localhost:7777", nil)
			// http://localhost:7777/debug/pprof/
			// go tool pprof --http=localhost:8800 /home/justdial/Desktop/test_workspace/api_load_testing_tool/__debug_bin ./profile
		}()
	}

	go func() {
		my_modules.InitBeforeBenchMarkStart()
		if gin_mode != "release" {
			// go tool pprof --http=localhost:8800 /home/justdial/Desktop/test_workspace/api_load_testing_tool/__debug_bin /home/justdial/Desktop/test_workspace/api_load_testing_tool/mem.pprof

			memprof, err := os.Create("mem.pprof")
			callback()
			if err == nil {
				pprof.WriteHeapProfile(memprof)
				memprof.Close()
			}
		} else{
			callback()
		}
		my_modules.OnBenchmarkEnd()
	}()
	server.RunServer(
		os.Getenv("DISABLE_COLOR") == "true",
		gin_mode,
	)
}
