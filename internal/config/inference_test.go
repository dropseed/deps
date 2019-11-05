package config

import "testing"

func TestInference(t *testing.T) {
	config, err := InferredConfigFromDir("./testdata/repo")
	if err != nil {
		t.Error(err)
	}
	if len(config.Dependencies) != 5 {
		t.FailNow()
	}
	expected := `version: 3
dependencies:
- type: python
  path: Pipfile
- type: python
  path: app/requirements.txt
- type: python
  path: app/requirements_test.txt
- type: php
  path: .
- type: js
  path: .
`
	dumped, err := config.DumpYAML()
	if err != nil {
		t.Error(err)
	}
	if expected != dumped {
		println(dumped)
		t.FailNow()
	}
}
