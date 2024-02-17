#include "api_req.h"



class api_req_async
{
private:
    void *data;
    int req_sent,req_timeout,req_finish;
    uv_loop_t *loop = nullptr;
    CURLM *curl_handle = nullptr;
    uv_timer_t timeout;
    void add_request_to_event_loop(request_input *req_input, response_data *response_ref, int debug);
    curl_context_t *create_curl_context(curl_socket_t sockfd);
    int start_timeout(CURLM *multi, long timeout_ms, void *userp);
    int handle_socket(CURL *easy, curl_socket_t s, int action, void *userp, void *socketp);

public:
    api_req_async();
    ~api_req_async();
    void loop_on_the_thread(request_input *req_inputs, response_data *response_ref, int total_requests, int debug);
    void (api_req_async::*on_timeout_ptr)(uv_timer_t *req);
};

// void loop_on_the_thread(request_input *req_inputs, response_data *response_ref, int total_requests, int debug);