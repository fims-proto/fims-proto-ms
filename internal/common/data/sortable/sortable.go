package sortable

type Sortable interface {
	IsSorted() bool
	Sorts() []Sort
}

type sortableImpl struct {
	sorts []Sort
}

// new

func Unsorted() Sortable {
	return New()
}

func New(sorts ...Sort) Sortable {
	return sortableImpl{sorts: sorts}
}

// impl

func (s sortableImpl) IsSorted() bool {
	return len(s.sorts) > 0
}

func (s sortableImpl) Sorts() []Sort {
	return s.sorts
}
