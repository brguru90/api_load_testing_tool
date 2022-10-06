import React from 'react'
import { useSelector } from 'react-redux'
import APITimes from './apis_time'
import styles from "./style.module.scss"

export default function APIPerSecondMetrics({index}) {
    const  iteration_ids= useSelector(state => state.metrics_data?.[index]?.iteration_data?.map(item => item?.iteration_id))
    return (
        <div className={styles['per_second_metrics']}>
            <h1 className={styles['title']}>API Per Second Metrics</h1>
            <div className={styles['per_second_metrics_set']}>
            {
                iteration_ids.map(id=>{
                    return <div key={id} className={styles['per_second_data']}>
                        <APITimes />
                    </div>
                })
            }
            </div>
        </div>
    )
}
