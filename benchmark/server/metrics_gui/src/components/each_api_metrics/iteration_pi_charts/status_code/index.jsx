import React, { useMemo } from 'react'
import ReactApexChart from "react-apexcharts"
import ApexCharts from "apexcharts"
import { useSelector } from 'react-redux'
import { chart_option } from "./chart_option"

export default function StatusCodePieChart({ APIindex }) {

  const _iteration_data = useSelector(state => {
    const iteration_data = state.metrics_data?.[APIindex]?.iteration_data
    if (iteration_data?.length) {
      return iteration_data
    }
    return []
  })

  const status_code_coverage = useMemo(() => {
    const StatusCodesInPerc = {}
    console.log("StatusCodePieChart _iteration_data", _iteration_data)
    _iteration_data.forEach(({ Status_code_in_percentage }) => {
      Object.entries(Status_code_in_percentage || {}).forEach(([code, perc]) => {
        if (StatusCodesInPerc[code] == undefined) {
          StatusCodesInPerc[code] = 0
        }
        StatusCodesInPerc[code] += perc
      });
    })

    Object.entries(StatusCodesInPerc).forEach(([code, total_perc]) => {
      StatusCodesInPerc[code] = Number((total_perc / _iteration_data?.length).toFixed(2))
    })

    console.log("StatusCodesInPerc", StatusCodesInPerc)
    return StatusCodesInPerc

  }, [_iteration_data?.length])


  const color = {
    "1xx": 'rgb(11, 123, 214)',
    "2xx": 'rgb(8, 201, 18)',
    "3xx": 'rgb(245, 186, 37)',
    "4xx": 'rgb(245, 41, 27)',
    "5xx": 'rgb(163, 26, 16)',
  }

  const _chart_option = Object.assign({}, chart_option)
  _chart_option.labels = Object.keys(status_code_coverage || {})
  _chart_option.fill = {
    colors: Object.keys(status_code_coverage || {}).map(val => {
      return color[String(val)[0] + "xx"]
    })
  }

  return (
    <div>
      <ReactApexChart options={_chart_option} series={Object.values(status_code_coverage || {})} type="pie" width={300} />
    </div>
  )
}
