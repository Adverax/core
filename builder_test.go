package core

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

type myObject struct {
	name string
	data map[string]interface{}
}

type myObjectBuilder struct {
	*Builder
	obj *myObject
}

func NewMyObjectBuilder() *myObjectBuilder {
	return &myObjectBuilder{
		Builder: NewBuilder("myObject"),
		obj:     &myObject{},
	}
}

func (that *myObjectBuilder) Name(name string) *myObjectBuilder {
	that.obj.name = name
	return that
}

func (that *myObjectBuilder) Data(data map[string]interface{}) *myObjectBuilder {
	that.obj.data = data
	return that
}

func (that *myObjectBuilder) Build() (*myObject, error) {
	if err := that.checkRequiredFields(that.obj); err != nil {
		return nil, err
	}

	return that.obj, nil
}

func (that *myObjectBuilder) checkRequiredFields(obj *myObject) error {
	that.RequiredField(obj.name, errFieldNameRequired)
	that.RequiredField(obj.data, errFieldDataRequired)
	return that.ResError()
}

var (
	errFieldNameRequired = errors.New("name is required")
	errFieldDataRequired = errors.New("data is required")
)

func TestBuilder(t *testing.T) {
	builder := NewMyObjectBuilder()
	obj, err := builder.Build()
	require.Nil(t, obj)
	require.NotNil(t, err)
	assert.True(t, builder.Errors().Contains(errFieldNameRequired))
	assert.True(t, builder.Errors().Contains(errFieldDataRequired))
}
