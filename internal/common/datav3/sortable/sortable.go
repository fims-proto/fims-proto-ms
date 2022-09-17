package sortable

type Sortable interface {
	IsSorted() bool
	Sorts() []Sort
}

type unsorted struct{}

type sortableImpl struct {
	sorts []Sort
}

// new

func Unsorted() Sortable {
	return unsorted{}
}

func New(sorts []Sort) Sortable {
	return sortableImpl{sorts: sorts}
}

// impl

func (u unsorted) IsSorted() bool {
	return false
}

func (u unsorted) Sorts() []Sort {
	return nil
}

func (s sortableImpl) IsSorted() bool {
	return true
}

func (s sortableImpl) Sorts() []Sort {
	return s.sorts
}
