import React from 'react'
import { useSelector } from 'react-redux'
import "./style.scss"

export default function APIMetricsOverview({ APIindex, url }) {

  const overview_data = useSelector(state => state.metrics_data?.[APIindex]?.all_data || {})

  return (
    <div className='api_metrics_overview'>
      <div className='title'>Overview:</div>
      <div className='tbl table-responsive'>
        <table className='table table-striped table-hover'>
          <tbody>
            <tr>
              <th>URL:</th>
              <td>{overview_data?.Url || url}</td>
            </tr>
            <tr>
              <th>API Success rate:</th>
              <td>{overview_data?.Status_code_in_percentage?.[200] || "-"}%</td>
            </tr>
            <tr>
              <th>Concurrent request:</th>
              <td>{overview_data?.Concurrent_request || "-"}</td>
            </tr>
            <tr>
              <th>Total request:</th>
              <td>{overview_data?.Total_number_of_request || "-"}</td>
            </tr>
            <tr>
              <th>Average time to connect API:</th>
              <td>{overview_data?.Avg_time_to_connect_api_in_sec || "-"} Seconds
              </td>
            </tr>
            <tr>
              <th>Average API response time:</th>
              <td>
                {overview_data?.Avg_time_to_complete_api_in_sec || "-"} Seconds
                {/* ({Math.round((1 / overview_data?.Avg_time_to_complete_api_in_sec)*overview_data?.Concurrent_request) || "-"} Req/Sec) */}
              </td>
            </tr>
            <tr>
              <th>Minimum API response time:</th>
              <td>{overview_data?.Min_time_to_complete_api_in_sec || "-"} Seconds</td>
            </tr>
            <tr>
              <th>Maximum API response time:</th>
              <td>{overview_data?.Max_time_to_complete_api_in_sec || "-"} Seconds</td>
            </tr>
            <tr>
              <th>Total time to complete all iteration:</th>
              <td>{overview_data?.Total_time_to_complete_all_apis_iteration_in_sec || "-"} Seconds
                ({overview_data?.Total_number_of_request / overview_data?.Total_time_to_complete_all_apis_iteration_in_sec || "-"} Req/Sec)
              </td>
            </tr>
            <tr>
              <th>Average request payload size:</th>
              <td>{overview_data?.Average_request_payload_size_in_bytes_in_all_iteration / 1024 || "-"} KB</td>
            </tr>
            <tr>
              <th>Average response payload size:</th>
              <td>{overview_data?.Average_response_payload_size_in_bytes_in_all_iteration / 1024 || "-"} KB</td>
            </tr>
            <tr>
              <th>Total operation time:</th>
              <td>{overview_data?.Total_operation_time_in_sec || "-"} Seconds</td>
            </tr>
          </tbody>
        </table>
      </div>


    </div>
  )
}
