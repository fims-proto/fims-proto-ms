package filterable

type Filterable interface {
	IsFiltered() bool
	Filters() []Filter
	AddFilter(f Filter)
}

type unfiltered struct{}

type filterableImpl struct {
	filters []Filter
}

// new

func Unfiltered() Filterable {
	return unfiltered{}
}

func New(filters ...Filter) Filterable {
	return &filterableImpl{filters: filters}
}

// impl

func (u unfiltered) IsFiltered() bool {
	return false
}

func (u unfiltered) Filters() []Filter {
	return nil
}

func (u unfiltered) AddFilter(Filter) {
	panic("cannot add filter into unfiltered")
}

func (f *filterableImpl) IsFiltered() bool {
	return true
}

func (f *filterableImpl) Filters() []Filter {
	return f.filters
}

func (f *filterableImpl) AddFilter(newFilter Filter) {
	exists := false
	for _, filter := range f.filters {
		if filter.Field().Equals(newFilter.Field()) &&
			filter.Operator() == newFilter.Operator() &&
			sliceEqual(filter.Values(), newFilter.Values()) {
			exists = true
		}
	}

	if !exists {
		f.filters = append(f.filters, newFilter)
	}
}

func sliceEqual(a, b []any) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
