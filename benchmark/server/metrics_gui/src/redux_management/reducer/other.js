const Reducer_2 = (state = {
    last_screen_scroll:null
}, action) => {
    switch (action.type) {
        case "SET_OTHER":
            return action.payload
        case "UPDATE_OTHER":
            return {...state,...action.payload}
        case "RESET_OTHER":
            return {}
        default:
            return state
    }
}

export default Reducer_2
