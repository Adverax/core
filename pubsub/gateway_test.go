package pubsub

import (
	"context"
	"github.com/Adverax/core/generic"
	"github.com/Adverax/core/json"
	"github.com/stretchr/testify/require"
	"testing"
)

type Receipt struct {
	Id string `json:"id"`
}

type Parcel struct {
	Id string `json:"id"`
}

type Bus struct {
	*ReceiptBus
	*ParcelBus
}

type ReceiptBus struct {
	OnCreated *PubSub[*Receipt]
}

type ParcelBus struct {
	OnCreated *PubSub[*Parcel]
}

func TestBus(t *testing.T) {
	bus := &Bus{
		ReceiptBus: &ReceiptBus{
			OnCreated: generic.Must(
				NewBuilder[*Receipt]().
					Subject("receipt.inserted").
					// PublisherMiddlewares(PublisherMiddlewareLogger[*Receipt](logger)).
					Build(),
			),
		},
		ParcelBus: &ParcelBus{
			OnCreated: generic.Must(
				NewBuilder[*Parcel]().
					Subject("parcel.inserted").
					Build(),
			),
		},
	}

	gateway := NewGateway(bus)
	err := gateway.Import(
		context.Background(),
		"receipt.inserted",
		json.RawMessage(`{"id":"1"}`),
	)
	require.NoError(t, err)
}
