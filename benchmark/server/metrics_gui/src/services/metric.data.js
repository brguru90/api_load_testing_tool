import {ManageWebSocket} from "./util"

const GetBenchmarkMetrics = (callback) => {
    if (!callback) return
    const manageWebSocket = new ManageWebSocket("/metrics/")
    let existing_data = []
    manageWebSocket.connect((raw_data) => {
        raw_data = JSON.parse(raw_data)
        // console.log("raw_data", raw_data)

        if (raw_data && Array.isArray(raw_data)) {
            let new_benchmark = true
            raw_data.forEach((raw_elem, raw_index) => {
                existing_data.forEach((elem, index) => {
                    if (
                        raw_elem.url == elem.url &&
                        raw_elem.process_uid == elem.process_uid
                    ) {
                        if (raw_elem.iteration_data) {
                            elem.iteration_data=[...elem.iteration_data,...raw_elem.iteration_data]
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

        console.log("existing_data", JSON.parse(JSON.stringify(existing_data)))
        callback(existing_data)
    })
}

export {GetBenchmarkMetrics}
