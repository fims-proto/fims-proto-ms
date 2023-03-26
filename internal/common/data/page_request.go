package data

import (
	"github/fims-proto/fims-proto-ms/internal/common/data/filterable"
	"github/fims-proto/fims-proto-ms/internal/common/data/pageable"
	"github/fims-proto/fims-proto-ms/internal/common/data/sortable"
)

type PageRequest interface {
	pageable.Pageable
	sortable.Sortable
	filterable.Filterable
	GetRawFilterable() filterable.Filterable
	AddAndFilterable(filterbal filterable.Filterable)
}

type pageRequestImpl struct {
	p pageable.Pageable
	s sortable.Sortable
	f filterable.Filterable
}

func NewPageRequest(p pageable.Pageable, s sortable.Sortable, f filterable.Filterable) PageRequest {
	return &pageRequestImpl{
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

func (p pageRequestImpl) Children() []filterable.Filterable {
	return p.f.Children()
}

func (p pageRequestImpl) FilterableType() filterable.FilterableType {
	return filterable.TypeRequest
}

func (p *pageRequestImpl) GetRawFilterable() filterable.Filterable {
	return p.f
}

func (p *pageRequestImpl) AddAndFilterable(fb filterable.Filterable) {
	if p.f.FilterableType() == filterable.TypeNONE {
		p.f = fb
	} else {
		newFilterable := filterable.NewFilterable(filterable.TypeAND, p.f, fb)
		p.f = newFilterable
	}
}
