package config

import (
	"io/ioutil"
	"testing"
)

func TestNormal(t *testing.T) {
	config := Config{
		EnvSettings: &EnvSettings{},
		Flags: &Flags{
			Branch: "test",
			Title:  "test title",
			Body:   "test body",
		},
	}
	title, _ := config.TitleFromConfig()
	if title != "test title" {
		t.Error("generated title does not match")
	}
	body, _ := config.BodyFromConfig()
	if body != "test body" {
		t.Error("generated body does not match")
	}
}

func TestSchemaTitle(t *testing.T) {
	schemaContent, _ := ioutil.ReadFile("../schema/test_data/single_dependency.json")
	config := Config{
		EnvSettings: &EnvSettings{},
		Flags: &Flags{
			Branch:             "test",
			DependenciesSchema: string(schemaContent),
			TitleFromSchema:    true,
			Body:               "test body",
		},
	}

	if title, _ := config.TitleFromConfig(); title != "Update pullrequest in / from 0.1.0 to 0.3.0" {
		t.Error("generated title does not match")
	}

	body, _ := config.BodyFromConfig()
	if body != "test body" {
		t.Error("generated body does not match")
	}
}

func TestSchemaBody(t *testing.T) {
	schemaContent, _ := ioutil.ReadFile("../schema/test_data/single_dependency.json")
	config := Config{
		EnvSettings: &EnvSettings{},
		Flags: &Flags{
			Branch:             "test",
			DependenciesSchema: string(schemaContent),
			Title:              "test title",
			BodyFromSchema:     true,
		},
	}

	if title, _ := config.TitleFromConfig(); title != "test title" {
		t.Error("generated title does not match")
	}

	body, _ := config.BodyFromConfig()
	expectedBody, _ := ioutil.ReadFile("../schema/test_data/single_body.txt")
	if body != string(expectedBody) {
		t.Error("generated body does not match")
	}
}

func TestSchemaTitleAndBody(t *testing.T) {
	schemaContent, _ := ioutil.ReadFile("../schema/test_data/single_dependency.json")
	config := Config{
		EnvSettings: &EnvSettings{},
		Flags: &Flags{
			Branch:             "test",
			DependenciesSchema: string(schemaContent),
			TitleFromSchema:    true,
			BodyFromSchema:     true,
		},
	}

	if title, _ := config.TitleFromConfig(); title != "Update pullrequest in / from 0.1.0 to 0.3.0" {
		t.Error("generated title does not match")
	}

	body, _ := config.BodyFromConfig()
	expectedBody, _ := ioutil.ReadFile("../schema/test_data/single_body.txt")
	if body != string(expectedBody) {
		t.Error("generated body does not match")
	}
}
