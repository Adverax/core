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
