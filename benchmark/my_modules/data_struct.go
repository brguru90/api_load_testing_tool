package my_modules

import (
	"net/http"
	"time"
)

// total_number_of_request: Total number of request
// concurrent_request: Number of parallel request per each iteration

type MessageType struct {
	UID                  int64
	Data                 APIData
	Time_to_complete_api int64
	Err                  error
	Res                  *http.Response
}

type BenchmarkData struct {
	Url                             string
	Status_code_in_percentage       map[int]float64
	Status_codes                    map[int]int64
	Concurrent_request              int64
	Total_number_of_request         int64
	Avg_time_to_connect_api         int64                     `json:"Avg_time_to_connect_api_in_millesec,omitempty"`
	Avg_time_to_receive_first_byte  int64                     `json:"Avg_time_to_receive_first_byte,omitempty"`
	Avg_time_to_complete_api        int64                     `json:"Avg_time_to_complete_api_in_millesec,omitempty"`
	Min_time_to_complete_api        int64                     `json:"Min_time_to_complete_api_in_millesec,omitempty"`
	Max_time_to_complete_api        int64                     `json:"Max_time_to_complete_api_in_millesec,omitempty"`
	Total_time_to_complete_all_apis int64                     `json:"Total_time_to_complete_all_apis_in_millesec,omitempty"`
	Benchmark_per_second_metric     []BenchMarkPerSecondCount `json:"Benchmark_per_second_metric,omitempty"`
	// it may not accurate because its not the size of whole raw request
	Average_request_payload_size  float64 `json:"Average_request_payload_size_in_bytes,omitempty"`
	Average_response_payload_size float64 `json:"Average_response_payload_size_in_bytes,omitempty"`

	Avg_time_to_connect_api_in_sec                   float64 `json:"Avg_time_to_connect_api_in_sec,omitempty"`
	Avg_time_to_complete_api_in_sec                  float64 `json:"Avg_time_to_complete_api_in_sec,omitempty"`
	Min_time_to_complete_api_in_sec                  float64 `json:"Min_time_to_complete_api_in_sec,omitempty"`
	Max_time_to_complete_api_in_sec                  float64 `json:"Max_time_to_complete_api_in_sec,omitempty"`
	Total_time_to_complete_all_apis_iteration_in_sec float64 `json:"Total_time_to_complete_all_apis_iteration_in_sec,omitempty"`
	Average_request_payload_size_in_all_iteration    float64 `json:"Average_request_payload_size_in_bytes_in_all_iteration,omitempty"`
	Average_response_payload_size_in_all_iteration   float64 `json:"Average_response_payload_size_in_bytes_in_all_iteration,omitempty"`
	Total_operation_time_in_sec                      float64 `json:"Total_operation_time_in_sec,omitempty"`
	IterationID                                      int64   `json:"iteration_id"`
}

type AdditionalAPIDetails struct {
	request_id                  int64
	request_sent                time.Time
	request_connected           time.Time
	request_receives_first_byte time.Time
	request_processed           time.Time
	request_payload_size        int
	response_payload_size       int
}

type BenchMarkPerSecondCount struct {
	From_time_duration                   time.Time
	To_time_duration                     time.Time
	Request_sent                         int64
	Request_connected                    int64
	Request_receives_first_byte          int64
	Request_processed                    int64
	Total_request_payload_size_in_bytes  float64
	Total_response_payload_size_in_bytes float64
}
