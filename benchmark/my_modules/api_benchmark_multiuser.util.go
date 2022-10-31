package my_modules

import (
	"fmt"
	"math"
	"net/http"
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

var BenchMarkEnded bool = false

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

	var each_iterations_data []BenchmarkData
	number_of_iteration := total_number_of_request / concurrent_request

	if total_number_of_request < concurrent_request {
		panic("total_number_of_request<concurrent_request")
	}

	if (total_number_of_request % concurrent_request) != 0 {
		panic("(total_number_of_request % concurrent_request)!=0")
	}

	var concurrent_req_wg, rh_concurrent_req_wg sync.WaitGroup
	var main_iter, rh_iteration_wg sync.WaitGroup

	var i, j, total_time_to_complete_api, avg_time_to_complete_api, avg_time_to_connect_api, avg_time_to_receive_first_byte_api int64
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

	var request_ahead_array []CreatedAPIRequestFormat
	for request_ahead := range requests_ahead {
		request_ahead_array = append(request_ahead_array, request_ahead)
	}
	requests_ahead = nil

	all_iteration_data := []AllIterationData{}
	for i = 0; i < number_of_iteration; i++ {
		all_iteration_data = append(all_iteration_data, AllIterationData{
			messages:           make(chan MessageType),
			additional_details: make(chan AdditionalAPIDetails, concurrent_request),
		})
	}
	main_iter.Add(1)
	go func() {
		defer main_iter.Done()
		iterations_start_time = time.Now()
		for i = 0; i < number_of_iteration; i++ {
			messages := &(all_iteration_data[i].messages)
			additional_details := &(all_iteration_data[i].additional_details)
			all_iteration_data[i].concurrent_req_start_time = time.Now()
			fmt.Printf("url=%s,i=%v\n", _url, i)

			// run whole parallel request routine in background
			// so that i can loop through channel data later
			// spin up all request parallally & wait for all those to finish
			concurrent_req_wg.Add(int(concurrent_request))
			// todo:
			// here benchmark server will distribute request to runner client
			// wait for all to complete
			// & then collect metrics
			for j = 0; j < concurrent_request; j++ {
				// fmt.Printf("%v-%v\n", i, j)
				go func(sub_iteration int64) {
					defer func() {
						request_ahead_array[sub_iteration] = CreatedAPIRequestFormat{
							req:          nil,
							err:          nil,
							payload:      nil,
							request_size: 0,
						}
						concurrent_req_wg.Done()
					}()

					data, time_to_complete_api, _, err := APIReq(&request_ahead_array[sub_iteration], response_interceptor, *additional_details)
					// fmt.Printf("finish APIReq\n")
					*messages <- MessageType{
						Data:                 data,
						Time_to_complete_api: time_to_complete_api,
						Err:                  err,
					}
					// fmt.Printf("finish channel\n")
				}((i * concurrent_request) + j)
			}
			concurrent_req_wg.Wait()
			// fmt.Println("all the parallel request finished")
			for ;len(*messages)>0;{
				time.Sleep(time.Millisecond*100)
			}
			close(*messages)
			for ;len(*additional_details)>0;{
				time.Sleep(time.Millisecond*100)
			}
			close(*additional_details)
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

			avg_time_to_complete_api = 0
			avg_time_to_connect_api = 0
			avg_time_to_receive_first_byte_api = 0
			min_time_to_complete_api := math.Inf(1)
			max_time_to_complete_api := 0.0
			status_codes := make(map[int]int64)
			// looping through channel data, whenever go routine finishes execution
			// fmt.Println("loop through msgs")
			for message := range *messages {
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
			concurrent_req_end_time := time.Now()

			avg_time_to_complete_api = avg_time_to_complete_api / concurrent_request
			avg_time_to_connect_api = avg_time_to_connect_api / concurrent_request
			avg_time_to_receive_first_byte_api = avg_time_to_receive_first_byte_api / concurrent_request

			status_code_in_percentage := make(map[int]float64)
			for status_code, occurrence := range status_codes {
				status_code_in_percentage[status_code] = (float64(occurrence) / float64(concurrent_request)) * 100
			}

			var avg_request_payload_size float64 = 0
			var avg_response_payload_size float64 = 0
			var additional_details_arr []AdditionalAPIDetails
			for additional_detail := range *additional_details {
				additional_details_arr = append(additional_details_arr, additional_detail)
				avg_request_payload_size += float64(additional_detail.request_payload_size)
				avg_response_payload_size += float64(additional_detail.response_payload_size)
			}
			if len(additional_details_arr) > 0 {
				avg_request_payload_size = avg_request_payload_size / float64(len(additional_details_arr))
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
			if track_iteration_time.Add(time.Second * 1).After(track_iteration_end_time) {
				track_iteration_time = track_iteration_end_time
			} else {
				track_iteration_time = track_iteration_time.Add(time.Second * 1)
			}
			var per_second_metrics []BenchMarkPerSecondCount
			// var _request_sent_in_sec_avg, _request_connected_in_sec_avg, _request_processed_in_sec_avg int64
			if len(additional_details_arr) > 0 {
				for {
					// fmt.Printf("prev_iteration_time=%v,track_iteration_time=%v\n",prev_iteration_time,track_iteration_time)
					var _request_sent_in_sec, _request_connected_in_sec, _request_received_first_byte_in_sec, _request_processed_in_sec int64
					var request_payload_size float64 = 0
					var response_payload_size float64 = 0
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
					if track_iteration_time.Add(time.Second * 1).After(track_iteration_end_time) {
						track_iteration_time = track_iteration_end_time
					} else {
						track_iteration_time = track_iteration_time.Add(time.Second * 1)
					}
				}
			}

			temp_data := BenchmarkData{
				Url:                             _url,
				Status_code_in_percentage:       status_code_in_percentage,
				Status_codes:                    status_codes,
				Concurrent_request:              concurrent_request,
				Avg_time_to_connect_api:         avg_time_to_connect_api,
				Avg_time_to_receive_first_byte:  avg_time_to_receive_first_byte_api,
				Avg_time_to_complete_api:        avg_time_to_complete_api,
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
			each_iterations_data = append(each_iterations_data, temp_data)
		}(&all_iteration_data[i], int64(i))
	}

	main_iter.Wait()
	all_iteration_data_collection_wg.Wait()

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
	return &each_iterations_data, &temp_data
}

func OnBenchmarkEnd() {
	BenchmarkMetricEvent.Emit(nil)
	BenchmarkMetricEvent.Dispose()
	BenchMarkEnded = true
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
