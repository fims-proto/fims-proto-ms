package errors

func ErrRecordNotFound() SlugErr {
	return NewNotFoundError(SlugRecordNotFound) // 404
}
