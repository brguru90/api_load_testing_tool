package my_modules

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

// total_number_of_request: Total number of request
// concurrent_request: Number of parallel request per each iteratio

type MessageType struct {
	data                 interface{}
	time_to_complete_api int64
	res                  *http.Response
	err                  error
}

type BenchmarkData struct {
	_url string
	concurrent_request int64
	avg_time_to_complete_api int64
	min_time_to_complete_api int64
	max_time_to_complete_api int64

}

func BenchmarkAPI(total_number_of_request int64, concurrent_request int64, _url string, method string, headers map[string]string, payload_obj map[string]interface{}) {

	var each_iterations_data []BenchmarkData
	number_of_iteration := total_number_of_request / concurrent_request

	if total_number_of_request < concurrent_request {
		panic("total_number_of_request<concurrent_request")
	}

	if (total_number_of_request % concurrent_request) != 0 {
		panic("(total_number_of_request % concurrent_request)!=0")
	}

	var wg sync.WaitGroup

	var i, j, avg_time_to_complete_api int64
	for i = 0; i < number_of_iteration; i++ {

		messages := make(chan MessageType)

		// spin up all request parallally & wait for all those to finish
		wg.Add(int(number_of_iteration))
		for j = 0; j < concurrent_request; j++ {
			go func() {
				defer wg.Done()
				data, time_to_complete_api, res, err := APIReq(_url, method, headers, payload_obj)
				messages <- MessageType{
					data:                 data,
					time_to_complete_api: time_to_complete_api,
					res:                  res,
					err:                  err,
				}
			}()
		}
		wg.Wait()
		close(messages)

		avg_time_to_complete_api = 0
		for message := range messages {
			avg_time_to_complete_api += message.time_to_complete_api
		}
		avg_time_to_complete_api=avg_time_to_complete_api/concurrent_request

		each_iterations_data = append(each_iterations_data, BenchmarkData{
			_url: _url,
			avg_time_to_complete_api:avg_time_to_complete_api,
		})

	}

}
