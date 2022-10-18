import React, { useMemo, useState } from 'react'
import "./style.scss"
import ReactApexChart from "react-apexcharts"
import { chart_option } from "./chart_option"
import { useSelector } from 'react-redux';
import ChartScrollbar from '../../common_components/chart_scrollbar';



export default function AllAPISReport() {

  const all_data = useSelector(state => {
    return state.metrics_data?.map(data => data.all_data)
  })



  const structure_data = (dt) => {
    return [
      {
        name: "total time to complete",
        data: dt.map(data => Number(data?.Total_time_to_complete_all_apis_iteration_in_sec).toFixed(2)),
      },
      {
        name: "average time to complete",
        data: dt.map(data => Number(data?.Avg_time_to_complete_api_in_sec).toFixed(2)),
      },
      {
        name: "average time to connect",
        data: dt.map(data => Number(data?.Avg_time_to_connect_api_in_sec).toFixed(2)),
      },
    ]
  }


  const chartData = useMemo(() => {
    const _chart_option = Object.assign({}, chart_option)
    _chart_option.xaxis.categories = all_data.map(data => data?.Url)
    return {
      series: structure_data(all_data),
      chart_option: _chart_option,
    }
  }, [all_data?.length])

  return (
    <div className={"all_apis_reports"}>
      AllAPISReport<br />
      <div>
        <ReactApexChart
          options={chartData.chart_option}
          series={chartData.series}
          type="line"
          height={400}
          className="benchmark_line_chart"
        />
      </div>
    </div>
  )
}
