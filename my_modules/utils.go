package my_modules

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"strings"
	"time"
)

var HTTPTimeout = time.Minute * 5

type ContextData struct {
	status_code     int
	payload         map[string]interface{}
	json_body       map[string]interface{}
	body            []byte
	time            int64
	time_to_connect int64
}

type APIData struct {
	url          string
	context      string
	context_data ContextData
}

func RandomBytes(size int) (blk []byte, err error) {
	blk = make([]byte, size)
	_, err = rand.Read(blk)
	return
}

func RandomString(size int) string {
	var r_err error = nil
	var _rand []byte
	if _rand, r_err = RandomBytes(size); r_err == nil {
		return base64.StdEncoding.EncodeToString(_rand)[:size]
	}
	return ""
}

func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", " ")
	err := encoder.Encode(t)
	return bytes.TrimRight(buffer.Bytes(), "\n"), err
}

func APIReq(
	_url string,
	method string,
	header map[string]string,
	payload_obj map[string]interface{},
	uid int64,
	request_interceptor func(req *http.Request, uid int64),
	response_interceptor func(resp *http.Response, uid int64),
	additional_detail_chan chan BenchMarkPerSecondDetail,
) (APIData, int64, *http.Response, error) {

	method = strings.ToUpper(method)
	additional_detail:=BenchMarkPerSecondDetail{
		request_id: uid,
	}

	payload, err := JSONMarshal(payload_obj)
	if err != nil {
		return APIData{
			url:     _url,
			context: "json.Marshal",
			context_data: ContextData{
				status_code: -1,
				payload:     payload_obj,
			},
		}, 0, nil, err
	}

	var req *http.Request
	switch method {
	case "GET":
		req, err = http.NewRequest(method, _url, nil)
	case "POST", "PUT", "DELETE":
		req, err = http.NewRequest(method, _url, bytes.NewBuffer(payload))
	default:
		req, err = http.NewRequest(method, _url, nil)
	}
	if err != nil {

		return APIData{
			url:     _url,
			context: "API request creation",
			context_data: ContextData{
				status_code: -1,
				payload:     payload_obj,
			},
		}, 0, nil, err
	}

	// req.Header.Add("Content-Type", "application/json")
	for key, value := range header {
		req.Header.Add(key, value)
	}

	if request_interceptor != nil {
		request_interceptor(req, uid)
	}

	start_time := time.Now()
	var connected_time time.Time = start_time
	additional_detail.request_sent=start_time

	// https://pkg.go.dev/net/http/httptrace@go1.18.2
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			connected_time = time.Now()
			additional_detail.request_connected=connected_time
			// fmt.Printf("Got Conn: %+v,\t%v\n", connInfo,connected_time.Sub(start_time).Milliseconds())
		},
		// DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
		// 	fmt.Printf("DNS Info: %+v\n", dnsInfo)
		// },
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))


	client := &http.Client{
		Timeout: HTTPTimeout,
	}

	resp, err := client.Do(req)
	end_time := time.Now()
	additional_detail.request_processed=end_time

	if additional_detail_chan!=nil{
		additional_detail_chan <- additional_detail
	}

	if err != nil {

		return APIData{
			url:     _url,
			context: "API request send",
			context_data: ContextData{
				status_code: -1,
				payload:     payload_obj,
			},
		}, end_time.Sub(start_time).Milliseconds(), nil, err
	}

	if response_interceptor != nil {
		response_interceptor(resp, uid)
	}

	defer resp.Body.Close()
	if payload_obj != nil {
		defer req.Body.Close()
	}
	json_body := make(map[string]interface{})
	var body []byte = nil
	if strings.Contains(resp.Header.Get("Content-Type"), "json") {
		json.NewDecoder(resp.Body).Decode(&json_body)
	} else {
		body, _ = ioutil.ReadAll(resp.Body)
		json_body = nil
	}

	return APIData{
		url:     _url,
		context: "API response",
		context_data: ContextData{
			status_code:     resp.StatusCode,
			payload:         payload_obj,
			json_body:       json_body,
			body:            body,
			time:            end_time.Sub(start_time).Milliseconds(),
			time_to_connect: connected_time.Sub(start_time).Milliseconds(),
		},
	}, end_time.Sub(start_time).Milliseconds(), resp, nil

}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
