package core

import (
	"regexp"
	"strings"
)

const (
	MatchFilterExact MatchFilterType = iota
	MatchFilterRegexp
	MatchFilterPrefix
	MatchFilterSuffix
)

type MatchFilterType int

type MatchFilterDef struct {
	Text string
	Type MatchFilterType
}

type MatchFilter interface {
	IsMatch(text string) bool
}

type matchFilterMulti []MatchFilter

func (fs matchFilterMulti) IsMatch(text string) bool {
	for _, re := range fs {
		if re.IsMatch(text) {
			return true
		}
	}
	return false
}

type matchFilterRegexp struct {
	re *regexp.Regexp
}

func (f *matchFilterRegexp) IsMatch(text string) bool {
	return f.re.Match([]byte(text))
}

type matchFilterExact struct {
	text string
}

func (f *matchFilterExact) IsMatch(text string) bool {
	return f.text == text
}

type matchFilterPrefix struct {
	text string
}

func (f *matchFilterPrefix) IsMatch(text string) bool {
	return strings.HasPrefix(text, f.text)
}

type matchFilterSuffix struct {
	text string
}

type matchFilterConst struct {
	allow bool
}

func (f *matchFilterConst) IsMatch(text string) bool {
	return f.allow
}

func (f *matchFilterSuffix) IsMatch(text string) bool {
	return strings.HasPrefix(text, f.text)
}

type matchFilterAllowDeny struct {
	allow MatchFilter
	deny  MatchFilter
}

func (f *matchFilterAllowDeny) IsMatch(text string) bool {
	if f.deny != nil && f.deny.IsMatch(text) {
		return false
	}
	return f.allow == nil || f.allow.IsMatch(text)
}

var (
	AllowFilter MatchFilter = &matchFilterConst{allow: true}
	DenyFilter  MatchFilter = &matchFilterConst{allow: false}
)
