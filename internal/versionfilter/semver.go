package versionfilter

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type SemverishFilter struct {
	*Version
	Operator string
}

const OP_ASTERISK = "*"
const OP_NEXTBEST = "-"
const OP_LTE = "<="
const OP_LT = "<"
const OP_GTE = ">="
const OP_GT = ">"
const OP_NEQ = "!="
const OP_CARET = "^"
const OP_TILDE = "~"

var operators = map[string]bool{
	OP_NEXTBEST: true,
	OP_ASTERISK: true,
	OP_LTE:      true,
	OP_LT:       true,
	OP_GTE:      true,
	OP_GT:       true,
	OP_NEQ:      true,
	OP_CARET:    true,
	OP_TILDE:    true,
}

var nonWordStartRegex = regexp.MustCompile("^\\W+")

func NewSemverishFilter(s string) (*SemverishFilter, error) {
	original := s

	operator, s, err := splitOperator(s)
	if err != nil {
		return nil, err
	}

	version, err := NewVersion(s)
	if err != nil {
		return nil, err
	}

	f := &SemverishFilter{
		Version:  version,
		Operator: operator,
	}
	f.Original = original // make sure original has operator on it

	if err := f.validateDotParts(); err != nil {
		return nil, err
	}

	return f, nil
}

func (f *SemverishFilter) validateDotParts() error {
	if f.Operator == OP_ASTERISK {
		// the only scenario where dotparts is allowed to be 1 empty string
		// (the only operator that doesn't have to precede something)
		return nil
	}

	for _, p := range f.DotParts {
		if p == "" {
			return errors.New("dot separated part is empty")
		}
		if p == YES {
			return nil
		}
		if lockRegex.MatchString(p) {
			return nil
		}

		if _, err := strconv.Atoi(p); err != nil {
			// You cannot use a string if there is a known Operator
			if f.Operator != "" {
				return err
			}
		}
	}

	return nil
}

func (f *SemverishFilter) Matching(versions []string, currentVersion string) ([]string, error) {
	matches, err := f.MatchingVersions(versions, currentVersion)
	if err != nil {
		return nil, err
	}
	return matches.ToStrings(), nil
}

func (f *SemverishFilter) MatchingVersions(versions []string, currentVersion string) (Versions, error) {
	if f.Operator == OP_NEXTBEST {
		return f.MatchingNextBestVersions(versions, currentVersion)
	}

	matches := Versions{}

	var current *Version
	if currentVersion != "" {
		var err error
		current, err = NewVersion(currentVersion)
		if err != nil {
			return nil, err
		}
	}

	for _, version := range versions {
		v, err := NewVersion(version)
		if err != nil {
			return nil, err
		}

		match, err := f.Match(v, current)
		if err != nil {
			return nil, err
		}

		if match {
			matches = append(matches, v)
		}
	}

	sort.Sort(matches)

	return matches, nil
}

func (f *SemverishFilter) MatchingNextBestVersions(versions []string, currentVersion string) (Versions, error) {
	if f.Build != "" || f.Prerelease != "" || f.PrereleaseBuild != "" {
		return nil, errors.New("next best filter cannot have pre-release or build fields")
	}

	var current *Version
	if currentVersion != "" {
		var err error
		current, err = NewVersion(currentVersion)
		if err != nil {
			return nil, err
		}
	}

	fResolved := f.ResolvedVersion(current)
	matches := Versions{}

	if fResolved.UsesYes() {
		// if Y is used, then compile a set of
		// (possible truncated and/or fake) versions that
		// and then we'll choose the first item of each
		versionsForKey := map[string]Versions{}

		lastYIndex := 0
		for i, p := range fResolved.DotParts {
			if p == YES {
				lastYIndex = i
			}
		}

		for _, vString := range versions {
			v, err := NewVersion(vString)
			if err != nil {
				// skip if not a valid version
				continue
			}

			// swap out the Y in the filter for what this version has
			keyParts := []string{}
			for i := 0; i <= lastYIndex; i++ {
				if fResolved.DotParts[i] == YES {
					keyParts = append(keyParts, v.DotParts[i])
				} else {
					keyParts = append(keyParts, fResolved.DotParts[i])
				}
			}
			key := strings.Join(keyParts, ".")
			versionsForKey[key] = append(versionsForKey[key], v)
		}

		// the first sorted version for each key is a match
		for _, keyVersions := range versionsForKey {
			sort.Sort(keyVersions)
			matches = append(matches, keyVersions[0])
		}
	} else {
		// if Y is not used, then we will only
		// ever get 1 or 0 results and can do it by just using a ^ filter
		if filter, err := NewSemverishFilter("^" + fResolved.Original); err == nil {
			if submatches, err := filter.MatchingVersions(versions, fResolved.Original); err != nil {
				return nil, err
			} else if len(submatches) > 0 {
				matches = append(matches, submatches[0])
			}
		} else {
			return nil, err
		}
	}

	// filter out everything at or below the current version
	anyFilter, err := NewSemverishFilter(OP_ASTERISK)
	if err != nil {
		panic(err)
	}
	matches, err = anyFilter.MatchingVersions(matches.ToStrings(), currentVersion)
	if err != nil {
		panic(err)
	}

	sort.Sort(matches)
	return matches, nil
}

