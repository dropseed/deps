package schema

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

type Dependencies struct {
	Lockfiles map[string]*Lockfile `json:"lockfiles,omitempty"`
	Manifests map[string]*Manifest `json:"manifests,omitempty"`
}

// NewDependenciesFromJSONPath loads Dependencies from a JSON file path
func NewDependenciesFromJSONPath(path string) (*Dependencies, error) {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return NewDependenciesFromJSONContent(fileContent)
}

// NewDependenciesFromJSONContent creates a Dependencies instance with Unmarshalled JSON data
func NewDependenciesFromJSONContent(content []byte) (*Dependencies, error) {
	deps := Dependencies{}
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&deps); err != nil {
		return nil, err
	}

	if err := deps.Validate(); err != nil {
		return nil, err
	}

	return &deps, nil
}

func (s *Dependencies) Validate() error {
	for _, lockfile := range s.Lockfiles {
		if err := lockfile.Validate(); err != nil {
			return err
		}
	}
	for _, manifest := range s.Manifests {
		if err := manifest.Validate(); err != nil {
			return err
		}
	}

	return nil
}
