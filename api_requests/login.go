package api_requests

import (
	"apis_load_test/benchmark/my_modules"
	"apis_load_test/benchmark/store"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)
type RequestSideSession struct {
	Cookies    []*http.Cookie
	CSRF_token string
}  
var RequestSession store.CredentialStore[RequestSideSession]

func LoginAsMultiUser(total_req int64, concurrent_req int64) interface{} {

	_url := "http://localhost:8000/api/login/"
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	RequestSession=RequestSession.NewCredentialStore(concurrent_req)

	fmt.Printf("len of cred ==> %v\n", len(*signup_credentials.CredentialStore_GetAllRefs()))

	payload_generator_callback := func(current_iteration int64) map[string]interface{} {
		if len(*signup_credentials.CredentialStore_GetAllRefs())>int(current_iteration % concurrent_req){
			return map[string]interface{}{
				"email": signup_credentials.CredentialStore_Get(current_iteration % concurrent_req).Email,
			}			
		}
		return map[string]interface{}{}
	}

	request_interceptor := func(req *http.Request, uid int64) {
		// fmt.Printf("request interceptor uid--> %v\n", uid)
	}

	response_interceptor := func(resp *http.Response, uid int64) {
		// fmt.Printf("response interceptor uid--> %v\n", uid)

		user_data := RequestSideSession{
			CSRF_token: resp.Header.Get("csrf_token"),
		}

		if len(resp.Cookies()) > 0 {
			if int64(RequestSession.CredentialStore_GetCount()) < concurrent_req && uid < concurrent_req {
				user_data.Cookies = resp.Cookies()
				RequestSession.CredentialStore_Append(user_data)
			}
		}
	}

	iteration_data, all_data := my_modules.BenchmarkAPIAsMultiUser(total_req, concurrent_req, _url, "post", headers, nil, payload_generator_callback, request_interceptor, response_interceptor)

	fmt.Println("bench mark on api finished")

	signup_credentials.Dispose()
	RequestSession.CredentialStore_WaitForAppend()
	RequestSession.CloseQ()
	fmt.Printf("total collected cookies %d\n", RequestSession.CredentialStore_GetCount())
	// fmt.Printf("collected cookies %v\n", *store.GetSessionsRefs())

	result := make(map[string]interface{})
	result[_url] = map[string]interface{}{
		"iteration_data": iteration_data,
		"all_data":       all_data,
	}
	return result
}
