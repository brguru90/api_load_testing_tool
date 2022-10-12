package my_modules

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
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

type CreatedAPIRequestFormat struct {
	req          *http.Request
	request_size int
	err          error
	payload      map[string]interface{}
	url          string
	uid          int64
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

func CreateAPIRequest(
	_url string,
	method string,
	header map[string]string,
	payload_obj map[string]interface{},
	uid int64,
	request_interceptor func(req *http.Request, uid int64),
) CreatedAPIRequestFormat {

	method = strings.ToUpper(method)

	payload, err := JSONMarshal(payload_obj)

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
		return CreatedAPIRequestFormat{
			req:          nil,
			request_size: -1,
			err:          err,
			payload:      payload_obj,
			url:          _url,
			uid:          uid,
		}
	}

	// req.Header.Add("Content-Type", "application/json")
	for key, value := range header {
		req.Header.Add(key, value)
	}

	if request_interceptor != nil {
		request_interceptor(req, uid)
	}

	var request_size int = 0

	reqDump, err := httputil.DumpRequestOut(req, true)
	if err == nil {
		request_size = len(reqDump)
	}

	return CreatedAPIRequestFormat{
		req:          req,
		request_size: request_size,
		err:          nil,
		payload:      payload_obj,
		url:          _url,
		uid:          uid,
	}

}

// make http request
// get the metrics like delay, payload size etc for particular request
func APIReq(
	api_request CreatedAPIRequestFormat,
	response_interceptor func(resp *http.Response, uid int64),
	additional_detail_chan chan AdditionalAPIDetails,
) (APIData, int64, *http.Response, error) {
	uid := api_request.uid
	additional_detail := AdditionalAPIDetails{
		request_id: uid,
	}

	if api_request.err != nil {
		return APIData{
			url:     api_request.url,
			context: "json.Marshal",
			context_data: ContextData{
				status_code: -1,
				payload:     api_request.payload,
			},
		}, 0, nil, api_request.err
	}

	start_time := time.Now()
	var connected_time time.Time = start_time
	additional_detail.request_sent = start_time
	additional_detail.request_payload_size = api_request.request_size

	// https://pkg.go.dev/net/http/httptrace@go1.18.2
	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			connected_time = time.Now()
			additional_detail.request_connected = connected_time
			// fmt.Printf("Got Conn: %+v,\t%v\n", connInfo,connected_time.Sub(start_time).Milliseconds())
		},
		GotFirstResponseByte: func(){
			additional_detail.request_received_first_byte=time.Now()
		},
		// DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
		// 	fmt.Printf("DNS Info: %+v\n", dnsInfo)
		// },
	}
	req := api_request.req.WithContext(httptrace.WithClientTrace(api_request.req.Context(), trace))

	client := &http.Client{
		Timeout: HTTPTimeout,
	}

	resp, err := client.Do(req)
	end_time := time.Now()
	additional_detail.request_processed = end_time

	if err != nil {
		return APIData{
			url:     api_request.url,
			context: "client.Do",
			context_data: ContextData{
				status_code: -1,
				payload:     api_request.payload,
			},
		}, 0, nil, err
	}

	// end of time difference calculation

	var response_size int = 0
	// if resp.Body!=nil && resp.ContentLength!=0{
	// 	var resp_body_copy bytes.Buffer
	// 	_, io_err := io.Copy(&resp_body_copy, resp.Body)
	// 	if io_err == nil {
	// 		response_size = len(resp_body_copy.Bytes())
	// 		reader := bytes.NewReader(resp_body_copy.Bytes())
	// 		resp.Body = ioutil.NopCloser(reader)
	// 	}
	// }
	// response_header_string:=""
	// if req.Header!=nil{
	// 	for key,value :=range resp.Header{
	// 		response_header_string+=fmt.Sprintf("%v: %v\n",key,strings.Join(value,","))
	// 	}
	// }
	// response_size+=len([]byte(response_header_string))

	respDump, err := httputil.DumpResponse(resp, true)
	if err == nil {
		response_size = len(respDump)
	}

	additional_detail.response_payload_size = response_size

	if additional_detail_chan != nil {
		additional_detail_chan <- additional_detail
	}

	if err != nil {
		return APIData{
			url:     api_request.url,
			context: "API request send",
			context_data: ContextData{
				status_code: -1,
				payload:     api_request.payload,
			},
		}, end_time.Sub(start_time).Milliseconds(), nil, err
	}

	if response_interceptor != nil {
		response_interceptor(resp, uid)
	}

	defer resp.Body.Close()
	if api_request.payload != nil {
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
		url:     api_request.url,
		context: "API response",
		context_data: ContextData{
			status_code:     resp.StatusCode,
			payload:         api_request.payload,
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
