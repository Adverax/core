package types

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// GetBooleanProperty is helper for get boolean property from the getter
func GetBooleanProperty(
	ctx context.Context,
	getter Getter,
	name string,
	defVal bool,
) (res bool, err error) {
	val, err := getter.GetProperty(ctx, name)
	if err != nil {
		if errors.Is(err, GetErrNoMatch()) {
			return defVal, nil
		}
		return
	}
	if val == nil {
		return defVal, nil
	}
	res, ok := Type.Boolean.TryCast(val)
	if ok {
		return
	}
	return false, fmt.Errorf("can not convert value %v into boolean with key %q", val, name)
}

// GetStringProperty is helper for get string property from the getter
func GetStringProperty(
	ctx context.Context,
	getter Getter,
	name string,
	defVal string,
) (res string, err error) {
	val, err := getter.GetProperty(ctx, name)
	if err != nil {
		if errors.Is(err, GetErrNoMatch()) {
			return defVal, nil
		}
		return
	}
	if val == nil {
		return defVal, nil
	}
	res, ok := Type.String.TryCast(val)
	if ok {
		return
	}
	return "", fmt.Errorf("can not convert value %v into string with key %q", val, name)
}

// GetIntegerProperty is helper for get integer property from the getter
func GetIntegerProperty(
	ctx context.Context,
	getter Getter,
	name string,
	defVal int64,
) (res int64, err error) {
	val, err := getter.GetProperty(ctx, name)
	if err != nil {
		if errors.Is(err, GetErrNoMatch()) {
			return defVal, nil
		}
		return
	}
	if val == nil {
		return defVal, nil
	}
	res, ok := Type.Integer.TryCast(val)
	if ok {
		return
	}
	return 0, fmt.Errorf("can not convert value %v into integer with key %q", val, name)
}

// GetFloatProperty is helper for get float property from the getter
func GetFloatProperty(
	ctx context.Context,
	getter Getter,
	name string,
	defVal float64,
) (res float64, err error) {
	val, err := getter.GetProperty(ctx, name)
	if err != nil {
		if errors.Is(err, GetErrNoMatch()) {
			return defVal, nil
		}
		return
	}
	if val == nil {
		return defVal, nil
	}
	res, ok := Type.Float.TryCast(val)
	if ok {
		return
	}
	return 0, fmt.Errorf("can not convert value %v into float with key %q", val, name)
}

// GetDurationProperty is helper for get duration property from the getter
func GetDurationProperty(
	ctx context.Context,
	getter Getter,
	name string,
	defVal time.Duration,
) (res time.Duration, err error) {
	val, err := getter.GetProperty(ctx, name)
	if err != nil {
		if errors.Is(err, GetErrNoMatch()) {
			return defVal, nil
		}
		return
	}
	if val == nil {
		return defVal, nil
	}
	res, ok := Type.Duration.TryCast(val)
	if ok {
		return
	}
	return 0, fmt.Errorf("can not convert value %v into duration with key %q", val, name)
}

// GetJsonProperty is helper for get duration property from the getter
func GetJsonProperty(
	ctx context.Context,
	getter Getter,
	name string,
	defVal json.RawMessage,
) (res json.RawMessage, err error) {
	val, err := getter.GetProperty(ctx, name)
	if err != nil {
		if errors.Is(err, GetErrNoMatch()) {
			return defVal, nil
		}
		return
	}
	if val == nil {
		return defVal, nil
	}

	res, ok := Type.Json.TryCast(val)
	if ok {
		return res, nil
	}

	return nil, fmt.Errorf("can not convert value %v into duration with key %q", val, name)
}
