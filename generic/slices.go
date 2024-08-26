package generic

// ToInterfaces converts a slice of any type to a slice of interfaces
func ToInterfaces[T any](values []T) []interface{} {
	result := make([]interface{}, len(values))
	for i, value := range values {
		result[i] = value
	}
	return result
}

// SkipNils skips nil values in a slice
func SkipNils[T any](values []T) []T {
	var result []T
	for _, value := range values {
		if !IsZeroVal(value) {
			result = append(result, value)
		}
	}
	return result
}

// Traverse calls a function for each item in a slice
func Traverse[T any, K comparable](items []T, action func(T) K) []K {
	values := make([]K, 0, len(items))
	var empty K
	for _, item := range items {
		value := action(item)
		if value != empty {
			values = append(values, value)
		}
	}
	return values
}

// KeysOf returns keys of a map
func KeysOf[T any, K comparable](m map[K]T) []K {
	keys := make([]K, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return keys
}

// ValuesOf returns values of a map
func ValuesOf[T any, K comparable](m map[K]T) []T {
	values := make([]T, 0, len(m))
	for _, value := range m {
		values = append(values, value)
	}
	return values
}

func IsStrings(values []interface{}) bool {
	for _, value := range values {
		if !IsString(value) {
			return false
		}
	}

	return true
}

func IsNumerics(values []interface{}) bool {
	for _, value := range values {
		if !IsNumeric(value) {
			return false
		}
	}
	return true
}

func IsFloats(values []interface{}) bool {
	for _, value := range values {
		if !IsFloat(value) {
			return false
		}
	}
	return true
}
