package pubsub

import (
	"adverax/core/json"
	"adverax/core/log"
	"context"
)

type PublisherHandler[T any] interface {
	Publish(ctx context.Context, event *Event[T])
}

type PublisherHandlerFunc[T any] func(ctx context.Context, event *Event[T])

func (fn PublisherHandlerFunc[T]) Publish(ctx context.Context, event *Event[T]) {
	fn(ctx, event)
}

type PublisherMiddleware[T any] func(handler PublisherHandler[T]) PublisherHandler[T]

func PublisherMiddlewareLogger[T any](logger log.Logger) PublisherMiddleware[T] {
	return func(handler PublisherHandler[T]) PublisherHandler[T] {
		return PublisherHandlerFunc[T](func(ctx context.Context, event *Event[T]) {
			ctx = logger.NewContext(ctx)

			e, _ := json.Marshal(event.entity)
			logger.WithFields(
				ctx,
				log.Fields{
					log.FieldKeyEntity:  log.EntityBus,
					log.FieldKeySubject: event.subject,
					log.FieldKeyAction:  log.ActionRequestSent,
					log.FieldKeyData:    string(e),
				},
			).Debug(ctx, "Event")

			handler.Publish(ctx, event)
		})
	}
}
