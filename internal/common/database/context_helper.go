package database

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ctxDBKey struct{}

func ReadDBFromContext(ctx context.Context) *gorm.DB {
	// if it's gin.Context passing through, get it from request context
	if ginCtx, ok := ctx.(*gin.Context); ok {
		return ginCtx.Request.Context().Value(ctxDBKey{}).(*gorm.DB)
	}
	return ctx.Value(ctxDBKey{}).(*gorm.DB)
}

func NewContextWithDB(parent context.Context, db *gorm.DB) context.Context {
	return context.WithValue(parent, ctxDBKey{}, db)
}
