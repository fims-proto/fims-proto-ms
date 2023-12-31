package errors

func ErrNoWhereUsed(args ...string) SlugErr {
	return NewSlugError("no-where-used", args)
}

// common

func ErrRecordNotFound() SlugErr {
	return NewSlugError("record-not-found")
}

// account

func ErrInvalidAccountClass(class string) SlugErr {
	return NewSlugError("invalid-account-class", class)
}

func ErrInvalidAccountGroup(class string, group string) SlugErr {
	return NewSlugError("invalid-account-group", class, group)
}

func ErrInvalidAccountNumber(number string) SlugErr {
	return NewSlugError("invalid-account-number", number)
}

// auxiliary account

func ErrInvalidAuxiliaryAccountKey(categoryKey, accountKey string) SlugErr {
	return NewSlugError("invalid-auxiliary-account-key", categoryKey, accountKey)
}

// period

func ErrPeriodClosed() SlugErr {
	return NewSlugError("period-closed")
}
