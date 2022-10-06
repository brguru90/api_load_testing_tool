import React, {Component} from "react"
import {Link} from "react-router-dom"

export default class Page2 extends Component {
    render() {
        return (
            <div>
                Page2
                <br />
                <Link to="/">view Page1</Link>
            </div>
        )
    }
}
