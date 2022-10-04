const proxy = require("http-proxy-middleware")
const http = require("http")
var keepAliveAgent = new http.Agent({keepAlive: true})
module.exports = function (app) {
    app.use(
        "/go_ws",
        proxy({
            target: "http://localhost:"+(process?.env?.SERVER_PORT || 7000),
            changeOrigin: true,
            agent: keepAliveAgent,
            ws: true,
            logLevel: "debug",
            // pathRewrite: {
            //     "^/ws": "/",
            // },
        })
    )

    app.use(
        "/api",
        proxy({
            target: "http://localhost:"+(process?.env?.SERVER_PORT || 7000),
            changeOrigin: true,
            agent: keepAliveAgent,
            // ws: true,
            xfwd: true,
            // pathRewrite: {
            //     "^/api": "/",
            // },
            logLevel: "debug",
        })
    )
}
