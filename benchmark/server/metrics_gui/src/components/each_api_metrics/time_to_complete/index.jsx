import React, { useEffect, useMemo, useRef, useState } from "react"
import ReactApexChart from "react-apexcharts"
import ApexCharts from "apexcharts"
import "./style.scss"
import { useSelector } from "react-redux"
import { chart_option } from "./chart_option"

export default function TimeToComplete({ index }) {

    const _iteration_data = useSelector(state => {
        const iteration_data = state.metrics_data?.[index]?.iteration_data
        if (iteration_data?.length) {
            return iteration_data
        }
        return []
    })


    const structure_data = (dt) => {
        return [
            {
                name: "total time to complete",
                data: dt.map(data => data?.Total_time_to_complete_all_apis_in_millesec),
            },
            {
                name: "average time to complete",
                data: dt.map(data => data?.Avg_time_to_complete_api_in_millesec),
            },
            {
                name: "average time to connect",
                data: dt.map(data => data?.Avg_time_to_connect_api_in_millesec),
            },
        ]
    }


    const [chartData, setChartData] = useState({
        series:structure_data(_iteration_data),
        chart_option:chart_option
    })
    useMemo(() => {
        chart_option.xaxis.categories=_iteration_data.map(data => data?.iteration_id+1)
        setChartData({
            series:structure_data(_iteration_data),
            chart_option:chart_option
        })
        // setSeries(() => {
        //     const s = structure_data(_iteration_data)
        //     ApexCharts.exec("realtime", "updateSeries", s)
        //     return s
        // })
    }, [_iteration_data?.length])






    const effectCalled = useRef(false)
    useEffect(() => {
        if (!effectCalled.current) {
            effectCalled.current = true
        }
    }, [])


    useEffect(() => {
        console.log(`Rendered: TimeToComplete index=${index}`)
    })

    return (
        <div className="ttc">
            <ReactApexChart
                options={chartData.chart_option}
                series={chartData.series}
                type="line"
                height={600}
                className="benchmark_line_chart"
            // width={600}
            />
        </div>
    )
}
