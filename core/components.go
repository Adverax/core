package core

import (
	"reflect"
	"sync"
)

type component struct {
	initializer func() error
	finalizer   func()
}

type componentCollection struct {
	mx   sync.Mutex
	list []*component
	ok   bool
}

func (c *componentCollection) append(component *component) {
	c.mx.Lock()
	defer c.mx.Unlock()

	c.list = append(c.list, component)
}

func (c *componentCollection) Init() {
	c.mx.Lock()
	defer c.mx.Unlock()

	for _, cc := range c.list {
		err := cc.initializer()
		if err != nil {
			c.Fail(err)
		}
	}
}

func (c *componentCollection) Done() {
	c.mx.Lock()
	defer c.mx.Unlock()

	for _, cc := range c.list {
		cc.finalizer()
	}
}

func (c *componentCollection) Fail(err error) {
	c.ok = false
	// log.Fatal(err)
}

var components = new(componentCollection)

func NewComponent[T interface{}](
	constructor func() T,
) func() T {
	return NewComponentEx(constructor, nil, nil)
}

func NewComponentEx[T interface{}](
	constructor func() T,
	initializer func(instance T) error,
	finalizer func(),
) func() T {
	var instance T

	return func() T {
		if isZeroVal(instance) {
			instance = constructor()
			if initializer != nil || finalizer != nil {
				components.append(&component{
					initializer: func() error {
						if initializer != nil {
							return initializer(instance)
						}
						return nil
					},
					finalizer: func() {
						if finalizer != nil {
							finalizer()
						}
					},
				})
			}
		}

		return instance
	}
}

// IsZeroVal check if any type is its zero value
func isZeroVal(x interface{}) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
