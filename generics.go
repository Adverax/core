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
