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
func ResolveTenantBySubdomain(tenantManager tenantManager) gin.HandlerFunc {
	if tenantManager == nil {
		panic("nil tenant manager")
	}
	return func(c *gin.Context) {
		hostParts := strings.Split(strings.Split(c.Request.Host, ":")[0], ".")

		log.Debug(c, "resolved subdoamin: %s", hostParts[0])

		db, err := tenantManager.GetDBConnBySubdomain(c, hostParts[0])
		if err != nil {
			panic(errors.Wrapf(err, "failed to get DB connection by subdomanin %s", hostParts[0]))
		}
		c.Set("db", db.WithContext(c))

		c.Next()
	}
}
