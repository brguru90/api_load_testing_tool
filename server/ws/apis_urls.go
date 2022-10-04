package ws

import (
	"apis_load_test/server/ws/views"
	"apis_load_test/server/ws/ws_modules"

	"github.com/gin-gonic/gin"
)

func InitWS(router *gin.RouterGroup) {
	router.GET(ws_modules.GetWsHandlers("metrics/",views.Metrics))
}
