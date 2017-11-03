package main

import (
	"io/ioutil"
	"testing"
)

func TestNormal(t *testing.T) {
	cf := ConfigFlags{
		branch: "test",
		title:  "test title",
		body:   "test body",
	}
	title, _ := cf.titleFromConfigFlags()
	if title != "test title" {
		t.Error("generated title does not match")
	}
	body, _ := cf.bodyFromConfigFlags()
	if body != "test body" {
		t.Error("generated body does not match")
	}
}

func TestSchemaTitle(t *testing.T) {
	schemaContent, _ := ioutil.ReadFile("../../internal/app/schema/test_data/single_dependency.json")
	cf := ConfigFlags{
		branch:             "test",
		dependenciesSchema: string(schemaContent),
		titleFromSchema:    true,
		body:               "test body",
	}

	if title, _ := cf.titleFromConfigFlags(); title != "Update pullrequest in / from 0.1.0 to 0.3.0" {
		t.Error("generated title does not match")
	}

	body, _ := cf.bodyFromConfigFlags()
	if body != "test body" {
		t.Error("generated body does not match")
	}
}

func TestSchemaBody(t *testing.T) {
	schemaContent, _ := ioutil.ReadFile("../../internal/app/schema/test_data/single_dependency.json")
	cf := ConfigFlags{
		branch:             "test",
		dependenciesSchema: string(schemaContent),
		title:              "test title",
		bodyFromSchema:     true,
	}

	if title, _ := cf.titleFromConfigFlags(); title != "test title" {
		t.Error("generated title does not match")
	}

	body, _ := cf.bodyFromConfigFlags()
	expectedBody, _ := ioutil.ReadFile("../../internal/app/schema/test_data/single_body.txt")
	if body != string(expectedBody) {
		t.Error("generated body does not match")
	}
}

func TestSchemaTitleAndBody(t *testing.T) {
	schemaContent, _ := ioutil.ReadFile("../../internal/app/schema/test_data/single_dependency.json")
	cf := ConfigFlags{
		branch:             "test",
		dependenciesSchema: string(schemaContent),
		titleFromSchema:    true,
		bodyFromSchema:     true,
	}

	if title, _ := cf.titleFromConfigFlags(); title != "Update pullrequest in / from 0.1.0 to 0.3.0" {
		t.Error("generated title does not match")
	}

	body, _ := cf.bodyFromConfigFlags()
	expectedBody, _ := ioutil.ReadFile("../../internal/app/schema/test_data/single_body.txt")
	if body != string(expectedBody) {
		t.Error("generated body does not match")
	}
}
