import React, { useEffect, useMemo, useRef, useState } from "react"
import {
    Chart as ChartJS,
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend,
} from 'chart.js';
import { Chart, Line } from 'react-chartjs-2';

import "./style.scss"
import { useSelector } from "react-redux"
import { chart_option } from "./chart_option"
import ChartScrollbar from "../../../common_components/chart_scrollbar/index.jsx"


ChartJS.register(
    CategoryScale,
    LinearScale,
    PointElement,
    LineElement,
    Title,
    Tooltip,
    Legend
);

export default function TimeToComplete({ index }) {

    const chartRef = useRef(null);
    const _iteration_data = useSelector(state => {
        const iteration_data = state.metrics_data?.[index]?.iteration_data
        if (iteration_data?.length) {
            return iteration_data
        }
        return []
    })



    const structure_data = (dt) => {
        return {
            labels: dt.map(data => data?.iteration_id + 1),
            datasets: [
                {
                    label: 'total time to complete',
                    data: dt.map(data => data?.Total_time_to_complete_all_apis_in_millesec),
                    borderColor: 'rgb(53, 162, 235)',
                    backgroundColor: 'rgba(53, 162, 235, 0.5)',
                    yAxisID: 'y',
                },
                {
                    label: 'average time to complete',
                    data: dt.map(data => data?.Avg_time_to_complete_api_in_millesec),
                    borderColor: 'rgb(8, 201, 18)',
                    backgroundColor: 'rgba(8, 201, 18, 0.5)',
                    yAxisID: 'y',
                },
                {
                    label: 'average time to connect',
                    data: dt.map(data => data?.Avg_time_to_connect_api_in_millesec),
                    borderColor: 'rgb(245, 186, 37)',
                    backgroundColor: 'rgba(245, 186, 37, 0.5)',
                    yAxisID: 'y',
                },
            ],
        };
    }



    const [max_items, set_max_items] = useState(10)
    const [pagination, set_pagination] = useState(0)
    let start_index = pagination * max_items;
    const [chartData, setChartData] = useState(structure_data(_iteration_data.slice(start_index, start_index + max_items)))
    useMemo(() => {
        start_index = pagination * max_items;
        setChartData(structure_data(_iteration_data.slice(start_index, start_index + max_items)))
        // setSeries(() => {
        //     const s = structure_data(_iteration_data)
        //     ApexCharts.exec("realtime", "updateSeries", s)
        //     return s
        // })
    }, [_iteration_data?.length, pagination,max_items])


    const onScroll = (page) => {
        console.log("page",page)
        set_pagination(Math.round(page))
    }

    // useEffect(() => {
    //     chartRef?.current?.update()
    // }, [chartData])

    





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
            <div>
               <label>
                Page Size: <input type="number" value={max_items} onChange={e => set_max_items(Math.max(4,e.target.value))} /> 
               </label>
                <Line
                    options={chart_option}
                    data={chartData}
                    height={200}
                    className="benchmark_line_chart"
                    ref={chartRef}
                />
            </div>
            <ChartScrollbar
                scroll_count={Math.ceil(_iteration_data?.length / max_items)}
                onScroll={onScroll}
                className="chart_scrollbar"
            />
        </div>
    )
}
