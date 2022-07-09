package data

import (
	"math"

	"github.com/pkg/errors"
)

type Page[T any] struct {
	Content          []T `json:"content"`
	Page             int `json:"page"`
	Size             int `json:"size"`
	Total            int `json:"total"`
	NumberOfElements int `json:"numberOfElements"`
}

func NewPage[T any](content []T, pageable Pageable, numberOfElements int) (Page[T], error) {
	if numberOfElements < 0 {
		return Page[T]{}, errors.Errorf("invalid numberOfElements %d", numberOfElements)
	}
	total := int(math.Ceil(float64(numberOfElements) / float64(pageable.Size())))

	return Page[T]{
		Content:          content,
		Page:             pageable.Page(),
		Size:             pageable.Size(),
		Total:            total,
		NumberOfElements: numberOfElements,
	}, nil
}
