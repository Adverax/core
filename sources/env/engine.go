package envSource

import (
	"os"
	"strings"
)

type Guard interface {
	IsSatisfied(text string) (key string, matched bool)
}

type Accumulator interface {
	Add(key, value string)
	Result() map[string]interface{}
}

type Consumer interface {
	Start() Accumulator
}

type Engine struct {
	guard    Guard
	consumer Consumer
}

func New(guard Guard, consumer Consumer) *Engine {
	return &Engine{
		guard:    guard,
		consumer: consumer,
	}
}

func (that *Engine) Fetch() (map[string]interface{}, error) {
	return that.fetch(os.Environ())
}

func (that *Engine) fetch(es []string) (map[string]interface{}, error) {
	accumulator := that.consumer.Start()

	for _, e := range es {
		ss := strings.Split(e, "=")
		if len(ss) < 2 {
			continue
		}

		key, ok := that.guard.IsSatisfied(ss[0])
		if !ok {
			continue
		}

		accumulator.Add(key, ss[1])
	}

	return accumulator.Result(), nil
}
