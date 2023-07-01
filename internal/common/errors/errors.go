package errors

func ErrNoWhereUsed(args ...string) SlugErr { return NewSlugError("no-where-used", args) }

func ErrRecordNotFound() SlugErr { return NewSlugError("record-not-found") }
