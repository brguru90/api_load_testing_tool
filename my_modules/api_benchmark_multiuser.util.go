package my_modules

import (
	"fmt"
	"math"
	"net/http"
	"sync"
	"time"
)

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

	var each_iterations_data []BenchmarkData
	number_of_iteration := total_number_of_request / concurrent_request

	if total_number_of_request < concurrent_request {
		panic("total_number_of_request<concurrent_request")
	}

	if (total_number_of_request % concurrent_request) != 0 {
		panic("(total_number_of_request % concurrent_request)!=0")
	}

	var concurrent_req_wg sync.WaitGroup
	var iteration_wg sync.WaitGroup

	var i, j, avg_time_to_complete_api, avg_time_to_connect_api int64

	iterations_start_time := time.Now()
	for i = 0; i < number_of_iteration; i++ {
		messages := make(chan MessageType)
		additional_details := make(chan BenchMarkPerSecondDetail, concurrent_request)
		fmt.Printf("url=%s,i=%v\n", _url, i)

		concurrent_req_start_time := time.Now()
		// run whole parallel request routine in background
		// so that i can loop through channel data later
		iteration_wg.Add(1)
		go func(main_iteration int64) {
			defer iteration_wg.Done()
			// spin up all request parallally & wait for all those to finish
			concurrent_req_wg.Add(int(concurrent_request))
			for j = 0; j < concurrent_request; j++ {
				// fmt.Printf("%v-%v\n", i, j)
				go func(sub_iteration int64) {
					defer concurrent_req_wg.Done()
					var api_payload map[string]interface{}
					if payload_obj == nil && payload_generator_callback != nil {
						api_payload = payload_generator_callback(sub_iteration)
					} else {
						api_payload = payload_obj
					}
					data, time_to_complete_api, res, err := APIReq(_url, method, headers, api_payload, sub_iteration, request_interceptor, response_interceptor, additional_details)
					// fmt.Printf("finish APIReq\n")
					messages <- MessageType{
						Data:                 data,
						Time_to_complete_api: time_to_complete_api,
						Res:                  res,
						Err:                  err,
					}
					// fmt.Printf("finish channel\n")
				}((main_iteration * concurrent_request) + j)
			}
			concurrent_req_wg.Wait()
			// fmt.Println("all the parallel request finished")
			close(messages)
			close(additional_details)
		}(i)

		avg_time_to_complete_api = 0
		avg_time_to_connect_api = 0
		min_time_to_complete_api := math.Inf(1)
		max_time_to_complete_api := 0.0
		status_codes := make(map[int]int64)
		// looping through channel data, whenever go routine finishes execution
		// fmt.Println("loop through msgs")
		for message := range messages {
			avg_time_to_complete_api += message.Time_to_complete_api
			avg_time_to_connect_api += message.Data.context_data.time_to_connect
			var cur_status_code int = message.Data.context_data.status_code
			status_codes[cur_status_code] += 1

			if float64(message.Time_to_complete_api) < min_time_to_complete_api {
				min_time_to_complete_api = float64(message.Time_to_complete_api)
			}

			if float64(message.Time_to_complete_api) > max_time_to_complete_api {
				max_time_to_complete_api = float64(message.Time_to_complete_api)
			}
		}
		// fmt.Println("finish loop through msgs")
		iteration_wg.Wait() // just wait here before doing next iteration
		// fmt.Println("finish wait")
		concurrent_req_end_time := time.Now()

		avg_time_to_complete_api = avg_time_to_complete_api / concurrent_request
		avg_time_to_connect_api = avg_time_to_connect_api / concurrent_request

		status_code_in_percentage := make(map[int]float64)
		for status_code, occurrence := range status_codes {
			status_code_in_percentage[status_code] = (float64(occurrence) / float64(concurrent_request)) * 100
		}

		var additional_details_arr []BenchMarkPerSecondDetail
		for additional_detail := range additional_details {
			additional_details_arr = append(additional_details_arr, additional_detail)
		}

		// var request_sent_in_sec_avg,request_connected_in_sec_avg,request_processed_in_sec_avg time.Time
		track_iteration_time := concurrent_req_start_time
		prev_iteration_time := concurrent_req_start_time
		track_iteration_time = track_iteration_time.Add(time.Second * 1)
		var per_second_metrics []BenchMarkPerSecondCount
		// var _request_sent_in_sec_avg, _request_connected_in_sec_avg, _request_processed_in_sec_avg int64
		if len(additional_details_arr) > 0 {
			for {
				// fmt.Printf("prev_iteration_time=%v,track_iteration_time=%v\n",prev_iteration_time,track_iteration_time)
				var _request_sent_in_sec, _request_connected_in_sec, _request_processed_in_sec int64
				for _, additional_detail := range additional_details_arr {
					if additional_detail.request_sent.After(prev_iteration_time) && additional_detail.request_sent.Before(track_iteration_time) {
						_request_sent_in_sec++
					}
					if additional_detail.request_connected.After(prev_iteration_time) && additional_detail.request_connected.Before(track_iteration_time) {
						_request_connected_in_sec++
					}
					if additional_detail.request_processed.After(prev_iteration_time) && additional_detail.request_processed.Before(track_iteration_time) {
						_request_processed_in_sec++
					}
				}
				per_second_metrics = append(per_second_metrics, BenchMarkPerSecondCount{
					From_time_duration: prev_iteration_time,
					To_time_duration:   track_iteration_time,
					Request_sent:       _request_sent_in_sec,
					Request_connected:  _request_connected_in_sec,
					Request_processed:  _request_processed_in_sec,
				})

				// iterate over elapsed time
				if track_iteration_time.After(concurrent_req_end_time) {
					break
				}

				prev_iteration_time = track_iteration_time
				track_iteration_time = track_iteration_time.Add(time.Second * 1)
			}
		}

		each_iterations_data = append(each_iterations_data, BenchmarkData{
			Url:                             _url,
			Status_code_in_percentage:       status_code_in_percentage,
			Status_codes:                    status_codes,
			Concurrent_request:              concurrent_request,
			Avg_time_to_connect_api:         avg_time_to_connect_api,
			Avg_time_to_complete_api:        avg_time_to_complete_api,
			Min_time_to_complete_api:        int64(min_time_to_complete_api),
			Max_time_to_complete_api:        int64(max_time_to_complete_api),
			Total_time_to_complete_all_apis: concurrent_req_end_time.Sub(concurrent_req_start_time).Milliseconds(),
			Benchmark_per_second_metric:     per_second_metrics,
		})

	}
	iterations_end_time := time.Now()

	avg_time_to_complete_api = 0
	avg_time_to_connect_api = 0
	min_time_to_complete_api := math.Inf(1)
	max_time_to_complete_api := 0.0
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
	}

	status_codes_in_perc := make(map[int]float64)
	for status_code, occurrence := range status_codes {
		status_codes_in_perc[status_code] = (float64(occurrence) / float64(total_number_of_request)) * 100.0
	}

	return &each_iterations_data, &BenchmarkData{
		Url:                             _url,
		Status_codes:                    status_codes,
		Status_code_in_percentage:       status_codes_in_perc,
		Total_number_of_request:         total_number_of_request,
		Concurrent_request:              concurrent_request,
		Avg_time_to_connect_api_in_sec:  (float64(avg_time_to_connect_api) / float64(total_number_of_request)) / 1000,
		Min_time_to_complete_api_in_sec: min_time_to_complete_api / 1000.0,
		Max_time_to_complete_api_in_sec: max_time_to_complete_api / 1000.0,
		Avg_time_to_complete_api_in_sec: (float64(avg_time_to_complete_api) / float64(total_number_of_request)) / 1000,
		Total_time_to_complete_all_apis_iteration_in_sec: float64(iterations_end_time.Sub(iterations_start_time).Milliseconds()) / 1000.0,
	}
}
