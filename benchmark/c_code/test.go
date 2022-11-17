package main

/*
#cgo linux pkg-config: libcurl
#cgo darwin LDFLAGS: -lcurl
#cgo windows LDFLAGS: -lcurl
#include <stdlib.h>
#include <string.h>
#include "api_req.h"
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -L${SRCDIR} api_req.so
*/
import "C"
import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"unsafe"
)

// https://stackoverflow.com/questions/65562170/go-pointer-stored-into-non-go-memory
// https://stackoverflow.com/questions/41492071/how-do-i-convert-a-go-array-of-strings-to-a-c-array-of-strings

var gchan chan int

//export goCallback
func goCallback(myid C.int) {
	fmt.Printf("go api call from c %d\n", int(myid))
	// gchan <- int(myid)
	// if int(myid)==-1{
	// 	close(gchan)
	// }
}

func callCAPI(s string) {
	cs := C.CString(s)
	defer C.free(unsafe.Pointer(cs))
	C.run_bulk_api_request(cs)
}

func check_error(err error) {
	if err != nil {
		panic(err)
	}
}


func parseHttpResponse(header string, _body string, req *http.Request) (*http.Response, error) {
	skip_string := "Transfer-Encoding: chunked\r\n"
	pos := strings.Index(header, skip_string)
	if pos >= 0 {
		pos_end := pos + len(skip_string)
		header = header[0:pos] + header[pos_end:]
	}
	r := header + _body
	body := bytes.NewBuffer([]byte(r))
	prefix := make([]byte, 7)
	n, err := io.ReadFull(body, prefix)
	if err != nil {
		panic("handler err")
	}
	// fmt.Println(n, err, string(prefix))
	if string(prefix[:n]) == "HTTP/2 " {
		// fix HTTP/2 proto
		return http.ReadResponse(bufio.NewReader(io.MultiReader(bytes.NewBufferString("HTTP/2.0 "), body)), req)
	} else {
		// other proto
		// return http.ReadResponse(bufio.NewReader(bytes.NewBuffer([]byte(r))), req)
		return http.ReadResponse(bufio.NewReader(io.MultiReader(bytes.NewBuffer(prefix[:n]), body)), req)
	}
}

func call_api() {
	total_requests:=2
	// url := "http://localhost:8000/api/hello/1?query=text"
	url := "http://guruinfo.epizy.com/edu.php"
	// url:="https://jsonplaceholder.typicode.com/posts"
	// url:="https://google.com/"
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte("some string")))
	req.Header.Add("some-header", "its value")
	req.Header.Add("some-header2", "its value2")
	check_error(err)
	c_url := C.CString(url)
	defer C.free(unsafe.Pointer(c_url))

	request_input := make([]C.struct_SingleRequestInput, total_requests)

	c_headers := C.malloc(C.size_t(len(req.Header)) * C.sizeof_struct_Headers)
	defer C.free(unsafe.Pointer(c_headers))
	headers_data := (*[1<<30 - 1]C.struct_Headers)(c_headers)
	var i int = 0
	for name, values := range req.Header {
		for _, value := range values {
			headers_data[i] = C.struct_Headers{
				header: C.CString(name + ": " + value),
			}
			i++
		}
	}

	Body, _ := req.GetBody()
	body, _ := ioutil.ReadAll(Body)

	for i = 0; i < total_requests; i++ {
		request_input[i] = C.struct_SingleRequestInput{
			url:         C.CString("http://localhost:8000/api/test/"+strconv.Itoa(i)),
			method:      C.CString(req.Method),
			headers:     (*C.struct_Headers)(c_headers),
			headers_len: C.int(len(req.Header)),
			body:        C.CString(string(body)),
		}
	}

	bulk_response_data := make([]C.struct_ResponseData, total_requests)
	ram_size_in_GB:=float64(C.sysconf(C._SC_PHYS_PAGES)*C.sysconf(C._SC_PAGE_SIZE))/(1024*1024)
	nor_of_thread:=math.Ceil(ram_size_in_GB/70)
	fmt.Println("Nor of threads",nor_of_thread);
	C.send_request_in_concurrently(&(request_input[0]), &(bulk_response_data[0]), C.int(total_requests), C.int(nor_of_thread),2)

	// for i = 0; i < total_requests; i++ {
	// 	// fmt.Println(i,C.GoString(bulk_response_data[i].response_body))
	// 	fmt.Println(int(bulk_response_data[i].status_code),C.GoString(bulk_response_data[i].response_body))
	// }

	// fmt.Println(int(response_data.status_code))

	// C.send_raw_request(&(request_input[0]), &response_data, 0)
	// fmt.Println(int(response_data.status_code))
	// resp, err := parseHttpResponse(C.GoString(response_data.response_header), C.GoString(response_data.response_body), nil)
	// body2, err := ioutil.ReadAll(resp.Body)
	// fmt.Println(string(body2), err)
}

func main() {
	gchan = make(chan int)
	// str1 := "hi how are you\n"
	// callCAPI(str1)

	call_api()

	// for v:=range gchan{
	// 	println(v);
	// }
	// fmt.Println("channel closed")
}
