package errors

import (
	"errors"
	"net/http"
)

const unknownErrorSlug = "unknown-error"

type slugErrResponse struct {
	Message string `json:"message,omitempty"`
	Slug    string `json:"slug,omitempty"`
}

// ErrorType classifies a SlugErr so the HTTP middleware can return the
// semantically correct status code instead of always returning 400.
type ErrorType int

const (
	ErrorTypeInvalidInput ErrorType = iota // HTTP 400 — business rule / validation failure
	ErrorTypeNotFound                      // HTTP 404 — resource does not exist
	ErrorTypeConflict                      // HTTP 409 — uniqueness / state conflict
	ErrorTypeInternal                      // HTTP 500 — unexpected infrastructure error
)

// HTTPStatus maps each ErrorType to its canonical HTTP status code.
func (t ErrorType) HTTPStatus() int {
	switch t {
	case ErrorTypeNotFound:
		return http.StatusNotFound // 404
	case ErrorTypeConflict:
		return http.StatusConflict // 409
	case ErrorTypeInternal:
		return http.StatusInternalServerError // 500
	default:
		return http.StatusBadRequest // 400
	}
}

type SlugErr struct {
	slug      string
	args      []any
	errorType ErrorType
}

func (s SlugErr) Error() string {
	return s.slug
}

func (s SlugErr) Is(target error) bool {
	var t SlugErr
	ok := errors.As(target, &t)
	if !ok {
		return false
	}
	return s.slug == t.Error()
}

// NewInvalidInputError creates a 500 internal error slug error.
func NewInvalidInputError(slug string, args ...any) SlugErr {
	return SlugErr{slug: slug, args: args, errorType: ErrorTypeInvalidInput}
}

// NewNotFoundError creates a 404 Not Found slug error.
func NewNotFoundError(slug string, args ...any) SlugErr {
	return SlugErr{slug: slug, args: args, errorType: ErrorTypeNotFound}
}

// NewConflictError creates a 409 Conflict slug error (e.g. unique-constraint violations).
func NewConflictError(slug string, args ...any) SlugErr {
	return SlugErr{slug: slug, args: args, errorType: ErrorTypeConflict}
}

// NewInternalError creates a 500 internal error slug error.
func NewInternalError(slug string, args ...any) SlugErr {
	return SlugErr{slug: slug, args: args, errorType: ErrorTypeInternal}
}
