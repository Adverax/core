package pubsub

import (
	"adverax/core"
	"adverax/core/json"
	"context"
	"reflect"
)

type Pin interface {
	Subject() string
	Attach(exporter Exporter)
	Import(ctx context.Context, maker string, entity []byte) error
}

type MatchFilter interface {
	IsMatch(text string) bool
}

type Publisher interface {
	Publish(ctx context.Context, subject string, entity json.RawMessage)
}

type Gateway struct {
	ps     map[string]Pin
	pub    Publisher
	filter MatchFilter
	id     string
}

func (that *Gateway) Import(
	ctx context.Context,
	subject string,
	entity json.RawMessage,
) error {
	ps, ok := that.ps[subject]
	if !ok {
		return nil
	}

	return ps.Import(ctx, that.id, entity)
}

func (that *Gateway) CanExport(
	ctx context.Context,
	maker, subject string,
) bool {
	if maker == that.id {
		return false
	}

	return that.filter.IsMatch(subject)
}

func (that *Gateway) Export(
	ctx context.Context,
	subject string,
	entity json.RawMessage,
) {
	that.pub.Publish(ctx, subject, entity)
}

func NewGateway(bus interface{}) *Gateway {
	ps := make(map[string]Pin)
	collectPins(bus, ps)

	gateway := &Gateway{id: core.NewGUID(), ps: ps}
	for _, pin := range ps {
		pin.Attach(gateway)
	}

	return gateway
}

func collectPins(bus interface{}, pins map[string]Pin) {
	val := reflect.ValueOf(bus)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		if field.Type().Implements(reflect.TypeOf((*Pin)(nil)).Elem()) {
			pin := field.Interface().(Pin)
			pins[pin.Subject()] = pin
			continue
		}
		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		if field.Kind() == reflect.Struct {
			if field.CanInterface() {
				collectPins(field.Interface(), pins)
			}
		}
	}
}
