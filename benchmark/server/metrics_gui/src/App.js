import React, { useEffect, useRef } from "react"
import "./App.scss"
import { BrowserRouter as Router, HashRouter, Switch, Route } from "react-router-dom"
import Dashboard from "./pages/dashboard.jsx"
import IterationsMetrics from "./pages/iterations_metrics.jsx"
import "antd/dist/antd.min.css"
import { useDispatch } from "react-redux"
import { GetBenchmarkMetrics } from "./services/metric.data"

export default function App() {

    const effectCalled = useRef(false)
    const dispatch = useDispatch()

    useEffect(() => {
        if (!effectCalled.current) {
            effectCalled.current = true
            GetBenchmarkMetrics((data) => {
                dispatch({
                    type: "SET_METRICS",
                    payload: [...data],
                })
            })                        
        }
    }, [])

    return (
        <div className="App">
            <Router>
                <Switch>
                    <HashRouter>
                        <Switch>
                            <Route path="/" exact component={Dashboard} />
                            <Route path="/iterations_metrics" component={IterationsMetrics} />
                        </Switch>
                    </HashRouter>
                </Switch>
            </Router>
        </div>
    )
}
