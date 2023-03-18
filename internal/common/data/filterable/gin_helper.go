package filterable

func NewFilterableFromQuery(filter string) (Filterable, error) {
	if filter == "" {
		return Unfiltered(), nil // TODO Lond create Unfiltered type
	}
	filterExpr := &FilterExpr{Buffer: filter}
	filterExpr.Init()
	filterExpr.Print()
	if err := filterExpr.Parse(); err != nil {
		return Unfiltered(), err
	}
	filterNode, err := filterExpr.ParseAsFilterable()
	if err != nil {
		return nil, err
	}
	return filterNode, nil
}
