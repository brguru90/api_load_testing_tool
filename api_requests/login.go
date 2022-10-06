package api_requests

import (
	"apis_load_test/benchmark/my_modules"
	"apis_load_test/benchmark/store"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

func getUserCredentialFromDB(_limit int64) []string {

	all_users_email := []string{}

	const (
		host     = "localhost"
		port     = 5432
		user     = "guru"
		password = "guru"
		dbname   = "jwt5"
	)

	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlconn)
	my_modules.CheckError(err)
	defer db.Close()

	err = db.Ping()
	my_modules.CheckError(err)

	rows, err := db.Query(fmt.Sprintf(`SELECT "email" FROM "users" LIMIT %d;`, _limit))
	my_modules.CheckError(err)
	defer rows.Close()

	for rows.Next() {
		var email string
		err = rows.Scan(&email)
		my_modules.CheckError(err)

		all_users_email = append(all_users_email, email)
	}
	my_modules.CheckError(err)

	return all_users_email
}

type Users struct {
	Users []map[string]interface{} `json:"users"`
}
type ResponseType struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Data   Users  `json:"data"`
}

func getUserCredentialFromAPI(_limit int64) []string {

	all_users_email := []string{}
	resp, err := http.Get("http://localhost:8000/api/all_users/")
	my_modules.CheckError(err)
	var json_body ResponseType
	// json_body := make(map[string]interface{})
	json.NewDecoder(resp.Body).Decode(&json_body)

	for _, user := range json_body.Data.Users {
		all_users_email = append(all_users_email, user["email"].(string))
	}

	return all_users_email
}

func LoginAsMultiUser(total_req int64, concurrent_req int64) interface{} {

	_url := "http://localhost:8000/api/login/"
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	// all_users_email := getUserCredentialFromAPI(concurrent_req)
	// if len(all_users_email) == 0 {
	// 	panic("user emails are empty")
	// }

	fmt.Printf("len of cred ==> %v\n", len(*store.LoginCredential_GetAll()))

	payload_generator_callback := func(current_iteration int64) map[string]interface{} {
		return map[string]interface{}{
			// "email": all_users_email[current_iteration],
			"email": store.LoginCredential_Get(current_iteration % concurrent_req).Email,
		}
	}

	request_interceptor := func(req *http.Request, uid int64) {
		// fmt.Printf("request interceptor uid--> %v\n", uid)
	}

	response_interceptor := func(resp *http.Response, uid int64) {
		// fmt.Printf("response interceptor uid--> %v\n", uid)

		user_data := store.RequestSideSession{
			CSRF_token: resp.Header.Get("csrf_token"),
		}

		if len(resp.Cookies()) > 0 {
			// condition check may not work all time since data is pushed concurrently`
			if int64(store.GetSessionsCount()) < concurrent_req && uid < concurrent_req {
				user_data.Cookies = resp.Cookies()
				store.AppendCSession(user_data)
			}
		}
	}

	iteration_data, all_data := my_modules.BenchmarkAPIAsMultiUser(total_req, concurrent_req, _url, "post", headers, nil, payload_generator_callback, request_interceptor, response_interceptor)

	fmt.Println("bench mark on api finished")

	store.RequestSideSession_WaitForAppend()
	fmt.Printf("total collected cookies %d\n", store.GetSessionsCount())
	// fmt.Printf("collected cookies %v\n", *store.GetSessionsRefs())

	result := make(map[string]interface{})
	result[_url] = map[string]interface{}{
		"iteration_data": iteration_data,
		"all_data":       all_data,
	}
	return result
}
