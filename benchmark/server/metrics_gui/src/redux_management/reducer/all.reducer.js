import metric_reducer from "./metrics.reducer"
import other from "./other.reducer"
import metrics_extra from "./metrics.extra.reducer"
import {combineReducers} from "redux"
import { persistReducer } from 'redux-persist'
import storage from 'redux-persist/lib/storage'


const rootPersistConfig = {
    key: 'root',
    storage: storage,
    blacklist: ['metrics_data','metrics_extra']
  }

const allReducer = combineReducers({
    metrics_data: metric_reducer,
    other,
    metrics_extra,
})
export default persistReducer(rootPersistConfig,allReducer)
