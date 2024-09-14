package core

import (
	uuid "github.com/satori/go.uuid"
)

// NewGUID - generate new guid
func NewGUID() string {
	id := uuid.NewV4()
	return id.String()
}
