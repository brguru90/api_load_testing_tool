import React from 'react'
import styles from "./style.module.scss"
import TimeToComplete from "../each_api_benchmark_summary/time_to_complete"


export default function APIMetrics({ index }) {
    return (
        <fieldset className={styles['api_metrics']}>
            <legend className={styles['api_metrics_legened']}>API URL-{index}...</legend>
            <div className={styles['overall_api_metrics']}>
                <div className={styles['overall_api_metric']}>
                    summary

                </div>
                <div className={styles['overall_api_metric']}>
                    each iteration request timings, line chart

                </div>
                <div className={styles['overall_api_metric']}>
                    each iteration status codes, line chart
                </div>
                <div className={styles['overall_api_metric']}>
                    pi charts
                </div>
            </div>
        </fieldset>
    )
}
