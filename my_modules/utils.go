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

func APIReq(_url string, method string, header map[string]string, payload_obj map[string]interface{}) (map[string]interface{}, int64, *http.Response, error) {

	method = strings.ToUpper(method)

	payload, err := JSONMarshal(payload_obj)
	if err != nil {
		return (map[string]interface{}{
			"url":          _url,
			"context":      "json.Marshal",
			"context_data": map[string]interface{}{},
		}), 0, nil, err
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
		return (map[string]interface{}{
			"url":     _url,
			"context": "API request creation",
			"context_data": map[string]interface{}{
				"status_code": -1,
				"payload":     payload,
			},
			"error": err,
		}), 0, nil, err
	}

	// req.Header.Add("Content-Type", "application/json")
	for key, value := range header {
		req.Header.Add(key, value)
	}

	client := &http.Client{}

	start_time := time.Now()
	resp, err := client.Do(req)
	end_time := time.Now()

	if err != nil {
		return (map[string]interface{}{
			"url":     _url,
			"context": "API request send",
			"context_data": map[string]interface{}{
				"status_code": -1,
				"payload":     payload,
			},
			"error": err,
		}), end_time.Sub(start_time).Milliseconds(), nil, err
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
	return (map[string]interface{}{
		"url":     _url,
		"context": "API request",
		"context_data": map[string]interface{}{
			"status_code": resp.StatusCode,
			"payload":     payload_obj,
			"json_body":   json_body,
			"body":        body,
			"time":        end_time.Sub(start_time).Milliseconds(),
		},
	}), end_time.Sub(start_time).Milliseconds(), resp, nil

}
