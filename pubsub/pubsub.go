package pubsub

import (
	"context"
	"sync"
)

type Handler[T any] func(ctx context.Context, msg T)

type PubSub[T any] interface {
	Close()
	Subscribe(ctx context.Context, handler Handler[T])
	SubscribeEx(ctx context.Context) <-chan T
	Publish(ctx context.Context, msg T)
}

type pubsub[T any] struct {
	mx     sync.RWMutex
	subs   []chan T
	closed bool
}

func (that *pubsub[T]) Subscribe(ctx context.Context, handler Handler[T]) {
	go func(ch <-chan T) {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-ch:
				handler(ctx, msg)
			}
		}
	}(that.SubscribeEx(ctx))
}

func (that *pubsub[T]) SubscribeEx(ctx context.Context) <-chan T {
	that.mx.Lock()
	defer that.mx.Unlock()

	ch := make(chan T, 1)
	that.subs = append(that.subs, ch)
	return ch
}

func (that *pubsub[T]) Publish(ctx context.Context, msg T) {
	that.mx.RLock()
	defer that.mx.RUnlock()

	if that.closed {
		return
	}

	for _, ch := range that.subs {
		go func(ch chan T) {
			ch <- msg
		}(ch)
	}
}

func (that *pubsub[T]) Close() {
	that.mx.Lock()
	defer that.mx.Unlock()

	if !that.closed {
		that.closed = true
		for _, ch := range that.subs {
			close(ch)
		}
	}
}

func New[T any]() PubSub[T] {
	return &pubsub[T]{}
}
