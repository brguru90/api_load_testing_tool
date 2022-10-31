import React, { useEffect, useMemo, useState } from 'react'
import styles from "./style.module.scss"
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
import { chart_option } from "./chart_option"
import ChartScrollbar from "../../../../common_components/chart_scrollbar/index.jsx"


ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);

export default function APITimes({ iteration }) {

  useEffect(() => {
    console.log(`Rendered: APITimes index=${iteration.iteration_id}`)
  })

  const structure_data = (_Benchmark_per_second_metric) => {
    
    return {
      labels: _Benchmark_per_second_metric.map(data => {
        const dt=new Date(data.To_time_duration)
        return dt.toLocaleTimeString()+":"+dt.getMilliseconds()
      }),
      datasets: [
        {
          label: 'Request sent',
          data: _Benchmark_per_second_metric.map(data => data?.Request_sent || 0),
          borderColor: 'rgb(53, 162, 235)',
          backgroundColor: 'rgba(53, 162, 235, 0.5)',
          yAxisID: 'y',
        },
        {
          label: 'Request connected',
          data: _Benchmark_per_second_metric.map(data => data?.Request_connected || 0),
          borderColor: 'rgb(245, 186, 27)',
          backgroundColor: 'rgba(245, 186, 27, 0.5)',
          yAxisID: 'y',
        },
        {
          label: 'Request receives first_byte',
          data: _Benchmark_per_second_metric.map(data => data?.Request_receives_first_byte || 0),
          borderColor: 'rgb(237, 112, 9)',
          backgroundColor: 'rgba(237, 112, 9, 0.5)',
          yAxisID: 'y',
        },
        {
          label: 'Request processed',
          data: _Benchmark_per_second_metric.map(data => data?.Request_processed || 0),
          borderColor: 'rgb(8, 201, 18)',
          backgroundColor: 'rgba(8, 201, 18, 0.5)',
          yAxisID: 'y',
        },
      ],
    };
  }

  const [max_items, set_max_items] = useState(10)
  const [pagination, set_pagination] = useState(0)
  let start_index = pagination * max_items;

  const chartData = useMemo(() => {
    start_index = pagination * max_items;
    if (iteration.Benchmark_per_second_metric?.length) {
      return structure_data([{
        To_time_duration: iteration.Benchmark_per_second_metric[0].From_time_duration
      }, ...iteration.Benchmark_per_second_metric].slice(start_index, start_index + max_items))

    }
    return structure_data([])

  }, [iteration?.Benchmark_per_second_metric?.length,pagination, max_items])


  const onScroll = (page) => {
    set_pagination(page)
  }

  return (
    <div className={styles["apis_time"]}>
      <div>
        <b>Iteration {iteration?.iteration_id+1}</b><br />
        {Math.round((1 / (iteration?.Avg_time_to_complete_api_in_millesec/1000))*iteration?.Concurrent_request) || "-"} Req/Sec<br /><br />
        <label>
          Page Size: <input type="number" value={max_items} onChange={e => set_max_items(Math.max(4, e.target.value))} />
        </label>
        <Line
          options={chart_option}
          data={chartData}
          height={200}
          className="benchmark_line_chart"
        />
      </div>
      <ChartScrollbar
        scroll_count={Math.ceil(iteration?.Benchmark_per_second_metric?.length / max_items)}
        onScroll={onScroll}
        className="chart_scrollbar"
      />


    </div>
  )
}
