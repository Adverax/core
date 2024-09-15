package pubsub

import (
	"github.com/adverax/core/json"
)

func parse[T any](data []byte) (T, error) {
	var result T
	err := json.Unmarshal(data, &result)
	return result, err
}
