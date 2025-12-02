package data

import (
	"fmt"
	"math"

	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
)

type Page[T any] interface {
	Content() []T
	PageNumber() int
	PageSize() int
	TotalPage() int
	NumberOfElements() int
}

type PageResponse[T any] struct {
	Content          []T `json:"content"`
	PageNumber       int `json:"pageNumber"`
	PageSize         int `json:"pageSize"`
	TotalPage        int `json:"totalPage"`
	NumberOfElements int `json:"numberOfElements"`
}

type pageImplWrapper[T any] struct {
	*PageResponse[T]
}

func NewPage[T any](content []T, p pageable.Pageable, numberOfElements int) (Page[T], error) {
	if numberOfElements < 0 {
		return pageImplWrapper[T]{}, fmt.Errorf("invalid numberOfElements %d", numberOfElements)
	}
	total := int(math.Ceil(float64(numberOfElements) / float64(p.PageSize())))

	return pageImplWrapper[T]{
		&PageResponse[T]{
			Content:          content,
			PageNumber:       p.PageNumber(),
			PageSize:         p.PageSize(),
			TotalPage:        total,
			NumberOfElements: numberOfElements,
		},
	}, nil
}

func (p pageImplWrapper[T]) Content() []T {
	return p.PageResponse.Content
}

func (p pageImplWrapper[T]) PageNumber() int {
	return p.PageResponse.PageNumber
}

func (p pageImplWrapper[T]) PageSize() int {
	return p.PageResponse.PageSize
}

func (p pageImplWrapper[T]) TotalPage() int {
	return p.PageResponse.TotalPage
}

func (p pageImplWrapper[T]) NumberOfElements() int {
	return p.PageResponse.NumberOfElements
}
