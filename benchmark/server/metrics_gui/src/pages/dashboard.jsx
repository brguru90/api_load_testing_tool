import React, {useEffect, useRef, useState} from "react"
import {Link} from "react-router-dom"
import Dashboard from "../components/dashboard/index.jsx"
import {useDispatch, useSelector} from "react-redux"
import {GetBenchmarkMetrics} from "../services/metric.data"
GetBenchmarkMetrics

export default function dashboard_page() {
    const effectCalled = useRef(false)

    const dispatch = useDispatch()

    useEffect(() => {
        if (!effectCalled.current) {
            effectCalled.current = true
            GetBenchmarkMetrics((data) => {
                dispatch({
                    type: "SET_METRICS",
                    payload: data,
                })
            })
        }
    }, [])

    return <Dashboard />
}
