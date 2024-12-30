package core

import (
	"encoding/json"
	"fmt"
)

type Enum[T comparable] struct {
	encoders map[T]string
	decoders map[string]T
}

func (that *Enum[T]) Encode(val string) (res T, err error) {
	if v, ok := that.decoders[val]; ok {
		return v, nil
	}
	return res, fmt.Errorf("unknown value %s of enum", val)
}

func (that *Enum[T]) Decode(val T) (string, error) {
	if v, ok := that.encoders[val]; ok {
		return v, nil
	}
	return "", fmt.Errorf("not a valid value of enum %v", val)
}

func (that *Enum[T]) EncodeOrDefault(val string, def T) T {
	if v, ok := that.decoders[val]; ok {
		return v
	}
	return def
}

func (that *Enum[T]) DecodeOrDefault(val T, def string) string {
	if v, ok := that.encoders[val]; ok {
		return v
	}
	return def
}

func (that *Enum[T]) Values() []T {
	res := make([]T, 0, len(that.encoders))
	for k := range that.encoders {
		res = append(res, k)
	}
	return res
}

func (that *Enum[T]) UnmarshalText(text []byte, val *T) error {
	v, err := that.Encode(string(text))
	if err != nil {
		return err
	}

	*val = v
	return nil
}

func (that *Enum[T]) MarshalText(val T) ([]byte, error) {
	s, err := that.Decode(val)
	if err != nil {
		return nil, err
	}

	return []byte(s), nil
}

func (that *Enum[T]) MarshalJSON(val T) ([]byte, error) {
	s, err := that.Decode(val)
	if err != nil {
		return nil, err
	}

	return json.Marshal(s)
}

func (that *Enum[T]) UnmarshalJSON(data []byte, val *T) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}

	v, err := that.Encode(s)
	if err != nil {
		return err
	}

	*val = v
	return nil
}

func NewEnum[T comparable](encoders map[T]string) *Enum[T] {
	decoders := make(map[string]T, len(encoders))
	for k, v := range encoders {
		decoders[v] = k
	}

	return &Enum[T]{
		encoders: encoders,
		decoders: decoders,
	}
}
