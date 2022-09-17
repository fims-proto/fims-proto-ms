package datav3

import (
	"math"

	"github/fims-proto/fims-proto-ms/internal/common/datav3/pageable"

	"github.com/pkg/errors"
)

type Page[T any] interface {
	Content() []T
	PageNumber() int
	PageSize() int
	TotalPage() int
	NumberOfElements() int
}

type pageImpl[T any] struct {
	Content          []T `json:"content"`
	PageNumber       int `json:"pageNumber"`
	PageSize         int `json:"pageSize"`
	TotalPage        int `json:"totalPage"`
	NumberOfElements int `json:"numberOfElements"`
}

type pageImplWrapper[T any] struct {
	*pageImpl[T]
}

func NewPage[T any](content []T, p pageable.Pageable, numberOfElements int) (Page[T], error) {
	if numberOfElements < 0 {
		return pageImplWrapper[T]{}, errors.Errorf("invalid numberOfElements %d", numberOfElements)
	}
	total := int(math.Ceil(float64(numberOfElements) / float64(p.PageSize())))

	return pageImplWrapper[T]{
		&pageImpl[T]{
			Content:          content,
			PageNumber:       p.PageNumber(),
			PageSize:         p.PageSize(),
			TotalPage:        total,
			NumberOfElements: numberOfElements,
		},
	}, nil
}

func (p pageImplWrapper[T]) Content() []T {
	return p.pageImpl.Content
}

func (p pageImplWrapper[T]) PageNumber() int {
	return p.pageImpl.PageNumber
}

func (p pageImplWrapper[T]) PageSize() int {
	return p.pageImpl.PageSize
}

func (p pageImplWrapper[T]) TotalPage() int {
	return p.pageImpl.TotalPage
}

func (p pageImplWrapper[T]) NumberOfElements() int {
	return p.pageImpl.NumberOfElements
}
