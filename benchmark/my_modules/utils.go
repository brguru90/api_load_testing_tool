package my_modules

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
	"strings"
	"time"
)

var HTTPTimeout = time.Minute * 5

var ShouldDumpRequestAndResponse bool = false

// Single transport reused across the whole process so connections + TLS sessions
// are pooled. Per-request http.Client allocation (old behaviour) forced new dial
// + new TLS handshake on every call -> the dominant overhead at high concurrency.
var DefaultTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}).DialContext,
	ForceAttemptHTTP2:     true, // HTTP/2 for HTTPS targets that ALPN-negotiate it
	MaxIdleConns:          0,    // unlimited
	MaxIdleConnsPerHost:   100000,
	MaxConnsPerHost:       0, // unlimited
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
	DisableCompression:    true,
}

// No Timeout on the client itself; per-request context honours a mutable HTTPTimeout.
var sharedClient = &http.Client{
	Transport: DefaultTransport,
}

type ContextData struct {
	status_code     int
	payload         map[string]interface{}
	json_body       map[string]interface{}
	body            []byte
	time            int64
	time_to_connect int64
	ttfb            int64
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

func JSONMarshal(t interface{}, indent bool) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	if indent {
		encoder.SetIndent("", " ")
	}
	err := encoder.Encode(t)
	return bytes.TrimRight(buffer.Bytes(), "\n"), err
}

// cheap header byte estimate, no serialization
func headerSize(h http.Header) int {
	n := 0
	for k, v := range h {
		for _, val := range v {
			n += len(k) + len(val) + 4 // ": " + "\r\n"
		}
	}
	return n
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

	payload, err := JSONMarshal(payload_obj, false)

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

	for key, value := range header {
		req.Header.Add(key, value)
	}

	if request_interceptor != nil {
		request_interceptor(req, uid)
	}

	var request_size int
	if ShouldDumpRequestAndResponse {
		// exact wire size via dump (heavy: full marshal). Opt-in only.
		reqDump, derr := httputil.DumpRequestOut(req, true)
		if derr == nil {
			request_size = len(reqDump)
		}
	} else {
		// cheap estimate: request line + headers + body length
		request_size = len(req.Method) + 1 + len(req.URL.RequestURI()) + 11 // "<METHOD> <uri> HTTP/1.1\r\n"
		request_size += headerSize(req.Header)
		request_size += len(payload)
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

// Run a single HTTP request, capture timing + size metrics, drain + close body
// so the connection returns to the pool. response_interceptor (if non-nil) runs
// before the body is drained so the caller can read it if they need to.
//
// additional_detail_chan: caller passes a *chan *AdditionalAPIDetails. The
// pointer is taken from additionalDetailsPool, populated, and sent on the
// channel. The receiver MUST call additionalDetailsPool.Put(ad) after use, or
// the pool grows unbounded. If chan is nil, this function puts it back itself.
func APIReq(
	api_request *CreatedAPIRequestFormat,
	response_interceptor func(resp *http.Response, uid int64),
	additional_detail_chan *chan *AdditionalAPIDetails,
) (APIData, int64, *http.Response, error) {
	uid := &api_request.uid

	ad := additionalDetailsPool.Get().(*AdditionalAPIDetails)
	*ad = AdditionalAPIDetails{request_id: *uid}

	// Single send-or-return path: the deferred handoff guarantees we never
	// leak the pooled struct, regardless of which return branch fires below.
	defer func() {
		if additional_detail_chan != nil {
			*additional_detail_chan <- ad
		} else {
			additionalDetailsPool.Put(ad)
		}
	}()

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
	var connected_time time.Time
	var ttfb time.Time

	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			connected_time = time.Now()
			ad.request_connected = connected_time
		},
		GotFirstResponseByte: func() {
			ttfb = time.Now()
			ad.request_receives_first_byte = ttfb
		},
	}

	ad.request_sent = time.Now()
	ad.request_payload_size = api_request.request_size

	ctx, cancel := context.WithTimeout(api_request.req.Context(), HTTPTimeout)
	defer cancel()
	req := api_request.req.WithContext(httptrace.WithClientTrace(ctx, trace))

	resp, err := sharedClient.Do(req)
	end_time := time.Now()
	ad.request_processed = end_time

	if err != nil {
		return APIData{
			url:     api_request.url,
			context: "client.Do",
			context_data: ContextData{
				status_code: -1,
				payload:     api_request.payload,
			},
		}, end_time.Sub(start_time).Milliseconds(), nil, err
	}

	if response_interceptor != nil {
		response_interceptor(resp, *uid)
	}

	// MUST drain + close so the connection goes back to the keep-alive pool.
	// Without this, MaxIdleConnsPerHost saturates and concurrency stalls.
	var bodyBytes int64
	if resp.Body != nil {
		bodyBytes, _ = io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}
	ad.response_payload_size = int(bodyBytes) + headerSize(resp.Header)
	end_time = time.Now()
	ad.request_processed = end_time

	return APIData{
		url:     api_request.url,
		context: "API response",
		context_data: ContextData{
			status_code:     resp.StatusCode,
			payload:         api_request.payload,
			time:            end_time.Sub(start_time).Milliseconds(),
			time_to_connect: connected_time.Sub(start_time).Milliseconds(),
			ttfb:            ttfb.Sub(start_time).Milliseconds(),
		},
	}, end_time.Sub(start_time).Milliseconds(), resp, nil
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}
