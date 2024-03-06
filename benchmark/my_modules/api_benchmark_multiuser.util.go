package my_modules

/*
#cgo CXXFLAGS: -std=gnu++17
#cgo linux pkg-config: libcurl
#cgo linux pkg-config: libuv
#cgo darwin LDFLAGS: -lcurl
#cgo windows LDFLAGS: -lcurl
#include <stdlib.h>
#include <string.h>
#include "api_req.h"
#cgo CFLAGS: -g -Wall
#cgo LDFLAGS: -lssl -lcrypto -lpthread -lm
*/
import "C"
import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httputil"
	"os"
	"reflect"
	"strconv"
	"sync/atomic"
	"unsafe"

	"sync"
	"time"

	"github.com/brguru90/api_load_testing_tool/benchmark/store"

	"github.com/google/uuid"
)

type BenchmarkMetricStreamInfo struct {
	UpdatedAt int64
	Data      interface{}
}

type BenchmarkMetricStruct struct {
	Url           string          `json:"url,omitempty"`
	ProcessUid    string          `json:"process_uid,omitempty"`
	AllData       BenchmarkData   `json:"all_data,omitempty"`
	IterationData []BenchmarkData `json:"iteration_data,omitempty"`
}

type AllIterationData struct {
	additional_details        chan AdditionalAPIDetails
	messages                  chan MessageType
	concurrent_req_start_time time.Time
	concurrent_req_end_time   time.Time
}

var BenchmarkMetricEvent = NewCustomEvent("benchmark_event")

func pushBenchMarkMetrics(data interface{}) {
	t := time.Now().UnixMilli()
	BenchmarkMetricEvent.Emit(BenchmarkMetricStreamInfo{
		Data:      data,
		UpdatedAt: t,
	})
	store.BenchmarkDataStore_Append(data, t)
}

var BenchMarkEnded atomic.Bool

func init() {
	BenchMarkEnded.Store(false)
}

type CGlobalAllIterationData struct {
	all_iteration_data  *[]AllIterationData
	request_ahead_array *[]CreatedAPIRequestFormat
}

var c_global_all_iteration_data map[string]CGlobalAllIterationData = map[string]CGlobalAllIterationData{}
var global_all_iteration_data sync.Mutex

func check_error(err error) {
	if err != nil {
		panic(err)
	}
}

