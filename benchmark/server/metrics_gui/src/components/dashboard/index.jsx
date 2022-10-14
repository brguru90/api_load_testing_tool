import React, { useEffect } from "react"
import APIMetrics from "../each_api_metrics/index.jsx"
import styles from "./style.module.scss"
import { useSelector } from "react-redux"


export default function Dashboard() {
    const test_cases_len = useSelector(state => state.metrics_data?.length || 0)

    useEffect(() => {
        console.log("Rendered: Dashboard")
    })


    return (
        <div className={styles["dashbaord"]}>
            <div className={styles["test_cases"]}>
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
