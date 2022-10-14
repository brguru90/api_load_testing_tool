import React, { useEffect, useMemo } from 'react'
import { useDispatch, useSelector } from 'react-redux'
import styles from "./style.module.scss"
import APIMetricsOverview from './overview'
import TimeToComplete from "./time_to_complete"
import StatusCodes from './status_codes'
import APIPerSecondMetrics from './each_iteration_in_detail'
import RequestTimingsPieChart from './iteration_pi_charts/request_timings'
import StatusCodePieChart from './iteration_pi_charts/status_code'
import TimeForEachIterationPieChart from './iteration_pi_charts/time_for_each_iteration'
import { useHistory } from "react-router-dom";


export default function APIMetrics({ APIindex }) {
    const history = useHistory();
    const dispatch = useDispatch()

    useEffect(() => {
        console.log("APIMetrics: Dashboard")
    })

    const [url, process_id] = useSelector(state => {
        const [url, process_id] = [state.metrics_data?.[APIindex]?.url, state.metrics_data?.[APIindex]?.process_uid]
        return [url || "no url", process_id || "no process_id"]
    }, () => APIindex != undefined)

    const navigateToIterationMetrics=(e)=>{
        e.preventDefault()
        dispatch({
            type: "SET_OTHER",
            payload: {
                last_screen_scroll:window.pageYOffset || document.documentElement.scrollTop
            },
        })
        history.push("iterations_metrics?api_index="+APIindex);
    }

    return (
        <fieldset className={styles['api_metrics']}>
            <legend className={styles['api_metrics_legened']}>{url} - {process_id}</legend>
            <div className={styles['overall_api_metrics']}>
                <div className={styles['overall_api_metric']}>
                    <APIMetricsOverview APIindex={APIindex} url={url} />
                </div>
                <div className={styles['overall_api_metric']}>
                    <TimeToComplete APIindex={APIindex} />
                </div>
                <div className={styles['overall_api_metric']}>
                    <StatusCodes APIindex={APIindex} />
                </div>
                <div className={`${styles['overall_api_metric']} ${styles['pie_charts']}`}>
                    <div className={styles['pie_chart']}>
                        <RequestTimingsPieChart APIindex={APIindex} />
                    </div>
                    <div className={styles['pie_chart']}>
                        <StatusCodePieChart APIindex={APIindex} />
                    </div>
                    <div className={styles['pie_chart']}>
                        <TimeForEachIterationPieChart APIindex={APIindex} />
                    </div>
                </div>
            </div>
            {/* <APIPerSecondMetrics APIindex={APIindex} /> */}

            <a href="" onClick={navigateToIterationMetrics}> view all iteration data</a> <br />
        </fieldset>
    )
}
