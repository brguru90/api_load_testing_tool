package benchmark

import (
	"apis_load_test/benchmark/my_modules"
	"apis_load_test/benchmark/server"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

func RunBenchmark(callback func()) {

	var gin_mode string = os.Getenv("GIN_MODE")
	var enable_profiling bool = os.Getenv("PROFILING")=="true"
	my_modules.ShouldDumpRequestAndResponse=os.Getenv("CALCULATE_PAYLOAD_SIZE")=="true"

	if enable_profiling{
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT)
		memprof, memErr := os.Create("mem.pprof")
		go func() {
			select {
			case sig := <-sigs:			
				fmt.Printf("Got %s signal. Collecting profile data\n", sig)	
				if memprof!=nil && memErr==nil{
					pprof.WriteHeapProfile(memprof)
					memprof.Close()
				}
				fmt.Printf("Done Collecting profile data\n", sig)	
				os.Exit(1)
			}
		}()

		go func() {
			http.ListenAndServe("localhost:7777", nil)
			// http://localhost:7777/debug/pprof/
			// go tool pprof --http=localhost:8800 /home/guruprasad/Desktop/test_workspace/api_load_testing_tool/__debug_bin ./profile
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
