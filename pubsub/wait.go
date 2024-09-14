package pubsub

import (
	"context"
	"sync"
)

type Waiter interface {
	Wait() error
}

type waitGroup struct {
	sync.WaitGroup
	chain Observer
	ctx   context.Context
}

func (that *waitGroup) Wait() error {
	done := make(chan struct{})
	go func() {
		that.WaitGroup.Wait()
		close(done)
	}()

	select {
	case <-that.ctx.Done():
		return that.ctx.Err()
	case <-done:
		return nil
	}
}

func (that *waitGroup) Capture(delta int) {
	that.WaitGroup.Add(delta)
	that.chain.Capture(delta)
}

func (that *waitGroup) Release() {
	that.WaitGroup.Done()
	that.chain.Release()
}

type dummyWaitGroup struct{}

func (that *dummyWaitGroup) Capture(delta int) {}
func (that *dummyWaitGroup) Release()          {}
func (that *dummyWaitGroup) Wait() error       { return nil }

var dummyObserver = &dummyWaitGroup{}
