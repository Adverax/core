package filters

import (
	"fmt"
	"regexp"
)

type FilterBuilder struct {
	allow []*MatchFilterDef
	deny  []*MatchFilterDef
}

func (that *FilterBuilder) Allow(typ MatchFilterType, text string) *FilterBuilder {
	that.allow = append(that.allow, &MatchFilterDef{Text: text, Type: typ})
	return that
}

func (that *FilterBuilder) Deny(typ MatchFilterType, text string) *FilterBuilder {
	that.deny = append(that.deny, &MatchFilterDef{Text: text, Type: typ})
	return that
}

func (that *FilterBuilder) Build() (MatchFilter, error) {
	allow2, err := that.newMatchFilter(that.allow)
	if err != nil {
		return nil, fmt.Errorf("NewMatchFilter: %w", err)
	}

	deny2, err := that.newMatchFilter(that.deny)
	if err != nil {
		return nil, fmt.Errorf("NewMatchFilter: %w", err)
	}

	return &matchFilterAllowDeny{
		allow: allow2,
		deny:  deny2,
	}, nil
}

func (that *FilterBuilder) newMatchFilter(
	defs []*MatchFilterDef,
) (filter MatchFilter, err error) {
	if len(defs) == 0 {
		return nil, nil
	}

	filters := make(matchFilterMulti, 0)
	for _, def := range defs {
		var filter MatchFilter
		switch def.Type {
		case MatchFilterRegexp:
			re, err := regexp.Compile(def.Text)
			if err != nil {
				return nil, fmt.Errorf("regexp.Compile: %w", err)
			}
			filter = &matchFilterRegexp{re: re}
		case MatchFilterPrefix:
			filter = &matchFilterPrefix{text: def.Text}
		case MatchFilterSuffix:
			filter = &matchFilterSuffix{text: def.Text}
		default:
			filter = &matchFilterExact{text: def.Text}
		}
		filters = append(filters, filter)
	}

	return filters, nil
}
