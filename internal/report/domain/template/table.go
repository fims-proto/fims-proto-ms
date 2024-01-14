package template

type Table struct {
	header Header
	items  []Item
}

type Header struct {
	text    string
	columns []Cell[string]
}

type Cell[T any] struct {
	key   string
	value T
}
