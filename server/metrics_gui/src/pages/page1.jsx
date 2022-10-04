import React, {useEffect, useRef, useState} from "react"
import {Link} from "react-router-dom"
import {GetBenchmarkMetrics} from "../services/metric.data"
import ReactApexChart from "react-apexcharts"
import ApexCharts from "apexcharts"

export default function Page1() {
    const effectCalled = useRef(false)
    const [series, setSeries] = useState([
        {
            name: "Request success rate",
            data: [10, 20, 30, 40],
        },
        {
            name: "Request proccessed per Second",
            data: [15, 25, 35, 45],
        },
    ])

    const chartOption = {
        chart: {
            id: "realtime",
            height: 350,
            type: "line",
            dropShadow: {
                enabled: true,
                color: "#000",
                top: 18,
                left: 7,
                blur: 10,
                opacity: 0.2,
            },
            animations: {
                enabled: true,
                easing: "easeout",
                speed: 1000,
                animateGradually: {
                    enabled: true,
                    delay: 150,
                },
                dynamicAnimation: {
                    enabled: true,
                    speed: 800,
                },
            },
            toolbar: {
                show: true,
                offsetX: -65,
                offsetY: -15,
                tools: {
                    download: true,
                    pan: true,
                    selection: true,
                    zoom: true,
                    zoomin: true,
                    zoomout: true,
                    reset: true,
                },
                autoSelected: "zoom",
            },
            zoom: {
                enabled: true,
            },
        },
        dataLabels: {
            enabled: true,
        },
        stroke: {
            curve: "straight",
        },
        title: {
            // text: 'Server load balancing',
            align: "left",
            margin: 40,
            offsetX: 0,
            offsetY: -20,
            floating: false,
            style: {
                fontSize: "16px",
                fontWeight: "bold",
                fontFamily: undefined,
                color: "#263238",
            },
        },
        grid: {
            borderColor: "#e7e7e7",
            row: {
                colors: ["#f3f3f3", "transparent"], // takes an array which will be repeated on columns
                opacity: 0.5,
            },
        },
        markers: {
            size: 4,
            strokeWidth: 0,
        },
        xaxis: {
            type: "category",
            title: {
                text: "Request per second (total request is 2000)",
                style: {
                    fontSize: "10px",
                    color: "#78787D",
                },
            },
        },
        yaxis: {
            title: {
                text: "Maximum reach in %",
                style: {
                    fontSize: "10px",
                    color: "#78787D",
                },
            },
            labels: {
                rotate: -45,
            },
        },

        legend: {
            position: "top",
            horizontalAlign: "right",
            floating: true,
            offsetY: -25,
            offsetX: -5,
        },
    }

    useEffect(() => {
        if (!effectCalled.current) {
            effectCalled.current = true
            GetBenchmarkMetrics(() => {})
            let count = 20
            const interval = setInterval(() => {
                if (--count < 0) {
                    clearInterval(interval)
                }
                setSeries((series) => {
                    series[0].data = [...series[0].data, Math.round(Math.random() * 100)]
                    series[1].data = [...series[1].data, Math.round(Math.random() * 100)]
                    ApexCharts.exec("realtime", "updateSeries", [
                        {
                            name: "Request success rate",
                            data: series[0].data,
                        },
                        {
                            name: "Request proccessed per Second",
                            data: series[1].data,
                        },
                    ])
                    return series
                })
            }, 1000)
        }
    }, [])

    return (
        <div>
            Page1 <br />
            <Link to="page2">view Page2</Link> <br />
            <ReactApexChart
                options={chartOption}
                series={series}
                type="line"
                height={600}
                // width={600}
            />
        </div>
    )
}
