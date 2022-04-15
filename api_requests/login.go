package api_requests

import (
	"apis_load_test/my_modules"
	"apis_load_test/store"
	"database/sql"
	"fmt"
	"net/http"
	"sync"

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

	rows, err := db.Query(fmt.Sprintf(`SELECT "email" FROM "users" LIMIT %d;`,_limit))
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

func LoginAsMultiUser() interface{} {
	var total_req int64=100
	var concurrent_req int64=10


	_url := "http://localhost:8000/api/login/"
	headers := map[string]string{
		"Content-Type": "application/json",
	}

	all_users_email := getUserCredentialFromDB(concurrent_req);
	if len(all_users_email) == 0 {
		panic("user emails are empty")
	}

	payload_generator_callback := func(current_iteration int64) map[string]interface{} {
		return map[string]interface{}{
			"email": all_users_email[current_iteration],
		}
	}

	request_interceptor := func(req *http.Request, uid int64) {
		fmt.Printf("request interceptor uid--> %v\n", uid)
	}

	var m sync.Mutex
	response_interceptor := func(resp *http.Response, uid int64) {
		fmt.Printf("response interceptor uid--> %v\n", uid)

		for _, cookie := range resp.Cookies() {
			if cookie.Name == "access_token" {
				m.Lock()
				if int64(len(*store.GetCookiesRefs()))<concurrent_req{
					store.AppendCookie(cookie)
				}
				m.Unlock()
			}
		}
	}

	iteration_data, all_data := my_modules.BenchmarkAPIAsMultiUser(total_req, concurrent_req, _url, "post", headers, nil, payload_generator_callback, request_interceptor, response_interceptor)

	fmt.Printf("total collected cookies %d\n",len(store.GetAllCookies()))
	fmt.Printf("total collected cookies %d\n",len(*store.GetCookiesRefs()))

	fmt.Println("bench mark on api finished")

	return map[string]interface{}{
		"iteration_data": iteration_data,
		"all_data":       all_data,
	}
}
