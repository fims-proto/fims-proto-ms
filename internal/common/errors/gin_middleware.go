package errors

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github/fims-proto/fims-proto-ms/internal/common/localization"
)

func ErrorHandler(localizer localization.Localizer) gin.HandlerFunc {
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
			if localize := localizer.Get(c.Request.Header.Get("Accept-Language"), slug, localizationArgs); localize != "" {
				message = localize
			}

			c.JSON(http.StatusBadRequest, slugErrResponse{
				Slug:    slug,
				Message: message,
			})
		}
	}
}
