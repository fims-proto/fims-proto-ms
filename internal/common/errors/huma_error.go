package errors

import (
	"encoding/json"
	"errors"
	"net/http"

	"github/fims-proto/fims-proto-ms/internal/common/localization"

	"github.com/danielgtaylor/huma/v2"
)

// InitHumaErrorHandler overrides huma.NewErrorWithContext to translate SlugErr
// into the project's {slug, message} JSON format with i18n support.
func InitHumaErrorHandler(localizer localization.Localizer) {
	huma.NewErrorWithContext = func(ctx huma.Context, status int, msg string, errs ...error) huma.StatusError {
		for _, err := range errs {
			var se SlugErr
			if asSlugErr(err, &se) {
				lang := ctx.Header("Accept-Language")
				localized := localizer.Get(lang, se.slug, se.args)
				if localized == "" {
					localized = se.slug
				}
				return &humaSlugError{
					status:  se.errorType.HTTPStatus(),
					slug:    se.slug,
					message: localized,
				}
			}
		}
		return huma.NewError(status, msg, errs...)
	}
}

// asSlugErr checks if err is or wraps a SlugErr and assigns it.
func asSlugErr(err error, target *SlugErr) bool {
	if err == nil {
		return false
	}
	if se, ok := errors.AsType[SlugErr](err); ok {
		*target = se
		return true
	}
	return false
}

type humaSlugError struct {
	status  int
	slug    string
	message string
}

func (e *humaSlugError) Error() string           { return e.message }
func (e *humaSlugError) GetStatus() int          { return e.status }
func (e *humaSlugError) GetHeaders() http.Header { return nil }

func (e *humaSlugError) MarshalJSON() ([]byte, error) {
	return json.Marshal(slugErrResponse{Slug: e.slug, Message: e.message})
}
