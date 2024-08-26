package core

import (
	"errors"
	"fmt"
)

// Builder is a base parent for building objects
type Builder struct {
	errors    *Errors
	component string
}

// ResError returns error list
func (that *Builder) ResError() error {
	if that.errors == nil {
		return nil
	}
	return that.errors.ResError()
}

// Errors returns error list
func (that *Builder) Errors() *Errors {
	that.EnsureErrors()
	return that.errors
}

// EnsureErrors ensures that errors list is initialized
func (that *Builder) EnsureErrors() {
	if that.errors == nil {
		that.errors = NewErrors()
	}
}

// RequiredField checks if field is not zero value
func (that *Builder) RequiredField(field interface{}, err error) {
	if IsZeroVal(field) {
		that.AddError(err)
	}
}

// AddError adds error
func (that *Builder) AddError(err error) {
	that.EnsureErrors()
	that.errors.AddError(fmt.Errorf("component %s: %w", that.component, err))
}

// AddErrorf adds error with formatted message
func (that *Builder) AddErrorf(format string, args ...interface{}) {
	err := errors.New(fmt.Sprintf(format, args...))
	that.AddError(err)
}

// NewBuilder creates new builder for slice of objects
func NewBuilder(component string) *Builder {
	return &Builder{
		component: component,
	}
}
