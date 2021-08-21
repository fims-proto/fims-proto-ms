package ginmiddleware

import (
	"context"
	"net/http"
	"net/url"
	"testing"

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

	c := gin.Context{
		Request: &http.Request{
			URL: localhost,
		},
	}

	ResolveTenantBySubdomain(mockTenantManager{})(&c)

	assert.Equal(t, "localhost", c.Value("db").(*gorm.DB).Error.Error())
}

func TestResolveTenantBySubdomain_remote(t *testing.T) {
	t.Parallel()

	remote, _ := url.Parse("https://some-domain.fims.com/test")

	c := gin.Context{
		Request: &http.Request{
			URL: remote,
		},
	}

	ResolveTenantBySubdomain(mockTenantManager{})(&c)

	assert.Equal(t, "remote", c.Value("db").(*gorm.DB).Error.Error())
}

type mockTenantManager struct{}

func (m mockTenantManager) GetDBConnBySubdomain(ctx context.Context, subdomain string) (*gorm.DB, error) {
	if subdomain == "localhost" {
		return db4Localhost, nil
	}
	if subdomain == "some-domain" {
		return db4Remote, nil
	}
	return nil, errors.New("not found")
}
