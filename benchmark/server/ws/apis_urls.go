package ws

import (
	"github.com/brguru90/api_load_testing_tool/benchmark/server/ws/views"

	"github.com/gin-gonic/gin"
)

func InitWS(router *gin.RouterGroup) {
	router.GET("metrics/", views.Metrics)
}
