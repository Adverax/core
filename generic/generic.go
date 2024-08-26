package generic

import (
	"encoding/json"
	"reflect"
)

var NULL interface{}

type Float interface {
	~float64 | ~float32
}

type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Numeric interface {
	Float | Integer
}

type Ordered interface {
	Integer | Float | ~string
}

// IsString check if a value is a string or a json.Number
func IsString(v interface{}) bool {
	switch v.(type) {
	case string:
	case json.Number:
	default:
		return false
	}
	return true
}

// IsNumeric check if a value is a numeric type
func IsNumeric(v interface{}) bool {
	switch v.(type) {
	case int:
	case int64:
	case int32:
	case int16:
	case int8:
	case uint:
	case uint64:
	case uint32:
	case uint16:
	case uint8:
	default:
		return false
	}
	return true
}

// IsFloat check if a value is a float type
func IsFloat(v interface{}) bool {
	switch v.(type) {
	case float64:
	case float32:
	default:
		return false
	}
	return true
}

func Max[T Numeric](a, b T) T {
	if a > b {
		return a
	}
	return b
}

func Min[T Numeric](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func If[T any](cond bool, a, b T) T {
	if cond {
		return a
	}
	return b
}

func Coalesce[T any](vals ...T) T {
	for _, val := range vals {
		if !IsZeroVal(val) {
			return val
		}
	}
	var v T
	return v
}

func CoalesceAny(vals ...any) interface{} {
	for _, val := range vals {
		if !IsZeroVal(val) {
			return val
		}
	}
	var empty interface{}
	return empty
}

func IndexOf[T comparable](vals []T, val T) int {
	for i, v := range vals {
		if v == val {
			return i
		}
	}
	return -1
}

func Append[T comparable](vals []T, val T) []T {
	if IndexOf(vals, val) != -1 {
		return vals
	}
	return append(vals, val)
}

func Compare[T Ordered](a, b T) int {
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}

// IsNil check if a value is nil
func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Map, reflect.Pointer, reflect.UnsafePointer, reflect.Interface, reflect.Slice:
		return val.IsNil()
	}
	return false
}

// IsZeroVal check if any type is its zero value
func IsZeroVal(x interface{}) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// Fetch returns the value of a pointer or the zero value of the type
func Fetch[T any](v *T) T {
	if v == nil {
		var empty T
		return empty
	}
	return *v
}

// Must panics if an error is not nil
func Must[T any](value T, err error) T {
	if err != nil {
		panic(err)
	}
	return value
}
