package ginmiddleware

import (
	"context"
	"github/fims-proto/fims-proto-ms/internal/common/log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type tenantManager interface {
	GetDBConnBySubdomain(ctx context.Context, subdomain string) (*gorm.DB, error)
}

// subdomain:
// https://random-domain.xxxhost.com -> random-domain
// http://localhost -> localhost
// http://127.0.0.1 -> localhost
func ResolveTenantBySubdomain(tenantManager tenantManager) gin.HandlerFunc {
	if tenantManager == nil {
		panic("nil tenant manager")
	}
	return func(c *gin.Context) {
		host := strings.Split(c.Request.Host, ":")[0]
		subdomain := strings.Split(host, ".")[0]
		if host == "127.0.0.1" {
			subdomain = "localhost"
		}

		log.Debug(c, "resolved subdoamin: %s", subdomain)

		db, err := tenantManager.GetDBConnBySubdomain(c, subdomain)
		if err != nil {
			panic(errors.Wrapf(err, "failed to get DB connection by subdomanin %s", subdomain))
		}
		c.Set("db", db.WithContext(c))
		c.Set("subdomain", subdomain)

		c.Next()
	}
}
