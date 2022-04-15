package my_modules

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type ContextData struct {
	status_code int
	payload     map[string]interface{}
	json_body   map[string]interface{}
	body        []byte
	time        int64
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
	request_interceptor func(req *http.Request),
	response_interceptor func(resp *http.Response),
) (APIData, int64, *http.Response, error) {

	method = strings.ToUpper(method)

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

	if request_interceptor!=nil{
		request_interceptor(req)
	}

	client := &http.Client{}

	start_time := time.Now()
	resp, err := client.Do(req)
	end_time := time.Now()

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


	if response_interceptor!=nil{
		response_interceptor(resp)
	}

	defer resp.Body.Close()
	defer req.Body.Close()
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
			status_code: resp.StatusCode,
			payload:     payload_obj,
			json_body:   json_body,
			body:        body,
			time:        end_time.Sub(start_time).Milliseconds(),
		},
	}, end_time.Sub(start_time).Milliseconds(), resp, nil

}
