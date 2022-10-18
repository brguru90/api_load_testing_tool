const MetricsExtraReducer = (state = {
    last_screen_scroll:null
}, action) => {
    switch (action.type) {
        case "SET_METRICS_EXTRA":
            return action.payload
        case "UPDATE_METRICS_EXTRA":
            return {...state,...action.payload}
        case "RESET_METRICS_EXTRA":
            return {}
        default:
            return state
    }
}

export default MetricsExtraReducer