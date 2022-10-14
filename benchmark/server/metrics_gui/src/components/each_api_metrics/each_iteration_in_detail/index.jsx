import React from 'react'
import { useSelector } from 'react-redux'
import { useHistory } from 'react-router-dom';
import APITimes from './apis_time'
import styles from "./style.module.scss"

export default function APIPerSecondMetrics({ APIindex }) {
    const history = useHistory();
    const [url, iteration_ids] = useSelector(state => [state.metrics_data?.[APIindex]?.url, state.metrics_data?.[APIindex]?.iteration_data?.map(item => item?.iteration_id)])

    const navigateToDashBoard = (e) => {
        e.preventDefault()
        history.push("")
    }

    return (
        <div className={styles['per_second_metrics']}>
            <h1 className={styles['title']}>API Per Second Metrics - {url}</h1>
            <div className={styles['per_second_metrics_set']}>
                {
                    iteration_ids?.map(id => {
                        return <div key={id} className={styles['per_second_data']}>
                            <APITimes />
                        </div>
                    })
                }
            </div>

            <a href="" onClick={navigateToDashBoard}>
                Back to dashboard
            </a>
        </div>
    )
}
