package core

import (
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type myObject struct {
	host  string `build:",required"`
	name  string `build:",required"`
	alias string `build:",required"`
}

type myBuilder struct {
	*Builder
	object *myObject
}

func newMyBuilder() *myBuilder {
	const defaultAlias = "default"

	return &myBuilder{
		Builder: NewBuilder("component"),
		object: &myObject{
			alias: defaultAlias,
		},
	}
}

func (that *myBuilder) WithHost(host string) *myBuilder {
	that.object.host = host
	return that
}

func (that *myBuilder) WithName(name string) *myBuilder {
	that.object.name = name
	return that
}

func (that *myBuilder) Build() (*myObject, error) {
	if err := that.checkRequiredFields(); err != nil {
		return nil, err
	}
	return that.object, nil
}

func (that *myBuilder) checkRequiredFields() error {
	that.RequiredField(that.object.name, errRequiredFieldName)
	that.RequiredField(that.object.host, errRequiredFieldHost)
	that.RequiredField(that.object.alias, errRequiredFieldAlias)
	return that.ResError()
}

var (
	errRequiredFieldName  = errors.New("required field name is required")
	errRequiredFieldHost  = errors.New("required field host is required")
	errRequiredFieldAlias = errors.New("required field alias is required")
)

type BuilderShould struct {
	suite.Suite
}

func TestBuilderShould(t *testing.T) {
	suite.Run(t, new(BuilderShould))
}

func (that *BuilderShould) TestBuildWithoutRequiredFields_mustBeFailed() {
	b := newMyBuilder()
	object, err := b.Build()
	that.Nil(object)
	that.NotNil(err)
	errs := b.Errors()
	that.NotNil(errs)
	that.True(errs.Contains(errRequiredFieldHost))
	that.True(errs.Contains(errRequiredFieldName))
}

func (that *BuilderShould) TestBuildWithRequiredFields_mustBeSucceeded() {
	b := newMyBuilder()
	object, err := b.
		WithName("component").
		WithHost("host").
		Build()
	that.Nil(err)
	that.NotNil(object)
	that.False(b.Errors().Contains(errRequiredFieldHost))
}
