package filterable

import (
	"fmt"
	"regexp"
	"strings"
)

func NewFilterableFromQuery(filter string) (Filterable, error) {
	if filter == "" {
		return Unfiltered(), nil
	}

	andSep := regexp.MustCompile(`\band\b|\bAND\b`)
	spaceSep := regexp.MustCompile(`\s+`)
	commaSep := regexp.MustCompile(`\s*,\s*`)

	conditions := andSep.Split(filter, -1)

	var filters []Filter

	for _, raw := range conditions {
		condition := strings.TrimSpace(raw)
		components := spaceSep.Split(condition, -1)
		if len(components) < 3 {
			return nil, fmt.Errorf("invalid filter %s", condition)
		}
		field := components[0]
		opt := components[1]
		valuesString := strings.Join(components[2:], " ")
		values := commaSep.Split(valuesString, -1)

		cleanValues := removeSingleQuote(values)
		f, err := NewFilter(field, opt, cleanValues...)
		if err != nil {
			return nil, err
		}
		filters = append(filters, f)
	}

	return New(filters...), nil
}

func removeSingleQuote(values []string) []string {
	singleQuoted := regexp.MustCompile(`^'(.+)'$`)
	for i := 0; i < len(values); i++ {
		if singleQuoted.MatchString(values[i]) {
			values[i] = values[i][1 : len(values[i])-1]
		}
	}
	return values
}
