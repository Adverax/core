package enums

import (
	"fmt"
	"reflect"
)

// Enum is abstract enumeration
type Enum[T any] struct {
	dict  map[string]int
	items []T
}

func (e *Enum[T]) Of(tType string) (val T, err error) {
	index, ok := e.dict[tType]
	if !ok {
		return val, ErrIllegalMember(reflect.TypeOf(val).Name(), tType)
	}

	return e.items[index], nil
}

func (e *Enum[T]) GetAll() []T {
	return e.items
}

// NewEnum creates new enumeration
func NewEnum[T ~string](values ...T) *Enum[T] {
	enum := &Enum[T]{
		dict:  make(map[string]int, len(values)),
		items: make([]T, len(values)),
	}
	for index, v := range values {
		enum.items[index] = v
		enum.dict[string(v)] = index
	}
	return enum
}

func ErrIllegalMember(aType, aItem string) error {
	return fmt.Errorf("illegal member %q for type %q", aItem, aType)
}
