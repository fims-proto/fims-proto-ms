package ginmiddleware

import (
	"context"
	"fmt"
	"strings"

	"github/fims-proto/fims-proto-ms/internal/common/database"

	"github/fims-proto/fims-proto-ms/internal/common/log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type tenantManager interface {
	GetDBConnBySubdomain(ctx context.Context, subdomain string) (*gorm.DB, error)
}

// ResolveTenantBySubdomain subdomain:
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

		log.Debug(c, "resolved subdomain: %s", subdomain)

		db, err := tenantManager.GetDBConnBySubdomain(c, subdomain)
		if err != nil {
			panic(fmt.Errorf("failed to get DB connection by subdomanin %s: %w", subdomain, err))
		}

		// set DB properties to request context
		db = db.WithContext(c.Request.Context())
		ctx := database.NewContextWithDB(c.Request.Context(), db)
		c.Request = c.Request.WithContext(ctx)

		c.Set("subdomain", subdomain)

		c.Next()
	}
}
