package versionfilter

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

type Version struct {
	Original        string
	DotParts        []string
	Build           string
	Prerelease      string
	PrereleaseBuild string
}

const PRERELEASE_SEPARATOR = "-"
const BUILD_SEPARATOR = "+"
const YES = "Y"

var lockRegex = regexp.MustCompile("^L\\d*$")
var vEqPrefixRegex = regexp.MustCompile("^v?=*?\\d")

func NewVersion(s string) (*Version, error) {
	f := &Version{
		Original: s,
	}

	s = strings.TrimSpace(s)

	if vEqPrefixRegex.MatchString(s) {
		s = strings.TrimLeft(s, "v=")
		s = strings.TrimSpace(s)
	}

	if nonWordStartRegex.MatchString(s) {
		return nil, fmt.Errorf("version cannot start with a non-word character: %s", s)
	}

	// if s == "" {
	// 	return nil, errors.New("version cannot be an empty string")
	// }

	prereleaseSplit := strings.Split(s, PRERELEASE_SEPARATOR)
	if len(prereleaseSplit) == 2 {
		s = prereleaseSplit[0]
		f.Prerelease = prereleaseSplit[1]
		buildSplit := strings.Split(f.Prerelease, BUILD_SEPARATOR)
		if len(buildSplit) == 2 {
			f.Prerelease = buildSplit[0]
			f.PrereleaseBuild = buildSplit[1]
		} else if len(buildSplit) != 1 {
			return nil, errors.New("invalid prerelease Build usage")
		}
	} else if len(prereleaseSplit) != 1 {
		return nil, errors.New("invalid prerelease usage")
	}

	buildSplit := strings.Split(s, BUILD_SEPARATOR)
	if len(buildSplit) == 2 {
		s = buildSplit[0]
		f.Build = buildSplit[1]
	} else if len(buildSplit) != 1 {
		return nil, errors.New("invalid build usage")
	}

	f.DotParts = strings.Split(s, ".")
	f.validateDotParts()

	return f, nil
}

func (v *Version) validateDotParts() error {
	for _, p := range v.DotParts {
		if p == "" {
			return errors.New("dot separated part is empty")
		}
	}

	return nil
}

func (v *Version) UsesLocks() bool {
	toCheck := []string{}
	toCheck = append(toCheck, v.DotParts...)
	toCheck = append(toCheck, v.Build, v.Prerelease, v.PrereleaseBuild)
	for _, p := range toCheck {
		if lockRegex.MatchString(p) {
			return true
		}
	}
	return false
}

func (v *Version) UsesYes() bool {
	toCheck := []string{}
	toCheck = append(toCheck, v.DotParts...)
	toCheck = append(toCheck, v.Build, v.Prerelease, v.PrereleaseBuild)
	for _, p := range toCheck {
		if p == YES {
			return true
		}
	}
	return false
}

func (v *Version) IsYesCompatible() bool {
	for _, s := range v.DotParts {
		if _, err := strconv.Atoi(s); err != nil {
			return false
		}
	}
	return true
}

func (f *Version) Compare(t *Version) (int, error) {
	if f.UsesLocks() || t.UsesLocks() {
		return 0, errors.New("cannot compare versions still containing lock variables")
	}

	comparisonValue := comparePartSlices(f.DotParts, t.DotParts)

	if comparisonValue == 0 && (f.Prerelease != "" || t.Prerelease != "") && f.Prerelease != YES {
		// if they are equal, compare prerelease strings
		fPrereleaseParts := strings.Split(f.Prerelease, ".")
		tPrereleaseParts := strings.Split(t.Prerelease, ".")
		comparisonValue = comparePartSlices(fPrereleaseParts, tPrereleaseParts)
	}

	return comparisonValue, nil
}

func comparePartSlices(a, b []string) int {
	fillSlices(&a, &b)

	comparisonValue := 0

	for i, aPart := range a {
		if aPart == YES {
			continue
		}

		bPart := b[i]

		aInt, aIntErr := strconv.Atoi(aPart)
		bInt, bIntErr := strconv.Atoi(bPart)

		if aIntErr != nil || bIntErr != nil {
			comparisonValue = strings.Compare(aPart, bPart)
		} else {
			if aInt > bInt {
				comparisonValue = 1
			} else if aInt < bInt {
				comparisonValue = -1
			}
		}

		if comparisonValue != 0 {
			// As soon as we find something inequal, we know where they fall
			break
		}
	}

	return comparisonValue
}

func (f *Version) matchesCaret(t *Version) (bool, error) {
	// Can only modify places after first non-zero
	firstNonZero := -1
	for i, p := range f.DotParts {
		if p != "0" {
			firstNonZero = i
			break
		}
	}

	if firstNonZero == -1 {
		return false, errors.New("no non-zero part found")
	}

	for i := 0; i <= firstNonZero; i++ {
		if f.DotParts[i] == YES {
			continue
		}
		if f.DotParts[i] != t.DotParts[i] {
			return false, nil
		}
	}

	// if they weren't asking for a prerelease, exclude all prereleases
	// (flip side -- if asked for any prerelease, than all are included, which is what node-semver does)
	if f.Prerelease == "" && t.Prerelease != "" {
		return false, nil
	}

	return true, nil
}

func (f *Version) matchesTilde(t *Version) (bool, error) {
	if len(f.DotParts) > 3 || len(t.DotParts) > 3 {
		return false, errors.New("can only use tilde operator on versions with 3 or fewer places")
	}

	fParts := f.DotParts
	tParts := t.DotParts
	fillSlices(&fParts, &tParts)

	// major range always has to be the same
	if len(f.DotParts) > 0 && fParts[0] != tParts[0] {
		return false, nil
	}

	// if anything after major was originally specified,
	// then minor has to be the same
	if len(f.DotParts) > 1 && fParts[1] != tParts[1] {
		return false, nil
	}

	// also has to be > than the filter
	// (ex. ~1.0.1 eliminates 1.0.0, 1.0.1)
	if comparePartSlices(fParts, tParts) >= 0 {
		return false, nil
	}

	// if they weren't asking for a prerelease, exclude all prereleases
	// (flip side -- if asked for any prerelease, than all are included, which is what node-semver does)
	if f.Prerelease == "" && t.Prerelease != "" {
		return false, nil
	}

	return true, nil
}
