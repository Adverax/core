package pubsub

import (
	"context"
	"github.com/google/uuid"
	"sync"
)

type Handler[T any] func(ctx context.Context, msg T)

type Subscriber[T any] struct {
	id string
	ch chan T
}

func (that *Subscriber[T]) ID() string {
	return that.id
}

func (that *Subscriber[T]) Channel() <-chan T {
	return that.ch
}

type PubSub[T any] struct {
	closed bool
	mx     sync.RWMutex
	subs   []*Subscriber[T]
}

func (that *PubSub[T]) Handler(ctx context.Context, handler Handler[T]) {
	go func(sub *Subscriber[T]) {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-sub.ch:
				handler(ctx, msg)
			}
		}
	}(that.Subscribe(ctx))
}

func (that *PubSub[T]) Subscribe(ctx context.Context) *Subscriber[T] {
	that.mx.Lock()
	defer that.mx.Unlock()

	sub := &Subscriber[T]{
		id: uuid.New().String(),
		ch: make(chan T, 1),
	}
	that.subs = append(that.subs, sub)
	return sub
}

func (that *PubSub[T]) Unsubscribe(ctx context.Context, subID string) {
	that.mx.Lock()
	defer that.mx.Unlock()

	for i, sub := range that.subs {
		if sub.id == subID {
			close(sub.ch)
			that.subs = append(that.subs[:i], that.subs[i+1:]...)
			break
		}
	}
}

func (that *PubSub[T]) Publish(ctx context.Context, msg T) {
	go func() {
		that.mx.RLock()
		defer that.mx.RUnlock()

		if that.closed {
			return
		}

		for _, sub := range that.subs {
			go func(sub *Subscriber[T]) {
				sub.ch <- msg
			}(sub)
		}
	}()
}

func (that *PubSub[T]) Close() {
	that.mx.Lock()
	defer that.mx.Unlock()

	if !that.closed {
		that.closed = true
		for _, sub := range that.subs {
			close(sub.ch)
		}
	}
}

func New[T any]() *PubSub[T] {
	return &PubSub[T]{}
}
