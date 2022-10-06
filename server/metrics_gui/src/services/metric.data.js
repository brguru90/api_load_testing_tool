import {ManageWebSocket} from "./util"

const GetBenchmarkMetrics = (callback) => {
    if (!callback) return
    const manageWebSocket = new ManageWebSocket("/metrics/")
    let existing_data = []
    manageWebSocket.connect((raw_data) => {
        // console.log("raw_data", raw_data,typeof(raw_data))
        // need to refine data in different categories
        // need to set ata in redux, so it can be accessed directly in chart component
        // since data will e updated on redux, no need of callback

        // callback(data)

        raw_data = JSON.parse(raw_data)
        console.log("raw_data", raw_data)

        if (raw_data && Array.isArray(raw_data)) {
            let new_benchmark = true
            raw_data.forEach((raw_elem, raw_index) => {
                existing_data.forEach((elem, index) => {
                    if (
                        raw_elem.url == elem.url &&
                        raw_elem.process_uid == elem.process_uid
                    ) {
                        if (raw_elem.iteration_data) {
                            elem.iteration_data.push(raw_elem.iteration_data)
                        }
                        if (raw_elem.all_data?.Url) {
                            elem.all_data = raw_elem.all_data
                        }
                        new_benchmark = false
                    }
                })
            })
            if (new_benchmark) {
                existing_data = [...existing_data, ...raw_data]
            }
        }

        console.log("existing_data", existing_data)
    })
}

export {GetBenchmarkMetrics}
