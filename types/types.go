package types

import (
	"context"
	"encoding/json"
	"github.com/adverax/core/generic"
	"time"
)

type StringType struct{}

func (that *StringType) Get(ctx context.Context, getter Getter, name string, defVal string) (res string, err error) {
	return GetStringProperty(ctx, getter, name, defVal)
}

func (that *StringType) TryCast(value interface{}) (string, bool) {
	return generic.ConvertToString(value)
}

func (that *StringType) Cast(v interface{}, defaults string) string {
	if vv, ok := that.TryCast(v); ok {
		return vv
	}
	return defaults
}

type IntegerType struct{}

func (that *IntegerType) Get(ctx context.Context, getter Getter, name string, defVal int64) (res int64, err error) {
	return GetIntegerProperty(ctx, getter, name, defVal)
}

func (that *IntegerType) TryCast(value interface{}) (int64, bool) {
	return generic.ConvertToInt64(value)
}

func (that *IntegerType) Cast(v interface{}, defaults int64) int64 {
	if vv, ok := that.TryCast(v); ok {
		return vv
	}
	return defaults
}

type FloatType struct{}

func (that *FloatType) Get(ctx context.Context, getter Getter, name string, defVal float64) (res float64, err error) {
	return GetFloatProperty(ctx, getter, name, defVal)
}

func (that *FloatType) TryCast(value interface{}) (float64, bool) {
	return generic.ConvertToFloat64(value)
}

func (that *FloatType) Cast(v interface{}, defaults float64) float64 {
	if vv, ok := that.TryCast(v); ok {
		return vv
	}
	return defaults
}

type BooleanType struct{}

func (that *BooleanType) Get(ctx context.Context, getter Getter, name string, defVal bool) (res bool, err error) {
	return GetBooleanProperty(ctx, getter, name, defVal)
}

func (that *BooleanType) TryCast(value interface{}) (bool, bool) {
	return generic.ConvertToBoolean(value)
}

func (that *BooleanType) Cast(v interface{}, defaults bool) bool {
	if vv, ok := that.TryCast(v); ok {
		return vv
	}
	return defaults
}

type DurationType struct{}

func (that *DurationType) Get(ctx context.Context, getter Getter, name string, defVal time.Duration) (res time.Duration, err error) {
	return GetDurationProperty(ctx, getter, name, defVal)
}

func (that *DurationType) TryCast(value interface{}) (time.Duration, bool) {
	return generic.ConvertToDuration(value)
}

func (that *DurationType) Cast(v interface{}, defaults time.Duration) time.Duration {
	if vv, ok := that.TryCast(v); ok {
		return vv
	}
	return defaults
}

type JsonType struct {
}

func (that *JsonType) Get(ctx context.Context, getter Getter, name string, defVal json.RawMessage) (res json.RawMessage, err error) {
	return GetJsonProperty(ctx, getter, name, defVal)
}

func (that *JsonType) TryCast(value interface{}) (json.RawMessage, bool) {
	return generic.ConvertToJson(value)
}

func (that *JsonType) Cast(v interface{}, defaults json.RawMessage) json.RawMessage {
	if vv, ok := that.TryCast(v); ok {
		return vv
	}
	return defaults
}
