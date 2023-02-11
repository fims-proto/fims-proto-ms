package errors

const unknownErrorSlug = "unknown-error"

type slugErrResponse struct {
	Message string `json:"message,omitempty"`
	Slug    string `json:"slug,omitempty"`
}

type SlugErr struct {
	slug string
	args []any
}

func (s SlugErr) Error() string {
	return s.slug
}

func NewSlugError(slug string, args ...any) SlugErr {
	return SlugErr{
		slug: slug,
		args: args,
	}
}
