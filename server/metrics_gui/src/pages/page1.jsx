import React, {useEffect, useRef} from "react"
import {Link} from "react-router-dom"

export default function Page1() {

    const effectCalled = useRef(false);

  


    let client=null
    const connect = async () => {        
        await fetch("/go_ws/")

        let loc = window.location,
            ws_url
        if (loc.protocol === "https:") {
            ws_url = "wss:"
        } else {
            ws_url = "ws:"
        }
        ws_url += "//" + loc.host + "/go_ws"
        client = new WebSocket(`${ws_url}/metrics/`)

        client.onopen = function () {
            console.log("[open] Connection established")
        }

        client.onmessage = function (event) {
            console.log(`[message] Data received from server: ${event.data}`)
        }

        client.onclose = function (event) {
            if (event.wasClean) {
                console.log(
                    `[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`
                )
            } else {
                console.log("[close] Connection died")
            }
        }

        client.onerror = function (error) {
            console.log(error)
        }

    }


    const send = (text) => {
        this.client.send(text)
    }

    const disconnect = () => {
        this.client.close()
    }

    useEffect(() => {

        if (!effectCalled.current) {
            effectCalled.current = true;
            connect()
        }       
        
    }, [])


    return (
        <div>
            Page1 <br />
            <Link to="page2">view Page2</Link> <br />
        </div>
    )
}