func send_concurrent_request_using_c_curl(main_iteration int64, concurrent_request int64, uuid string) {
	var i int = 0
	var j int64 = 0
	global_all_iteration_data.Lock()
	request_ahead_array := c_global_all_iteration_data[uuid].request_ahead_array
	all_iteration_data := c_global_all_iteration_data[uuid].all_iteration_data
	global_all_iteration_data.Unlock()

	request_input := make([]C.struct_SingleRequestInput, concurrent_request)
	for j = 0; j < concurrent_request; j++ {
		var sub_iteration int64 = (main_iteration * concurrent_request) + j
		req := (*request_ahead_array)[sub_iteration].req

		cookies := ""
		for _, cookie := range req.Cookies() {
			// fmt.Printf("%s\n\n",cookie.String())
			cookies += cookie.String() + "; "
		}

		c_headers := C.malloc(C.size_t(len(req.Header)) * C.sizeof_struct_Headers)
		defer C.free(unsafe.Pointer(c_headers))
		headers_data := (*[1<<30 - 1]C.struct_Headers)(c_headers)
		i = 0
		for name, values := range req.Header {
			for _, value := range values {
				// if name == "cookie" {
				// 	cookies += value + ";"
				// 	continue
				// }
				headers_data[i] = C.struct_Headers{
					header: C.CString(name + ": " + value),
				}
				i++
			}
		}

		// fmt.Printf("cookies=%v\n\n", cookies)

		var body []byte = nil
		if req.Body != http.NoBody {
			Body, err := req.GetBody()
			check_error(err)
			if err == nil && Body != nil {
				body, err = ioutil.ReadAll(Body)
				check_error(err)
			}
		}

		if body != nil {
			request_input[j] = C.struct_SingleRequestInput{
				uid:             C.CString(strconv.FormatInt(((*request_ahead_array)[sub_iteration].uid), 10)),
				url:             C.CString(req.URL.String()),
				method:          C.CString(req.Method),
				headers:         (*C.struct_Headers)(c_headers),
				headers_len:     C.int(len(req.Header)),
				cookies:         C.CString(cookies),
				time_out_in_sec: 5 * 60,
				body:            C.CString(string(body)),
			}
		} else {
			request_input[j] = C.struct_SingleRequestInput{
				uid:             C.CString(strconv.FormatInt(((*request_ahead_array)[sub_iteration].uid), 10)),
				url:             C.CString(req.URL.String()),
				method:          C.CString(req.Method),
				headers:         (*C.struct_Headers)(c_headers),
				headers_len:     C.int(len(req.Header)),
				cookies:         C.CString(cookies),
				time_out_in_sec: 5 * 60,
			}
		}
	}

	bulk_response_data := make([]C.struct_ResponseData, concurrent_request)
	// ram_size_in_GB := float64(C.sysconf(C._SC_PHYS_PAGES)*C.sysconf(C._SC_PAGE_SIZE)) / (1024 * 1024)
	// nor_of_thread := math.Ceil(ram_size_in_GB / 50)
	// // fmt.Println("Nor of threads", nor_of_thread)
	(*all_iteration_data)[main_iteration].concurrent_req_start_time = time.Now()
	// C.send_request_in_concurrently(&(request_input[0]), &(bulk_response_data[0]), C.struct_AdditionalDetails{
	// 	uuid:           C.CString(uuid),
	// 	total_requests: C.int(concurrent_request),
	// 	total_threads:  C.int(nor_of_thread),
	// }, 0)
	C.send_request_in_concurrently(&(request_input[0]), &(bulk_response_data[0]), C.int(concurrent_request), 0)

	(*all_iteration_data)[main_iteration].concurrent_req_end_time = time.Now()

	for j = 0; j < concurrent_request; j++ {
		data := bulk_response_data[j]
		var sub_iteration int64 = (main_iteration * concurrent_request) + j
		additional_detail := AdditionalAPIDetails{
			request_id:                  int64(*data.uid),
			request_sent:                time.UnixMicro(int64(data.before_connect_time_microsec)),
			request_connected:           time.UnixMicro(int64(data.connected_at_microsec)),
			request_receives_first_byte: time.UnixMicro(int64(data.first_byte_at_microsec)),
			request_processed:           time.UnixMicro(int64(data.after_response_time_microsec)),
		}
		api_data := APIData{
			url:     (*request_ahead_array)[sub_iteration].req.URL.String(),
			context: "API response",
			context_data: ContextData{
				status_code: int(data.status_code),
				payload:     (*request_ahead_array)[sub_iteration].payload,
				// json_body:       json_body,
				// body:            body,
				time:            (float64(int64(data.total_time_microsec)) / 1000),
				time_to_connect: (float64(int64(data.connect_time_microsec)) / 1000),
				ttfb:            (float64(int64(data.time_to_first_byte_microsec)) / 1000),
			},
		}
		resp, err := parseHttpResponse(C.GoString(data.response_header), C.GoString(data.response_body), nil, additional_detail.request_id)
		if err != nil {
			resp = nil
		}
		messages := MessageType{
			UID:                  (*request_ahead_array)[sub_iteration].uid,
			Data:                 api_data,
			Time_to_complete_api: float64(int64(data.total_time_microsec)) / 1000,
			Err:                  fmt.Errorf("Curl error code %d", strconv.Itoa(int(data.err_code))),
			Res:                  resp,
		}

		(*all_iteration_data)[main_iteration].additional_details <- additional_detail
		(*all_iteration_data)[main_iteration].messages <- messages
		if (*request_ahead_array)[sub_iteration].req.Body != http.NoBody {
			(*request_ahead_array)[sub_iteration].req.Body.Close()
		}
		(*request_ahead_array)[sub_iteration] = CreatedAPIRequestFormat{
			req:          nil,
			err:          nil,
			payload:      nil,
			request_size: 0,
		}
	}

}

