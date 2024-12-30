package core

import (
	"context"
	uuid "github.com/satori/go.uuid"
)

func NewGUID() string {
	return uuid.NewV4().String()
}

type nonCancelableContext struct {
	context.Context
}

func NewNonCancelableContext(ctx context.Context) context.Context {
	return &nonCancelableContext{Context: ctx}
}

func (that *nonCancelableContext) Done() <-chan struct{} {
	return dummyDone
}

var dummyDone = make(chan struct{})

type Aggregator[T any] interface {
	Aggregate(ctx context.Context, value T) error
}

type AggregatorFunc[T any] func(context.Context, T) error

func (fn AggregatorFunc[T]) Aggregate(ctx context.Context, value T) error {
	return fn(ctx, value)
}
