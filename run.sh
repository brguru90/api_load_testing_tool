export GIN_MODE=release
export DISABLE_COLOR=true

# CALCULATE_PAYLOAD_SIZE require more memory
export CALCULATE_PAYLOAD_SIZE=false

go build -v -o ./benchmark.bin && ./benchmark.bin