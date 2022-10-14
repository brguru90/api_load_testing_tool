import metric_reducer from "./metrics.reducer"
import other from "./other"
import {combineReducers} from "redux"
import { persistReducer } from 'redux-persist'
import storage from 'redux-persist/lib/storage'


const rootPersistConfig = {
    key: 'root',
    storage: storage,
    blacklist: ['metrics_data']
  }

const allReducer = combineReducers({
    metrics_data: metric_reducer,
    other: other,
})
export default persistReducer(rootPersistConfig,allReducer)
