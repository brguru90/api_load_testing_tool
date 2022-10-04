package ws_modules

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/olahol/melody.v1"
)

// when every time creating new melody/websocker object,  requires a new unique variable
// & its difficult to group it
// so using bellow function we can easily map handler to the api in apis_urls



type WsHandlers struct {
	M *melody.Melody
	OnConnect func(s *melody.Session)
	OnMessage func(s *melody.Session, msg []byte)
	OnDisconnect func(s *melody.Session)
}

var AllHandler map[string]WsHandlers = make(map[string]WsHandlers)

func GetWsHandlers(url string, handler_callback func(*WsHandlers)) (string,func(c *gin.Context)) {
	if AllHandler[url].M == nil {
		var WS WsHandlers
		WS.M = melody.New()

		handler_callback(&WS)
		
		WS.M.HandleMessage(func(s *melody.Session, msg []byte) {
			if WS.OnMessage!=nil{
				WS.OnMessage(s,msg)
			}
		})
		WS.M.HandleConnect(func(s *melody.Session) {
			if WS.OnConnect!=nil{
				WS.OnConnect(s)
			}		
		})
		WS.M.HandleDisconnect(func(s *melody.Session) {
			if WS.OnDisconnect!=nil{
				WS.OnDisconnect(s)
			}			
		})
		AllHandler[url] = WS
	}
	return url,func(c *gin.Context) {
		AllHandler[url].M.HandleRequest(c.Writer, c.Request)
	}
}
