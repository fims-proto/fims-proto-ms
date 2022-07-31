package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/pkg/errors"
)

func ErrorHandler(bundle *i18n.Bundle) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		message := c.Errors.String() // default error message

		var slug string
		var localizationArgs []any
		if len(c.Errors) == 1 {
			// there should only be 1 error in the stack
			ginErr := c.Errors.Last()

			se, ok := errors.Cause(ginErr.Err).(SlugErr)
			if ok {
				slug = se.slug
				localizationArgs = se.args
			} else {
				slug = unknownErrorSlug
			}
		} else if len(c.Errors) > 1 {
			// if multiple errors in gin, then it's unknown to us
			slug = unknownErrorSlug
		}

		if slug != "" {
			localize, _ := i18n.NewLocalizer(bundle, c.Request.Header.Get("Accept-Language")).
				Localize(&i18n.LocalizeConfig{
					MessageID:    slug,
					TemplateData: localizationArgs,
				})
			if localize != "" {
				message = localize
			}

			c.JSON(http.StatusInternalServerError, slugErrResponse{
				Slug:    slug,
				Message: message,
			})
		}
	}
}
