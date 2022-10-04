package views

import (
	"apis_load_test/my_modules"
	"apis_load_test/server/ws/ws_modules"
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
			result, err := json.MarshalIndent(my_modules.BenchmarkMetricArray, "", "  ")
			if err == nil {
				w.M.Broadcast([]byte(result))
			}
			for data := range my_modules.BenchmarkMetricStream {
				result, err := json.MarshalIndent([]map[string]interface{}{data}, "", "  ")
				if err == nil {
					w.M.Broadcast([]byte(result))
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
