package data

import "github.com/pkg/errors"

type Pageable interface {
	Page() int
	Size() int
	Offset() int
	Sorts() []Sort
	Chooses() []Choose
}

type pageRequest struct {
	page    int
	size    int
	offset  int
	sorts   []Sort
	chooses []Choose
}

func NewPageRequest(page, size int, sortFields map[string]string, chooseFields []string) (Pageable, error) {
	if page < 1 {
		return nil, errors.New("zero page number. page number starts with 1")
	}
	if size < 1 {
		return nil, errors.New("zero page size")
	}

	offset := (page - 1) * size

	var sorts []Sort
	for field, order := range sortFields {
		sort, err := newSort(field, order)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create sorts request")
		}
		sorts = append(sorts, sort)
	}

	var chooses []Choose
	for _, field := range chooseFields {
		choose, err := newChoose(field)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create chooses request")
		}
		chooses = append(chooses, choose)
	}

	return pageRequest{
		page:    page,
		size:    size,
		offset:  offset,
		sorts:   sorts,
		chooses: chooses,
	}, nil
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
