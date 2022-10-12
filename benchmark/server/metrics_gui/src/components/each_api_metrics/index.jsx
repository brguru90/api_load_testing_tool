import React, { useEffect, useMemo } from 'react'
import { useSelector } from 'react-redux'
import styles from "./style.module.scss"
import APIMetricsOverview from './overview'
import TimeToComplete from "./time_to_complete"
import StatusCodes from './status_codes'
import APIPerSecondMetrics from './each_iteration_in_detail'


export default function APIMetrics({ index }) {
    useEffect(() => {
        console.log("APIMetrics: Dashboard")
    })

    const [url, process_id] = useSelector(state => {
        const [url, process_id] = [state.metrics_data?.[index]?.url, state.metrics_data?.[index]?.process_uid]
        return [url || "no url", process_id || "no process_id"]
    }, () => index != undefined)

    return (
        <fieldset className={styles['api_metrics']}>
            <legend className={styles['api_metrics_legened']}>{url} - {process_id}</legend>
            <div className={styles['overall_api_metrics']}>
                <div className={styles['overall_api_metric']}>
                    <APIMetricsOverview index={index} url={url} />
                </div>
                <div className={styles['overall_api_metric']}>
                    <TimeToComplete index={index}/>
                </div>
                <div className={styles['overall_api_metric']}>
                    <StatusCodes index={index}/>
                </div>
                <div className={styles['overall_api_metric']}>
                    pi charts
                </div>
            </div>
            <APIPerSecondMetrics index={index} />
        </fieldset>
    )
}
