package json

import (
	"encoding/json"
	"errors"
	"github.com/adverax/core/types"
)

type Number = json.Number

type RawMessage = json.RawMessage

// Boolean is a custom type for JSON booleans.
type Boolean bool

func (b *Boolean) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"true"`, `true`, `"1"`, `1`:
		*b = true
		return nil
	case `"false"`, `false`, `"0"`, `0`, `""`:
		*b = false
		return nil
	default:
		return errors.New("CustomBool: parsing \"" + string(data) + "\": unknown value")
	}
}

func (b *Boolean) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(*b)
	return data, err
}

// String is a custom type for JSON strings.
type String string

func (s *String) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `true`:
		*s = "true"
		return nil
	case `false`:
		*s = "false"
		return nil
	default:
		if val, ok := types.Type.String.TryCast(json.RawMessage(data)); ok {
			*s = String(val)
			return nil
		}
		return errors.New("CustomBool: parsing \"" + string(data) + "\": unknown value")
	}
}

func (s *String) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(*s)
	return data, err
}

// Logical is a custom type for JSON.
type Logical int

func (b *Logical) UnmarshalJSON(data []byte) error {
	switch string(data) {
	case `"true"`, `true`, `"1"`, `1`:
		*b = 1
		return nil
	case `"false"`, `false`, `"0"`, `0`:
		*b = -1
		return nil
	default:
		*b = 0
		return nil
	}
}

func (b *Logical) MarshalJSON() ([]byte, error) {
	data, err := json.Marshal(*b)
	return data, err
}
