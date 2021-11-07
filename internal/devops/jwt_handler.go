package devops

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func InitJwtHandler(r *gin.RouterGroup) {
	r.GET("jwt", func(c *gin.Context) {
		token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
		c.String(http.StatusOK, token)
	})
}
