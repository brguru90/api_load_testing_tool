import React, { useEffect, useRef } from 'react'
import { useSelector } from 'react-redux'
import { useHistory } from 'react-router-dom';
import APITimes from './apis_time'
import styles from "./style.module.scss"

export default function APIPerSecondMetrics({ APIindex }) {
    const history = useHistory();
    const [url, iteration_ids] = useSelector(state => [state.metrics_data?.[APIindex]?.url, state.metrics_data?.[APIindex]?.iteration_data?.map(item => item?.iteration_id)])

    const navigateToDashBoard = (e) => {
        e.preventDefault()
        history.push({
            pathname: '/',
            state: { restoreScroll: true }
        })
    }


    const effectCalled = useRef(false)
    useEffect(() => {
        if (!effectCalled.current) {
            effectCalled.current = true
            history.listen(location => {
                if (history.action === 'POP' && location.pathname=="/") {
                    console.log(window.location.pathname,location.pathname)
                    history.replace({
                        pathname: '/',
                        state: { restoreScroll: true }
                    });
                }
            })
        }
    })

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
