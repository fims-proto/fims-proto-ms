package data

import (
	"math"

	"github.com/pkg/errors"
)

type Page[T any] struct {
	Content          []T  `json:"content"`
	Page             int  `json:"page"`
	Size             int  `json:"size"`
	Total            int  `json:"total"`
	NumberOfElements int  `json:"numberOfElements"`
	IsFirst          bool `json:"isFirst"`
	IsLast           bool `json:"isLast"`
}

func NewPage[T any](content []T, page, size, numberOfElements int) (Page[T], error) {
	if page < 1 {
		return Page[T]{}, errors.Errorf("invalid page number %d", page)
	}
	if size < 1 {
		return Page[T]{}, errors.Errorf("invalid page size %d", size)
	}
	if numberOfElements < 0 {
		return Page[T]{}, errors.Errorf("invalid numberOfElements %d", numberOfElements)
	}
	total := int(math.Ceil(float64(numberOfElements) / float64(size)))
	isFirst := page == 1
	isLast := page == total

	return Page[T]{
		Content:          content,
		Page:             page,
		Size:             size,
		Total:            total,
		NumberOfElements: numberOfElements,
		IsFirst:          isFirst,
		IsLast:           isLast,
	}, nil
}