func carray2slice(array *C.int, len int) []C.int {
	var list []C.int
	sliceHeader := (*reflect.SliceHeader)((unsafe.Pointer(&list)))
	sliceHeader.Cap = len
	sliceHeader.Len = len
	sliceHeader.Data = uintptr(unsafe.Pointer(array))
	return list
}

//export response_callback_from_c
func response_callback_from_c(arr_len C.int, response_data []C.struct_ResponseData, uuid *C.char) {
	fmt.Println(C.GoString(uuid))
	// fmt.Printf("%s",uuid)
	fmt.Println(arr_len)
}

func send_concurrent_request(i int64, concurrent_request int64, uuid string) {

	if os.Getenv("USING_C_CURL") == "true" {
		println("Warning:  using libcurl & libuv")
		// runtime.LockOSThread()
		send_concurrent_request_using_c_curl(i, concurrent_request, uuid)
		// runtime.UnlockOSThread()
		return
	}

	global_all_iteration_data.Lock()
	request_ahead_array := c_global_all_iteration_data[uuid].request_ahead_array
	all_iteration_data := c_global_all_iteration_data[uuid].all_iteration_data
	global_all_iteration_data.Unlock()

	var concurrent_req_wg sync.WaitGroup
	var j int64
	concurrent_req_wg.Add(int(concurrent_request))
	(*all_iteration_data)[i].concurrent_req_start_time = time.Now()
	for j = 0; j < concurrent_request; j++ {
		go func(main_iteration int64, sub_iteration int64) {
			defer func() {
				(*request_ahead_array)[sub_iteration] = CreatedAPIRequestFormat{
					req:          nil,
					err:          nil,
					payload:      nil,
					request_size: 0,
				}
			}()

			data, time_to_complete_api, resp, additional_details, err := APIReq(&(*request_ahead_array)[sub_iteration], concurrent_request)
			// fmt.Printf("finish APIReq\n")
			(*all_iteration_data)[main_iteration].additional_details <- additional_details
			(*all_iteration_data)[main_iteration].messages <- MessageType{
				UID:                  (*request_ahead_array)[sub_iteration].uid,
				Data:                 data,
				Time_to_complete_api: time_to_complete_api,
				Err:                  err,
				Res:                  resp,
			}
			concurrent_req_wg.Done()
		}(i, (i*concurrent_request)+j)
	}
	concurrent_req_wg.Wait()
	(*all_iteration_data)[i].concurrent_req_end_time = time.Now()
}