func (f *SemverishFilter) Match(t *Version, currentVersion *Version) (bool, error) {
	// Filter out anything below the current version
	if currentVersion != nil {
		v, err := t.Compare(currentVersion)
		if err != nil {
			return false, err
		}
		if v < 1 {
			return false, nil
		}
	}

	if f.Operator == OP_ASTERISK {
		return true, nil
	}

	if f.UsesYes() && !t.IsYesCompatible() {
		// if the version isn't compatible with the Y mask,
		// then silently throw it out (ex. python syntax of 1.2.1dev)
		return false, nil
	}

	fResolved := f.ResolvedVersion(currentVersion)

	if f.Operator == OP_CARET {
		return fResolved.matchesCaret(t)
	}

	if f.Operator == OP_TILDE {
		return fResolved.matchesTilde(t)
	}

	comparisonValue, err := fResolved.Compare(t)
	if err != nil {
		return false, err
	}

	if comparisonValue == 0 && (f.Operator == "" || f.Operator == OP_GTE || f.Operator == OP_LTE) {
		return true, nil
	}

	if comparisonValue != 0 && f.Operator == OP_NEQ {
		return true, nil
	}

	if comparisonValue == -1 && (f.Operator == OP_GT || f.Operator == OP_GTE) {
		return true, nil
	}

	if comparisonValue == 1 && (f.Operator == OP_LT || f.Operator == OP_LTE) {
		return true, nil
	}

	return false, nil
}

func (f *SemverishFilter) ResolvedVersion(currentVersion *Version) *Version {
	resolved := &Version{
		Original:        f.Original,
		DotParts:        f.DotParts,
		Build:           f.Build,
		Prerelease:      f.Prerelease,
		PrereleaseBuild: f.PrereleaseBuild,
	}
	if currentVersion != nil {
		for i, p := range resolved.DotParts {
			if currentVersion != nil && len(currentVersion.DotParts) >= i+1 {
				resolved.DotParts[i] = resolveVariable(p, currentVersion.DotParts[i], true)
			} else {
				resolved.DotParts[i] = resolveVariable(p, "", true)
			}
		}

		resolved.Build = resolveVariable(resolved.Build, currentVersion.Build, true)
		resolved.Prerelease = resolveVariable(resolved.Prerelease, currentVersion.Prerelease, true)
		resolved.PrereleaseBuild = resolveVariable(resolved.PrereleaseBuild, currentVersion.PrereleaseBuild, true)

		asString := strings.Join(resolved.DotParts, ".")
		if resolved.Build != "" {
			asString = asString + BUILD_SEPARATOR + resolved.Build
		}
		if resolved.Prerelease != "" {
			asString = asString + PRERELEASE_SEPARATOR + resolved.Prerelease
		}
		if resolved.PrereleaseBuild != "" {
			asString = asString + BUILD_SEPARATOR + resolved.PrereleaseBuild
		}
		resolved.Original = asString

	} else {
		for i, p := range resolved.DotParts {
			resolved.DotParts[i] = resolveVariable(p, "", false)
		}
		resolved.Build = resolveVariable(resolved.Build, "", false)
		resolved.Prerelease = resolveVariable(resolved.Prerelease, "", false)
		resolved.PrereleaseBuild = resolveVariable(resolved.PrereleaseBuild, "", false)
		_, vString, err := splitOperator(resolved.Original)
		if err != nil {
			panic(err)
		}
		resolved.Original = strings.TrimSpace(vString)
	}

	return resolved
}

func resolveVariable(a, b string, hasCurrentVersion bool) string {
	if lockRegex.MatchString(a) {
		if !hasCurrentVersion {
			panic(errors.New("cannot use L without a current version"))
		}
		lockIncrement := a[1:]
		if lockIncrement == "" {
			return b
		}
		incrementInt, incErr := strconv.Atoi(lockIncrement)
		if incErr != nil {
			panic(incErr)
		}
		bInt, bErr := strconv.Atoi(b)
		if bErr != nil {
			return b
		}
		return strconv.Itoa(incrementInt + bInt)
	}
	return a
}

func splitOperator(s string) (string, string, error) {
	s = strings.TrimSpace(s)
	operator := nonWordStartRegex.FindString(s)
	if operator != "" {
		if ok := operators[operator]; !ok {
			return "", "", fmt.Errorf("%s is an unknown operator", operator)
		}

		s = s[len(operator):]
		s = strings.TrimSpace(s)
	}

	return operator, s, nil
}
