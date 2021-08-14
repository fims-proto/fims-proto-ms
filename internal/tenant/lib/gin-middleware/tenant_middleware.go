package ginmiddleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
)

type tenantService interface {
	ReadTenantIdBySubdomain(ctx context.Context, subdomain string) (uuid.UUID, error)
}

func ResolveTenantBySubdomain(tenantService tenantService) gin.HandlerFunc {
	if tenantService == nil {
		panic("nil tenant service")
	}
	return func(c *gin.Context) {
		hostParts := strings.Split(strings.Split(c.Request.URL.Host, ":")[0], ".")
		tenantId, err := tenantService.ReadTenantIdBySubdomain(c, hostParts[0])
		if err != nil {
			panic(errors.Wrapf(err, "failed to read tenant id by subdomanin %s", hostParts[0]))
		}
		c.Set("tenant_id", tenantId)

		c.Next()
	}
}
