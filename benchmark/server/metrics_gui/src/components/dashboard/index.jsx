import React, { useEffect, useRef, useState } from "react"
import { Link } from "react-router-dom"
import APIMetrics from "../each_api_metrics/layout"
import styles from "./style.module.scss"
import { useDispatch, useSelector } from "react-redux"


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
                            .map((e,i) => <APIMetrics key={i} index={i} />)
                    }
                </div>
            </div>
            <Link to="page2">view Page2</Link> <br />
        </div>
    )
}
