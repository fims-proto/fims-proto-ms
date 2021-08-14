package ginmiddleware

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	uuid4Localhost = uuid.New()
	uuid4Remote    = uuid.New()
)

func TestResolveTenantBySubdomain_localhost(t *testing.T) {
	t.Parallel()

	localhost, _ := url.Parse("http://localhost:3000/test")

	c := gin.Context{
		Request: &http.Request{
			URL: localhost,
		},
	}

	ResolveTenantBySubdomain(mockTenantService{})(&c)

	assert.Equal(t, uuid4Localhost, c.Value("tenant_id").(uuid.UUID))
}

func TestResolveTenantBySubdomain_remote(t *testing.T) {
	t.Parallel()

	remote, _ := url.Parse("https://some-domain.fims.com/test")

	c := gin.Context{
		Request: &http.Request{
			URL: remote,
		},
	}

	ResolveTenantBySubdomain(mockTenantService{})(&c)

	assert.Equal(t, uuid4Remote, c.Value("tenant_id").(uuid.UUID))
}

type mockTenantService struct{}

func (m mockTenantService) ReadTenantIdBySubdomain(ctx context.Context, subdomain string) (uuid.UUID, error) {
	if subdomain == "localhost" {
		return uuid4Localhost, nil
	}
	if subdomain == "some-domain" {
		return uuid4Remote, nil
	}
	return uuid.Nil, errors.New("not found")
}
