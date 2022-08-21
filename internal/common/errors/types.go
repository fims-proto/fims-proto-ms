package errors

const unknownErrorSlug = "unknown-error"

type slugErrResponse struct {
	Message string `json:"message,omitempty"`
	Slug    string `json:"slug,omitempty"`
}

type SlugErr struct {
	slug  string
	error string
	args  []any
}

func (s SlugErr) Error() string {
	return s.error
}

func NewSlugError(slug, error string, args ...any) SlugErr {
	return SlugErr{
		slug:  slug,
		error: error,
		args:  args,
	}
}
