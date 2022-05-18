package config

import (
	"fmt"
	"testing"

	"github.com/dropseed/deps/pkg/schema"
)

func compareToYAML(t *testing.T, config *Config, expected string) {
	dumped, err := config.DumpYAML()
	if err != nil {
		t.Error(err)
	}
	if expected != dumped {
		fmt.Print(dumped)
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
      settings:
        github_labels:
        - django
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

func TestFilterSettings(t *testing.T) {
	// create a schema with a manifest update
	deps, err := schema.NewDependenciesFromJSONPath("../runner/testdata/single_dependency.json")
	if err != nil {
		t.Error(err)
	}

	depConfig := &Dependency{
		Settings: Settings{
			"github_labels": []string{"test"},
		},
		ManifestUpdates: ManifestUpdates{
			Filters: []*Filter{
				{
					Name: "pullrequest",
					Settings: Settings{
						"github_labels": []string{"pullrequest"},
					},
				},
				{
					Name: ".*",
					Settings: Settings{
						"github_labels": []string{"all"},
					},
				},
			},
		},
	}

	value := depConfig.GetSettingForSchema("github_labels", deps)
	if len(value.([]string)) != 1 {
		t.FailNow()
	}
	if value.([]string)[0] != "pullrequest" {
		t.Error("expected pullrequest, got ", value.([]string)[0])
	}
}
