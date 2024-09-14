package pubsub

import (
	"adverax/core/json"
	"adverax/core/log"
	"context"
)

type SubscriberMiddleware[T any] func(handler Handler[T]) Handler[T]

func WithRecover[T any](logger log.Logger) SubscriberMiddleware[T] {
	return func(handler Handler[T]) Handler[T] {
		return HandlerFunc[T](func(ctx context.Context, event *Event[T]) {
			defer func() {
				if r := recover(); r != nil {
					logger.Errorf(ctx, "panic recovered: %v", r)
				}
			}()

			handler.Handle(ctx, event)
		})
	}
}

func WithLogger[T any](logger log.Logger) SubscriberMiddleware[T] {
	return func(handler Handler[T]) Handler[T] {
		return HandlerFunc[T](func(ctx context.Context, event *Event[T]) {
			e, _ := json.Marshal(event.entity)
			logger.WithFields(
				ctx,
				log.Fields{
					log.FieldKeyEntity:  log.EntityBus,
					log.FieldKeySubject: event.subject,
					log.FieldKeyAction:  log.ActionRequestReceived,
					log.FieldKeyData:    string(e),
				},
			).Debug(ctx, "Event")

			handler.Handle(ctx, event)
		})
	}
}

func makeSubscriberHandler[T any](
	handler Handler[T],
	middlewares []SubscriberMiddleware[T],
) Handler[T] {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
