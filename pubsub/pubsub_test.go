package pubsub

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type notification struct {
	Subject string `json:"subject"`
	Message string `json:"message"`
}

func TestSync(t *testing.T) {
	ps, err := newTestPubSub[*notification]()
	require.NoError(t, err)

	ctx := context.Background()
	var event *Event[*notification]
	_ = ps.SubscribeHandlerFunc(ctx, func(ctx context.Context, e *Event[*notification]) {
		time.Sleep(10 * time.Millisecond)
		event = e
	})

	notification := &notification{
		Subject: "subject",
		Message: "hello",
	}
	err = ps.Publish(ctx, notification).Wait()
	require.NoError(t, err)

	require.NotNil(t, event)
	assert.Equal(t, notification, event.Entity())
}

func newTestPubSub[T any]() (*PubSub[T], error) {
	return NewBuilder[T]().
		Subject("test").
		Middlewares(WithRecover[T](nil)).
		Build()
}
