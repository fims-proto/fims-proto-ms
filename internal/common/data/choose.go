package data

type Choose interface {
	Field() string
}

type chooseRequest struct {
	field string
}

func newChoose(field string) (Choose, error) {
	fieldSnakeCase, err := toSnakeCase(field)
	if err != nil {
		return nil, err
	}

	return chooseRequest{
		field: fieldSnakeCase,
	}, nil
}

func (c chooseRequest) Field() string {
	return c.field
}
