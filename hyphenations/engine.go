package hyphenations

import (
	"regexp"
	"strings"
)

// Библиотека автоматической расстановки переносов
// * Русский - https://github.com/kozachenko/jQuery-Russian-Hyphenation
// * Английский -https://github.com/bramstein/hypher
// Пример подключены оба плагина - http://jsfiddle.net/yoowordpress/e5cormh3/1/

const (
	statusInvalid = iota
	statusValid
	statusRedo
	statusSuccess
)

type state struct {
	s int
	p int
	l int
}

type Engine struct {
	hyphen rune
	rules  []*regexp.Regexp
}

// Hyphenate returns text with all available hyphens.
// Hyph is a string that will be inserted between parts of the word.
func (that *Engine) Hyphenate(text string, hyph string) string {
	for _, r := range that.rules {
		text = r.ReplaceAllString(text, "$1"+hyph+"$2")
	}
	return text
}

// Split returns slice of lines, where length of each of them is less or equal param width.
func (that *Engine) Split(
	text string,
	width int,
) (lines []string) {
	raw := that.Hyphenate(text, string(that.hyphen))
	hyph := []rune(raw)
	cur := state{s: 0, p: 0, l: 0}
	var valid = cur
	for len(hyph) > 0 {
		st, status := that.resolve(hyph, width, cur)
		if status == statusValid || status == statusInvalid {
			cur = st
			if status == statusValid {
				valid = st
			}
		} else {
			if status == statusSuccess {
				valid = state{s: 0, p: len(hyph), l: cur.l}
			}
			line := hyph[:valid.p]
			line = []rune(strings.TrimSpace(string(line)))
			var ln = len(line)
			if len(line) == 0 {
				break
			}
			if line[ln-1] == that.hyphen {
				line = append(line[:ln-1], '-')
			}
			line2 := strings.ReplaceAll(string(line), string(that.hyphen), "")
			lines = append(lines, line2)
			if status == statusSuccess {
				break
			}
			hyph = hyph[valid.p:]
			cur = state{s: 0, p: 0, l: 0}
			valid = cur
		}
	}
	return lines
}

func (that *Engine) resolve(
	hyph []rune,
	width int,
	cur state,
) (next state, status int) {
	if cur.l > width {
		// Достигнута граничная ширина
		status = statusRedo
		return
	}

	if len(hyph) == cur.p {
		// Весь текст обработан
		status = statusSuccess
		return
	}

	var ch = hyph[cur.p]
	if ch == that.hyphen {
		if cur.l == width {
			status = statusRedo
			return
		}
	} else {
		cur.l++
	}

	if ch == ' ' || ch == that.hyphen {
		status = statusValid
	}

	return state{
		s: cur.s,
		p: cur.p + 1,
		l: cur.l,
	}, status
}
