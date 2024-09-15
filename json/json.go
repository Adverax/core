package json

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adverax/core"
	"github.com/adverax/core/generic"
	jsoniter "github.com/json-iterator/go"
	"io"
)

const (
	TypeUnknown Type = 0
	TypeArray   Type = 1
	TypeObject  Type = 2
)

type Type int

// Config is a custom JSON configuration.
var Config = jsoniter.Config{
	EscapeHTML:             false,
	SortMapKeys:            false,
	ValidateJsonRawMessage: false,
	UseNumber:              true,
}.Froze()

// ConfigSorted is a custom JSON configuration.
var ConfigSorted = jsoniter.Config{
	EscapeHTML:             false,
	SortMapKeys:            true,
	ValidateJsonRawMessage: false,
	UseNumber:              true,
}.Froze()

// Unmarshal is a helper function to unmarshal a JSON document.
func Unmarshal(raw []byte, value interface{}) error {
	err := Config.Unmarshal(raw, &value)
	if err != nil && raw != nil {
		return fmt.Errorf("Config.Unmarshal: %w", err)
	}
	return nil
}

// Marshal is a helper function to marshal a JSON document.
func Marshal(value interface{}) ([]byte, error) {
	return Config.Marshal(value)
}

// MarshalIndent is a helper function to marshal a JSON document with indentation.
func MarshalIndent(value interface{}) ([]byte, error) {
	return json.MarshalIndent(value, "", "  ")
}

// Update is a helper function to update a JSON document.
// It takes a document and a function that takes a Map and returns an error.
// Example: Update(doc, func(m Map) error { m["key"] = "value"; return nil })
func Update(
	doc []byte,
	actions ...func(Map) error,
) ([]byte, error) {
	m, err := NewMap(doc)
	if err != nil {
		return nil, fmt.Errorf("NewMapFromJson: %w", err)
	}

	for _, action := range actions {
		err = action(m)
		if err != nil {
			return nil, fmt.Errorf("action: %w", err)
		}
	}

	data, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("Marshal: %w", err)
	}
	return data, nil
}

// UpdateAll is a helper function to update a list of JSON documents.
func UpdateAll(
	raw []byte,
	action func(Map) error,
) ([]byte, error) {
	var rows []RawMessage
	err := Unmarshal(raw, &rows)
	if err != nil {
		return nil, fmt.Errorf("Unmarshal: %w", err)
	}

	for i, row := range rows {
		rows[i], err = Update(row, action)
		if err != nil {
			return nil, fmt.Errorf("Update: %w", err)
		}
	}

	return Marshal(rows)
}

// UpdateValues is a helper function to update a JSON document with a map of values.
// It takes a document and a map of fields to update.
// Example: UpdateValues(doc, map[string]interface{}{"key": "value"})
func UpdateValues(
	doc []byte,
	fields map[string]interface{},
) ([]byte, error) {
	return Update(doc, func(m Map) error {
		m.ExpandBy(NewMapFromStruct(fields))
		return nil
	})
}

// Get is a helper function to extract a value from a JSON document.
// Example: Get(doc, GetInteger("key", 0))
func Get[T any](
	doc []byte,
	getter func(Map) (T, error),
) (val T, err error) {
	m, err := NewMap(doc)
	if err != nil {
		return val, fmt.Errorf("NewMapFromJson: %w", err)
	}

	val, err = getter(m)
	if err != nil {
		return val, fmt.Errorf("getter: %w", err)
	}

	return val, nil
}

