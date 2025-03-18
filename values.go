package core

import (
	"context"
	"encoding/json"
	"time"
)

type Boolean interface {
	Get(ctx context.Context) (bool, error)
}

type Integer interface {
	Get(ctx context.Context) (int64, error)
}

type Float interface {
	Get(ctx context.Context) (float64, error)
}

type String interface {
	Get(ctx context.Context) (string, error)
}

type Duration interface {
	Get(ctx context.Context) (time.Duration, error)
}

type Json interface {
	Get(ctx context.Context) (json.RawMessage, error)
}

type Strings interface {
	Get(ctx context.Context) ([]string, error)
}

type Enumeration[T comparable] interface {
	Get(ctx context.Context) (T, error)
}

type RWBoolean interface {
	Boolean
	Set(ctx context.Context, value bool) error
}

type RWInteger interface {
	Integer
	Set(ctx context.Context, value int64) error
}

type RWFloat interface {
	Float
	Set(ctx context.Context, value float64) error
}

type RWString interface {
	String
	Set(ctx context.Context, value string) error
}

type RWStrings interface {
	Strings
	Set(ctx context.Context, value []string) error
}

type RWDuration interface {
	Duration
	Set(ctx context.Context, value time.Duration) error
}

type RWJson interface {
	Json
	Set(ctx context.Context, value json.RawMessage) error
}
