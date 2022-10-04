package views

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Hello_api(c *gin.Context) {
	c.String(http.StatusOK, "Welcome hello")
}
