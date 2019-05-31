package versionfilter

import "regexp"

type RegexFilter struct {
	regex *regexp.Regexp
}

func NewRegexFilter(s string) (*RegexFilter, error) {
	regex, err := regexp.Compile(s)
	if err != nil {
		return nil, err
	}
	f := &RegexFilter{
		regex: regex,
	}
	return f, nil
}

func (f *RegexFilter) Matching(versions []string, currentVersion string) ([]string, error) {
	matches := []string{}
	for _, version := range versions {
		if f.regex.MatchString(version) {
			matches = append(matches, version)
		}
	}
	return matches, nil
}
