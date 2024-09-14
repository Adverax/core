package log

import (
	"io"
)

type Constructor func(output io.Writer, level Level) (Logger, error)

type Factory struct {
	logs        map[Level]Logger
	output      io.WriteCloser
	constructor Constructor
}

func (that *Factory) Close() error {
	return that.output.Close()
}

func (that *Factory) NewLogger(level Level) (Logger, error) {
	if l, ok := that.logs[level]; ok {
		return l, nil
	}

	l, err := that.constructor(that.output, level)
	if err != nil {
		return nil, err
	}
	that.logs[level] = l
	return l, nil
}

func NewLoggerFactory(output io.WriteCloser, constructor Constructor) *Factory {
	if constructor == nil {
		constructor = DefaultConstructor
	}

	return &Factory{
		constructor: constructor,
		logs:        make(map[Level]Logger),
		output:      output,
	}
}

var DefaultConstructor = func(output io.Writer, level Level) (Logger, error) {
	l, err := NewLogrusBuilder().
		Output(output).
		Level(level).
		Formatter(NewDefaultFormatter()).
		Build()
	if err != nil {
		return nil, err
	}
	return NewLogger(l, nil), nil
}
