package pubsub

import (
	"context"
	"sync"
)

type Observer interface {
	Capture(delta int)
	Release()
}

type Handler[T any] interface {
	Handle(ctx context.Context, event *Event[T])
}

type HandlerFunc[T any] func(ctx context.Context, event *Event[T])

func (fn HandlerFunc[T]) Handle(ctx context.Context, event *Event[T]) {
	fn(ctx, event)
}

type Subscriber[T any] interface {
	Handler[T]
	ID() string
	Close(ctx context.Context)
}

type PubSub[T any] struct {
	mx          sync.RWMutex
	subs        []Subscriber[T]
	observer    Observer
	publisher   PublisherHandler[T]
	exporters   ExportHub[T]
	middlewares []SubscriberMiddleware[T]
	subject     string
	closed      bool
}

func (that *PubSub[T]) Subject() string {
	return that.subject
}

func (that *PubSub[T]) Close(ctx context.Context) {
	that.mx.Lock()
	defer that.mx.Unlock()

	if !that.closed {
		that.closed = true
		for _, sub := range that.subs {
			sub.Close(ctx)
		}
	}
}

func (that *PubSub[T]) Subscribe(ctx context.Context, sub Subscriber[T]) {
	that.mx.Lock()
	defer that.mx.Unlock()

	handler := makeSubscriberHandler[T](sub, that.middlewares)
	wrapper := &wrapperSubscription[T]{Subscriber: sub, handler: handler}

	that.subs = append(that.subs, wrapper)
}

func (that *PubSub[T]) SubscribeHandler(ctx context.Context, handler Handler[T]) Subscriber[T] {
	sub := NewSubscription[T](handler)
	that.Subscribe(ctx, sub)
	return sub
}

func (that *PubSub[T]) SubscribeHandlerFunc(ctx context.Context, fn HandlerFunc[T]) Subscriber[T] {
	return that.SubscribeHandler(ctx, fn)
}

func (that *PubSub[T]) SubscribeChannel(ctx context.Context, cap int) Subscriber[T] {
	sub := NewChannelSubscription[T](cap)
	that.Subscribe(ctx, sub)
	return sub
}

func (that *PubSub[T]) Unsubscribe(ctx context.Context, subID string) {
	that.mx.Lock()
	defer that.mx.Unlock()

	for i, sub := range that.subs {
		if sub.ID() == subID {
			sub.Close(ctx)
			that.subs[i] = nil
			that.subs = append(that.subs[:i], that.subs[i+1:]...)
			break
		}
	}
}

func (that *PubSub[T]) Import(
	ctx context.Context,
	maker string,
	entity []byte,
) error {
	e, err := parse[T](entity)
	if err != nil {
		return err
	}

	event := &Event[T]{
		ctx:      ctx,
		subject:  that.subject,
		entity:   e,
		observer: dummyObserver,
		maker:    maker,
	}

	that.publish(ctx, event)
	return nil
}

func (that *PubSub[T]) Publish(ctx context.Context, entity T) Waiter {
	wg := &waitGroup{ctx: ctx, chain: that.observer}

	event := &Event[T]{
		ctx:      ctx,
		subject:  that.subject,
		entity:   entity,
		observer: wg,
	}

	that.publish(ctx, event)

	return wg
}

func (that *PubSub[T]) publish(ctx context.Context, event *Event[T]) {
	event.Capture(1)
	go that.post(ctx, event)
}

func (that *PubSub[T]) post(ctx context.Context, event *Event[T]) {
	defer event.Release()

	that.mx.RLock()
	defer that.mx.RUnlock()

	if that.closed {
		return
	}

	that.publisher.Publish(ctx, event)

	if that.exporters != nil {
		that.exporters.Export(ctx, event)
	}
}

func (that *PubSub[T]) publishEvent(
	ctx context.Context,
	event *Event[T],
) {
	event.Capture(len(that.subs))
	for _, sub := range that.subs {
		go sub.Handle(ctx, event)
	}
}

func (that *PubSub[T]) Attach(exporter Exporter) {
	that.mx.Lock()
	defer that.mx.Unlock()

	if that.exporters == nil {
		that.exporters = NewExporters[T]()
	}

	that.exporters.Attach(exporter)
}