// run the http request concurrently with set of iteration
// collect the metric & calculate the metric for concurrent request & for all iteration
// have the callback to generate payload, intercept request & response
func BenchmarkAPIAsMultiUser(
	total_number_of_request int64,
	concurrent_request int64,
	_url string, method string,
	headers map[string]string,
	payload_obj map[string]interface{},
	payload_generator_callback func(current_iteration int64) map[string]interface{},
	request_interceptor func(req *http.Request, uid int64),
	response_interceptor func(resp *http.Response, uid int64),
) (*[]BenchmarkData, *BenchmarkData) {

	// Todo:
	// need to make an wrapper
	// which distribute set of loads/calculated number of concurrent request's size & payload details
	// to a group of connected benchmark runner client
	// benchmark runner client can be connected by websocket/gRPC protocol
	// later metrics from all client will be sent back to benchmark server
	// & should also support as standalone benchmark tool when no runner configured/connected
	// connection should be established from runner client in private network to single publicly hosted server
	process_uuid := uuid.New().String()

	// var each_iterations_data []BenchmarkData
	number_of_iteration := total_number_of_request / concurrent_request
	each_iterations_data := make([]BenchmarkData, number_of_iteration)

	if total_number_of_request < concurrent_request {
		panic("total_number_of_request<concurrent_request")
	}

	if (total_number_of_request % concurrent_request) != 0 {
		panic("(total_number_of_request % concurrent_request)!=0")
	}

	var rh_concurrent_req_wg sync.WaitGroup
	var rh_iteration_wg sync.WaitGroup

	var i, j int64
	var iterations_start_time, iterations_end_time time.Time

	requests_ahead := make(chan CreatedAPIRequestFormat, number_of_iteration*concurrent_request)

	for i = 0; i < number_of_iteration; i++ {
		rh_iteration_wg.Add(1)
		go func(main_iteration int64) {
			defer rh_iteration_wg.Done()
			rh_concurrent_req_wg.Add(int(concurrent_request))
			for j = 0; j < concurrent_request; j++ {
				go func(sub_iteration int64) {
					defer rh_concurrent_req_wg.Done()
					var api_payload map[string]interface{}
					if payload_obj == nil && payload_generator_callback != nil {
						api_payload = payload_generator_callback(sub_iteration)
					} else {
						api_payload = payload_obj
					}
					requests_ahead <- CreateAPIRequest(_url, method, headers, api_payload, sub_iteration, request_interceptor)
					// fmt.Printf("requests_ahead %d->%d<%d\n",(main_iteration * concurrent_request) + j,len(requests_ahead),number_of_iteration*concurrent_request)
				}((main_iteration * concurrent_request) + j)
			}
			rh_concurrent_req_wg.Wait()

		}(i)
		rh_iteration_wg.Wait()
	}
	// fmt.Printf("requests_ahead %d<%d\n",len(requests_ahead),number_of_iteration*concurrent_request)
	close(requests_ahead)

	// var request_ahead_array []CreatedAPIRequestFormat
	request_ahead_array := make([]CreatedAPIRequestFormat, len(requests_ahead))
	ri := 0
	for request_ahead := range requests_ahead {
		request_ahead_array[ri] = request_ahead
		// request_ahead_array = append(request_ahead_array, request_ahead)
		ri++
	}
	requests_ahead = nil

	// all_iteration_data := []AllIterationData{}
	all_iteration_data := make([]AllIterationData, number_of_iteration)
	for i = 0; i < number_of_iteration; i++ {
		all_iteration_data[i] = AllIterationData{
			messages:           make(chan MessageType),
			additional_details: make(chan AdditionalAPIDetails, concurrent_request),
		}
		// all_iteration_data = append(all_iteration_data, AllIterationData{
		// 	messages:           make(chan MessageType),
		// 	additional_details: make(chan AdditionalAPIDetails, concurrent_request),
		// })
	}

	global_all_iteration_data.Lock()
	c_global_all_iteration_data[process_uuid] = CGlobalAllIterationData{
		all_iteration_data:  &all_iteration_data,
		request_ahead_array: &request_ahead_array,
	}
	global_all_iteration_data.Unlock()

	var send_req_wg sync.WaitGroup
	go func() {
		defer send_req_wg.Done()
		send_req_wg.Add(1)
		iterations_start_time = time.Now()
		for i = 0; i < number_of_iteration; i++ {
			messages := &(all_iteration_data[i].messages)
			// additional_details := &(all_iteration_data[i].additional_details)
			fmt.Printf("url=%s,i=%v\n", _url, i)
			send_concurrent_request(i, concurrent_request, process_uuid)

			close(*messages)
			close(all_iteration_data[i].additional_details)
		}
		iterations_end_time = time.Now()
	}()

	var all_iteration_data_collection_wg sync.WaitGroup
	all_iteration_data_collection_wg.Add(len(all_iteration_data))
	for i := 0; i < len(all_iteration_data); i++ {
		go func(cur_iteration_data *AllIterationData, _i int64) {
			defer all_iteration_data_collection_wg.Done()
			messages := &cur_iteration_data.messages
			additional_details := &cur_iteration_data.additional_details

			var avg_time_to_complete_api, avg_time_to_connect_api, avg_time_to_receive_first_byte_api float64

			avg_time_to_complete_api = 0
			avg_time_to_connect_api = 0
			avg_time_to_receive_first_byte_api = 0
			min_time_to_complete_api := math.Inf(1)
			max_time_to_complete_api := 0.0
			status_codes := make(map[int]int64)
			// looping through channel data, whenever go routine finishes execution
			// fmt.Println("loop through msgs")

			var avg_request_payload_size float64 = 0
			var avg_response_payload_size float64 = 0

			for message := range *messages {

				if message.Res != nil && response_interceptor != nil {
					response_interceptor(message.Res, message.UID)
				}

				response_size := 0
				if message.Res != nil {
					if ShouldDumpRequestAndResponse {
						respDump, err := httputil.DumpResponse(message.Res, true)
						if err == nil {
							response_size = len(respDump)
							respDump = nil
						}
					}
					message.Res.Body.Close()
				}

				avg_response_payload_size += float64(response_size)

				avg_time_to_complete_api += message.Time_to_complete_api
				avg_time_to_connect_api += message.Data.context_data.time_to_connect
				avg_time_to_receive_first_byte_api += message.Data.context_data.ttfb
				var cur_status_code int = message.Data.context_data.status_code
				status_codes[cur_status_code] += 1

				if float64(message.Time_to_complete_api) < min_time_to_complete_api {
					min_time_to_complete_api = float64(message.Time_to_complete_api)
				}

				if float64(message.Time_to_complete_api) > max_time_to_complete_api {
					max_time_to_complete_api = float64(message.Time_to_complete_api)
				}
			}

			concurrent_req_start_time := cur_iteration_data.concurrent_req_start_time
			concurrent_req_end_time := cur_iteration_data.concurrent_req_end_time

			avg_time_to_complete_api = avg_time_to_complete_api / float64(concurrent_request)
			avg_time_to_connect_api = avg_time_to_connect_api / float64(concurrent_request)
			avg_time_to_receive_first_byte_api = avg_time_to_receive_first_byte_api / float64(concurrent_request)

			status_code_in_percentage := make(map[int]float64)
			for status_code, occurrence := range status_codes {
				status_code_in_percentage[status_code] = (float64(occurrence) / float64(concurrent_request)) * 100
			}

			var additional_details_arr []AdditionalAPIDetails
			for additional_detail := range *additional_details {
				additional_details_arr = append(additional_details_arr, additional_detail)
				avg_request_payload_size += float64(additional_detail.request_payload_size)
			}
			if avg_request_payload_size >= 0 {
				avg_request_payload_size = avg_request_payload_size / float64(len(additional_details_arr))
			}

			if avg_response_payload_size >= 0 {
				avg_response_payload_size = avg_response_payload_size / float64(len(additional_details_arr))
			}

			// concurrent_req_start_time & end time will not be perfectly accurate
			// better to take
			// min(additional_detail.request_sent)
			// max(additional_detail.request_processed)
			track_iteration_start_time := concurrent_req_start_time
			track_iteration_end_time := concurrent_req_end_time
			if len(additional_details_arr) > 0 {

				track_iteration_start_time := additional_details_arr[0].request_sent
				track_iteration_end_time := additional_details_arr[0].request_processed
				for _, val := range additional_details_arr {
					if val.request_sent.Before(track_iteration_start_time) {
						track_iteration_start_time = val.request_sent
					}
					if val.request_processed.After(track_iteration_end_time) {
						track_iteration_end_time = val.request_processed
					}
				}
			}
			track_iteration_time := track_iteration_start_time
			prev_iteration_time := track_iteration_start_time.Add(time.Nanosecond * -1)
			page_size := 8
			k_nano := (track_iteration_end_time.Sub(track_iteration_start_time) / time.Duration(page_size)).Nanoseconds()
			time_frame_size := time.Duration(math.Round(float64(k_nano/100)) * 100) // rounding up to nearest decimal
			if time_frame_size < (100 * time.Nanosecond) {
				time_frame_size = 100 * time.Nanosecond
			}
			// time_frame_size := time.Millisecond * 200 // decrease to increase the detail of iteration data
			if track_iteration_time.Add(time_frame_size).After(track_iteration_end_time) {
				track_iteration_time = track_iteration_end_time
			} else {
				// add 1 sec gap, if diff between track_iteration_time & track_iteration_end_time is greater than 1 sec
				track_iteration_time = track_iteration_time.Add(time_frame_size)
			}
			var per_second_metrics []BenchMarkPerSecondCount
			// var _request_sent_in_sec_avg, _request_connected_in_sec_avg, _request_processed_in_sec_avg int64
			if len(additional_details_arr) > 0 {
				for {
					// fmt.Printf("prev_iteration_time=%v,track_iteration_time=%v\n",prev_iteration_time,track_iteration_time)
					var _request_sent_in_sec, _request_connected_in_sec, _request_received_first_byte_in_sec, _request_processed_in_sec int64
					var request_payload_size float64 = 0
					var response_payload_size float64 = 0

					// filtering each request in the constructed time boundary
					for _, additional_detail := range additional_details_arr {

						if additional_detail.request_sent.After(prev_iteration_time) && (additional_detail.request_sent.Equal(track_iteration_time) || additional_detail.request_sent.Before(track_iteration_time)) {
							_request_sent_in_sec++
							request_payload_size += float64(additional_detail.request_payload_size)
						}
						if additional_detail.request_connected.After(prev_iteration_time) && (additional_detail.request_connected.Equal(track_iteration_time) || additional_detail.request_connected.Before(track_iteration_time)) {
							_request_connected_in_sec++
						}
						if additional_detail.request_receives_first_byte.After(prev_iteration_time) && (additional_detail.request_receives_first_byte.Equal(track_iteration_time) || additional_detail.request_receives_first_byte.Before(track_iteration_time)) {
							_request_received_first_byte_in_sec++
						}
						if additional_detail.request_processed.After(prev_iteration_time) && (additional_detail.request_processed.Equal(track_iteration_time) || additional_detail.request_processed.Before(track_iteration_time)) {
							_request_processed_in_sec++
							response_payload_size += float64(additional_detail.response_payload_size)
						}
					}
					per_second_metrics = append(per_second_metrics, BenchMarkPerSecondCount{
						From_time_duration:                   prev_iteration_time,
						To_time_duration:                     track_iteration_time,
						Request_sent:                         _request_sent_in_sec,
						Request_connected:                    _request_connected_in_sec,
						Request_processed:                    _request_processed_in_sec,
						Request_receives_first_byte:          _request_received_first_byte_in_sec,
						Total_request_payload_size_in_bytes:  request_payload_size,
						Total_response_payload_size_in_bytes: response_payload_size,
					})

					// iterate over elapsed time
					if track_iteration_time.After(track_iteration_end_time) || track_iteration_time.Equal(track_iteration_end_time) {
						break
					}
					prev_iteration_time = track_iteration_time
					if track_iteration_time.Add(time_frame_size).After(track_iteration_end_time) {
						track_iteration_time = track_iteration_end_time
					} else {
						track_iteration_time = track_iteration_time.Add(time_frame_size)
					}
				}
			}

			temp_data := BenchmarkData{
				Url:                             _url,
				Status_code_in_percentage:       status_code_in_percentage,
				Status_codes:                    status_codes,
				Concurrent_request:              concurrent_request,
				Avg_time_to_connect_api:         int64(avg_time_to_connect_api),
				Avg_time_to_receive_first_byte:  int64(avg_time_to_receive_first_byte_api),
				Avg_time_to_complete_api:        int64(avg_time_to_complete_api),
				Min_time_to_complete_api:        int64(min_time_to_complete_api),
				Max_time_to_complete_api:        int64(max_time_to_complete_api),
				Average_request_payload_size:    avg_request_payload_size,
				Average_response_payload_size:   avg_response_payload_size,
				Total_time_to_complete_all_apis: concurrent_req_end_time.Sub(concurrent_req_start_time).Milliseconds(),
				Benchmark_per_second_metric:     per_second_metrics,
				IterationID:                     _i,
			}
			result := BenchmarkMetricStruct{
				Url:           temp_data.Url,
				IterationData: []BenchmarkData{temp_data},
				ProcessUid:    process_uuid,
			}
			pushBenchMarkMetrics(result)
			// each_iterations_data = append(each_iterations_data, temp_data)
			each_iterations_data[_i] = temp_data
		}(&all_iteration_data[i], int64(i))
	}

	all_iteration_data_collection_wg.Wait()

	var total_time_to_complete_api, avg_time_to_complete_api, avg_time_to_connect_api int64

	avg_time_to_complete_api = 0
	avg_time_to_connect_api = 0
	min_time_to_complete_api := math.Inf(1)
	max_time_to_complete_api := 0.0
	avg_request_payload_size := 0.0
	avg_response_payload_size := 0.0
	status_codes := make(map[int]int64)
	for _, _each_iterations_data := range each_iterations_data {
		for status_code, occurrence := range _each_iterations_data.Status_codes {
			status_codes[status_code] += occurrence
		}
		if float64(_each_iterations_data.Min_time_to_complete_api) < min_time_to_complete_api {
			min_time_to_complete_api = float64(_each_iterations_data.Min_time_to_complete_api)
		}

		if float64(_each_iterations_data.Max_time_to_complete_api) > max_time_to_complete_api {
			max_time_to_complete_api = float64(_each_iterations_data.Max_time_to_complete_api)
		}
		avg_time_to_complete_api += _each_iterations_data.Avg_time_to_complete_api
		avg_time_to_connect_api += _each_iterations_data.Avg_time_to_connect_api
		avg_request_payload_size += _each_iterations_data.Average_request_payload_size
		avg_response_payload_size += _each_iterations_data.Average_response_payload_size
		total_time_to_complete_api += _each_iterations_data.Total_time_to_complete_all_apis
	}

	status_codes_in_perc := make(map[int]float64)
	for status_code, occurrence := range status_codes {
		status_codes_in_perc[status_code] = (float64(occurrence) / float64(total_number_of_request)) * 100.0
	}

	send_req_wg.Wait()

	temp_data := BenchmarkData{
		Url:                                              _url,
		Status_codes:                                     status_codes,
		Status_code_in_percentage:                        status_codes_in_perc,
		Total_number_of_request:                          total_number_of_request,
		Concurrent_request:                               concurrent_request,
		Avg_time_to_connect_api_in_sec:                   (float64(avg_time_to_connect_api) / float64(number_of_iteration)) / 1000,
		Min_time_to_complete_api_in_sec:                  min_time_to_complete_api / 1000.0,
		Max_time_to_complete_api_in_sec:                  max_time_to_complete_api / 1000.0,
		Avg_time_to_complete_api_in_sec:                  (float64(avg_time_to_complete_api) / float64(number_of_iteration)) / 1000,
		Average_request_payload_size_in_all_iteration:    avg_request_payload_size / float64(number_of_iteration),
		Average_response_payload_size_in_all_iteration:   avg_response_payload_size / float64(number_of_iteration),
		Total_time_to_complete_all_apis_iteration_in_sec: float64(total_time_to_complete_api) / 1000.0,
		Total_operation_time_in_sec:                      float64(iterations_end_time.Sub(iterations_start_time).Milliseconds()) / 1000.0,
		IterationID:                                      -1,
	}
	result := BenchmarkMetricStruct{
		Url:           temp_data.Url,
		AllData:       temp_data,
		IterationData: nil,
		ProcessUid:    process_uuid,
	}
	pushBenchMarkMetrics(result)
	// runtime.GC()

	global_all_iteration_data.Lock()
	c_global_all_iteration_data[process_uuid] = CGlobalAllIterationData{}
	global_all_iteration_data.Unlock()

	return &each_iterations_data, &temp_data
}

func OnBenchmarkEnd() {
	BenchmarkMetricEvent.Emit(nil)
	BenchmarkMetricEvent.Dispose()
	BenchMarkEnded.Store(true)
	fmt.Println("*** Benchmark completed ***")
}

func InitBeforeBenchMarkStart() {
	f := func(data []interface{}, cur_data interface{}) []interface{} {
		_matched := false
		for key, dt := range data {
			temp := dt.(BenchmarkMetricStruct)
			if temp.Url == cur_data.(BenchmarkMetricStruct).Url && temp.ProcessUid == cur_data.(BenchmarkMetricStruct).ProcessUid {
				if cur_data.(BenchmarkMetricStruct).IterationData != nil {
					temp.IterationData = append(temp.IterationData, cur_data.(BenchmarkMetricStruct).IterationData...)
				} else {
					temp.AllData = cur_data.(BenchmarkMetricStruct).AllData
				}
				data[key] = temp
				_matched = true
			}
		}
		if !_matched {
			data = append(data, cur_data)
		}
		return data
	}
	store.BenchmarkDataStore_ManualAppendFromQ(&f)
}
