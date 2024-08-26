package types

import (
	"context"
	"encoding/json"
	"errors"
	"time"
)

// Getter is abstract property props
type Getter interface {
	GetProperty(ctx context.Context, name string) (interface{}, error)
}

// Setter is abstract property setter
type Setter interface {
	SetProperty(ctx context.Context, name string, value interface{}) error
}

type Typer[T any] interface {
	Get(ctx context.Context, getter Getter, name string, defVal T) (res T, err error)
	TryCast(value interface{}) (T, bool)
	Cast(value interface{}, defVal T) T
}

type Types struct {
	Boolean  Typer[bool]
	Integer  Typer[int64]
	Float    Typer[float64]
	String   Typer[string]
	Duration Typer[time.Duration]
	Json     Typer[json.RawMessage]
}

var Type = &Types{
	Boolean:  &BooleanType{},
	Integer:  &IntegerType{},
	Float:    &FloatType{},
	String:   &StringType{},
	Duration: &DurationType{},
	Json:     &JsonType{},
}

var GetErrNoMatch = func() error {
	return errNoMatch
}

var errNoMatch = errors.New("no match")
