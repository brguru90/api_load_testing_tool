#include "api_req.h"

thread_local uv_loop_t *loop;
thread_local CURLM *curl_handle;
thread_local uv_timer_t timeout;

void loop_on_the_thread(request_input *req_inputs, response_data *response_ref, int total_requests, int debug);