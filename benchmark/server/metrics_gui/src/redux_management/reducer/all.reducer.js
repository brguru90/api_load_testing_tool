import metric_reducer from "./metrics.reducer"
import reducer_2 from "./reducer2"
import {combineReducers} from "redux"

const allReducer = combineReducers({
    metrics_data: metric_reducer,
    data_2: reducer_2,
})
export default allReducer
