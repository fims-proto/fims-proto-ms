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

func (s SlugErr) Is(target error) bool {
	t, ok := target.(SlugErr)
	if !ok {
		return false
	}
	return s.slug == t.Error()
}

func NewSlugError(slug string, args ...any) SlugErr {
	return SlugErr{
		slug: slug,
		args: args,
	}
}
