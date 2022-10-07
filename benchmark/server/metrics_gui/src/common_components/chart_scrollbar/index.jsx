import React from 'react'
import styles from "./style.module.scss"

export default function ChartScrollbar({className="", scroll_count, onScroll = e => console.log(e) }) {
    return (
        <div
            className={styles["chart_scrollbar_parent"]+" "+className}
            style={{ "--scroll-count": scroll_count }}
            onScroll={e => {
                let a=e?.target?.scrollLeft
                let b = e?.target?.scrollWidth - e?.target?.offsetWidth;
                onScroll( scroll_count*(a/b),e,scroll_count)
            }}
        >
            <div className={styles["chart_scrollbar"]}>
            </div>
        </div>
    )
}
