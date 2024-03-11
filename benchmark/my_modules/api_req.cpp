#include "api_req_helper.hpp"

// git clone https://github.com/curl/curl.git
// https://gist.github.com/nolim1t/126991/ae3a7d36470d2a81190339fbc78821076a4059f7
// https://github.com/ppelleti/https-example/blob/master/https-client.c
// https://stackoverflow.com/questions/40303354/how-to-make-an-https-request-in-c
// https://stackoverflow.com/questions/22077802/simple-c-example-of-doing-an-http-post-and-consuming-the-response
// https://stackoverflow.com/questions/62387069/golang-parse-raw-http-2-response
// https://curl.se/libcurl/c/sendrecv.html

#ifdef _WIN32
#include <Windows.h>
struct timezone
{
    int tz_minuteswest;
    int tz_dsttime;
};

int gettimeofday(struct timeval *tv, struct timezone *tz)
{
    if (tv)
    {
        FILETIME filetime; /* 64-bit value representing the number of 100-nanosecond intervals since January 1, 1601 00:00 UTC */
        ULARGE_INTEGER x;
        ULONGLONG usec;
        static const ULONGLONG epoch_offset_us = 11644473600000000ULL; /* microseconds betweeen Jan 1,1601 and Jan 1,1970 */

#if _WIN32_WINNT >= _WIN32_WINNT_WIN8
        GetSystemTimePreciseAsFileTime(&filetime);
#else
        GetSystemTimeAsFileTime(&filetime);
#endif
        x.LowPart = filetime.dwLowDateTime;
        x.HighPart = filetime.dwHighDateTime;
        usec = x.QuadPart / 10 - epoch_offset_us;
        tv->tv_sec = (time_t)(usec / 1000000ULL);
        tv->tv_usec = (long)(usec % 1000000ULL);
    }
    if (tz)
    {
        TIME_ZONE_INFORMATION timezone;
        GetTimeZoneInformation(&timezone);
        tz->tz_minuteswest = timezone.Bias;
        tz->tz_dsttime = 0;
    }
    return 0;
}
#endif

long long get_current_time()
{
    struct timeval tv;
    gettimeofday(&tv, NULL);
    return (((long long)tv.tv_sec) * 1e6) + (tv.tv_usec);
}

static size_t response_writer(void *data, size_t size, size_t nmemb, void *userp)
{
    size_t realsize = size * nmemb;
    struct memory *mem = (struct memory *)userp;

    char *ptr = (char *)realloc(mem->data, mem->size + realsize + 1);
    if (ptr == NULL)
        return 0; /* out of memory! */

    mem->data = ptr;
    memcpy(&(mem->data[mem->size]), data, realsize);
    mem->size += realsize;
    mem->data[mem->size] = 0;

    return realsize;
}


// uv_loop_t *loop;
// CURLM *curl_handle;
// uv_timer_t timeout;
thread_local uv_loop_t *loop;
thread_local CURLM *curl_handle;
thread_local uv_timer_t timeout;

static curl_context_t *create_curl_context(curl_socket_t sockfd)
{
    curl_context_t *context;

    context = (curl_context_t *)malloc(sizeof(*context));

    context->sockfd = sockfd;

    uv_poll_init_socket(loop, &context->poll_handle, sockfd);
    context->poll_handle.data = context;

    return context;
}

static void curl_close_cb(uv_handle_t *handle)
{
    curl_context_t *context = (curl_context_t *)handle->data;
    free(context);
}

static void destroy_curl_context(curl_context_t *context)
{
    uv_close((uv_handle_t *)&context->poll_handle, curl_close_cb);
}

