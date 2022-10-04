const proxy = require("http-proxy-middleware")
const { createProxyMiddleware } = require('http-proxy-middleware');

const http = require("http")
var keepAliveAgent = new http.Agent({ keepAlive: true })

const filter = function (pathname, req) {
    console.log("pathname",pathname)
    return pathname?.match('^/go_ws/');
};


module.exports = function (app) {

    app.use(
        createProxyMiddleware("/go_ws/",{
            target: "http://localhost:" + (process?.env?.SERVER_PORT || 7000),
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
        createProxyMiddleware("/api",{
            target: "http://localhost:" + (process?.env?.SERVER_PORT || 7000),
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
