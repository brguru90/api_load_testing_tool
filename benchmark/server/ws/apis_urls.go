package ws

import (
	"apis_load_test/benchmark/server/ws/views"

	"github.com/gin-gonic/gin"
)

func InitWS(router *gin.RouterGroup) {
	router.GET("metrics/", views.Metrics)
}