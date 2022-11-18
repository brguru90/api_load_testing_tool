export GIN_MODE=release
export DISABLE_COLOR=true
# export USING_C_CURL=true

# CALCULATE_PAYLOAD_SIZE require more memory
export CALCULATE_PAYLOAD_SIZE=true

go build -v -o ./benchmark.bin && ./benchmark.bin