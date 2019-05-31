package versionfilter

type VersionFilter interface {
	Matching(versions []string, currentVersion string) ([]string, error)
}

// NewVersionFilter returns new filter struct if a match is found
func NewVersionFilter(s string) VersionFilter {
	if f, err := NewBooleanFilter(s); err == nil {
		return f
	}

	if f, err := NewSemverishFilter(s); err == nil {
		return f
	}

	if f, err := NewRegexFilter(s); err == nil {
		return f
	}

	return nil
}
