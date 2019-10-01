package config

import (
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

func TestFull(t *testing.T) {
	config, err := NewConfigFromPath("./testdata/v2_full.yml")
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
    - name: django
    - name: .*
`
	compareToYAML(t, config, expected)
}

func TestMinimal(t *testing.T) {
	config, err := NewConfigFromPath("./testdata/v2_minimal.yml")
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
}

func TestConfigFromMap(t *testing.T) {
	m := map[string]interface{}{
		"version": Version,
	}
	config, err := newConfigFromMap(m)

	if err != nil {
		t.Error(err)
	}

	if config.Version != Version {
		t.FailNow()
	}
}
