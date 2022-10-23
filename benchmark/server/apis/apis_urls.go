package apis

import (
	"github.com/brguru90/api_load_testing_tool/benchmark/server/apis/views"

	"github.com/gin-gonic/gin"
)

// only the functions whose initial letter is upper case only those can be exportable from package
func InitApis(router *gin.RouterGroup) {
	router.GET("hello/", views.Hello_api)
}
