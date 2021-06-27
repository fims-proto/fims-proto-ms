package domain

import (
	"errors"
)

type Sob struct {
	id          string
	name        string
	description string
}

func NewSob(id, name, description string) (*Sob, error) {
	if id == "" {
		return nil, errors.New("empty sob id")
	}
	if name == "" {
		return nil, errors.New("empty sob name")
	}

	return &Sob{
		id:          id,
		name:        name,
		description: description,
	}, nil
}

func (s Sob) Id() string {
	return s.id
}

func (s Sob) Name() string {
	return s.name
}

func (s Sob) Description() string {
	return s.description
}
