import React, { useEffect, useMemo } from 'react'
import { useSelector } from 'react-redux'
import ReactApexChart from "react-apexcharts"
import { chart_option } from "./chart_option"


export default function TimeForEachIterationPieChart({ APIindex }) {

  const _iteration_data = useSelector(state => {
    const iteration_data = state.metrics_data?.[APIindex]?.iteration_data
    if (iteration_data?.length) {
      return iteration_data
    }
    return []
  })


  const duration_of_iterations = useMemo(() => {
    return _iteration_data.map(iter => Number((iter.Total_time_to_complete_all_apis_in_millesec/1000).toFixed(2)))
  }, [_iteration_data?.length])
  const _chart_option = Object.assign({}, chart_option)
  _chart_option.labels = _iteration_data?.map(({iteration_id})=>iteration_id+1) || []

  useEffect(() => {
    console.log(`Rendered: StatusCodePieChart index=${APIindex}`)
  })

  return (
    <div>
      <ReactApexChart options={_chart_option} series={duration_of_iterations || []} type="pie" width={500} height={500} />
    </div>
  )
}
