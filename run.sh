export GIN_MODE=release
export DISABLE_COLOR=true

go build -v -o ./benchmark.bin && ./benchmark.bin