package errors

type ObjectNotFoundErr struct {
	slugErr SlugErr
}

func (o ObjectNotFoundErr) Error() string {
	return o.slugErr.slug
}

func NewObjectNotFoundErr(object string) ObjectNotFoundErr {
	return ObjectNotFoundErr{
		slugErr: SlugErr{
			slug: object + "-notFound",
			args: nil,
		},
	}
}
