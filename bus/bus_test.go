package bus

import (
	"context"
	"testing"
)

type MyBus struct {
	Master Channel[bool]
}

func NewBus(logger Logger) *MyBus {
	return &MyBus{
		Master: NewChannel[bool]("master", logger),
	}
}

func TestNewChannel(t *testing.T) {
	ctx := context.Background()
	bus := NewBus(nil)
	bus.Master.Subscribe(
		ctx,
		func(ctx context.Context, notification bool) error {
			return nil
		},
	)
}
