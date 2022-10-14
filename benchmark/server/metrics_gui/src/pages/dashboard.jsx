import React, { useEffect, useRef, useState } from "react"
import Dashboard from "../components/dashboard/index.jsx"
import { useSelector } from "react-redux"
import { useLocation } from 'react-router-dom';

export default function dashboard_page() {
    const last_scroll_position = useSelector(state => state.other?.last_screen_scroll)
    const {state} = useLocation();

    useEffect(() => {
        if (state?.restoreScroll && last_scroll_position != null) {
            window.scroll({
                top: last_scroll_position,
                left: 0,
                behavior: 'smooth'
            });
        }
    }, [])

    return <Dashboard />
}
