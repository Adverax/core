package core

import (
	"fmt"
)

// Builder is a base parent for building objects
type Builder struct {
	errors    Errors
	component string
}

// ResError returns error list
func (that *Builder) ResError() error {
	return that.errors.ResError()
}

// Errors returns error list
func (that *Builder) Errors() *Errors {
	return &that.errors
}

// RequiredField checks if field is not zero value
func (that *Builder) RequiredField(field interface{}, err error) {
	if IsZeroValue(field) {
		that.AddError(err)
	}
}

// AddError adds error
func (that *Builder) AddError(err error) {
	that.errors.AddError(fmt.Errorf("component %s: %w", that.component, err))
}

// NewBuilder creates new builder
func NewBuilder(component string) *Builder {
	return &Builder{
		component: component,
	}
}
