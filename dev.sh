export GIN_MODE=debug
export DISABLE_COLOR=false
export CALCULATE_PAYLOAD_SIZE=true
export USING_C_CURL=true

go build -v -o ./benchmark.bin && gdb ./benchmark.bin