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
    int headers_len;
    char *body;
} request_input;

typedef struct ResponseData
{
    char *uid;
    char *response_header;
    char *response_body;
    // time_t before_connect_time; //long int
    long long before_connect_time_microsec;
    long long connected_at_microsec;
    long long first_byte_at_microsec;
    long long finish_at_microsec;
    long connect_time_microsec;
    long time_to_first_byte_microsec;
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

struct memory
{
    char *data;
    size_t size;
};

// extern void response_callback_from_c(int loop_size, response_data res_data[], char *uuid);
void *loop_on_the_thread(void *data);
void send_request_in_concurrently(request_input *req_inputs, response_data *response_ref, additional_details _additional_details, int debug);
response_data send_raw_request(request_input *req_input, response_data *response_ref, int debug);