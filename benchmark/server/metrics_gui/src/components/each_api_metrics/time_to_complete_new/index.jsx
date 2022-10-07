import React, { useEffect, useMemo, useRef, useState } from "react"
import ReactApexChart from "react-apexcharts"
import ApexCharts from "apexcharts"
import "./style.scss"
import { useSelector } from "react-redux"
import { chart_option } from "./chart_option"
import ChartScrollbar from "../../../common_components/chart_scrollbar/index.jsx"

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


    const max_items=8
    const [pagination, set_pagination] = useState(0)
    let start_index=pagination*max_items;
    const [chartData, setChartData] = useState({
        series: structure_data(_iteration_data.slice(start_index,start_index+max_items)),
        chart_option: chart_option
    })
    useMemo(() => {
        start_index=pagination*max_items;
        chart_option.xaxis.categories = _iteration_data.slice(start_index,start_index+max_items).map(data => data?.iteration_id + 1)
        setChartData({
            series: structure_data(_iteration_data.slice(start_index,start_index+max_items)),
            chart_option: chart_option
        })
        // setChartData(() => {
        //     const s = structure_data(_iteration_data.slice(start_index,start_index+max_items))
        //     ApexCharts.exec("realtime", "updateSeries", s)
        //     return {
        //         series:s,
        //         chart_option: chart_option
        //     }
        // })
    }, [_iteration_data?.length,pagination])

    const effectCalled = useRef(false)
    useEffect(() => {
        if (!effectCalled.current) {
            effectCalled.current = true
        }
    }, [])


    useEffect(() => {
        console.log(`Rendered: TimeToComplete index=${index}`)
    })

    const onScroll=(page)=>{
        set_pagination(page)
    }

    return (
        <div className="ttc">
            <ReactApexChart
                options={chartData.chart_option}
                series={chartData.series}
                type="line"
                height={400}
                // width={chartData?.chart_option?.xaxis?.categories?.length*50}
                className="benchmark_line_chart"
            // width={600}
            />
            <ChartScrollbar
                scroll_count={Math.ceil(_iteration_data?.length / max_items)}
                onScroll={onScroll}
            />

            {JSON.stringify(chartData.chart_option.xaxis.categories)}
        </div>
    )
}
