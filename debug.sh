
export GIN_MODE=debug
export DISABLE_COLOR=false
export PROFILING=true
export CALCULATE_PAYLOAD_SIZE=false
export USING_C_CURL=true
# export GOGC=off

PID_LIST=""

function beforeExit() {
    echo;
    echo "statrting Profiling...";
    trap - SIGINT
    kill -s SIGINT $PID_LIST    
    sleep 2
    # analyse generated profiling data
    go tool pprof --http=localhost:8800 ./debug.bin ./mem.pprof
    echo "Benchmark Exited";
}


rm -rf ./mem.pprof ./cpu.pprof
go build -v -o ./debug.bin
./debug.bin &

PID_LIST+=" $!";

echo "execution started, press ctrl+c to end execution & view profiling data"

echo "PIDs=$PID_LIST"
trap beforeExit SIGINT
wait $PID_LIST

echo "Exited";
