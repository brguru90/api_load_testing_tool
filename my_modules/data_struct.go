package my_modules

import "net/http"

// total_number_of_request: Total number of request
// concurrent_request: Number of parallel request per each iteration

type MessageType struct {
	Data                 APIData
	Time_to_complete_api int64
	Res                  *http.Response
	Err                  error
}

type BenchmarkData struct {
	Url                             string
	Status_code_in_percentage       map[int]float64
	Status_codes                    map[int]int64
	Concurrent_request              int64
	Total_number_of_request         int64
	Avg_time_to_connect_api         int64 `json:"Avg_time_to_connect_api_in_millesec,omitempty"`
	Avg_time_to_complete_api        int64 `json:"Avg_time_to_complete_api_in_millesec,omitempty"`
	Min_time_to_complete_api        int64 `json:"Min_time_to_complete_api_in_millesec,omitempty"`
	Max_time_to_complete_api        int64 `json:"Max_time_to_complete_api_in_millesec,omitempty"`
	Total_time_to_complete_all_apis int64 `json:"Total_time_to_complete_all_apis_in_millesec,omitempty"`

	Avg_time_to_connect_api_in_sec                   float64 `json:"Avg_time_to_connect_api_in_sec,omitempty"`
	Avg_time_to_complete_api_in_sec                  float64 `json:"Avg_time_to_complete_api_in_sec,omitempty"`
	Min_time_to_complete_api_in_sec                  float64 `json:"Min_time_to_complete_api_in_sec,omitempty"`
	Max_time_to_complete_api_in_sec                  float64 `json:"Max_time_to_complete_api_in_sec,omitempty"`
	Total_time_to_complete_all_apis_iteration_in_sec float64 `json:"Total_time_to_complete_all_apis_iteration_in_sec,omitempty"`
}
