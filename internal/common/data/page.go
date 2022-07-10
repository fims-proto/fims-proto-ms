package data

import (
	"math"

	"github.com/pkg/errors"
)

type Page[T any] interface {
	Content() []T
	Page() int
	Size() int
	Total() int
	NumberOfElements() int
}

type pageImpl[T any] struct {
	Content          []T `json:"content"`
	Page             int `json:"page"`
	Size             int `json:"size"`
	Total            int `json:"total"`
	NumberOfElements int `json:"numberOfElements"`
}

type pageImplWrapper[T any] struct {
	*pageImpl[T]
}

func NewPage[T any](content []T, pageable Pageable, numberOfElements int) (Page[T], error) {
	if numberOfElements < 0 {
		return pageImplWrapper[T]{}, errors.Errorf("invalid numberOfElements %d", numberOfElements)
	}
	total := int(math.Ceil(float64(numberOfElements) / float64(pageable.Size())))

	return pageImplWrapper[T]{
		&pageImpl[T]{
			Content:          content,
			Page:             pageable.Page(),
			Size:             pageable.Size(),
			Total:            total,
			NumberOfElements: numberOfElements,
		},
	}, nil
}

func (p pageImplWrapper[T]) Content() []T {
	return p.pageImpl.Content
}

func (p pageImplWrapper[T]) Page() int {
	return p.pageImpl.Page
}

func (p pageImplWrapper[T]) Size() int {
	return p.pageImpl.Size
}

func (p pageImplWrapper[T]) Total() int {
	return p.pageImpl.Total
}

func (p pageImplWrapper[T]) NumberOfElements() int {
	return p.pageImpl.NumberOfElements
}
