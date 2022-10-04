package views

import (
	"apis_load_test/my_modules"
	"apis_load_test/server/ws/ws_modules"
	"apis_load_test/store"
	"encoding/json"
	"fmt"

	"gopkg.in/olahol/melody.v1"
)

var MetricsStream chan interface{}
var AccumulatedMetrics map[string]interface{}

func Metrics(w *ws_modules.WsHandlers) {

	w.OnConnect = func(s *melody.Session) {
		fmt.Println("on Connect")
		w.M.Broadcast([]byte("connection success : initial message from server"))
		go func() {
			temp_data, info := store.GeneralStore_GetAllWithInfo()
			_temp_data, _info := *temp_data, *info

			result, err := json.MarshalIndent(_temp_data, "", "  ")
			if err == nil {
				w.M.Broadcast([]byte(result))
			}
			for _stream := range my_modules.BenchmarkMetricStream {
				if _stream.UpdatedAt > _info.UpdatedAt {
					result, err := json.MarshalIndent([]map[string]interface{}{_stream.Data}, "", "  ")
					if err == nil {
						w.M.Broadcast([]byte(result))
					}
				}
			}
		}()

	}

	w.OnMessage = func(s *melody.Session, msg []byte) {
		fmt.Println("on Receiving message")
		switch string(msg) {
		case "hi":
			w.M.Broadcast([]byte("bi"))
		case "GM":
			w.M.Broadcast([]byte("GN"))
		default:
			w.M.Broadcast(msg)
		}
	}

	w.OnDisconnect = func(s *melody.Session) {
		fmt.Println("on Disconnect")
	}

}