static void add_request_to_event_loop(request_input *req_input, response_data *response_ref, int debug)
{
    response_ref->uid = req_input->uid;
    response_ref->status_code = -2;
    response_ref->debug = debug;
    if (debug > 0)
    {
        printf("debug_level%d\n", debug);
        printf("%s\n", req_input->url);
    }
    if (debug > 2)
    {
        printf("cookies=>%s\n\n", req_input->cookies);
    }
    if (debug > 2)
    {
        printf("header count=%d\n\n", req_input->headers_len);
        for (int i = 0; i < req_input->headers_len; i++)
        {
            printf("header=>%s\n\n", req_input->headers[i].header);
        }
        printf("body=>%s\n\n", req_input->body);
    }

    CURL *curl;
    curl = curl_easy_init();
    struct curl_slist *header_list = NULL;
    if (req_input->headers_len > 0)
    {
        for (int i = 0; i < req_input->headers_len; i++)
        {
            header_list=curl_slist_append(header_list, req_input->headers[i].header);
        }
    }
    curl_easy_setopt(curl, CURLOPT_PRIVATE, response_ref);
    curl_easy_setopt(curl, CURLOPT_VERBOSE, debug > 3 ? 1L : 0);
    // curl_easy_setopt(curl, CURLOPT_VERBOSE, 1L);
    curl_easy_setopt(curl, CURLOPT_URL, req_input->url);
    curl_easy_setopt(curl, CURLOPT_SSL_VERIFYSTATUS, 0);
    curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);
    curl_easy_setopt(curl, CURLOPT_COOKIE, req_input->cookies);
    curl_easy_setopt(curl, CURLOPT_HTTPHEADER, header_list);
    curl_easy_setopt(curl, CURLOPT_USERAGENT, "cgo benchmark tool");
    if (req_input->body != NULL)
    {
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, req_input->body);
    }
    curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, req_input->method);
    curl_easy_setopt(curl, CURLOPT_TIMEOUT, req_input->time_out_in_sec);

    response_ref->Resp_header = {0};
    response_ref->Resp_body = {0};
    // from response
    curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_writer);
    curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&response_ref->Resp_body);
    curl_easy_setopt(curl, CURLOPT_HEADERDATA, &response_ref->Resp_header);
    curl_easy_setopt(curl, CURLOPT_HEADERFUNCTION, response_writer);

    curl_multi_add_handle(curl_handle, curl);
    response_ref->before_connect_time_microsec = get_current_time();
    if (debug > 1)
    {
        printf("request added to event loop: %s\n", req_input->url);
    }
}

static void on_request_complete(CURLMcode resp)
{
    char *done_url;
    CURLMsg *message;
    int pending;
    CURL *easy_handle;
    response_data *response_ref;

    while ((message = curl_multi_info_read(curl_handle, &pending)))
    {
        switch (message->msg)
        {
        case CURLMSG_DONE:
        {
            /* Do not use message data after calling curl_multi_remove_handle() and
           curl_easy_cleanup(). As per curl_multi_info_read() docs:
           "WARNING: The data the returned pointer points to will not survive
           calling curl_multi_cleanup, curl_multi_remove_handle or
           curl_easy_cleanup." */
            easy_handle = message->easy_handle;
            curl_easy_getinfo(easy_handle, CURLINFO_EFFECTIVE_URL, &done_url);
            // printf("\n---done_url=%s\n", done_url);
            curl_easy_getinfo(easy_handle, CURLINFO_PRIVATE, &response_ref);
            // if (res != CURLE_OK)
            // {
            //     response_ref->status_code = -1;
            //     response_ref->err_code = res;
            // }
            response_ref->after_response_time_microsec = get_current_time();

            curl_off_t start = -1, connect = -1, total = -1;
            struct memory body = {0}, header = {0};
            // from response
            curl_easy_setopt(easy_handle, CURLOPT_WRITEFUNCTION, response_writer);
            curl_easy_setopt(easy_handle, CURLOPT_WRITEDATA, (void *)&body);
            curl_easy_setopt(easy_handle, CURLOPT_HEADERDATA, &header);
            curl_easy_setopt(easy_handle, CURLOPT_HEADERFUNCTION, response_writer);

            curl_easy_getinfo(easy_handle, CURLINFO_RESPONSE_CODE, &response_ref->status_code);
            CURLcode res = curl_easy_getinfo(easy_handle, CURLINFO_CONNECT_TIME_T, &connect);
            if (CURLE_OK != res)
            {
                connect = -1;
            }
            res = curl_easy_getinfo(easy_handle, CURLINFO_STARTTRANSFER_TIME_T, &start);
            if (CURLE_OK != res)
            {
                start = -1;
            }
            res = curl_easy_getinfo(easy_handle, CURLINFO_TOTAL_TIME_T, &total);
            if (CURLE_OK != res)
            {
                total = -1;
            }


            response_ref->connect_time_microsec = connect;
            response_ref->time_to_first_byte_microsec = start;
            response_ref->total_time_from_curl_microsec = total;
            response_ref->total_time_microsec = (response_ref->after_response_time_microsec - response_ref->before_connect_time_microsec);
            response_ref->connected_at_microsec = response_ref->before_connect_time_microsec + connect;
            response_ref->first_byte_at_microsec = response_ref->before_connect_time_microsec + start;
            response_ref->finish_at_microsec = response_ref->before_connect_time_microsec + response_ref->total_time_microsec;

            response_ref->response_header = response_ref->Resp_header.data;
            response_ref->response_body = response_ref->Resp_body.data;


            if(response_ref->status_code==0){
                printf("-- Failed request --\n");
                printf("Resp_header=%s\n",response_ref->Resp_header.data);
                printf("Resp_body=%s\n",response_ref->Resp_body.data);
                printf("seconds to connect=%lf,ttfb=%lf\n",  connect / 1e6, start / 1e6);
            }

            if (response_ref->debug > 2)
            {
                printf("status_code=%d\n", response_ref->status_code);
                printf("before_connect_time_microsec=%lld,after_response_time_microsec=%lld,seconds to connect=%lf,ttfb=%lf,total=%lf.total2=%lld\n", response_ref->before_connect_time_microsec, response_ref->after_response_time_microsec, connect / 1e6, start / 1e6, total / 1e6, response_ref->after_response_time_microsec - response_ref->before_connect_time_microsec);
            }
            if (response_ref->debug > 3)
            {
                printf("%s\n%s\n", header.data, body.data);
            }

            curl_multi_remove_handle(curl_handle, easy_handle);
            curl_easy_cleanup(easy_handle);

            break;
        }

        default:
        {
            fprintf(stderr, "CURLMSG default\n");
            break;
        }
        }
    }

    // free(message);
}

