package authentication

import (
	"github/fims-proto/fims-proto-ms/internal/common/log"

	"github.com/gin-gonic/gin"
)

func Authn() gin.HandlerFunc {
	return func(c *gin.Context) {
		// before request
		log.Debug(c, "check authentication against URL: %s", c.Request.RequestURI)

		c.Next()

		// after request
	}
}
