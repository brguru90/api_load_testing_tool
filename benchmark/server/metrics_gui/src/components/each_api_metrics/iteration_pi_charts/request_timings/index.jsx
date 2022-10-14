import React, { useEffect, useMemo } from 'react'
import useExtractIteration from '../../util/useExtractIteration'
import ReactApexChart from "react-apexcharts"
import { chart_option } from "./chart_option"
import "./style.scss"

export default function RequestTimingsPieChart({ APIindex }) {

  const _iteration_data = useExtractIteration({ APIindex })


  const requestTimings = useMemo(() => {
    const requestTimings = {
      connecting: 0,
      processing: 0
    }
    _iteration_data.forEach(iter => {
      requestTimings.connecting += iter["Avg_time_to_connect_api_in_millesec"] || 0
      requestTimings.processing += iter["Avg_time_to_complete_api_in_millesec"] || 0
    });
    requestTimings.connecting /= _iteration_data?.length
    requestTimings.processing /= _iteration_data?.length
    return requestTimings
  }, [_iteration_data?.length])


  useEffect(() => {
    console.log(`Rendered: RequestTimingsPieChart index=${APIindex}`)
  })


  return (
    <div className='request_timing_pi_chart'>
      <ReactApexChart options={chart_option} series={[requestTimings.connecting, requestTimings.processing]} type="pie" width={300} />
    </div>
  )
}
