package core

import (
	"context"
	"sync"
)

type Waiter interface {
	Wait()
	WaitWithContext(ctx context.Context) error
}

type Control interface {
	Enter()
	Leave()
}

type WaitGroup interface {
	Waiter
	Control
}

type waitGroup struct {
	wg sync.WaitGroup
}

func (that *waitGroup) Enter() {
	that.wg.Add(1)
}

func (that *waitGroup) Leave() {
	that.wg.Done()
}

func (that *waitGroup) Wait() {
	that.wg.Wait()
}

func (that *waitGroup) WaitWithContext(ctx context.Context) error {
	done := make(chan struct{}, 1)

	go func() {
		that.wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func NewWaitGroup() WaitGroup {
	return &waitGroup{}
}

type DummyWaitGroup struct {
}

func (that *DummyWaitGroup) Enter() {
}

func (that *DummyWaitGroup) Leave() {
}

func (that *DummyWaitGroup) Wait() {
}

func (that *DummyWaitGroup) WaitWithContext(ctx context.Context) error {
	return nil
}

var dummyWaitGroup = new(DummyWaitGroup)

func NewDummyWaitGroup() WaitGroup {
	return dummyWaitGroup
}
