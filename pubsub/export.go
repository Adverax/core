package pubsub

import (
	"adverax/core/json"
	"context"
	"sync"
)

type Exporter interface {
	CanExport(ctx context.Context, maker, subject string) bool
	Export(ctx context.Context, subject string, entity json.RawMessage)
}

type ExportHub[T any] interface {
	Export(ctx context.Context, event *Event[T])
	Attach(exporter Exporter)
}

type Exporters[T any] struct {
	mx        sync.Mutex
	exporters []Exporter
}

func NewExporters[T any]() *Exporters[T] {
	return &Exporters[T]{}
}

func (that *Exporters[T]) Attach(exporter Exporter) {
	that.mx.Lock()
	defer that.mx.Unlock()

	that.exporters = append(that.exporters, exporter)
}

func (that *Exporters[T]) Detach(exporter Exporter) {
	that.mx.Lock()
	defer that.mx.Unlock()

	for i, e := range that.exporters {
		if e == exporter {
			that.exporters = append(that.exporters[:i], that.exporters[i+1:]...)
			break
		}
	}
}

func (that *Exporters[T]) Export(ctx context.Context, event *Event[T]) {
	that.mx.Lock()
	defer that.mx.Unlock()

	var exporters []Exporter
	for _, exporter := range that.exporters {
		if exporter.CanExport(ctx, event.maker, event.subject) {
			exporters = append(exporters, exporter)
		}
	}

	if len(exporters) == 0 {
		return
	}

	data, err := json.Marshal(event.entity)
	if err != nil {
		return
	}

	event.Capture(len(that.exporters))
	for _, exporter := range exporters {
		go func(exporter Exporter) {
			defer event.Release()
			exporter.Export(ctx, event.subject, data)
		}(exporter)
	}
}
