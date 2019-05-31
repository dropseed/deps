package versionfilter

import (
	"regexp"
	"testing"
)

func expectFilterResults(t *testing.T, filter string, versions []string, current string, expected []string) {
	f := NewVersionFilter(filter)
	if f == nil {
		t.Errorf("version filter parse failed")
		t.FailNow()
	}
	matches, err := f.Matching(versions, current)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}

	if len(matches) != len(expected) {
		t.Errorf("matches differ from expected: %v %v", matches, expected)
		t.FailNow()
	}
	for i, m := range matches {
		if m != expected[i] {
			t.Errorf("match differs from expected at %d: %v %v", i, m, expected[i])
			t.FailNow()
		}
	}
}

func TestRegexParse(t *testing.T) {
	filter := NewVersionFilter(".*")
	if filter == nil {
		t.FailNow()
	}
	if filter.(*RegexFilter).regex.String() != regexp.MustCompile(".*").String() {
		t.FailNow()
	}

	bad := NewVersionFilter(".*[")
	if bad != nil {
		t.FailNow()
	}
}

func TestAsterisk(t *testing.T) {
	filter := "*"
	versions := []string{"1.8", "1.8.1", "1.8.2", "1.9", "1.9.1", "1.10", "nightly"}
	current := "1.9"
	expected := []string{"1.9.1", "1.10", "nightly"}
	expectFilterResults(t, filter, versions, current, expected)
}

// def test_specitemmask_asterisk():
//     s = SpecItemMask("*")
//     assert(Spec("*") == s.spec)
//     assert(Version("0.0.1") in s.spec)
//     assert(Version("0.1.1-alpha") in s.spec)
//
//
// def test_specitemmask_asterisk_forms():
//     s = SpecItemMask(" *")
//     assert s.kind == "*"
//
//     s = SpecItemMask(" *  ")
//     assert s.kind == "*"
//
//     s = SpecItemMask("*")
//     assert s.kind == "*"
//
//     with pytest.raises(ValueError):
//         SpecItemMask("bad*")
//
//     with pytest.raises(ValueError):
//         SpecItemMask("* bad")
//
//
// def test_specitemmask_lock1():
//     s = SpecItemMask("L.0.0", current_version=Version("1.0.0"))
//     assert(Spec("1.0.0") == s.spec)
//
//
// def test_specitemmask_lock2():
//     s = SpecItemMask("L.L.0", current_version=Version("1.8.0"))
//     assert(Spec("1.8.0") == s.spec)
//
//
// def test_specitemmask_lock3():
//     s = SpecItemMask("L.L.L", current_version=Version("1.8.3"))
//     assert(Spec("1.8.3") == s.spec)
//
//
// def test_specitemmask_lock4():
//     s = SpecItemMask("L1.L.L", current_version=Version("1.8.3"))
//     assert(Spec("2.8.3") == s.spec)
//
//
// def test_specitemmask_lock5():
//     s = SpecItemMask("L1.L999.L", current_version=Version("1.8.3"))
//     assert(Spec("2.1007.3") == s.spec)
//
//
// def test_specitemmask_lock6():
//     with pytest.raises(ValueError):
//         SpecItemMask("L-1.L999.L", current_version=Version("1.8.3"))
//
//
// def test_specitemmask_yes1():
//     s = SpecItemMask("Y.Y.0", current_version=Version("1.8.3"))
//     assert(Spec("*") == s.spec)
//
//
// def test_specitemmask_yes2():
//     s = SpecItemMask("L.Y.0", current_version=Version("1.8.3"))
//     assert(Spec("*") == s.spec)
//
//
// def test_specitemmask_yes3():
//     s = SpecItemMask("L.L.Y", current_version=Version("1.8.3"))
//     assert(Spec("*") == s.spec)
//
//
// def test_specitemmask_modifiers_1():
//     s = SpecItemMask(">1.0.0")
//     assert(Spec(">1.0.0") == s.spec)
//
//
// def test_coerceable_version():
//     s = SpecItemMask("1")
//     assert(Spec("1") == s.spec)
//
//
// def test_specmask():
//     s = SpecMask("1.0.0")
//     assert(Spec("1.0.0") == s.specs[0].spec)
//
//
// def test_specmask_one_or():
//     s = SpecMask("1.0.0 || 2.0.0")
//     assert(Spec("1.0.0") == s.specs[0].spec)
//     assert(Spec("2.0.0") == s.specs[1].spec)
//
//
// def test_specmask_multi_ors():
//     s = SpecMask("1.0.0 || 2.0.0 || 3.0.0 || 4.0.0")
//     assert(Spec("1.0.0") == s.specs[0].spec)
//     assert(Spec("2.0.0") == s.specs[1].spec)
//     assert(Spec("3.0.0") == s.specs[2].spec)
//     assert(Spec("4.0.0") == s.specs[3].spec)
//
//
// def test_specmask_one_and():
//     s = SpecMask("1.0.0 && 2.0.0")
//     assert(Spec("1.0.0") == s.specs[0].spec)
//     assert(Spec("2.0.0") == s.specs[1].spec)
//
//
// def test_specmask_multi_ands():
//     s = SpecMask("1.0.0 && 2.0.0 && 3.0.0 && 4.0.0")
//     assert(Spec("1.0.0") == s.specs[0].spec)
//     assert(Spec("2.0.0") == s.specs[1].spec)
//     assert(Spec("3.0.0") == s.specs[2].spec)
//     assert(Spec("4.0.0") == s.specs[3].spec)

