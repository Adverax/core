package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
)

type Logger interface {
	Error(ctx context.Context, msg string)
	Info(ctx context.Context, msg string)
}

type pubsubWithLogger[T any] struct {
	PubSub[T]
	logger  Logger
	subject string
}

func (that *pubsubWithLogger[T]) Subscribe(ctx context.Context, handler Handler[T]) {
	that.logger.Info(ctx, fmt.Sprintf("SUBSCRIBE %s", that.subject))
	that.PubSub.Subscribe(ctx, handler)
}

func (that *pubsubWithLogger[T]) SubscribeEx(ctx context.Context) <-chan T {
	that.logger.Info(ctx, fmt.Sprintf("SUBSCRIBE %s", that.subject))
	return that.PubSub.SubscribeEx(ctx)
}

func (that *pubsubWithLogger[T]) Publish(ctx context.Context, msg T) {
	data, err := json.Marshal(msg)
	if err != nil {
		that.logger.Error(ctx, err.Error())
	}
	that.logger.Info(ctx, fmt.Sprintf("EVENT %s: %s", that.subject, string(data)))
	that.PubSub.Publish(ctx, msg)
}

func NewPubSubWithLogger[T any](subject string, logger Logger) PubSub[T] {
	return &pubsubWithLogger[T]{
		PubSub:  New[T](),
		logger:  logger,
		subject: subject,
	}
}
