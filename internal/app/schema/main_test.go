package schema

import (
	"io/ioutil"
	"testing"
)

func generateTitleFromFilename(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	schema, err := NewSchemaFromString(content)
	if err != nil {
		return "", err
	}

	title, err := schema.GenerateTitleFromSchema()
	if err != nil {
		return "", err
	}

	return title, nil
}

func TestMalformedJSON(t *testing.T) {
	_, err := NewSchemaFromString([]byte("{not a json}"))
	if err == nil {
		t.Fail()
	}
}

func TestGenerateTitleFromSchemaWithSingleDependency(t *testing.T) {
	title, err := generateTitleFromFilename("./test_data/single_dependency.json")
	if err != nil {
		t.Error(err)
	}
	if title != "Update pullrequest in / from 0.1.0 to 0.3.0" {
		t.Error("Title does not match expected: ", title)
	}
}

func TestGenerateTitleFromSchemaWithTwoDependencies(t *testing.T) {
	title, err := generateTitleFromFilename("./test_data/two_dependencies.json")
	if err != nil {
		t.Error(err)
	}
	if title != "Update 2 dependencies from go, pip" {
		t.Error("Title does not match expected: ", title)
	}
}

func TestGenerateTitleFromSchemaNoDependencies(t *testing.T) {
	title, err := generateTitleFromFilename("./test_data/no_dependencies.json")
	if title != "" {
		t.Fail()
	}
	if err == nil {
		t.Fail()
	}
}

func generateBodyFromFilename(filename string) (string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	schema, err := NewSchemaFromString(content)
	if err != nil {
		return "", err
	}

	body, err := schema.GenerateBodyFromSchema()
	if err != nil {
		return "", err
	}

	return body, nil
}

func TestGenerateBodyFromSchemaWithSingleDependency(t *testing.T) {
	body, err := generateBodyFromFilename("./test_data/single_dependency.json")
	if err != nil {
		t.Error(err)
	}
	expected, err := ioutil.ReadFile("./test_data/single_body.txt")
	if err != nil {
		panic(err)
	}
	if body != string(expected) {
		t.Error("Body does not match expected: ", body)
	}
}

func TestGenerateBodyFromSchemaWithTwoDependencies(t *testing.T) {
	body, err := generateBodyFromFilename("./test_data/two_dependencies.json")
	if err != nil {
		t.Error(err)
	}
	expected, err := ioutil.ReadFile("./test_data/two_body.txt")
	if err != nil {
		panic(err)
	}
	if body != string(expected) {
		t.Error("Body does not match expected: ", body)
	}
}

func TestGenerateBodyFromSchemaNoDependencies(t *testing.T) {
	body, err := generateBodyFromFilename("./test_data/no_dependencies.json")
	if body != "" {
		t.Fail()
	}
	if err == nil {
		t.Fail()
	}
}
