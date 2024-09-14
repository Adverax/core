package pubsub

import (
	"context"
)

type Event[T any] struct {
	ctx      context.Context
	observer Observer
	subject  string
	maker    string
	entity   T
}

func (that *Event[T]) Capture(delta int) {
	that.observer.Capture(delta)
}

func (that *Event[T]) Release() {
	that.observer.Release()
}

func (that *Event[T]) Context() context.Context {
	return that.ctx
}

func (that *Event[T]) Subject() string {
	return that.subject
}

func (that *Event[T]) Entity() T {
	return that.entity
}
