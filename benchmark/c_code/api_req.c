#include "api_req.h"

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

#define NUM_THREAD 5

void *goCallback_wrap(void *vargp)
{
    int *myid = (int *)vargp;
    printf("tid=%d\n", *myid);
    // goCallback(*myid);
    pthread_exit(NULL);
}

void run_bulk_api_request(char *s)
{
    printf("%s\n", s);

    int i, nor_of_thread;
    // pthread_t threads[NUM_THREAD];
    printf("Enter number of thread\n");
    scanf("%d", &nor_of_thread);
    printf("Entered %d\n", nor_of_thread);
    pthread_t *threads = malloc(sizeof(pthread_t) * nor_of_thread);
    pthread_t tid;

    for (i = 0; i < nor_of_thread; i++)
    {
        pthread_create(&threads[i], NULL, goCallback_wrap, (void *)&threads[i]);
    }

    for (i = 0; i < nor_of_thread; i++)
    {
        pthread_join(threads[i], NULL);
    }

    // printf("Main thread exiting...\n");
    // pthread_exit(NULL);

    goCallback(-1);
}

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

    char *ptr = realloc(mem->data, mem->size + realsize + 1);
    if (ptr == NULL)
        return 0; /* out of memory! */

    mem->data = ptr;
    memcpy(&(mem->data[mem->size]), data, realsize);
    mem->size += realsize;
    mem->data[mem->size] = 0;

    return realsize;
}

void send_raw_request(request_input *req_input, response_data *response_ref, int debug)
{
    response_data res_data;
    res_data.status_code = -2;

    switch (debug)
    {
    case 2:
        printf("header count=%d\n\n", req_input->headers_len);
        printf("%s\n\n", req_input->headers[0].header);
        printf("body=%s\n\n", req_input->body);
    case 1:
        printf("debug_level%d\n", debug);
        printf("%s\n", req_input->url);
        break;

    default:
        break;
    }

    CURL *curl;
    CURLcode res;
    curl_global_init(CURL_GLOBAL_ALL);
    curl = curl_easy_init();
    struct curl_slist *header_list = NULL;
    if (req_input->headers_len > 0)
    {
        for (int i = 0; i < req_input->headers_len; i++)
        {
            curl_slist_append(header_list, req_input->headers[i].header);
        }
    }
    // struct curl_slist *header_list = curl_slist_append(NULL, "Content-Type: text/html");
    if (curl)
    {
        long response_code;
        curl_off_t start = -1, connect = -1, total = -1;
        struct memory body = {0}, header = {0};
        // to request
        curl_easy_setopt(curl, CURLOPT_VERBOSE, debug > 2 ? 1L : 0);
        curl_easy_setopt(curl, CURLOPT_URL, req_input->url);
        curl_easy_setopt(curl, CURLOPT_SSL_VERIFYSTATUS, 0);
        curl_easy_setopt(curl, CURLOPT_FOLLOWLOCATION, 1L);
        curl_easy_setopt(curl, CURLOPT_HTTPHEADER, header_list);
        curl_easy_setopt(curl, CURLOPT_USERAGENT, "cgo benchmark tool");
        curl_easy_setopt(curl, CURLOPT_POSTFIELDS, req_input->body);
        curl_easy_setopt(curl, CURLOPT_CUSTOMREQUEST, req_input->method);

        // from response
        curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, response_writer);
        curl_easy_setopt(curl, CURLOPT_WRITEDATA, (void *)&body);
        curl_easy_setopt(curl, CURLOPT_HEADERDATA, &header);
        curl_easy_setopt(curl, CURLOPT_HEADERFUNCTION, response_writer);

        // need to map all inputs like
        // Header
        // query param
        // method
        // body
        // blob files
        // keep alive
        // time out
        // chunked, see?
        // UA

        res_data.before_connect_time_microsec = get_current_time();
        res = curl_easy_perform(curl);
        /* Check for errors */
        if (res != CURLE_OK)
        {
            if (debug > 1)
            {
                fprintf(stderr, "curl_easy_perform() failed: %s\n",
                        curl_easy_strerror(res));
            }
            response_code = -1;
        }
        curl_easy_getinfo(curl, CURLINFO_RESPONSE_CODE, &response_code);
        res = curl_easy_getinfo(curl, CURLINFO_CONNECT_TIME_T, &connect);
        if (CURLE_OK != res)
        {
            connect = -1;
        }
        res = curl_easy_getinfo(curl, CURLINFO_STARTTRANSFER_TIME_T, &start);
        if (CURLE_OK != res)
        {
            start = -1;
        }
        res = curl_easy_getinfo(curl, CURLINFO_TOTAL_TIME_T, &total);
        if (CURLE_OK != res)
        {
            total = -1;
        }
        if (debug > 1)
        {
            printf("status_code=%ld\n", response_code);
            printf("before_connect_time_microsec=%lld,seconds to connect=%lf,ttfb=%lf,total=%lf\n",res_data.before_connect_time_microsec, connect / 1e6, start / 1e6, total / 1e6);
        }
        if (debug > 2)
        {
            printf("%s\n%s\n", header.data, body.data);
        }

        res_data.status_code = response_code;
        res_data.connect_time_microsec = connect;
        res_data.time_at_first_byte_microsec = start;
        res_data.total_time_microsec = total;
        res_data.connected_at_microsec=res_data.before_connect_time_microsec+connect;
        res_data.first_byte_at_microsec=res_data.before_connect_time_microsec+start;
        res_data.finish_at_microsec=res_data.before_connect_time_microsec+total;

        res_data.response_header = header.data;
        res_data.response_body = body.data;

        curl_easy_cleanup(curl);
    }
    curl_global_cleanup();
    *response_ref = res_data;
}

