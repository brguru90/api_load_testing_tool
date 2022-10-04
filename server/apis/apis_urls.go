package apis

import (
	"apis_load_test/server/apis/views"

	"github.com/gin-gonic/gin"
)

// only the functions whose initial letter is upper case only those can be exportable from package
func InitApis(router *gin.RouterGroup) {
	router.GET("hello/", views.Hello_api)
}
