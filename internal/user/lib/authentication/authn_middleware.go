package authentication

import (
	"context"
	"fmt"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"net/http"

	"github.com/gin-gonic/gin"
	kratos "github.com/ory/kratos-client-go"
)

const kratosCookieKey = "ory_kratos_session"

type tenantManager interface {
	GetKratosClientBySubdomain(ctx context.Context, subdomain string) (*kratos.APIClient, error)
}

func Authn(tenantManager tenantManager) gin.HandlerFunc {
	if tenantManager == nil {
		panic("nil tenantManager")
	}
	return func(c *gin.Context) {
		// before request
		log.Debug(c, "check authentication against URL: %s", c.Request.RequestURI)

		apiClient, err := tenantManager.GetKratosClientBySubdomain(c, c.Value("subdomain").(string))
		if err != nil {
			log.Err(c, err, "failed to get kratos client by subdomain %s", c.Value("subdomain").(string))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "Internal Server Error"})
			return
		}

		sessionKey, err := c.Cookie(kratosCookieKey)
		if err != nil {
			log.Err(c, err, "No cookie named %s found", kratosCookieKey)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
			return
		}

		session, resp, err := apiClient.V0alpha1Api.ToSession(c).Cookie(fmt.Sprintf("%s=%s", kratosCookieKey, sessionKey)).Execute()
		if resp != nil && resp.StatusCode == http.StatusUnauthorized {
			log.Err(c, err, "Not Unauthorized")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"Message": "Unauthorized"})
			return
		} else if err != nil {
			log.Err(c, err, "error when validate session, full response: %+v", resp)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"Message": "Internal Server Error"})
			return
		}

		// log for now, will determine what to do next
		log.Debug(c, "session validated: %+v", session)

		c.Next()
		// after request
	}
}
