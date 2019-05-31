package versionfilter

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type BooleanFilter struct {
	Operator   string
	Subfilters []VersionFilter
}

const OP_AND = "&&"
const OP_OR = "||"

func NewBooleanFilter(s string) (*BooleanFilter, error) {
	andParts := strings.Split(s, OP_AND)
	orParts := strings.Split(s, OP_OR)

	if len(andParts) == 1 && len(orParts) == 1 {
		return nil, fmt.Errorf("does not contain \"%s\" or \"%s\"", OP_AND, OP_OR)
	}

	if len(andParts) > 1 && len(orParts) > 1 {
		return nil, fmt.Errorf("cannot contain both \"%s\" and \"%s\"", OP_AND, OP_OR)
	}

	f := &BooleanFilter{}

	if len(andParts) > 1 {
		f.Operator = OP_AND
		for _, part := range andParts {
			subfilter := NewVersionFilter(part)
			if subfilter != nil {
				f.Subfilters = append(f.Subfilters, subfilter)
			} else {
				return nil, fmt.Errorf("invalid filter %s", part)
			}
		}
	}

	if len(orParts) > 1 {
		f.Operator = OP_OR
		for _, part := range orParts {
			subfilter := NewVersionFilter(part)
			if subfilter != nil {
				f.Subfilters = append(f.Subfilters, subfilter)
			} else {
				return nil, fmt.Errorf("invalid filter %s", part)
			}
		}
	}

	return f, nil
}

func (f *BooleanFilter) Matching(versions []string, currentVersion string) ([]string, error) {
	matchesSet := map[string]int{}

	for _, sub := range f.Subfilters {
		submatches, err := sub.Matching(versions, currentVersion)
		if err != nil {
			return nil, err
		}
		for _, sm := range submatches {
			matchesSet[sm] = matchesSet[sm] + 1
		}
	}

	matches := []string{}
	for k, v := range matchesSet {
		if f.Operator == OP_AND {
			// intersection (item appeared in every subfilter)
			if v == len(f.Subfilters) {
				matches = append(matches, k)
			}
		} else if f.Operator == OP_OR {
			// union (item appeared in any subfilter)
			matches = append(matches, k)
		} else {
			return nil, errors.New("unknown operator")
		}
	}

	// try to return them sorted
	versionMatches, err := NewVersions(matches)
	if err == nil {
		sort.Sort(versionMatches)
		return versionMatches.ToStrings(), nil
	}

	return matches, nil
}
