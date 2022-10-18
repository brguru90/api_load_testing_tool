import React, { useEffect } from "react"
import APIMetrics from "../each_api_metrics/index.jsx"
import styles from "./style.module.scss"
import { useSelector } from "react-redux"
import AllAPISReport from "../all_apis_report/index.jsx"

export default function Dashboard() {
    const test_cases_len = useSelector(state => state.metrics_data?.length || 0)
    const benchmark_finished = useSelector(state => state.metrics_extra?.benchmark_finished)

    useEffect(() => {
        console.log("Rendered: Dashboard")
    })


    return (
        <div className={styles["dashbaord"]}>
            <div className={styles["test_cases"]}>
                <div className={styles["test_case"]}>
                    {
                        benchmark_finished && <AllAPISReport />
                    }
                </div>
                <div className={styles["test_case"]}>
                    {
                        Array.from({ length: test_cases_len }, (_, i) => i + 1)
                            .map((e, i) => <APIMetrics key={i} APIindex={i} />)
                    }
                </div>
            </div>
        </div>
    )
}
