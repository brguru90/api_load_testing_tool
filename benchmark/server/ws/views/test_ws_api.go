package views

import (
	"github.com/brguru90/api_load_testing_tool/benchmark/my_modules"
	"github.com/brguru90/api_load_testing_tool/benchmark/store"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

func Metrics(c *gin.Context) {

	M := melody.New()
	M.HandleMessage(func(s *melody.Session, msg []byte) {
		fmt.Println("on Receiving message")
		switch string(msg) {
		case "hi":
			M.Broadcast([]byte("bi"))
		case "GM":
			M.Broadcast([]byte("GN"))
		default:
			M.Broadcast(msg)
		}
	})
	M.HandleConnect(func(s *melody.Session) {
		fmt.Println("on Connect")
		go func() {
			temp_data, info := store.BenchmarkDataStore_GetAllWithInfo()
			_temp_data, _info := *temp_data, *info

			result, err := json.MarshalIndent(_temp_data, "", "  ")
			if err == nil {
				M.Broadcast([]byte(result))
				if my_modules.BenchMarkEnded{
					store.BenchmarkDataStore_CloseQ()
					M.Close()
					return
				}
			}
			t2 := func(data interface{}) {
				if data==nil{
					store.BenchmarkDataStore_CloseQ()
					M.Close()
					return
				}

				_stream := data.(my_modules.BenchmarkMetricStreamInfo)
				if _stream.UpdatedAt > _info.UpdatedAt {
					result, err := json.MarshalIndent([]interface{}{_stream.Data}, "", "  ")
					if err == nil {
						M.Broadcast([]byte(result))
					}
				}
			}
			my_modules.BenchmarkMetricEvent.OnEvent(&t2)
		}()

	})
	M.HandleDisconnect(func(s *melody.Session) {
		fmt.Println("on Disconnect")
	})

	M.HandleRequest(c.Writer, c.Request)

}
