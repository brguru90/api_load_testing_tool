import React, { useState } from 'react'
import styles from "./style.module.scss"

export default function ChartScrollbar({ className = "", scroll_count, onScroll = e => console.log(e) }) {
    const [currentRange, setCurrentRange] = useState(0)
    const scrollFractions = 100
    return (
        // <div
        //     className={styles["chart_scrollbar_parent"]+" "+className}
        //     style={{ "--scroll-count": scroll_count }}
        //     onScroll={e => {
        //         let a=e?.target?.scrollLeft
        //         let b = e?.target?.scrollWidth - e?.target?.offsetWidth;
        //         onScroll( scroll_count*(a/b),e,scroll_count)
        //     }}
        // >
        //     <div className={styles["chart_scrollbar"]}>
        //     </div>
        // </div>
        <div className={styles["chart_scrollbar_parent2"]+" "+className}>
            <input type="range" min="0" max={scroll_count * scrollFractions} value={currentRange} step="1" className={styles["chart_scrollbar2"]} onChange={e => {
                console.log("e", e)
                setCurrentRange(e.target.value)
                onScroll(e.target.value / scrollFractions)
            }} /><span className={styles["scroll_count"]}>{Math.round(currentRange/ scrollFractions)}/{scroll_count}</span>
        </div>
    )
}
