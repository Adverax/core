package jsonFetcher

import (
	"bytes"
	"encoding/json"
)

type Fetcher interface {
	Fetch() ([]byte, error)
}

type Engine struct {
	fetcher Fetcher
}

func New(fetcher Fetcher) *Engine {
	return &Engine{
		fetcher: fetcher,
	}
}

func (that *Engine) Fetch() (map[string]interface{}, error) {
	data := make(map[string]interface{})

	source, err := that.fetcher.Fetch()
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(bytes.NewBuffer(source))
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
