class ManageWebSocket {
    client = null
    path = null

    constructor(path) {
        if (path && !path?.startsWith("/")) {
            path = "/" + path
        }
        this.path = path
    }

    connect = (
        callback = (...e) => {
            console.log("unhandled onmessage", e)
        }
    ) => {
        let loc = window.location,
            ws_url
        if (loc.protocol === "https:") {
            ws_url = "wss:"
        } else {
            ws_url = "ws:"
        }
        ws_url += "//" + loc.host + "/go_ws"
        this.client = new WebSocket(`${ws_url}${this.path}`)

        this.client.onopen = function () {
            console.log("[open] Connection established")
        }

        this.client.onmessage = function (event) {
            console.log(`[message] Data received from server: ${event.data}`)
            callback()
        }

        this.client.onclose = function (event) {
            if (event.wasClean) {
                console.log(
                    `[close] Connection closed cleanly, code=${event.code} reason=${event.reason}`
                )
            } else {
                console.log("[close] Connection died")
            }
        }

        this.client.onerror = function (error) {
            console.log(error)
        }
    }

    send = (text) => {
        this.client.send(text)
    }

    disconnect = () => {
        this.client.close()
        this.client = null
    }
}

export {ManageWebSocket}
