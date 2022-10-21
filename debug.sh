
export GIN_MODE=debug
export DISABLE_COLOR=false
export PROFILING=true

PID_LIST=""

function beforeExit() {
    echo;
    echo "statrting Profiling...";
    trap - SIGINT
    kill -s SIGINT $PID_LIST    
    go tool pprof --http=localhost:8800 ./debug.bin ./mem.pprof
    echo "Benchmark Exited";
}


rm -rf ./mem.pprof
go build -v -o ./debug.bin
./debug.bin &

PID_LIST+=" $!";
echo "PIDs=$PID_LIST"
trap beforeExit SIGINT
wait $PID_LIST

echo "Exited";
