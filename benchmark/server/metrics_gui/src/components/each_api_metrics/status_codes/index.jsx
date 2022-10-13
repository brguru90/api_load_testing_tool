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

export default function StatusCodes({ APIindex }) {

    const chartRef = useRef(null);
    const _iteration_data = useSelector(state => {
        const iteration_data = state.metrics_data?.[APIindex]?.iteration_data
        if (iteration_data?.length) {
            return iteration_data
        }
        return []
    })



    const structure_data = (dt) => {
        const datasets = {}
        dt.forEach((data, index) => {
            Object.entries(data.Status_codes).forEach(([key, value]) => {
                if (datasets[key] == undefined) {
                    datasets[key] = {}
                }
                datasets[key][index] = value
            })
        })
        return {
            labels: dt.map(data => data?.iteration_id + 1),
            datasets: Object.entries(datasets).map(([label, iterations]) => {
                label=String(label)
                const color = {
                    "1xx": {
                        borderColor: 'rgb(245, 186, 37)',
                        backgroundColor: 'rgba(245, 186, 37, 0.5)'
                    },
                    "2xx": {
                        borderColor: 'rgb(8, 201, 18)',
                        backgroundColor: 'rgba(8, 201, 18, 0.5)'
                    },
                    "3xx": {
                        borderColor: 'rgb(245, 186, 37)',
                        backgroundColor: 'rgba(245, 186, 37, 0.5)'
                    },
                    "4xx": {
                        borderColor: 'rgb(245, 41, 27)',
                        backgroundColor: 'rgba(245, 41, 27)'
                    },
                    "5xx": {
                        borderColor: 'rgb(163, 26, 16)',
                        backgroundColor: 'rgba(163, 26, 16)'
                    },
                }
               
                return {
                    label: label,
                    data: dt.map((_, iteration) => iterations[iteration] || 0),
                    yAxisID: 'y',
                    ...color[label[0]+"xx"]
                }
            })
        };
    }

    const [max_items, set_max_items] = useState(10)
    const [pagination, set_pagination] = useState(0)
    let start_index = pagination * max_items;
    const [chartData, setChartData] = useState(structure_data(_iteration_data.slice(start_index, start_index + max_items)))
    useMemo(() => {
        start_index = pagination * max_items;
        setChartData(structure_data(_iteration_data.slice(start_index, start_index + max_items)))
    }, [_iteration_data?.length, pagination, max_items])


    const onScroll = (page) => {
        set_pagination(page)
    }

    useEffect(() => {
        console.log(`Rendered: StatusCodes index=${APIindex}`)
    })

    return (
        <div className="status_codes">
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
