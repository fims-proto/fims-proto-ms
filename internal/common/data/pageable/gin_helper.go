package pageable

func NewPageableFromQuery(page, size int) (Pageable, error) {
	return New(page, size)
}
