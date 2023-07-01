package core

import (
	"reflect"
	"sync"
)

type components struct {
	mx   sync.Mutex
	list []interface{}
}

func (c *components) append(component interface{}) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.list = append(c.list, component)
}

func (c *components) Init() {
	c.mx.Lock()
	defer c.mx.Unlock()

	for _, component := range c.list {
		if initializer, ok := component.(ComponentWithInitialization); ok {
			initializer.InitComponent()
		}
	}
}

func (c *components) Done() {
	c.mx.Lock()
	defer c.mx.Unlock()

	for _, component := range c.list {
		if closer, ok := component.(ComponentWithFinalization); ok {
			closer.CloseComponent()
		}
	}
}

var Components components

type ComponentWithInitialization interface {
	InitComponent()
}

type ComponentWithFinalization interface {
	CloseComponent()
}

func NewComponent[T interface{}](builder func() T) func() T {
	var instance T

	return func() T {
		if isZeroVal(instance) {
			instance = builder()
			Components.append(instance)
		}

		return instance
	}
}

// IsZeroVal check if any type is its zero value
func isZeroVal(x interface{}) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
