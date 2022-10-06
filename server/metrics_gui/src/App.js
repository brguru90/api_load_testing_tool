import React, {Component} from "react"
import "./App.scss"
import {BrowserRouter as Router, HashRouter, Switch, Route} from "react-router-dom"
import Dashboard from "./pages/dashboard.jsx"
import page2 from "./pages/page2.jsx"
import "antd/dist/antd.min.css"

export default class App extends Component {
    render() {
        return (
            <div className="App">
                <Router>
                    <Switch>
                        <HashRouter>
                            <Switch>
                                <Route path="/" exact component={Dashboard} />
                                <Route path="/page2" exact component={page2} />
                            </Switch>
                        </HashRouter>
                    </Switch>
                </Router>
            </div>
        )
    }
}
