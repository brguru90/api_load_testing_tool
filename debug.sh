
export GIN_MODE=debug
export DISABLE_COLOR=false
export PROFILING=true
export CALCULATE_PAYLOAD_SIZE=false
# export GOGC=off

# Default 7000 is held by macOS AirPlay Receiver (ControlCenter). Override here.
export SERVER_PORT="${SERVER_PORT:-7100}"

PID_LIST=""

function beforeExit() {
    echo;
    echo "statrting Profiling...";
    trap - SIGINT
    kill -s SIGINT $PID_LIST    
    sleep 2
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
