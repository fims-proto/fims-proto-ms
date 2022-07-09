package data

import "github.com/pkg/errors"

type Pageable interface {
	IsPaged() bool
	Page() int
	Size() int
	Offset() int
	Sorts() []Sort
	Chooses() []Choose
	Filters() []Filter
}

type pageRequest struct {
	page    int
	size    int
	offset  int
	sorts   []Sort
	chooses []Choose
	filters []Filter
}

func newPageRequest(page, size int, sorts []Sort, chooses []Choose, filters []Filter) (Pageable, error) {
	if page < 1 {
		return nil, errors.New("zero page number. page number starts with 1")
	}
	if size < 1 {
		return nil, errors.New("zero page size")
	}

	offset := (page - 1) * size

	return pageRequest{
		page:    page,
		size:    size,
		offset:  offset,
		sorts:   sorts,
		chooses: chooses,
		filters: filters,
	}, nil
}

func (p pageRequest) IsPaged() bool {
	return true
}

func (p pageRequest) Page() int {
	return p.page
}

func (p pageRequest) Size() int {
	return p.size
}

func (p pageRequest) Offset() int {
	return p.offset
}

func (p pageRequest) Sorts() []Sort {
	return p.sorts
}

func (p pageRequest) Chooses() []Choose {
	return p.chooses
}

func (p pageRequest) Filters() []Filter {
	return p.filters
}

type unpaged struct{}

func Unpaged() Pageable {
	return unpaged{}
}

func (u unpaged) IsPaged() bool {
	return false
}

func (u unpaged) Page() int {
	return 0
}

func (u unpaged) Size() int {
	return 0
}

func (u unpaged) Offset() int {
	return 0
}

func (u unpaged) Sorts() []Sort {
	return nil
}

func (u unpaged) Chooses() []Choose {
	return nil
}

func (u unpaged) Filters() []Filter {
	return nil
}
