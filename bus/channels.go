package bus

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

type Logger interface {
	Error(ctx context.Context, msg string)
	Info(ctx context.Context, msg string)
}

type Subscriber[T any] func(ctx context.Context, notification T) error

type subscribers[T any] []Subscriber[T]

type Notifier[T any] interface {
	Notify(ctx context.Context, notification T)
}

type Registrar[T any] interface {
	Subscribe(ctx context.Context, subscriber Subscriber[T])
}

type Channel[T any] interface {
	Publish(ctx context.Context, notification T)
	Notifier[T]
	Registrar[T]
}

type channel[T any] struct {
	sync.Mutex
	subject     string
	logger      Logger
	subscribers subscribers[T]
}

func (c *channel[T]) Subscribe(ctx context.Context, action Subscriber[T]) {
	c.logger.Info(ctx, fmt.Sprintf("NEW SUBSCRIBER: %s", c.subject))
	c.subscribe(ctx, action)
}

func (c *channel[T]) subscribe(ctx context.Context, action Subscriber[T]) {
	c.Lock()
	defer c.Unlock()

	c.subscribers = append(c.subscribers, action)
}

func (c *channel[T]) Notify(ctx context.Context, notification T) {
	c.Publish(ctx, notification)
}

func (c *channel[T]) Publish(ctx context.Context, notification T) {
	data, err := json.Marshal(notification)
	if err != nil {
		c.logger.Error(ctx, err.Error())
	}
	c.logger.Info(ctx, fmt.Sprintf("EVENT %s: %s", c.subject, string(data)))

	c.publish(ctx, notification)
}

func (c *channel[T]) publish(ctx context.Context, notification T) {
	c.Lock()
	defer c.Unlock()

	for _, subscriber := range c.subscribers {
		err := subscriber(ctx, notification)
		if err != nil {
			c.logger.Info(ctx, err.Error())
		}
	}
}

func NewChannel[T any](subject string, logger Logger) Channel[T] {
	return &channel[T]{
		subject:     subject,
		logger:      logger,
		subscribers: make(subscribers[T], 4),
	}
}
