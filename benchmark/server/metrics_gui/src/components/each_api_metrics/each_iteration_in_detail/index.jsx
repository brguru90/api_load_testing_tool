import React, { useEffect, useRef } from 'react'
import { useSelector } from 'react-redux'
import { useHistory } from 'react-router-dom';
import APITimes from './apis_time'
import styles from "./style.module.scss"

export default function APIPerSecondMetrics({ APIindex }) {
    const history = useHistory();
    const [url, iterations] = useSelector(state => [state.metrics_data?.[APIindex]?.url, state.metrics_data?.[APIindex]?.iteration_data || []])

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
            window.scroll({
                top: 0,
                left: 0
            });
            history.listen(location => {
                if (history.action === 'POP' && location.pathname == "/") {
                    console.log(window.location.pathname, location.pathname)
                    history.replace({
                        pathname: '/',
                        state: { restoreScroll: true }
                    });
                }
            })
        }
    })

    useEffect(() => {
        console.log(`Rendered: APIPerSecondMetrics index=${APIindex}`)
    })

    return (
        <div className={styles['per_second_metrics_dashboard']}>
            <div className={styles['per_second_metrics']}>
                <h1 className={styles['title']}>                    
                    <a href="" onClick={navigateToDashBoard}>
                        &lt;&lt; Back &nbsp;&nbsp;
                    </a>
                    API Per Second Metrics - {url}</h1>
                <div className={styles['per_second_metrics_set']}>
                    {
                        iterations?.map((iter) => {
                            return <div key={iter.iteration_id} className={styles['per_second_data']}>
                                <APITimes iteration={iter} />
                            </div>
                        })
                    }
                </div>
            </div>

        </div>
    )
}
