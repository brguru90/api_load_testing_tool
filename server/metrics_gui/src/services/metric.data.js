import {ManageWebSocket} from "./util"

const GetBenchmarkMetrics = (callback) => {
    if (!callback) return
    const manageWebSocket = new ManageWebSocket("/metrics/")
    manageWebSocket.connect((raw_data) => {
        console.log("raw_data", raw_data)
        // need to refine data in different categories
        // need to set ata in redux, so it can be accessed directly in chart component
        // since data will e updated on redux, no need of callback

        // callback(data)
    })
}

export {GetBenchmarkMetrics}
