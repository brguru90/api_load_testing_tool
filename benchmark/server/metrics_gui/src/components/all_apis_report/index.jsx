import React, { useMemo, useState } from 'react'
import "./style.scss"
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
import { useSelector } from 'react-redux';
import ChartScrollbar from '../../common_components/chart_scrollbar';

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend
);


export default function AllAPISReport() {

  const all_data = useSelector(state => {
    return state.metrics_data?.map(data => data.all_data)
  })

  const structure_data = (dt) => {
    return {
      labels: dt.map(data => {
        let _url=[]
        for(let i=0;i<=data?.Url?.length;i+=20){
          _url.push(data?.Url?.slice(i,i+20))
        }
        return _url;
      }),
      datasets: [
        {
          label: 'total time to complete complete APIs',
          data: dt.map(data => data?.Total_time_to_complete_all_apis_iteration_in_sec),
          borderColor: 'rgb(53, 162, 235)',
          backgroundColor: 'rgba(53, 162, 235, 0.5)',
          yAxisID: 'y',
        },
        {
          label: 'average time to complete APIs',
          data: dt.map(data => data?.Avg_time_to_complete_api_in_sec),
          borderColor: 'rgb(8, 201, 18)',
          backgroundColor: 'rgba(8, 201, 18, 0.5)',
          yAxisID: 'y',
        },
        {
          label: 'average time to connect APIs',
          data: dt.map(data => data?.Avg_time_to_connect_api_in_sec),
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
  const [chartData, setChartData] = useState(structure_data(all_data.slice(start_index, start_index + max_items)))
  useMemo(() => {
    start_index = pagination * max_items;
    console.log("start_index", start_index)
    setChartData(structure_data(all_data.slice(start_index, start_index + max_items)))
  }, [all_data?.length, pagination, max_items])

  return (
    <div className={"all_apis_reports"}>
      AllAPISReport<br />
      <div>
        <label>
          Page Size: <input type="number" value={max_items} onChange={e => set_max_items(Math.max(4, e.target.value))} />
        </label>
        <Line
          options={chart_option}
          data={chartData}
          height={150}
          className="benchmark_line_chart"
        />
      </div>
      <ChartScrollbar
        scroll_count={Math.ceil(all_data?.length / max_items)}
        onScroll={page => set_pagination(page)}
        className="chart_scrollbar"
      />
    </div>
  )
}
