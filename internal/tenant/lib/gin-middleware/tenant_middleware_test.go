package ginmiddleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github/fims-proto/fims-proto-ms/internal/common/database"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var (
	db4Localhost = &gorm.DB{
		Error:     errors.New("localhost"),
		Config:    &gorm.Config{},
		Statement: &gorm.Statement{},
	}
	db4Remote = &gorm.DB{
		Error:     errors.New("remote"),
		Config:    &gorm.Config{},
		Statement: &gorm.Statement{},
	}
)

func TestResolveTenantBySubdomain_localhost(t *testing.T) {
	t.Parallel()

	localhost, _ := url.Parse("http://localhost:3000/test")

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = &http.Request{
		Host: "localhost:3000",
		URL:  localhost,
	}

	ResolveTenantBySubdomain(mockTenantManager{})(c)

	assert.Equal(t, "localhost", database.ReadDBFromContext(c).Error.Error())
}

func TestResolveTenantBySubdomain_127_0_0_1(t *testing.T) {
	t.Parallel()

	localhost, _ := url.Parse("http://127.0.0.1:3000/test")

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = &http.Request{
		Host: "127.0.0.1:3000",
		URL:  localhost,
	}

	ResolveTenantBySubdomain(mockTenantManager{})(c)

	assert.Equal(t, "localhost", database.ReadDBFromContext(c).Error.Error())
}

func TestResolveTenantBySubdomain_remote(t *testing.T) {
	t.Parallel()

	remote, _ := url.Parse("https://some-domain.fims.com/test")

	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = &http.Request{
		Host: "some-domain.fims.com",
		URL:  remote,
	}

	ResolveTenantBySubdomain(mockTenantManager{})(c)

	assert.Equal(t, "remote", database.ReadDBFromContext(c.Request.Context()).Error.Error())
}

type mockTenantManager struct{}

func (m mockTenantManager) GetDBConnBySubdomain(_ context.Context, subdomain string) (*gorm.DB, error) {
	if subdomain == "localhost" {
		return db4Localhost, nil
	}
	if subdomain == "some-domain" {
		return db4Remote, nil
	}
	return nil, errors.New("not found")
}
