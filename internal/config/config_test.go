package config

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func compareToYAML(t *testing.T, config *Config, expected string) {
	dumped, err := config.DumpYAML()
	if err != nil {
		t.Error(err)
	}
	if expected != dumped {
		print(dumped)
		t.FailNow()
	}
}

func compareToJSON(t *testing.T, config *Config, path string) {
	config.Compile()
	j, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		t.Error(err)
	}
	f, err := ioutil.ReadFile(path)
	if err != nil {
		t.Error(err)
	}
	if string(j)+"\n" != string(f) {
		println(string(j))
		t.FailNow()
	}
}

func TestFull(t *testing.T) {
	config, err := NewConfigFromPath("./testdata/v2_full.yml", nil)
	if err != nil {
		t.Error(err)
	}
	if len(config.Dependencies) != 2 {
		t.FailNow()
	}
	expected := `version: 3
dependencies:
- type: js
- type: python
  path: requirements.txt
  settings:
    ok: true
  lockfile_updates:
    enabled: false
  manifest_updates:
    enabled: true
    filters:
    - name: django-braces
      enabled: false
    - name: django-.*
      versions: Y.Y.Y
    - name: django
      versions: L.Y.Y
    - name: .*
      versions: Y.Y.Y
`
	compareToYAML(t, config, expected)
	compareToJSON(t, config, "./testdata/v2_full.json")
}

func TestMinimal(t *testing.T) {
	config, err := NewConfigFromPath("./testdata/v2_minimal.yml", nil)
	if err != nil {
		t.Error(err)
		return
	}
	expected := `version: 3
dependencies:
- type: python
  path: requirements.txt
`
	compareToYAML(t, config, expected)
	compareToJSON(t, config, "./testdata/v2_minimal.json")
}

func TestMinimalVariables(t *testing.T) {
	python := map[string]interface{}{
		"path": "requirements.txt",
		"type": "python",
	}
	configVars := map[string]interface{}{
		"deps": []map[string]interface{}{python},
	}
	config, err := NewConfigFromPath("./testdata/v2_minimal_variables.yml", configVars)
	if err != nil {
		t.Error(err)
		return
	}
	expected := `version: 3
dependencies:
- type: python
  path: requirements.txt
`
	compareToYAML(t, config, expected)
	compareToJSON(t, config, "./testdata/v2_minimal_variables.json")
}

func TestVariableInSettings(t *testing.T) {
	configVars := map[string]interface{}{
		"labels": []string{"deps", "automerge"},
	}
	config, err := NewConfigFromPath("./testdata/v2_variable_in_settings.yml", configVars)
	if err != nil {
		t.Error(err)
		return
	}
	expected := `version: 3
dependencies:
- type: python
  settings:
    github_labels:
    - deps
    - automerge
`
	compareToYAML(t, config, expected)
	compareToJSON(t, config, "./testdata/v2_variable_in_settings.json")
}

func TestConfigFromMap(t *testing.T) {
	m := map[string]interface{}{
		"version": Version,
	}
	config, err := newConfigFromMap(m, nil)

	if err != nil {
		t.Error(err)
	}

	if config.Version != Version {
		t.FailNow()
	}
}

func TestMissingVariable(t *testing.T) {
	m := map[string]interface{}{
		"version":      Version,
		"dependencies": "$test",
	}
	config, err := newConfigFromMap(m, nil)
	if err.Error() != `1 error(s) decoding:

* error decoding 'dependencies': $test variable not found` {
		t.Fail()
	}
	if config != nil {
		t.Fail()
	}
}
