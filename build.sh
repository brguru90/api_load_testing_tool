export LD_LIBRARY_PATH=$PWD/benchmark/my_modules/

rm -rf ./*.so
gcc -fPIC -shared benchmark/my_modules/api_req.c -o api_req.so
# gcc -fPIC -c *.c
# gcc -shared -o api_req.so api_req.o
go run main.go
# go build -o a.bin && ./a.bin