import React, { useEffect, useMemo } from 'react'
import ReactApexChart from "react-apexcharts"
import { chart_option } from "./chart_option"
import "./style.scss"
import useExtractIteration from '../../util/useExtractIteration'

export default function TimeForEachIterationPieChart({ APIindex }) {

  const _iteration_data = useExtractIteration({APIindex})



  const duration_of_iterations = useMemo(() => {
    return _iteration_data.map(iter => Number((iter.Total_time_to_complete_all_apis_in_millesec/1000).toFixed(2)))
  }, [_iteration_data?.length])
  const _chart_option = Object.assign({}, chart_option)
  _chart_option.labels = _iteration_data?.map(({iteration_id})=>iteration_id+1) || []

  useEffect(() => {
    console.log(`Rendered: StatusCodePieChart index=${APIindex}`)
  })

  return (
    <div className='time_for_each_duration_pie_chart'>
      <ReactApexChart options={_chart_option} series={duration_of_iterations || []} type="pie" width={400} />
    </div>
  )
}
