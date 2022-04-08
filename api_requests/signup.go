package api_requests

import (
	"apis_load_test/my_modules"
	"bytes"
	"encoding/json"
	"net/http"
	"strconv"
)

func SignUp() (string,int, string, error) {
	_url:="http://localhost:8000/api/sign_up/"

	payload, err := json.Marshal(map[string]interface{}{
		"email":       my_modules.RandomString(100) + "@gmail.com",
		"name":        my_modules.RandomString(20),
		"description": my_modules.RandomString(100),
	})
	if err != nil {
		return _url,-1, "json", err
	}

	resp, err := http.Post(
		_url,
		"application/json",
		bytes.NewBuffer(payload),
	)
	if err != nil {
		return _url,-1, "Post", err
	}

	defer resp.Body.Close()
	payload2, err :=strconv.Unquote(string(payload))
	if err != nil {
		return _url,-1, "Unquote", err
	}
	return _url,resp.StatusCode, string(payload2), nil
}
