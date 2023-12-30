package filterable

func NewFilterableFromQuery(filter string) (Filterable, error) {
	if filter == "" {
		return Unfiltered(), nil // TODO Lind create Unfiltered type
	}
	filterExpr := &FilterAST{Buffer: filter}
	if err := filterExpr.Init(); err != nil {
		return Unfiltered(), err
	}

	if err := filterExpr.Parse(); err != nil {
		return Unfiltered(), err
	}
	filterNode, err := filterExpr.ParseAsFilterable()
	if err != nil {
		return nil, err
	}
	return filterNode, nil
}
