package exception

const unknownErrorSlug = "unknown-error"

type slugErr interface {
	Slug() string
	Args() []any
}

type slugErrResponse struct {
	Message string `json:"message"`
	Slug    string `json:"slug"`
}
