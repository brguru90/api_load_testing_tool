const MetricsReducer = (state = null, action) => {
    switch (action.type) {
        case "SET_METRICS":
            return action.payload
        case "RESET_METRICS":
            return {}
        default:
            return state
    }
}

export default MetricsReducer
