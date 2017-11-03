package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// Schema stores a dependencies-schema
type Schema struct {
	content []byte // the original JSON encoded string
	data    map[string]interface{}
}

// NewSchemaFromString creates a Schema instance with Unmarshalled JSON data
func NewSchemaFromString(content []byte) (*Schema, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, err
	}
	return &Schema{content: content, data: data}, nil
}

// GenerateTitleFromSchema generates a title string from the dependencies schema
func (s *Schema) GenerateTitleFromSchema() (string, error) {
	dependencies, _ := s.data["dependencies"].([]interface{})

	if len(dependencies) < 1 {
		return "", errors.New("Must have at least 1 dependency")
	}

	if len(dependencies) == 1 {
		dep := dependencies[0].(map[string]interface{})
		name, _ := dep["name"].(string)
		path, _ := dep["path"].(string)
		installed, _ := dep["installed"].(map[string]interface{})["version"].(string)
		highest, _ := dep["available"].([]interface{})[0].(map[string]interface{})["version"].(string)
		return fmt.Sprintf("Update %v in %v from %v to %v", name, path, installed, highest), nil
	}

	// create a "set" of sources
	sources := make(map[string]bool)
	for _, dep := range dependencies {
		source, _ := dep.(map[string]interface{})["source"].(string)
		sources[source] = true
	}

	// get the keys remaining
	sourceNames := []string{}
	for k := range sources {
		sourceNames = append(sourceNames, k)
	}

	sort.Strings(sourceNames)

	// TODO if > 2 items, put an "and " in front of the last one

	return fmt.Sprintf("Update %v dependencies from %v", len(dependencies), strings.Join(sourceNames, ", ")), nil
}

// GenerateBodyFromSchema generates a body string from the dependencies schema
func (s *Schema) GenerateBodyFromSchema() (string, error) {
	dependencies, _ := s.data["dependencies"].([]interface{})

	if len(dependencies) < 1 {
		return "", errors.New("Must have at least 1 dependency")
	}

	if len(dependencies) == 1 {
		body, err := getContentForDependency(dependencies[0].(interface{}))
		body = body[:len(body)-6] // remove the final "\n---\n"
		return body, err
	}

	summary := "## Overview\n\nThe following dependencies have been updated:\n\n"
	body := "## Details\n\n"

	for _, dep := range dependencies {
		depSummary, err := getSummaryForDependency(dep)
		if err != nil {
			return "", err
		}
		summary = summary + fmt.Sprintf("- %v\n", depSummary)

		depBody, err := getContentForDependency(dep)
		if err != nil {
			return "", err
		}
		body = body + depBody
	}

	body = body[:len(body)-6] // remove the final "\n---\n"

	return summary + "\n" + body, nil
}

func getSummaryForDependency(dep interface{}) (string, error) {
	name, _ := dep.(map[string]interface{})["name"].(string)
	path, _ := dep.(map[string]interface{})["path"].(string)
	installed, _ := dep.(map[string]interface{})["installed"].(map[string]interface{})["version"]
	available, _ := dep.(map[string]interface{})["available"]
	highest, _ := available.([]interface{})[0].(map[string]interface{})["version"].(string)
	return fmt.Sprintf("%v in `%v` from `%v` to `%v`", name, path, installed, highest), nil
}

func getContentForDependency(dep interface{}) (string, error) {
	// TODO add notes

	// there's got to be a better way to traverse this...
	name, _ := dep.(map[string]interface{})["name"].(string)
	path, _ := dep.(map[string]interface{})["path"].(string)
	source, _ := dep.(map[string]interface{})["source"].(string)
	installed, _ := dep.(map[string]interface{})["installed"].(map[string]interface{})["version"]
	available, _ := dep.(map[string]interface{})["available"]
	highest, _ := available.([]interface{})[0].(map[string]interface{})["version"].(string)
	subject := fmt.Sprintf("[Dependencies.io](https://www.dependencies.io) has updated %v (a %v dependency in `%v`) from `%v` to `%v`.", name, source, path, installed, highest)

	content := ""

	for _, v := range available.([]interface{}) {
		version := v.(map[string]interface{})["version"].(string)

		vContent, ok := v.(map[string]interface{})["content"].(string)
		if !ok {
			vContent = "_No content found. Please open an issue at https://github.com/dependencies-io/support if you think this content could have been found._"
		}
		content += fmt.Sprintf("\n\n<details>\n<summary>%v</summary>\n\n%v\n\n</details>", version, vContent)
	}

	return subject + content + "\n\n---\n\n", nil
}
