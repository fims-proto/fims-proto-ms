package domain

func (c IdentifierConfiguration) IsMatchProperties(objectsToMatch map[string]string) bool {
	propertyMatchers := make(map[string]string)
	for _, matcher := range c.propertyMatchers {
		propertyMatchers[matcher.Name()] = matcher.Value()
	}

	for name, value := range objectsToMatch {
		targetValue, ok := propertyMatchers[name]
		if !ok || targetValue != value {
			return false
		}
	}

	return true
}
