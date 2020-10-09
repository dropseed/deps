package schemaext

import (
	"io/ioutil"
	"testing"

	"github.com/dropseed/deps/pkg/schema"
)

func generateTitleFromFilename(filename string) (string, error) {
	dependencies, err := schema.NewDependenciesFromJSONPath(filename)
	if err != nil {
		return "", err
	}

	return TitleForDeps(dependencies), nil
}

func TestMalformedJSON(t *testing.T) {
	_, err := schema.NewDependenciesFromJSONContent([]byte("{not a json}"))
	if err == nil {
		t.FailNow()
	}
}

func TestGenerateTitleWithSingleDependency(t *testing.T) {
	title, err := generateTitleFromFilename("./testdata/single_dependency.json")
	if err != nil {
		t.Error(err)
	}
	if title != "Update pullrequest from 0.1.0 to 0.3.0" {
		t.Error("Title does not match expected: ", title)
	}
}

func TestGenerateTitleWithSingleDependencyNoManifestName(t *testing.T) {
	title, err := generateTitleFromFilename("./testdata/single_dependency_empty_manifest.json")
	if err != nil {
		t.Error(err)
	}
	if title != "Update dropseed/pullrequest from 0.1.0 to 0.3.0" {
		t.Error("Title does not match expected: ", title)
	}
}

// func TestGenerateTitleWithTwoDependencies(t *testing.T) {
// 	title, err := generateTitleFromFilename("./testdata/two_dependencies.json")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if title != "Update 2 dependencies from go, pip" {
// 		t.Error("Title does not match expected: ", title)
// 	}
// }

// func TestGenerateTitleNoDependencies(t *testing.T) {
// 	title, err := generateTitleFromFilename("./testdata/no_dependencies.json")
// 	if title != "" {
// 		t.FailNow()
// 	}
// 	if err == nil {
// 		t.FailNow()
// 	}
// }

// func TestGenerateTitleWithOneLockfile(t *testing.T) {
// 	title, err := generateTitleFromFilename("./testdata/single_lockfile.json")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if title != "Update yarn.lock" {
// 		t.Error("Title does not match expected: ", title)
// 	}
// }

// func TestGenerateTitleWithTwoLockfiles(t *testing.T) {
// 	title, err := generateTitleFromFilename("./testdata/two_lockfiles.json")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if title != "Update lockfiles: composer.lock, yarn.lock" {
// 		t.Error("Title does not match expected: ", title)
// 	}
// }

// func TestGenerateTitleWithLockfilesAndManifests(t *testing.T) {
// 	title, err := generateTitleFromFilename("./testdata/lockfiles_and_manifests.json")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if title != "Update 1 lockfile and 1 manifest" {
// 		t.Error("Title does not match expected: ", title)
// 	}
// }

func generateBodyFromFilename(filename string) (string, error) {
	dependencies, err := schema.NewDependenciesFromJSONPath(filename)
	if err != nil {
		return "", err
	}
	return DescriptionForDeps(dependencies), nil
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

func TestGenerateBodyWithSingleDependencyEmptyManifest(t *testing.T) {
	body, err := generateBodyFromFilename("./testdata/single_dependency_empty_manifest.json")
	if err != nil {
		t.Error(err)
	}
	expected, err := ioutil.ReadFile("./testdata/single_body_empty_manifest.txt")
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

// func TestGenerateBodyWithTwoLockfiles(t *testing.T) {
// 	body, err := generateBodyFromFilename("./testdata/two_lockfiles.json")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	expected, err := ioutil.ReadFile("./testdata/two_lockfiles.txt")
// 	if err != nil {
// 		panic(err)
// 	}
// 	if body != string(expected) {
// 		t.Error("Body does not match expected: ", body)
// 	}
// }

// func TestGenerateBodyWithLockfilesAndManifests(t *testing.T) {
// 	body, err := generateBodyFromFilename("./testdata/lockfiles_and_manifests.json")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	expected, err := ioutil.ReadFile("./testdata/lockfiles_and_manifests.txt")
// 	if err != nil {
// 		panic(err)
// 	}
// 	if body != string(expected) {
// 		t.Error("Body does not match expected: ", body)
// 	}
// }

// func TestGetIDConsistency(t *testing.T) {
// 	dependencies, err := NewDependenciesFromJSONPath("./testdata/lockfiles_and_manifests.json")
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	initialID := dependencies.getUpdateID()
// 	for index := 0; index < 100; index++ {
// 		testID := dependencies.getUpdateID()
// 		if initialID != testID {
// 			t.FailNow()
// 		}
// 	}
// }
