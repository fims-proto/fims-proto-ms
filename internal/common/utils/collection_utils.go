package utils

type void struct{}

var empty void

func SliceToMap[E any, K comparable, V any](
	elements []E,
	keyMapper func(element E) K,
	valueMapper func(element E) V,
) map[K]V {
	return SliceToMapMerge(elements, keyMapper, valueMapper, func(_ V, replacement V) V { return replacement })
}

func SliceToMapMerge[E any, K comparable, V any](
	elements []E,
	keyMapper func(element E) K,
	valueMapper func(element E) V,
	mergeFn func(existing V, replacement V) V,
) map[K]V {
	result := make(map[K]V)
	for _, element := range elements {
		existing, ok := result[keyMapper(element)]
		if ok {
			result[keyMapper(element)] = mergeFn(existing, valueMapper(element))
		} else {
			result[keyMapper(element)] = valueMapper(element)
		}
	}
	return result
}

func SliceToSet[E any, K comparable](
	elements []E,
	keyMapper func(element E) K,
) map[K]void {
	// a set is a map that has keys but with empty values
	return SliceToMap(elements, keyMapper, func(element E) void { return empty })
}

func MapToKeySlice[K comparable, V any](m map[K]V) []K {
	keys := make([]K, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