func TestMixedBooleanWillAssert(t *testing.T) {
	_, err := NewBooleanFilter("1.0.0 && 2.0.0 || 3.0.0")
	if err.Error() != "cannot contain both \"&&\" and \"||\"" {
		t.FailNow()
	}
}

// func TestSpecmaskMatch(t *testing.T) {
// 	s, err := NewSemverishFilter("1.0.0")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	match, err := s.MatchStrings("1.0.0", "")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if match != true {
// 		t.FailNow()
// 	}
//
// 	match, err = s.MatchStrings("1.0.1", "")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if match != false {
// 		t.FailNow()
// 	}
// }

// def test_specmask_contains():
//     filter := "1.0.0"
//     s = SpecMask(mask)
//     assert("1.0.0" in s)
//     assert("1.0.1" not in s)
//
//
// def test_partial_versions_1():
//     filter := "1.0"
//     s = SpecItemMask(mask)
//     assert(s.spec == Spec("1.0"))
//
//
// def test_partial_versions_2():
//     filter := "1"
//     s = SpecItemMask(mask)
//     assert(s.spec == Spec("1"))
//
//
// def test_partial_versions_3():
//     filter := "L"
//     current := "1"
//     s = SpecItemMask(mask, current_version)
//     assert(s.spec == Spec("==1"))
//
//
func TestPartialVersions_4(t *testing.T) {
	filter := "L.Y"
	versions := []string{"1.8", "1.8.1", "1.8.2", "1.9", "1.9.1", "1.10", "nightly"}
	current := "1.8"
	expected := []string{"1.9", "1.10"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestPartialVersions_5(t *testing.T) {
	filter := "L.Y.Y"
	versions := []string{"1.8", "1.8.1", "1.8.2", "1.9", "1.9.1", "1.10", "nightly"}
	current := "1.8"
	expected := []string{"1.8.1", "1.8.2", "1.9", "1.9.1", "1.10"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestReadmeExampleSemver(t *testing.T) {
	filter := "L.Y.Y"
	versions := []string{"1.8.0", "1.8.1", "1.8.2", "1.9.0", "1.9.1", "1.10.0", "nightly"}
	current := "1.9.0"
	expected := []string{"1.9.1", "1.10.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestReadmeExampleRegex(t *testing.T) {
	filter := "^night"
	versions := []string{"1.8.0", "1.8.1", "1.8.2", "1.9.0", "1.9.1", "1.10.0", "nightly"}
	current := ""
	expected := []string{"nightly"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestMajorUpdatesOnly_1(t *testing.T) {
	filter := "Y.0.0"
	versions := []string{"1.8.0", "1.8.1", "1.8.2", "1.9.0", "1.9.1", "1.10.0", "2.0.0", "2.0.1"}
	current := "1.9.0"
	expected := []string{"2.0.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestMajorUpdatesOnly_2(t *testing.T) {
	filter := "Y.0.0"
	versions := []string{"1.8.0", "1.8.1", "1.8.2", "1.9.0", "1.9.1", "1.10.0", "2.0.1", "2.0.2"}
	current := "1.9.0"
	expected := []string{}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestMinorUpdates_1(t *testing.T) {
	filter := "Y.Y.0"
	versions := []string{"1.8.0", "1.8.1", "1.8.2", "1.9.0", "1.9.1", "1.10.0", "2.0.0", "2.0.1"}
	current := "1.8.0"
	expected := []string{"1.9.0", "1.10.0", "2.0.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestMinorUpdates_2(t *testing.T) {
	filter := "L.Y.0"
	versions := []string{"1.8.0", "1.8.1", "1.8.2", "1.9.0", "1.9.1", "1.10.0", "2.0.0", "2.0.1"}
	current := "1.8.0"
	expected := []string{"1.9.0", "1.10.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestAllUpdates_1(t *testing.T) {
	filter := "Y.Y.Y"
	versions := []string{"1.8.0", "1.8.1", "1.8.2", "1.9.0", "1.9.1", "1.10.0", "2.0.0", "2.0.1"}
	current := "1.8.0"
	expected := []string{"1.8.1", "1.8.2", "1.9.0", "1.9.1", "1.10.0", "2.0.0", "2.0.1"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestExplicitMajorUpdatesOnly_1(t *testing.T) {
	filter := "2.0.0"
	versions := []string{"1.8.0", "1.8.1", "1.8.2", "1.9.0", "1.9.1", "1.10.0", "2.0.0", "2.0.1"}
	current := "1.9.0"
	expected := []string{"2.0.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestPythonFilter(t *testing.T) {
	filter := "Y.Y.Y"
	versions := []string{"0.1.0", "1.1.0", "1.2.1.dev0", "1.2.dev0", "1.2.post0", "1.2.0a1"}
	current := ""
	expected := []string{"0.1.0", "1.1.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestDjangoConfigExample_1(t *testing.T) {
	filter := "1.8.Y"
	versions := []string{"1.8.0", "1.8.1", "1.8.2", "1.9.0", "1.9.1", "1.10.0", "2.0.0", "2.0.1"}
	current := ""
	expected := []string{"1.8.0", "1.8.1", "1.8.2"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestDjangoConfigExample_2(t *testing.T) {
	filter := "1.8.Y || 1.10.Y"
	versions := []string{"1.8.0", "1.8.1", "1.8.2", "1.9.0", "1.9.1", "1.10.0", "2.0.0", "2.0.1"}
	current := ""
	expected := []string{"1.8.0", "1.8.1", "1.8.2", "1.10.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestDjangoCurrentExample_1(t *testing.T) {
	filter := "Y.Y.0 || L.L.Y"
	versions := []string{"1.8.0", "1.8.1", "1.8.2", "1.9.0", "1.9.1", "1.10.0", "2.0.0", "2.0.1"}
	current := "1.8.1"
	expected := []string{"1.8.2", "1.9.0", "1.10.0", "2.0.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestModifierExample_1(t *testing.T) {
	filter := ">1.8.0 && <2.0.0"
	versions := []string{"1.8.0", "1.8.1", "1.8.2", "1.9.0", "1.9.1", "1.10.0", "2.0.0", "2.0.1"}
	current := "1.8.1"
	expected := []string{"1.8.2", "1.9.0", "1.9.1", "1.10.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestPrerelease_1(t *testing.T) {
	filter := "L.Y.Y"
	versions := []string{"0.9.5", "1.0.0-alpha.e2", "1.0.0-alpha.12", "1.0.0-alpha.58"}
	current := "0.9.5"
	expected := []string{}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestPrerelease_2(t *testing.T) {
	filter := "L.Y.Y-Y"
	versions := []string{"0.9.5", "1.0.0-alpha.e2", "1.0.0-alpha.12", "1.0.0-alpha.58"}
	current := "0.9.5"
	expected := []string{}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestPrerelease_3(t *testing.T) {
	filter := "L.Y.Y-Y"
	versions := []string{"0.9.5", "1.0.0-alpha.e2", "1.0.0-alpha.12", "1.0.0-alpha.58", "0.9.6-alpha.ef"}
	current := "0.9.5"
	expected := []string{"0.9.6-alpha.ef"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestPrerelease_4(t *testing.T) {
	filter := "L.Y.Y"
	versions := []string{"0.9.5", "1.0.0-alpha.e2", "1.0.0-alpha.12", "1.0.0-alpha.58", "0.9.6", "1.0.0"}
	current := "0.9.5"
	expected := []string{"0.9.6"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestPrerelease_5(t *testing.T) {
	filter := "Y.Y.Y"
	versions := []string{"0.9.5", "1.0.0-alpha.e2", "1.0.0-alpha.12", "1.0.0-alpha.58", "0.9.6", "1.0.0"}
	current := ""
	expected := []string{"0.9.5", "0.9.6", "1.0.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestPrerelease_6(t *testing.T) {
	filter := "Y.Y.Y-Y"
	versions := []string{"0.9.5", "1.0.0-alpha.e2", "1.0.0-alpha.12", "1.0.0-alpha.58", "0.9.6", "1.0.0"}
	current := ""
	expected := []string{"0.9.5", "0.9.6", "1.0.0", "1.0.0-alpha.12", "1.0.0-alpha.58", "1.0.0-alpha.e2"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestPrereleaseMatching(t *testing.T) {
	filter := "L.L.Y-alpine"
	versions := []string{"3.6", "3.6-alpine", "3.6-onbuild", "3.6.1", "3.6.1-alpine", "3.6.1-alpine3.6", "3.6.1-onbuild"}
	current := "3.6-alpine"
	expected := []string{"3.6.1-alpine"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestPrereleaseMatching_2(t *testing.T) {
	filter := "L.L.Y-alpine3.6"
	versions := []string{"3.6", "3.6-alpine", "3.6-alpine3.6", "3.6-onbuild", "3.6.1", "3.6.1-alpine", "3.6.1-alpine3.6", "3.6.1-onbuild"}
	current := "3.6-alpine3.6"
	expected := []string{"3.6.1-alpine3.6"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestPrereleaseLock(t *testing.T) {
	filter := "L.L.Y-L"
	versions := []string{"3.6", "3.6-alpine", "3.6-onbuild", "3.6.1", "3.6.1-alpine", "3.6.1-alpine3.6", "3.6.1-onbuild"}
	current := "3.6-alpine"
	expected := []string{"3.6.1-alpine"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestPrereleaseLock_2(t *testing.T) {
	filter := "L.L.Y-L"
	versions := []string{"3.6", "3.6-alpine", "3.6-alpine3.6", "3.6-onbuild", "3.6.1", "3.6.1-alpine", "3.6.1-alpine3.6", "3.6.1-onbuild"}
	current := "3.6-alpine3.6"
	expected := []string{"3.6.1-alpine3.6"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestVPrefixOnVersions(t *testing.T) {
	filter := "L.L.Y"
	versions := []string{"v0.9.5", "v0.9.6", "v1.0.0"}
	current := "0.9.5"
	expected := []string{"v0.9.6"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestVPrefixOnCurrentVersion(t *testing.T) {
	filter := "L.L.Y"
	versions := []string{"0.9.5", "0.9.6", "1.0.0"}
	current := "v0.9.5"
	expected := []string{"0.9.6"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestEqPrefixOnVersions(t *testing.T) {
	filter := "L.L.Y"
	versions := []string{"=0.9.5", "=0.9.6", "=1.0.0"}
	current := "0.9.5"
	expected := []string{"=0.9.6"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestEqPrefixOnCurrentVersion(t *testing.T) {
	filter := "L.L.Y"
	versions := []string{"0.9.5", "0.9.6", "1.0.0"}
	current := "=0.9.5"
	expected := []string{"0.9.6"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestVAndEqPrefixOnCurrentVersion(t *testing.T) {
	filter := "L.L.Y"
	versions := []string{"0.9.5", "0.9.6", "1.0.0"}
	current := "v=0.9.5"
	expected := []string{"0.9.6"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestEqAndEqPrefixOnCurrentVersion(t *testing.T) {
	filter := "L.L.Y"
	versions := []string{"0.9.5", "0.9.6", "1.0.0"}
	current := "==0.9.5"
	expected := []string{"0.9.6"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestCaret(t *testing.T) {
	filter := "^1.0.0"
	versions := []string{"1.0.0", "1.0.1", "1.1.0", "1.2.0-alpha", "2.0.0", "2.0.0-beta"}
	current := "1.0.0"
	expected := []string{"1.0.1", "1.1.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestCaretBabelCliExample(t *testing.T) {
	filter := "^L.L.L"
	versions := []string{
		"6.24.0",
		"7.0.0-alpha.3",
		"7.0.0-alpha.4",
		"7.0.0-alpha.6",
		"7.0.0-alpha.7",
		"6.24.1",
		"7.0.0-alpha.8",
		"7.0.0-alpha.9",
		"7.0.0-alpha.10",
		"7.0.0-alpha.11",
		"7.0.0-alpha.12",
		"7.0.0-alpha.14",
	}
	current := "6.24.0"
	expected := []string{"6.24.1"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestSemverCaret(t *testing.T) {
	filter, err := NewSemverishFilter("^1.0.0")
	if err != nil {
		t.Error(err)
	}

	matches := []string{"1.1.0"}
	notMatches := []string{"1.1.0-alpha", "2.0.0", "2.0.0-alpha"}
	for _, s := range matches {
		v, err := NewVersion(s)
		if err != nil {
			t.Error(err)
		}
		match, _ := filter.Version.matchesCaret(v)
		if !match {
			t.FailNow()
		}
	}
	for _, s := range notMatches {
		v, err := NewVersion(s)
		if err != nil {
			t.Error(err)
		}
		match, _ := filter.Version.matchesCaret(v)
		if match {
			t.FailNow()
		}
	}
}

func TestSemverCaretWithZeroMajor(t *testing.T) {
	filter, err := NewSemverishFilter("^0.3.0")
	if err != nil {
		t.Error(err)
	}

	matches := []string{"0.3.1"}
	notMatches := []string{"0.4.0", "1.3.0"}
	for _, s := range matches {
		v, err := NewVersion(s)
		if err != nil {
			t.Error(err)
		}
		match, _ := filter.Version.matchesCaret(v)
		if !match {
			t.FailNow()
		}
	}
	for _, s := range notMatches {
		v, err := NewVersion(s)
		if err != nil {
			t.Error(err)
		}
		match, _ := filter.Version.matchesCaret(v)
		if match {
			t.FailNow()
		}
	}
}

func TestSemverCaretWithLocks(t *testing.T) {
	filter := "^L.L.L"
	versions := []string{"0.3.0", "0.3.1", "0.4.0", "1.0.0", "1.3.0"}
	current := "0.3.0"
	expected := []string{"0.3.1"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestTildeMajorOnly(t *testing.T) {
	filter := "~1"
	versions := []string{"1.0.0", "1.0.1", "1.1.0", "1.2.0-alpha", "2.0.0", "2.0.0-beta"}
	current := "1.0.0"
	expected := []string{"1.0.1", "1.1.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestTildeMajorMinor(t *testing.T) {
	filter := "~1.0"
	versions := []string{"1.0.0", "1.0.1", "1.1.0", "1.2.0-alpha", "2.0.0", "2.0.0-beta"}
	current := "1.0.0"
	expected := []string{"1.0.1"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestTildeMajorMinorPatch(t *testing.T) {
	filter := "~1.0.0"
	versions := []string{"1.0.0", "1.0.1", "1.1.0", "1.2.0-alpha", "2.0.0", "2.0.0-beta"}
	current := "1.0.0"
	expected := []string{"1.0.1"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestTildeMinor(t *testing.T) {
	filter := "~1.1"
	versions := []string{"1.0.0", "1.0.1", "1.1.1", "1.2.0", "1.2.0-alpha", "2.0.0", "2.0.0-beta"}
	current := "1.0.0"
	expected := []string{"1.1.1"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestTildeMinorPatch(t *testing.T) {
	filter := "~1.1.0"
	versions := []string{"1.0.0", "1.0.1", "1.1.1", "1.2.0", "1.2.0-alpha", "2.0.0", "2.0.0-beta"}
	current := "1.0.0"
	expected := []string{"1.1.1"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestTildePatch(t *testing.T) {
	filter := "~1.0.1"
	versions := []string{"1.0.0", "1.0.1", "1.0.2", "1.1.0", "1.2.0-alpha", "2.0.0", "2.0.0-beta"}
	current := "1.0.0"
	expected := []string{"1.0.2"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestSemverishFilterSortedReturn(t *testing.T) {
	filter := "Y.Y.Y"
	versions := []string{"0.9.6", "1.0.0", "0.9.5"}
	current := ""
	expected := []string{"0.9.5", "0.9.6", "1.0.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

//
//
// def test_valid_version_parsing_1():
//     assert(Version("0.0.1") == _parse_semver("0.0.1"))
//     assert(Version("0.0.1-dev0") == _parse_semver("0.0.1-dev0"))
//     assert(Version("0.0.1-dev0.build0") == _parse_semver("0.0.1-dev0.build0"))
//     assert(Version("0.0.1+something") == _parse_semver("0.0.1+something"))
//
//
// def test_invalid_version_parsing_1():
//     with pytest.raises(InvalidSemverError):
//         _parse_semver("0.0.1.build0")  # invalid build string
//
//
// def test_next_best_specitemmask():
//     s = SpecItemMask("-1.0.0")
//     assert(Spec("1.0.0") == s.spec)
//     assert s.has_next_best
//
//
func TestNextBestSpecitemmaskMatchingVersionsLiteral1(t *testing.T) {
	filter := "-1.0.0"
	versions := []string{
		"1.0.1",
		"2.0.1",
	}
	current := ""
	expected := []string{"1.0.1"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestNextBestSpecitemmaskMatchingVersionsLiteral2(t *testing.T) {
	filter := "-1.0.0"
	versions := []string{
		"1.1.0",
		"2.0.1",
	}
	current := ""
	expected := []string{"1.1.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestNextBestSpecitemmaskMatchingVersionsLiteral3(t *testing.T) {
	filter := "-1.0.0"
	versions := []string{
		"1.0.0",
		"2.0.1",
	}
	current := ""
	expected := []string{}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestNextBestSpecitemmaskMatchingVersionsLock1(t *testing.T) {
	filter := "-L.0.0"
	versions := []string{
		"1.0.1",
		"2.0.1",
	}
	current := "1.0.0"
	expected := []string{"1.0.1"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestNextBestSpecitemmaskMatchingVersionsLock2(t *testing.T) {
	filter := "-L.0.0"
	versions := []string{
		"1.1.0",
		"2.0.1",
	}
	current := "1.0.0"
	expected := []string{"1.1.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestNextBestSpecitemmaskMatchingVersionsLock3(t *testing.T) {
	filter := "-L.0.0"
	versions := []string{
		"1.0.0",
		"2.0.1",
	}
	current := "1.0.0"
	expected := []string{}
	expectFilterResults(t, filter, versions, current, expected)
}

//
//
// def test_get_next_best_versions1():
//     y = YesVersion("*", "Y.0.0")
//     versions := []string{
//         "1.0.0",
//         "1.0.1",
//         "2.0.1",
//         "2.1.0"
//     ]
//     versions := []string{_parse_semver(x) for x in versions]
//
//     result = y.get_next_best_versions(versions)
//     assert(1 == len(result))
//     assert(_parse_semver("2.0.0") in result)
//
//
// def test_get_next_best_versions2():
//     y = YesVersion("*", "Y.Y.0")
//     versions := []string{
//         "1.0.0",
//         "1.0.1",
//         "1.1.1",
//         "2.0.1",
//         "2.1.1"
//     ]
//     versions := []string{_parse_semver(x) for x in versions]
//
//     result = y.get_next_best_versions(versions)
//     assert(3 == len(result))
//     assert(_parse_semver("1.1.0") in result)
//     assert(_parse_semver("2.0.0") in result)
//     assert(_parse_semver("2.1.0") in result)
//
//
// def test_get_next_best_versions3():
//     y = YesVersion("*", "1.0.0")
//     versions := []string{
//         "1.0.0",
//         "1.0.1",
//         "1.1.1",
//         "2.0.1",
//         "2.1.1"
//     ]
//     versions := []string{_parse_semver(x) for x in versions]
//
//     result = y.get_next_best_versions(versions)
//     assert(0 == len(result))
//
//
// def test_get_next_best_versions4():
//     y = YesVersion("*", "Y.Y.Y")
//     versions := []string{
//         "1.0.0",
//         "1.0.1",
//         "1.0.5",
//         "1.1.1",
//         "2.0.1",
//         "2.1.1"
//     ]
//     versions := []string{_parse_semver(x) for x in versions]
//
//     result = y.get_next_best_versions(versions)
//     assert(3 == len(result))
//     assert(_parse_semver("1.0.2") in result)
//     assert(_parse_semver("1.0.3") in result)
//     assert(_parse_semver("1.0.4") in result)
//
//
func TestNextBestSpecitemmaskMatchingVersionsYes1(t *testing.T) {
	filter := "-Y.0.0"
	versions := []string{
		"1.0.1",
		"2.0.1",
	}
	current := "1.0.0"
	expected := []string{"1.0.1", "2.0.1"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestNextBestSpecitemmaskMatchingVersionsYes2(t *testing.T) {
	filter := "-Y.0.0"
	versions := []string{
		"1.1.0",
		"1.2.0",
		"1.3.0",
		"2.0.1",
		"3.1.2",
	}
	current := "1.0.0"
	expected := []string{"1.1.0", "2.0.1", "3.1.2"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestNextBestSpecitemmaskMatchingVersionsYes3(t *testing.T) {
	filter := "-Y.0.0"
	versions := []string{
		"1.0.0",
		"1.1.0",
		"1.2.0",
		"1.3.0",
		"2.0.1",
		"3.1.2",
	}
	current := "1.0.0"
	expected := []string{"2.0.1", "3.1.2"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestNextBestSpecitemmaskMatchingVersionsYes4(t *testing.T) {
	filter := "-Y.0.0"
	versions := []string{
		"1.0.0",
		"2.0.1",
	}
	current := "1.0.0"
	expected := []string{"2.0.1"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestNextBestSpecitemmaskMatchingVersionsYes5(t *testing.T) {
	filter := "-Y.Y.0"
	versions := []string{
		"1.0.0",
		"1.1.0",
		"1.1.1",
		"1.2.1",
		"2.0.1",
	}
	current := "1.0.0"
	expected := []string{"1.1.0", "1.2.1", "2.0.1"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestNextBestSpecitemmaskMatchingVersionsYes6(t *testing.T) {
	filter := "-Y.0.0"
	versions := []string{
		"1.0.0",
		"1.1.0",
		"1.1.1",
		"1.2.1",
		"2.0.1",
	}
	current := "1.2.1"
	expected := []string{"2.0.1"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestNextBestSpecitemmaskMatchingVersionsYes7(t *testing.T) {
	filter := "-1.Y.0"
	versions := []string{
		"1.0.0",
		"1.1.0",
		"1.1.1",
		"1.2.1",
		"2.0.1",
	}
	current := "1.0.1"
	expected := []string{"1.1.0", "1.2.1"}
	expectFilterResults(t, filter, versions, current, expected)
}

//
//
// def test_next_best_specitemmask_with_range_1():
//     """Mixing semver range operators and next_best matching is not allowed"""
//     filter := "-^1.0.0"
//     versions := []string{
//         "1.0.0",
//     ]
//     with pytest.raises(ValueError):
//         VersionFilter.semver_filter(mask, versions)
//
//
// def test_next_best_example():
//     versions := []string{"1.0.0", "2.0.0", "3.0.1"}
//     current := "2.0.0"
//     assert VersionFilter.semver_filter("Y.0.0", versions, current_version) == []
//     # but with next_best ...
//     assert VersionFilter.semver_filter("-Y.0.0", versions, current_version) == ["3.0.1"}
//
//
func TestGreaterThanNextMajor(t *testing.T) {
	filter := ">=L1.0.0"
	versions := []string{"1.0.0", "1.0.1", "1.1.0", "1.2.0", "2.0.0", "2.0.1", "3.0.0"}
	current := "1.0.0"
	expected := []string{"2.0.0", "2.0.1", "3.0.0"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestInNextMajor(t *testing.T) {
	filter := "L1.Y.Y"
	versions := []string{"1.0.0", "1.0.1", "1.1.0", "1.2.0", "2.0.0", "2.0.1", "3.0.0"}
	current := "1.0.0"
	expected := []string{"2.0.0", "2.0.1"}
	expectFilterResults(t, filter, versions, current, expected)
}

func TestValidateOnly(t *testing.T) {
	validMasks := []string{"L1.Y.Y", "1.0.0", "1.0.0", "1.0", "1", "L", "L.Y", "L.Y.Y", "L.Y.Y", "Y.0.0", "Y.0.0", "Y.Y.0",
		"L.Y.0", "Y.Y.Y", "2.0.0", "Y.Y.Y", "1.8.Y", "1.8.Y || 1.10.Y", "Y.Y.0 || L.L.Y",
		">1.8.0 && <2.0.0", "L.Y.Y", "L.Y.Y-Y", "L.Y.Y-Y", "L.Y.Y", "Y.Y.Y", "Y.Y.Y-Y", "L.L.Y-alpine",
		"L.L.Y-alpine3.6", "L.L.Y-L", "L.L.Y-L", "L.L.Y", "L.L.Y", "L.L.Y", "L.L.Y", "L.L.Y", "^1.0.0",
		"^L.L.L", "^L.L.L", "-1.0.0", "-1.0.0", "-1.0.0", "-L.0.0", "-L.0.0", "-L.0.0", "-Y.0.0", "-Y.0.0",
		"-Y.0.0", "-Y.0.0", "-Y.Y.0", "-Y.0.0", ">=L1.0.0", "L1.Y.Y", "L.L.L-L"}

	for _, m := range validMasks {
		v := NewVersionFilter(m)
		if v == nil {
			t.Error("not valid")
		}
	}

	// invalidMasks := []string{"-^1.1.1", "a", "?.?.?", "", "YY.0.0", "LL.0.0"}
	// for _, m := range invalidMasks {
	// 	if _, err := NewSemverishFilter(m); err == nil {
	// 		t.Error("not expected to be valid")
	// 	}
	// }
}
