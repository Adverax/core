package log

import (
	"context"
	"errors"
	"fmt"
	"github.com/Adverax/core/log/logrus"
	"strings"
)

type contextLoggerType = int

const contextLoggerKey contextLoggerType = 0

type ContextMode string

const (
	ContextModeNone        ContextMode = ""
	ContextModeTransparent ContextMode = "transparent"
	ContextModeOpaque      ContextMode = "opaque"
)

type staticBehavior struct {
	logger Logger
}

func (that *staticBehavior) NewContext(ctx context.Context) context.Context {
	return ctx
}

func (that *staticBehavior) Resolve(ctx context.Context) Logger {
	return that.logger
}

type opaqueBehavior struct {
	logger Logger
}

func (that *opaqueBehavior) NewContext(ctx context.Context) context.Context {
	return NewContext(ctx, that.logger)
}

func (that *opaqueBehavior) Resolve(ctx context.Context) Logger {
	return getLogger(ctx, that.logger)
}

type transparentBehavior struct {
	logger Logger
}

func (that *transparentBehavior) NewContext(ctx context.Context) context.Context {
	logger := getLogger(ctx, nil)
	if logger == nil {
		return NewContext(ctx, that.logger)
	}

	return ctx
}

func (that *transparentBehavior) Resolve(ctx context.Context) Logger {
	return getLogger(ctx, that.logger)
}

func NewContextLogger(logger Logger, mode ContextMode) Logger {
	var behavior Behavior

	switch mode {
	case ContextModeTransparent:
		behavior = &transparentBehavior{logger: logger}
	case ContextModeOpaque:
		behavior = &opaqueBehavior{logger: logger}
	default:
		behavior = &staticBehavior{logger: logger}
	}

	return NewCustomLogger(logger, behavior)
}

type String interface {
	Get(ctx context.Context) (string, error)
}

type ContextLoggerOptions struct {
	Level   String
	Context String
}

type LevelBuilder interface {
	NewLogger(level Level) (Logger, error)
}

type ContextLoggerFactory struct {
	factory  LevelBuilder
	defaults ContextLoggerOptions
	warnings Logger
	errors   Logger
}

func NewContextLoggerFactory(
	factory LevelBuilder,
	warnings Logger,
	errors Logger,
	defOptions ContextLoggerOptions,
) *ContextLoggerFactory {
	return &ContextLoggerFactory{
		factory:  factory,
		warnings: warnings,
		errors:   errors,
		defaults: defOptions,
	}
}

func (that *ContextLoggerFactory) Warnings() Logger {
	return that.warnings
}

func (that *ContextLoggerFactory) Errors() Logger {
	return that.errors
}

func (that *ContextLoggerFactory) NewLogger(
	ctx context.Context,
	options ContextLoggerOptions,
) (Logger, error) {
	level, err := options.Level.Get(ctx)
	if err != nil {
		return nil, err
	}

	if level == "" {
		level, err = that.defaults.Level.Get(ctx)
		if err != nil {
			return nil, err
		}
	}

	context, err := options.Context.Get(ctx)
	if err != nil {
		return nil, err
	}
	if context == "" {
		context, err = that.defaults.Context.Get(ctx)
		if err != nil {
			return nil, err
		}
	}

	aLevel, err := logrus.ParseLevel(level)
	if err != nil {
		return nil, err
	}

	aContext, err := parseContextMode(context)
	if err != nil {
		return nil, err
	}

	return that.Build(aLevel, aContext)
}

func (that *ContextLoggerFactory) NewHub(ctx context.Context, options ContextLoggerOptions) (IHub, error) {
	logger, err := that.NewLogger(ctx, options)
	if err != nil {
		return nil, err
	}

	return NewHub(logger, that.warnings, that.errors), nil
}

func (that *ContextLoggerFactory) Build(level Level, context ContextMode) (Logger, error) {
	logger, err := that.factory.NewLogger(level)
	if err != nil {
		return nil, err
	}
	return NewContextLogger(logger, context), nil
}

func getLogger(ctx context.Context, defVal Logger) Logger {
	val := ctx.Value(contextLoggerKey)
	if l, ok := val.(Logger); ok {
		return l
	}

	return defVal
}

func parseContextMode(mode string) (ContextMode, error) {
	switch strings.ToLower(mode) {
	case "none":
		return ContextModeNone, nil
	case string(ContextModeNone):
		return ContextModeNone, nil
	case string(ContextModeTransparent):
		return ContextModeTransparent, nil
	case string(ContextModeOpaque):
		return ContextModeOpaque, nil
	default:
		return "", ErrInvalidContextMode
	}
}

var (
	ErrInvalidContextMode = errors.New("invalid context mode")
)

// Resolve returns logger from context
func Resolve(ctx context.Context) Logger {
	log := getLogger(ctx, nil)
	if log == nil {
		panic(fmt.Errorf("logger not found in context: %v", ctx))
	}

	return log
}

// NewContext returns new context with logger
func NewContext(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, contextLoggerKey, logger)
}
