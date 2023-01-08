package filterable

type Filterable interface {
	IsFiltered() bool
	Filters() []Filter
	AddFilter(f Filter)
}

type filterableImpl struct {
	filters []Filter
}

// new

func Unfiltered() Filterable {
	return New()
}

func New(filters ...Filter) Filterable {
	return &filterableImpl{filters: filters}
}

// impl

func (f *filterableImpl) IsFiltered() bool {
	return len(f.filters) > 0
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
