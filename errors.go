package core

import (
	"errors"
	"strings"
)

type Errors struct {
	errors []error
}

// NewErrors creates a new Errors object.
func NewErrors(es ...error) *Errors {
	res := &Errors{
		errors: make([]error, 0),
	}
	for _, e := range es {
		res.AddError(e)
	}
	return res
}

// Check is a helper function to check multiple errors.
func (that *Errors) Check(errs ...error) error {
	for _, err := range errs {
		if err != nil {
			that.AddError(err)
		}
	}
	return that.ResError()
}

// ResError returns an error if there are any errors present.
func (that *Errors) ResError() error {
	if that.IsEmpty() {
		return nil
	}

	return that
}

// Error returns a string representation of the errors.
func (that *Errors) Error() string {
	if len(that.errors) == 0 {
		return ""
	}

	errStrings := make([]string, 0)
	for _, e := range that.errors {
		errStrings = append(errStrings, e.Error())
	}
	return strings.Join(errStrings, "\n")
}

// AddError adds an error to the Errors object.
func (that *Errors) AddError(err error) {
	that.errors = append(that.errors, err)
}

// AddErrors adds multiple errors to the Errors object.
func (that *Errors) AddErrors(errs *Errors) {
	that.errors = append(that.errors, errs.errors...)
}

// IsEmpty returns true if there are no errors present.
func (that *Errors) IsEmpty() bool {
	return len(that.errors) == 0
}

// IsPresent returns true if there are errors present.
func (that *Errors) IsPresent() bool {
	return len(that.errors) > 0
}

// Count returns the number of errors.
func (that *Errors) Count() int {
	return len(that.errors)
}

// Contains returns true if the error is present in the Errors object.
func (that *Errors) Contains(err error) bool {
	for _, e := range that.errors {
		if errors.Is(e, err) {
			return true
		}
	}

	return false
}

// Result returns an error if there are any errors present.
func (that *Errors) Result() error {
	if that.IsPresent() {
		return that
	}

	return nil
}

// Check is a helper function to check multiple errors.
func Check(errs ...error) error {
	es := NewErrors()
	return es.Check(errs...)
}