static void curl_perform(uv_poll_t *req, int status, int events)
{
    int running_handles;
    int flags = 0;
    curl_context_t *context;
    CURLMcode res;

    if (events & UV_READABLE)
        flags |= CURL_CSELECT_IN;
    if (events & UV_WRITABLE)
        flags |= CURL_CSELECT_OUT;

    context = (curl_context_t *)req->data;

    res = curl_multi_socket_action(curl_handle, context->sockfd, flags,
                                   &running_handles);

    on_request_complete(res);
}

static void on_timeout(uv_timer_t *req)
{
    int running_handles;
    CURLMcode res;
    // printf("2curl_handle=%ld\n",(long)curl_handle);
    res = curl_multi_socket_action(curl_handle, CURL_SOCKET_TIMEOUT, 0,
                                   &running_handles);
    on_request_complete(res);
}

static int start_timeout(CURLM *multi, long timeout_ms, void *userp)
{
    // printf("timeout_ms->%ld\n", timeout_ms);
    if (timeout_ms < 0)
    {
        uv_timer_stop(&timeout);
    }
    else
    {
        if (timeout_ms == 0)
            timeout_ms = 1; /* 0 means directly call socket_action, but we will do it
                               in a bit */
        uv_timer_start(&timeout, on_timeout, timeout_ms, 0);
    }
    return 0;
}

static int handle_socket(CURL *easy, curl_socket_t s, int action, void *userp,
                         void *socketp)
{
    curl_context_t *curl_context;
    int events = 0;

    switch (action)
    {
    case CURL_POLL_IN:
    case CURL_POLL_OUT:
    case CURL_POLL_INOUT:
        curl_context = socketp ? (curl_context_t *)socketp : create_curl_context(s);

        curl_multi_assign(curl_handle, s, (void *)curl_context);

        if (action != CURL_POLL_IN)
            events |= UV_WRITABLE;
        if (action != CURL_POLL_OUT)
            events |= UV_READABLE;

        uv_poll_start(&curl_context->poll_handle, events, curl_perform);
        break;
    case CURL_POLL_REMOVE:
        if (socketp)
        {
            uv_poll_stop(&((curl_context_t *)socketp)->poll_handle);
            destroy_curl_context((curl_context_t *)socketp);
            curl_multi_assign(curl_handle, s, NULL);
        }
        break;
    default:
        abort();
    }
    return 0;
}

void loop_on_the_thread(request_input *req_inputs, response_data *response_ref, int total_requests, int debug)
{
    loop = uv_default_loop();

    if (curl_global_init(CURL_GLOBAL_ALL))
    {
        fprintf(stderr, "Could not init curl\n");
        return;
    }

    uv_timer_init(loop, &timeout);

    curl_handle = curl_multi_init();
    if(debug>0) printf("curl_handle=%ld\n",(long)curl_handle);
    curl_multi_setopt(curl_handle, CURLMOPT_MAX_HOST_CONNECTIONS, 100); // if the number of connection increased then server may fail to respond, for now fixing it to 100
    // & as i see for less connection server responds fast
    curl_multi_setopt(curl_handle, CURLMOPT_MAX_PIPELINE_LENGTH, total_requests);
    curl_multi_setopt(curl_handle, CURLMOPT_SOCKETFUNCTION, handle_socket);
    curl_multi_setopt(curl_handle, CURLMOPT_TIMERFUNCTION, start_timeout);
    for (int i = 0; i < total_requests; i++)
    {
        add_request_to_event_loop(&(req_inputs[i]), &(response_ref[i]), debug);
    }
    uv_run(loop, UV_RUN_DEFAULT);
    curl_multi_cleanup(curl_handle);
}