void *loop_on_the_thread(void *data)
{
    thread_data *td = (thread_data *)data;
    for (int i = td->th_pool_data.start_index; i <= td->th_pool_data.end_index; i++)
    {
        send_raw_request(&(td->req_inputs_ptr[i]), &(td->response_ref_ptr[i]), td->debug_flag);
    }
}

void send_request_in_concurrently(request_input *req_inputs, response_data *response_ref, int total_requests, int total_threads, int debug)
{

    int num_of_threads = total_requests >= total_threads ? total_threads : total_requests;
    int max_work_on_thread = floor((float)total_requests / num_of_threads);
    int left_out_work = total_requests % num_of_threads;

    printf("total_requests=%d,total_threads=%d,num_of_threads=%d,max_work_on_thread=%d,left_out_work=%d\n", total_requests, total_threads, num_of_threads, max_work_on_thread, left_out_work);

    thread_pool_data proc_data[left_out_work == 0 ? num_of_threads : num_of_threads + 1];

    for (int p = 0; p < num_of_threads; p++)
    {
        proc_data[p].start_index = p * max_work_on_thread;
        proc_data[p].end_index = (proc_data[p].start_index + max_work_on_thread) - 1;
        proc_data[p].full_index = false;
    }
    if (left_out_work > 0)
    {
        proc_data[num_of_threads].start_index = num_of_threads * max_work_on_thread;
        proc_data[num_of_threads].end_index = (proc_data[num_of_threads].start_index + left_out_work) - 1;
        proc_data[num_of_threads].full_index = false;
    }

    int thread_size = (left_out_work == 0 ? num_of_threads : num_of_threads + 1);
    pthread_t *threads = malloc(sizeof(pthread_t) * thread_size);
    thread_data *threads_data = malloc(sizeof(thread_data) * thread_size);

    for (int i = 0; i < thread_size; i++)
    {
        threads_data[i].req_inputs_ptr = req_inputs;
        threads_data[i].response_ref_ptr = response_ref;
        threads_data[i].debug_flag = debug;
        threads_data[i].thread_id = i;
        threads_data[i].th_pool_data = proc_data[i];
    }

    for (int p = 0; p < thread_size; p++)
    {
        if (pthread_create(&threads[p], NULL, loop_on_the_thread, (void *)&threads_data[p]) != 0)
        {
            perror("pthread_create() error");
            exit(1);
        }
    }

    for (int i = 0; i < thread_size; i++)
    {
        pthread_join(threads[i], NULL);
    }
    free(threads);
    free(threads_data);
}