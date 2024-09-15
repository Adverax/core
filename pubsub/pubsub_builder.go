package pubsub

import (
	"fmt"
	"github.com/Adverax/core"
)

type Builder[T any] struct {
	*core.Builder
	pubsub *PubSub[T]
	pm     []PublisherMiddleware[T]
}

func NewBuilder[T any]() *Builder[T] {
	return &Builder[T]{
		Builder: core.NewBuilder("PubSub"),
		pubsub:  &PubSub[T]{},
	}
}

func (that *Builder[T]) Subject(subject string) *Builder[T] {
	that.pubsub.subject = subject
	return that
}

func (that *Builder[T]) Observer(observer Observer) *Builder[T] {
	that.pubsub.observer = observer
	return that
}

func (that *Builder[T]) Exporter(hub ExportHub[T]) *Builder[T] {
	that.pubsub.exporters = hub
	return that
}

func (that *Builder[T]) PublisherMiddlewares(middlewares ...PublisherMiddleware[T]) *Builder[T] {
	that.pm = append(that.pm, middlewares...)
	return that
}

func (that *Builder[T]) Middlewares(middlewares ...SubscriberMiddleware[T]) *Builder[T] {
	that.pubsub.middlewares = append(that.pubsub.middlewares, middlewares...)
	return that
}

func (that *Builder[T]) Build() (*PubSub[T], error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}

	if err := that.updateDefaultFields(); err != nil {
		return nil, err
	}

	that.pubsub.publisher = makePublisherHandler[T](
		PublisherHandlerFunc[T](that.pubsub.publishEvent),
		that.pm,
	)

	return that.pubsub, nil
}

func (that *Builder[T]) checkRequiredFields() error {
	that.RequiredField(that.pubsub.subject, ErrFieldSubjectIsRequired)

	return that.ResError()
}

func (that *Builder[T]) updateDefaultFields() error {
	if that.pubsub.observer == nil {
		that.pubsub.observer = dummyObserver
	}

	return that.ResError()
}

var (
	ErrFieldSubjectIsRequired = fmt.Errorf("Field 'subject' is required")
)

func makePublisherHandler[T any](
	handler PublisherHandler[T],
	middlewares []PublisherMiddleware[T],
) PublisherHandler[T] {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}
	return handler
}
