#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <fcntl.h>

#include <pthread.h>
#include <time.h>
#include <math.h>
#include <sys/mman.h>
#include <sys/types.h>
#include <sys/shm.h>
#include <sys/ipc.h>
#include <sys/wait.h>

#ifndef __cplusplus
#include <stdbool.h>
#endif

#include <curl/curl.h>
#include <uv.h>

typedef struct curl_context_s
{
    uv_poll_t poll_handle;
    curl_socket_t sockfd;
} curl_context_t;

typedef struct CurlHandlers
{
    uv_loop_t *loop;
    CURLM *curl_handle;
    uv_timer_t timeout;
} curl_handlers_t;

typedef struct Headers
{
    char *header;
} headers_type;

typedef struct AdditionalDetails
{
    char *uuid;
    int total_requests;
    int total_threads;
} additional_details;

typedef struct SingleRequestInput
{
    char *uid;
    char *url;
    char *method;
    headers_type *headers;
    char *cookies;
    int headers_len;
    char *body;
    int time_out_in_sec;
} request_input;


struct memory
{
    char *data;
    size_t size;
};

typedef struct ResponseData
{
    int debug;
    char *uid;
    char *response_header;
    char *response_body;
    struct memory Resp_header;
    struct memory Resp_body;
    // time_t before_connect_time; //long int
    long long before_connect_time_microsec;
    long long after_response_time_microsec;
    long long connected_at_microsec;
    long long first_byte_at_microsec;
    long long finish_at_microsec;
    long connect_time_microsec;
    long time_to_first_byte_microsec;
    long total_time_from_curl_microsec;
    long total_time_microsec;
    int status_code;
    int err_code;
} response_data;

typedef struct ThreadPoolData
{
    int start_index;
    int end_index;
    pid_t pid;
    char *uuid;
    bool full_index;
} thread_pool_data;
typedef struct ThreadData
{
    request_input *req_inputs_ptr;
    response_data *response_ref_ptr;
    int thread_id;
    int debug_flag;
    thread_pool_data th_pool_data;
} thread_data;


#ifdef __cplusplus
extern "C"
{
#endif

    void send_request_in_concurrently(request_input *req_inputs, response_data *response_ref, int total_requests, int debug);
    // void send_raw_request(request_input *req_input, response_data *response_ref, int debug);

#ifdef __cplusplus
}
#endif