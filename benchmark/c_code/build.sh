export LD_LIBRARY_PATH=${PWD}
# export GOFLAGS="-count=1"

rm -rf ./*.so ./*.o ./*.a
# gcc  -static  -Wall curl.c -o libcurl.a  -nostartfiles -g -DCURL_STATICLIB -lssl -lcrypto -lpthread -L$PWD/curl/lib -I$PWD/curl/include 
# gcc -fPIC  -shared curl.c -o libcurl.a -lssl -lcrypto -lpthread -L$PWD/curl/lib -I$PWD/curl/include 
# gcc -fPIC  -Wall -shared api_req.c -o api_req.so -lssl -lcrypto -lpthread -lcurl

gcc -fPIC -shared api_req.c -o api_req.so -lssl -lcrypto -lpthread  
# gcc -fPIC -c *.c
# gcc -shared -o api_req.so api_req.o
go run main.go
# go build -o a.bin && ./a.bin