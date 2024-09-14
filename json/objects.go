package json

import (
	"adverax/core/types"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type AtomGetter interface {
	Get(ctx context.Context) (interface{}, error)
}

type AtomSetter interface {
	Set(ctx context.Context, value interface{}) error
}

// GetPropertyEx is a helper function to extract a value from an object
// It is wrapper for JsonGetProperty with panic handling.
// Example: GetPropertyEx(ctx, object, "key")
func GetPropertyEx(
	ctx context.Context,
	object interface{},
	key string,
) (val interface{}, err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = fmt.Errorf("jsonGetProperty: %v", e)
		}
	}()
	return GetProperty(ctx, object, key)
}

// SetPropertyEx is a helper function to set a value in an object.
// It is wrapper for JsonSetProperty with panic handling.
// Example: SetPropertyEx(ctx, object, "key", "value")
func SetPropertyEx(
	ctx context.Context,
	object interface{},
	key string,
	val interface{},
) (err error) {
	defer func() {
		e := recover()
		if e != nil {
			err = fmt.Errorf("jsonSetProperty: %v", e)
		}
	}()
	return SetProperty(ctx, object, key, val)
}

// GetProperty is a helper function to extract a value from an object.
// It takes an object and a key and returns the value. It can panic.
// Example: GetProperty(ctx, object, "key")
func GetProperty(
	ctx context.Context,
	object interface{},
	key string,
) (interface{}, error) {
	if strings.HasPrefix(key, "$.") {
		keys := strings.Split(key[2:], ".")
		for _, k := range keys {
			var err error
			object, err = getProperty(ctx, object, k)
			if err != nil {
				return nil, fmt.Errorf("GetProperty %q: %w", key, err)
			}
		}

		return object, nil
	}

	var err error
	object, err = getProperty(ctx, object, key)
	if err != nil {
		return nil, fmt.Errorf("GetProperty %q: %w", key, err)
	}
	return object, nil
}

func getProperty(
	ctx context.Context,
	object interface{},
	key string,
) (interface{}, error) {
	if getter, ok := object.(types.Getter); ok {
		val, err := getter.GetProperty(ctx, key)
		if err == nil {
			return val, nil
		}

		if !errors.Is(err, types.GetErrNoMatch()) {
			return nil, fmt.Errorf("GetParam: %w", err)
		}
	}

	v := reflect.ValueOf(object).Elem()
	// It's possible we can cache this, which is why precompute all these ahead of time.
	findJsonName := func(t reflect.StructTag) (string, error) {
		if jt, ok := t.Lookup("json"); ok {
			return strings.Split(jt, ",")[0], nil
		}
		return "", fmt.Errorf("tag provided does not define a json tag %s", key)
	}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag
		jname, _ := findJsonName(tag)
		if jname == key {
			val := v.Field(i).Interface()
			if getter, ok := val.(AtomGetter); ok {
				val, err := getter.Get(ctx)
				if err != nil {
					return nil, fmt.Errorf("Get atom %q: %w", key, err)
				}
				return val, nil
			}
			if v.Kind() == reflect.Struct {
				return v.Field(i).Addr().Interface(), nil
			}
			return val, nil
		}
	}

	return nil, types.GetErrNoMatch()
}

// SetProperty is a helper function to set a value in an object.
// It takes an object, a key and a value and sets the value. It can panic.
// Example: SetProperty(ctx, object, "key", "value")
func SetProperty(
	ctx context.Context,
	object interface{},
	key string,
	value interface{},
) error {
	if strings.HasPrefix(key, "$.") {
		keys := strings.Split(key[2:], ".")
		for i, k := range keys {
			if i == len(keys)-1 {
				return setProperty(ctx, object, k, value)
			}
			var err error
			object, err = getProperty(ctx, object, k)
			if err != nil {
				return fmt.Errorf("GetProperty %q: %w", key, err)
			}
		}

		return nil
	}

	return setProperty(ctx, object, key, value)
}

func setProperty(
	ctx context.Context,
	object interface{},
	key string,
	value interface{},
) error {
	if setter, ok := object.(types.Setter); ok {
		err := setter.SetProperty(ctx, key, value)
		if err == nil {
			return nil
		}

		if !errors.Is(err, types.GetErrNoMatch()) {
			return fmt.Errorf("GetParam: %w", err)
		}
	}

	v := reflect.ValueOf(object).Elem()
	if !v.CanAddr() {
		return fmt.Errorf("cannot assign to the item passed, item must be a pointer in order to assign")
	}
	// It's possible we can cache this, which is why precompute all these ahead of time.
	findJsonName := func(t reflect.StructTag) (string, error) {
		if jt, ok := t.Lookup("json"); ok {
			return strings.Split(jt, ",")[0], nil
		}
		return "", fmt.Errorf("tag provided does not define a json tag, %s", key)
	}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag
		jname, _ := findJsonName(tag)
		if jname == key {
			field := v.Field(i)
			val := field.Interface()
			if setter, ok := val.(AtomSetter); ok {
				err := setter.Set(ctx, value)
				if err != nil {
					return fmt.Errorf("set atom %q: %v", key, value)
				}
				return nil
			}
			field.Set(reflect.ValueOf(value).Convert(field.Type()))
			return nil
		}
	}

	return types.GetErrNoMatch()
}

// ImportProperties is a helper function to import properties into an object.
// It takes an object and a getter and sets the properties.
// Example: ImportProperties(ctx, object, getter)
func ImportProperties(
	ctx context.Context,
	object interface{},
	getter types.Getter,
) error {
	v := reflect.ValueOf(object).Elem()
	// It's possible we can cache this, which is why precompute all these ahead of time.
	findJsonName := func(key string, t reflect.StructTag) (string, error) {
		if jt, ok := t.Lookup("json"); ok {
			return strings.Split(jt, ",")[0], nil
		}
		return "", fmt.Errorf("tag provided does not define a json tag, %s", key)
	}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag
		key := typeField.Name
		jname, _ := findJsonName(key, tag)
		if jname != "" && jname != "-" {
			val, err := getter.GetProperty(ctx, jname)
			if val != nil && err == nil {
				err := SetPropertyEx(ctx, object, jname, val)
				if err != nil && !errors.Is(err, types.GetErrNoMatch()) {
					return fmt.Errorf("JsonSetProperty %q: %w", jname, err)
				}
			}
		}
	}

	return nil
}

// EnumProperties is a helper function to export properties from an object.
// It takes an object and returns a map of properties.
// Example: ExportProperties(ctx, object)
func EnumProperties(
	object interface{},
) (list []string, err error) {
	v := reflect.ValueOf(object).Elem()
	// It's possible we can cache this, which is why precompute all these ahead of time.
	findJsonName := func(key string, t reflect.StructTag) (string, error) {
		if jt, ok := t.Lookup("json"); ok {
			return strings.Split(jt, ",")[0], nil
		}
		return "", fmt.Errorf("tag provided does not define a json tag, %s", key)
	}
	for i := 0; i < v.NumField(); i++ {
		typeField := v.Type().Field(i)
		tag := typeField.Tag
		key := typeField.Name
		jname, _ := findJsonName(key, tag)
		if jname != "" && jname != "-" {
			list = append(list, jname)
		}
	}

	return list, nil
}
