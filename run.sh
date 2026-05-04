export GIN_MODE=release
export DISABLE_COLOR=true

# CALCULATE_PAYLOAD_SIZE require more memory
export CALCULATE_PAYLOAD_SIZE=true

# Default 7000 is held by macOS AirPlay Receiver (ControlCenter). Override here.
export SERVER_PORT="${SERVER_PORT:-7100}"

# Buffered logging: open log file once, write entries via 1 MiB bufio.Writer,
# flush on benchmark end. Drops per-call open/seek/close overhead. Set to false
# (or unset) for the default per-call durable mode.
export BUFFERED_LOGGING="${BUFFERED_LOGGING:-true}"

go build -v -o ./benchmark.bin && ./benchmark.bin