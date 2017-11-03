package main

import (
	"errors"
	"flag"

	"github.com/dependencies-io/pullrequest/internal/app/schema"
)

// ConfigFlags store the go flags used by the user
type ConfigFlags struct {
	branch             string
	title              string
	body               string
	dependenciesSchema string
	titleFromSchema    bool
	bodyFromSchema     bool
}

func parseFlags() *ConfigFlags {
	branch := flag.String("branch", "", "branch that pull request will be created from")
	title := flag.String("title", "", "pull request title")
	body := flag.String("body", "", "pull request body")

	dependenciesSchema := flag.String("dependencies-schema", "", "dependencies.io schema for the dependencies being acted on")
	titleFromSchema := flag.Bool("title-from-schema", false, "automatically generate the title from the dependencies-schema")
	bodyFromSchema := flag.Bool("body-from-schema", false, "automatically generate the body from the dependencies-schema")

	flag.Parse()

	return &ConfigFlags{
		branch:             *branch,
		title:              *title,
		body:               *body,
		dependenciesSchema: *dependenciesSchema,
		titleFromSchema:    *titleFromSchema,
		bodyFromSchema:     *bodyFromSchema,
	}
}

func (cf *ConfigFlags) titleFromConfigFlags() (string, error) {
	// user supplied title overrides all
	if cf.title != "" {
		return cf.title, nil
	}

	if cf.dependenciesSchema != "" {
		schema, err := schema.NewSchemaFromString([]byte(cf.dependenciesSchema))
		if err != nil {
			return "", err
		}

		if cf.titleFromSchema {
			title, err := schema.GenerateTitleFromSchema()
			if err != nil {
				return "", err
			}

			return title, nil
		}
	}

	return "", errors.New("\"title\" is required or needs the appropriate settings to be generated from the dependencies-schema")
}

func (cf *ConfigFlags) bodyFromConfigFlags() (string, error) {
	// user supplied body overrides all
	if cf.body != "" {
		return cf.body, nil
	}

	if cf.dependenciesSchema != "" {
		schema, err := schema.NewSchemaFromString([]byte(cf.dependenciesSchema))
		if err != nil {
			return "", err
		}

		if cf.bodyFromSchema {
			body, err := schema.GenerateBodyFromSchema()
			if err != nil {
				return "", err
			}

			return body, nil
		}
	}

	return "", errors.New("\"body\" is required or needs the appropriate settings to be generated from the dependencies-schema")
}
