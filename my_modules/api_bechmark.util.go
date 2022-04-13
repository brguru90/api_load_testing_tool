package my_modules

import (
	"math"
	"net/http"
	"sync"
)

// total_number_of_request: Total number of request
// concurrent_request: Number of parallel request per each iteratio

type MessageType struct {
	Data                 APIData
	Time_to_complete_api int64
	Res                  *http.Response
	Err                  error
}

type BenchmarkData struct {
	Url                       string
	Status_code_in_percentage map[int]float64
	Status_codes              map[int]int64
	Concurrent_request        int64
	Total_number_of_request   int64
	Avg_time_to_complete_api  int64 `json:"Avg_time_to_complete_api_in_millesec,omitempty"`
	Min_time_to_complete_api  int64 `json:"Min_time_to_complete_api_in_millesec,omitempty"`
	Max_time_to_complete_api  int64 `json:"Max_time_to_complete_api_in_millesec,omitempty"`


	Avg_time_to_complete_api_in_sec  float64 `json:"Avg_time_to_complete_api_in_sec"`
	Min_time_to_complete_api_in_sec  float64 `json:"Min_time_to_complete_api_in_sec,omitempty"`
	Max_time_to_complete_api_in_sec  float64 `json:"Max_time_to_complete_api_in_sec,omitempty"`
}

func BenchmarkAPI(total_number_of_request int64, concurrent_request int64, _url string, method string, headers map[string]string, payload_obj map[string]interface{},payload_generator_callback  func() map[string]interface{}) (*[]BenchmarkData, *BenchmarkData) {

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

	var i, j, avg_time_to_complete_api int64
	for i = 0; i < number_of_iteration; i++ {
		messages := make(chan MessageType)
		// fmt.Printf("i=%v\n", i)

		// run whole parallel request routine in background
		// so that i can loop through channel data later
		iteration_wg.Add(1)
		go func() {
			defer iteration_wg.Done()
			// spin up all request parallally & wait for all those to finish
			concurrent_req_wg.Add(int(concurrent_request))
			for j = 0; j < concurrent_request; j++ {
				// fmt.Printf("%v-%v\n", i, j)
				go func() {
					defer concurrent_req_wg.Done()
					var api_payload map[string]interface{}
					if payload_obj==nil{
						api_payload=payload_generator_callback()
					} else {
						api_payload=payload_obj
					}
					data, time_to_complete_api, res, err := APIReq(_url, method, headers, api_payload)
					// fmt.Printf("finish APIReq\n")
					messages <- MessageType{
						Data:                 data,
						Time_to_complete_api: time_to_complete_api,
						Res:                  res,
						Err:                  err,
					}
					// fmt.Printf("finish channel\n")
				}()
			}
			concurrent_req_wg.Wait()
			// fmt.Println("all the parallel request finished")
			close(messages)
		}()

		avg_time_to_complete_api = 0
		min_time_to_complete_api := math.Inf(1)
		max_time_to_complete_api := 0.0
		status_codes := make(map[int]int64)
		// looping through channel data, whenever go routine finishes execution
		// fmt.Println("loop through msgs")
		for message := range messages {
			avg_time_to_complete_api += message.Time_to_complete_api
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
		avg_time_to_complete_api = avg_time_to_complete_api / concurrent_request

		status_code_in_percentage := make(map[int]float64)
		for status_code, occurrence := range status_codes {
			status_code_in_percentage[status_code] = (float64(occurrence) / float64(concurrent_request)) * 100
		}

		each_iterations_data = append(each_iterations_data, BenchmarkData{
			Url:                       _url,
			Status_code_in_percentage: status_code_in_percentage,
			Status_codes:              status_codes,
			Concurrent_request:        concurrent_request,
			Avg_time_to_complete_api:  avg_time_to_complete_api,
			Min_time_to_complete_api:  int64(min_time_to_complete_api),
			Max_time_to_complete_api:  int64(max_time_to_complete_api),
		})

	}

	avg_time_to_complete_api = 0
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
	}

	return &each_iterations_data, &BenchmarkData{
		Url:                      _url,
		Status_codes:             status_codes,
		Total_number_of_request:  total_number_of_request,
		Concurrent_request:       concurrent_request,
		Min_time_to_complete_api_in_sec: min_time_to_complete_api/1000.0,
		Max_time_to_complete_api_in_sec: max_time_to_complete_api/1000.0,
		Avg_time_to_complete_api_in_sec: (float64(avg_time_to_complete_api) / float64(total_number_of_request))/1000,
	}
}
