package core

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
		if !IsZeroValue(val) {
			return val
		}
	}
	var v T
	return v
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