// Set is a helper function to update a JSON document.
// It takes a document and a list of setter functions.
// Example: Set(doc, Override("key", "value"), JsonDefault("key2", "value2"))
func Set(
	doc []byte,
	setter ...func(doc Map) error,
) ([]byte, error) {
	res, err := Update(
		doc,
		func(recData Map) error {
			for _, set := range setter {
				if err := set(recData); err != nil {
					return fmt.Errorf("setter: %w", err)
				}
			}
			return nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("JsonUpdate: %w", err)
	}
	return res, nil
}

// Override is a helper function to override a value in a JSON document.
// See Set for more information.
func Override(key string, val interface{}) func(Map) error {
	return func(doc Map) error {
		return doc.SetProperty(context.Background(), key, val)
	}
}

// Default is a helper function to set a default value in a JSON document.
// See Set for more information.
func Default(key string, val interface{}) func(Map) error {
	return func(doc Map) error {
		value, err := doc.GetProperty(context.Background(), key)
		if err != nil {
			if errors.Is(err, core.ErrNoMatch) {
				doc[key] = val
				return nil
			}
			return fmt.Errorf("GetProperty: %w", err)
		}

		if generic.IsZeroVal(value) {
			return doc.SetProperty(context.Background(), key, val)
		}
		return nil
	}
}

// If is a helper function to conditionally execute an action.
// See Set for more information.
func If(cond bool, action func(Map) error) func(Map) error {
	if cond {
		return action
	}
	return dummyAction
}

// GetInteger is a helper function to extract an integer from a JSON document.
// See Get for more information.
func GetInteger(key string, defVal int64) func(Map) (int64, error) {
	return func(doc Map) (int64, error) {
		return doc.GetInteger(context.Background(), key, defVal)
	}
}

// GetFloat is a helper function to extract a float from a JSON document.
// See JsonGet for more information.
func GetFloat(key string, defVal float64) func(Map) (float64, error) {
	return func(doc Map) (float64, error) {
		return doc.GetFloat(context.Background(), key, defVal)
	}
}

// GetString is a helper function to extract a string from a JSON document.
// See Get for more information.
func GetString(key string, defVal string) func(Map) (string, error) {
	return func(doc Map) (string, error) {
		return doc.GetString(context.Background(), key, defVal)
	}
}

// GetBoolean is a helper function to extract a boolean from a JSON document.
// See Get for more information.
func GetBoolean(key string, defVal bool) func(Map) (bool, error) {
	return func(doc Map) (bool, error) {
		return doc.GetBoolean(context.Background(), key, defVal)
	}
}

// GetAny is a helper function to extract a any value from a JSON document.
// See Get for more information.
func GetAny(key string) func(Map) (any, error) {
	return func(doc Map) (any, error) {
		return doc.GetProperty(context.Background(), key)
	}
}

// Remove is a helper function to remove a value from a JSON document.
func Remove(key string) func(Map) error {
	return func(doc Map) error {
		delete(doc, key)
		return nil
	}
}

// Ensure is a helper function to ensure a JSON document is not empty.
// It takes a document and returns an empty object if the document is empty.
func Ensure(doc []byte) []byte {
	if IsEmpty(doc) {
		return Empty
	}
	return doc
}

// TypeOf is a helper function to determine the type of a JSON document.
// It takes a document and returns the type (array, object).
// Example: TypeOf(doc)
func TypeOf(in io.Reader) (Type, error) {
	dec := json.NewDecoder(in)
	// Get just the first valid JSON token from input
	t, err := dec.Token()
	if err != nil {
		return TypeUnknown, err
	}
	if d, ok := t.(json.Delim); ok {
		// The first token is a delimiter, so this is an array or an object
		switch d {
		case '[':
			return TypeArray, nil
		case '{':
			return TypeObject, nil
		default: // ] or }
			return TypeUnknown, errors.New("Unexpected delimiter")
		}
	}
	return TypeUnknown, errors.New("Input does not represent a JSON object or array")
}

// AsArray is a helper function to ensure a JSON document is an array.
func AsArray(data []byte) ([]byte, error) {
	isArray, err := IsArray(data)
	if err != nil {
		return nil, fmt.Errorf("JsonIsArray: %w", err)
	}
	if isArray {
		return data, nil
	}

	return []byte(fmt.Sprintf("[%s]", string(data))), nil
}

// IsArray is a helper function to determine if a JSON document is an array.
func IsArray(raw []byte) (bool, error) {
	if len(raw) == 0 {
		return false, nil
	}
	tp, err := TypeOf(bytes.NewBuffer(raw))
	if err != nil {
		return false, fmt.Errorf("JsonType: %w", err)
	}
	return tp == TypeArray, nil
}

// IsObject is a helper function to determine if a JSON document is an object.
func IsObject(raw []byte) (bool, error) {
	if len(raw) == 0 {
		return false, nil
	}
	tp, err := TypeOf(bytes.NewBuffer(raw))
	if err != nil {
		return false, fmt.Errorf("JsonType: %w", err)
	}
	return tp == TypeObject, nil
}

// Empty is empty document
var Empty = []byte("{}")

func IsEmpty(doc []byte) bool {
	return len(doc) <= 2
}

// CoalesceString is a helper function to coalesce a string from multiple documents.
func CoalesceString(ctx context.Context, path string, defVal string, docs ...Map) (string, error) {
	for _, doc := range docs {
		if doc == nil {
			continue
		}
		val, err := doc.GetString(ctx, path, "")
		if err != nil {
			return val, fmt.Errorf("GetString: %w", err)
		}
	}

	return defVal, nil
}

// RestoreString is a helper function to restore a string from a document.
func RestoreString(ctx context.Context, path string, docs ...Map) error {
	val, err := CoalesceString(ctx, path, "", docs...)
	if err != nil {
		return fmt.Errorf("JsonCoalesceString: %w", err)
	}

	if val != "" {
		err = docs[0].SetString(ctx, path, val)
		if err != nil {
			return fmt.Errorf("SetString: %w", err)
		}
	}

	return nil
}

func Normalize(data []byte) ([]byte, error) {
	m, err := NewMap(data)
	if err != nil {
		return nil, err
	}

	return MarshalIndent(m)
}

func IsEqual(a, b []byte) bool {
	aa, err := Normalize(a)
	if err != nil {
		return false
	}

	bb, err := Normalize(b)
	if err != nil {
		return false
	}

	return string(aa) == string(bb)
}

// Merge is a helper function to merge JSON documents.
func Merge(values ...[]byte) ([]byte, error) {
	res := Map{}
	for i, raw := range values {
		m, err := NewMap(raw)
		if err != nil {
			return nil, fmt.Errorf("NewMap[%d]: %w", i, err)
		}
		if i == 0 {
			res = m
			continue
		}
		res.ExpandBy(m)
	}
	return res.Json()
}

var (
	dummyAction = func(Map) error { return nil }
)
