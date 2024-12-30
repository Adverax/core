package core

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type color int

func (that color) String() string {
	return colors.DecodeOrDefault(that, "unknown")
}

func (that *color) UnmarshalText(text []byte) error {
	return colors.UnmarshalText(text, that)
}

func (that color) MarshalText() ([]byte, error) {
	return colors.MarshalText(that)
}

const (
	red color = iota
	green
	blue
)

var colors = NewEnum(
	map[color]string{
		red:   "red",
		green: "green",
		blue:  "blue",
	},
)

type EnumShould struct {
	suite.Suite
}

func TestEnumShould(t *testing.T) {
	suite.Run(t, new(EnumShould))
}

func (that *EnumShould) TestEncode() {
	clr, err := colors.Encode("green")
	that.NoError(err)
	that.Equal(green, clr)
}

func (that *EnumShould) TestDecode() {
	clr, err := colors.Decode(green)
	that.NoError(err)
	that.Equal("green", clr)
}

func (that *EnumShould) TestValues() {
	that.Equal(3, len(colors.Values()))
}
