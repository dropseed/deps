package schema

import (
	"io/ioutil"
	"os"
	"testing"
)

func generateTitleFromFilename(filename string) (string, error) {
	dependencies, err := NewDependenciesFromJSONPath(filename)
	if err != nil {
		return "", err
	}

	title, err := dependencies.GenerateTitle()
	if err != nil {
		return "", err
	}

	return title, nil
}

func TestMalformedJSON(t *testing.T) {
	_, err := NewDependenciesFromJSONContent([]byte("{not a json}"))
	if err == nil {
		t.FailNow()
	}
}

func TestGenerateTitleWithSingleDependency(t *testing.T) {
	title, err := generateTitleFromFilename("./testdata/single_dependency.json")
	if err != nil {
		t.Error(err)
	}
	if title != "Update pullrequest in / from 0.1.0 to 0.3.0" {
		t.Error("Title does not match expected: ", title)
	}
}

func TestGenerateTitleWithTwoDependencies(t *testing.T) {
	title, err := generateTitleFromFilename("./testdata/two_dependencies.json")
	if err != nil {
		t.Error(err)
	}
	if title != "Update 2 dependencies from go, pip" {
		t.Error("Title does not match expected: ", title)
	}
}

func TestGenerateTitleNoDependencies(t *testing.T) {
	title, err := generateTitleFromFilename("./testdata/no_dependencies.json")
	if title != "" {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
}

func TestGenerateTitleWithOneLockfile(t *testing.T) {
	title, err := generateTitleFromFilename("./testdata/single_lockfile.json")
	if err != nil {
		t.Error(err)
	}
	if title != "Update yarn.lock" {
		t.Error("Title does not match expected: ", title)
	}
}

func TestGenerateTitleWithTwoLockfiles(t *testing.T) {
	title, err := generateTitleFromFilename("./testdata/two_lockfiles.json")
	if err != nil {
		t.Error(err)
	}
	if title != "Update lockfiles: composer.lock, yarn.lock" {
		t.Error("Title does not match expected: ", title)
	}
}

func TestGenerateTitleWithLockfilesAndManifests(t *testing.T) {
	title, err := generateTitleFromFilename("./testdata/lockfiles_and_manifests.json")
	if err != nil {
		t.Error(err)
	}
	if title != "Update 1 lockfile and 1 manifest" {
		t.Error("Title does not match expected: ", title)
	}
}

func generateRelatedPRTitleSearchFromFilename(filename string) (string, error) {
	dependencies, err := NewDependenciesFromJSONPath(filename)
	if err != nil {
		return "", err
	}

	title, err := dependencies.GenerateRelatedPRTitleSearch()
	if err != nil {
		return "", err
	}

	return title, nil
}

func TestGenerateRelatedPRTitleSearchWithSingleDependency(t *testing.T) {
	title, err := generateRelatedPRTitleSearchFromFilename("./testdata/single_dependency.json")
	if err != nil {
		t.Error(err)
	}
	if title != "Update pullrequest in / from 0.1.0 to " {
		t.Error("Related PR title search does not match expected: ", title)
	}
}

func TestGenerateRelatedPRTitleSearchWithTwoDependencies(t *testing.T) {
	title, err := generateRelatedPRTitleSearchFromFilename("./testdata/two_dependencies.json")
	if err == nil || title != "" {
		t.FailNow()
	}
}

func TestGenerateRelatedPRTitleSearchNoDependencies(t *testing.T) {
	title, err := generateRelatedPRTitleSearchFromFilename("./testdata/no_dependencies.json")
	if err == nil || title != "" {
		t.FailNow()
	}
}

func TestGenerateRelatedPRTitleWithOneLockfile(t *testing.T) {
	title, err := generateRelatedPRTitleSearchFromFilename("./testdata/single_lockfile.json")
	if err != nil {
		t.Error(err)
	}
	if title != "Update yarn.lock" {
		t.Error("Title does not match expected: ", title)
	}
}

func TestGenerateRelatedPRTitleWithTwoLockfiles(t *testing.T) {
	title, err := generateRelatedPRTitleSearchFromFilename("./testdata/two_lockfiles.json")
	if err != nil {
		t.Error(err)
	}
	if title != "Update lockfiles: composer.lock, yarn.lock" {
		t.Error("Title does not match expected: ", title)
	}
}

func TestGenerateRelatedPRTitleWithLockfilesAndManifests(t *testing.T) {
	title, err := generateRelatedPRTitleSearchFromFilename("./testdata/lockfiles_and_manifests.json")
	if err == nil || title != "" {
		t.FailNow()
	}
}

func generateBodyFromFilename(filename string) (string, error) {
	dependencies, err := NewDependenciesFromJSONPath(filename)
	if err != nil {
		return "", err
	}

	os.Setenv("SETTING_PULLREQUEST_VERSIONS_API_DISABLED", "true")
	body, err := dependencies.GenerateBody()
	if err != nil {
		return "", err
	}

	return body, nil
}

func TestGenerateBodyWithSingleDependency(t *testing.T) {
	body, err := generateBodyFromFilename("./testdata/single_dependency.json")
	if err != nil {
		t.Error(err)
	}
	expected, err := ioutil.ReadFile("./testdata/single_body.txt")
	if err != nil {
		panic(err)
	}
	if body != string(expected) {
		t.Error("Body does not match expected: ", body)
	}
}

func TestGenerateBodyWithTwoDependencies(t *testing.T) {
	body, err := generateBodyFromFilename("./testdata/two_dependencies.json")
	if err != nil {
		t.Error(err)
	}
	expected, err := ioutil.ReadFile("./testdata/two_body.txt")
	if err != nil {
		panic(err)
	}
	if body != string(expected) {
		t.Error("Body does not match expected: ", body)
	}
}

func TestGenerateBodyNoDependencies(t *testing.T) {
	body, err := generateBodyFromFilename("./testdata/no_dependencies.json")
	if body != "" {
		t.FailNow()
	}
	if err == nil {
		t.FailNow()
	}
}

func TestGenerateBodyWithOneLockfile(t *testing.T) {
	body, err := generateBodyFromFilename("./testdata/single_lockfile.json")
	if err != nil {
		t.Error(err)
	}
	expected, err := ioutil.ReadFile("./testdata/single_lockfile.txt")
	if err != nil {
		panic(err)
	}
	if body != string(expected) {
		t.Error("Body does not match expected: ", body)
	}
}

func TestGenerateBodyWithTwoLockfiles(t *testing.T) {
	body, err := generateBodyFromFilename("./testdata/two_lockfiles.json")
	if err != nil {
		t.Error(err)
	}
	expected, err := ioutil.ReadFile("./testdata/two_lockfiles.txt")
	if err != nil {
		panic(err)
	}
	if body != string(expected) {
		t.Error("Body does not match expected: ", body)
	}
}

func TestGenerateBodyWithLockfilesAndManifests(t *testing.T) {
	body, err := generateBodyFromFilename("./testdata/lockfiles_and_manifests.json")
	if err != nil {
		t.Error(err)
	}
	expected, err := ioutil.ReadFile("./testdata/lockfiles_and_manifests.txt")
	if err != nil {
		panic(err)
	}
	if body != string(expected) {
		t.Error("Body does not match expected: ", body)
	}
}
