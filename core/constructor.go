package core

import "reflect"

func Constructor[T interface{}](builder func() T) func() T {
	var instance T

	return func() T {
		if isZeroVal(instance) {
			instance = builder()
		}

		return instance
	}
}

// IsZeroVal check if any type is its zero value
func isZeroVal(x interface{}) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
