#include "api_req_helper.hpp"

void send_request_in_concurrently(request_input *req_inputs, response_data *response_ref, int total_requests, int debug)
{
    loop_on_the_thread(req_inputs, response_ref, total_requests, debug);
}
