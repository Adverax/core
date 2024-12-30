package core

import "context"

type Handler[T any] interface {
	Handle(ctx context.Context, entity T) error
}

type HandlerFunc[T any] func(ctx context.Context, entity T) error

func (fn HandlerFunc[T]) Handle(ctx context.Context, entity T) error {
	return fn(ctx, entity)
}
