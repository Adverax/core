package hyphenations

import (
	"github.com/adverax/core"
	"regexp"
)

type Builder struct {
	*core.Builder
	engine  *Engine
	charA   string
	charV   string
	charN   string
	charX   string
	charAen string
	charVen string
	charNen string
	charXen string
}

func NewBuilder() *Builder {
	return &Builder{
		Builder: core.NewBuilder("hyphenations"),
		charA:   "[абвгдеёжзийклмнопрстуфхцчшщъыьэюяґєії]",
		charV:   "[аеёиоуыэюяєії]",
		charN:   "[бвгджзклмнпрстфхцчшщґ]",
		charX:   "[йъь]",
		charAen: "[abcdefghijklmnopqrstuvwxyz]",
		charVen: "[aeiou]",
		charNen: "[bcdfghjklmnpqrstvwxyz]",
		charXen: "[wyz]",
		engine: &Engine{
			hyphen: '\xAD',
		},
	}
}

func (that *Builder) CharA(charA string) *Builder {
	that.charA = charA
	return that
}

func (that *Builder) CharV(charV string) *Builder {
	that.charV = charV
	return that
}

func (that *Builder) CharN(charN string) *Builder {
	that.charN = charN
	return that
}

func (that *Builder) CharX(charX string) *Builder {
	that.charX = charX
	return that
}

func (that *Builder) CharAen(charAen string) *Builder {
	that.charAen = charAen
	return that
}

func (that *Builder) CharVen(charVen string) *Builder {
	that.charVen = charVen
	return that
}

func (that *Builder) CharNen(charNen string) *Builder {
	that.charNen = charNen
	return that
}

func (that *Builder) CharXen(charXen string) *Builder {
	that.charXen = charXen
	return that
}

func (that *Builder) Hyphen(hyphen rune) *Builder {
	that.engine.hyphen = hyphen
	return that
}

func (that *Builder) Build() (*Engine, error) {
	const flags = "(?i)"

	that.engine.rules = []*regexp.Regexp{
		regexp.MustCompile(flags + "(" + that.charX + ")(" + that.charA + that.charA + ")"),
		regexp.MustCompile(flags + "(" + that.charV + ")(" + that.charV + that.charA + ")"),
		regexp.MustCompile(flags + "(" + that.charV + that.charN + ")(" + that.charN + that.charV + ")"),
		regexp.MustCompile(flags + "(" + that.charN + that.charV + ")(" + that.charN + that.charV + ")"),
		regexp.MustCompile(flags + "(" + that.charV + that.charN + ")(" + that.charN + that.charN + that.charV + ")"),
		regexp.MustCompile(flags + "(" + that.charV + that.charN + that.charN + ")(" + that.charN + that.charN + that.charV + ")"),

		regexp.MustCompile(flags + "(" + that.charXen + ")(" + that.charAen + that.charAen + ")"),
		regexp.MustCompile(flags + "(" + that.charVen + ")(" + that.charVen + that.charAen + ")"),
		regexp.MustCompile(flags + "(" + that.charVen + that.charNen + ")(" + that.charNen + that.charVen + ")"),
		regexp.MustCompile(flags + "(" + that.charNen + that.charVen + ")(" + that.charNen + that.charVen + ")"),
		regexp.MustCompile(flags + "(" + that.charVen + that.charNen + ")(" + that.charNen + that.charNen + that.charVen + ")"),
		regexp.MustCompile(flags + "(" + that.charVen + that.charNen + that.charNen + ")(" + that.charNen + that.charNen + that.charVen + ")"),
	}

	return that.engine, nil
}
