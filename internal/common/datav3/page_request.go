package datav3

import (
	"github/fims-proto/fims-proto-ms/internal/common/datav3/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/datav3/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/datav3/sortable"
)

type PageRequest interface {
	pageable.Pageable
	sortable.Sortable
	filterable.Filterable
}

type pageRequestImpl struct {
	p pageable.Pageable
	s sortable.Sortable
	f filterable.Filterable
}

func NewPageRequest(p pageable.Pageable, s sortable.Sortable, f filterable.Filterable) PageRequest {
	return pageRequestImpl{
		p: p,
		s: s,
		f: f,
	}
}

// impl

func (p pageRequestImpl) IsPaged() bool {
	return p.p.IsPaged()
}

func (p pageRequestImpl) PageNumber() int {
	return p.p.PageNumber()
}

func (p pageRequestImpl) PageSize() int {
	return p.p.PageSize()
}

func (p pageRequestImpl) Offset() int {
	return p.p.Offset()
}

func (p pageRequestImpl) IsSorted() bool {
	return p.s.IsSorted()
}

func (p pageRequestImpl) Sorts() []sortable.Sort {
	return p.s.Sorts()
}

func (p pageRequestImpl) IsFiltered() bool {
	return p.f.IsFiltered()
}

func (p pageRequestImpl) Filters() []filterable.Filter {
	return p.f.Filters()
}

func (p pageRequestImpl) AddFilter(f filterable.Filter) {
	p.f.AddFilter(f)
}
