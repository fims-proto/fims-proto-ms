package datasource

import (
	"strings"

	"github/fims-proto/fims-proto-ms/internal/common/log"

	"github.com/gin-gonic/gin"
)

// ResolveSubdomain parse subdomain from URL and set into context
// https://random-domain.xxxhost.com -> random-domain
// http://localhost or http://127.0.0.1 -> localhost
func ResolveSubdomain() gin.HandlerFunc {
	return func(c *gin.Context) {
		host := strings.Split(c.Request.Host, ":")[0]
		subdomain := strings.Split(host, ".")[0]
		if host == "127.0.0.1" {
			subdomain = "localhost"
		}
		log.Debug(c, "resolved subdomain: %s", subdomain)
		c.Set("subdomain", subdomain)

		c.Next()
	}
}
