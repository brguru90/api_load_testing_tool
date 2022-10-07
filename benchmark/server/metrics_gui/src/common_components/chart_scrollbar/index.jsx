import React from 'react'
import styles from "./style.module.scss"

export default function ChartScrollbar({scroll_count,onScroll=e=>console.log(e)}) {
    return (
        <div className={styles["chart_scrollbar_parent"]} style={{"--scroll-count":scroll_count}} onScroll={e=>onScroll(e,scroll_count)}>
            <div className={styles["chart_scrollbar"]}>
            </div>
        </div>
    )
}
