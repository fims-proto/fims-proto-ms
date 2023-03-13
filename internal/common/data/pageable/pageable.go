package pageable

import (
	"github.com/pkg/errors"
)

type Pageable interface {
	IsPaged() bool
	PageNumber() int
	PageSize() int
	Offset() int
}

type unpaged struct{}

type pageableImpl struct {
	page   int
	size   int
	offset int
}

// new

func Unpaged() Pageable {
	return unpaged{}
}

func New(page, size int) (Pageable, error) {
	if page < 1 {
		return nil, errors.New("zero page number. page number starts with 1")
	}
	if size < 1 {
		return nil, errors.New("zero page size")
	}

	offset := (page - 1) * size

	return pageableImpl{
		page:   page,
		size:   size,
		offset: offset,
	}, nil
}

// impl

func (u unpaged) IsPaged() bool {
	return false
}

func (u unpaged) PageNumber() int {
	return 0
}

func (u unpaged) PageSize() int {
	return 0
}

func (u unpaged) Offset() int {
	return 0
}

func (p pageableImpl) IsPaged() bool {
	return true
}

func (p pageableImpl) PageNumber() int {
	return p.page
}

func (p pageableImpl) PageSize() int {
	return p.size
}

func (p pageableImpl) Offset() int {
	return p.offset
}
