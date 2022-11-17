export LD_LIBRARY_PATH=${PWD}/benchmark/c_code/
export GIN_MODE=debug
export DISABLE_COLOR=false
export CALCULATE_PAYLOAD_SIZE=true

rm -rf ./benchmark/c_code/*.so ./benchmark/c_code/*.o ./benchmark/c_code/*.a

gcc -fPIC -shared ./benchmark/c_code/api_req.c -o api_req.so -lssl -lcrypto -lpthread  

go run main.go