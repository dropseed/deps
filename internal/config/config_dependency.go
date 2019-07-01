package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/dropseed/deps/internal/env"
)

// Dependency is a path + type in dependencies.yml
type Dependency struct {
	Type string `mapstructure:"type" yaml:"type" json:"type"`
	// Version         string                 `mapstructure:"version,omitempty" yaml:"version,omitempty" json:"version,omitempty"`
	Path            string                 `mapstructure:"path,omitempty" yaml:"path,omitempty" json:"path,omitempty"`
	Env             map[string]string      `mapstructure:"env,omitempty" yaml:"env,omitempty" json:"env"`
	Settings        map[string]interface{} `mapstructure:"settings,omitempty" yaml:"settings,omitempty" json:"settings"`
	LockfileUpdates LockfileUpdates        `mapstructure:"lockfile_updates,omitempty" yaml:"lockfile_updates,omitempty" json:"lockfile_updates,omitempty"`
	ManifestUpdates ManifestUpdates        `mapstructure:"manifest_updates,omitempty" yaml:"manifest_updates,omitempty" json:"manifest_updates,omitempty"`
}

func (dependency *Dependency) Compile() {
	dependency.Path = strings.Trim(dependency.Path, "/")
	if dependency.Path == "" {
		dependency.Path = "."
	}

	// set defaults
	if dependency.Settings == nil {
		dependency.Settings = map[string]interface{}{}
	}
	if dependency.LockfileUpdates.Enabled == nil {
		t := true
		dependency.LockfileUpdates.Enabled = &t
	}
	if dependency.ManifestUpdates.Enabled == nil {
		t := true
		dependency.ManifestUpdates.Enabled = &t
	}

	// if no filters then set the default 1
	if len(dependency.ManifestUpdates.Filters) == 0 {
		defaultFilter := &Filter{
			Name: ".*",
		}
		dependency.ManifestUpdates.Filters = append(dependency.ManifestUpdates.Filters, defaultFilter)
	}
	for _, filter := range dependency.ManifestUpdates.Filters {
		if filter.Enabled == nil {
			t := true
			filter.Enabled = &t
		}
		if filter.Group == nil {
			t := false
			filter.Group = &t
		}
	}
}

func (dependency *Dependency) Environ() ([]string, error) {
	environ := os.Environ()

	for k, v := range dependency.Settings {
		environString, err := env.SettingToEnviron(k, v)
		if err != nil {
			return nil, err
		}
		environ = append(environ, environString)
	}

	for k, v := range dependency.Env {
		environ = append(environ, fmt.Sprintf("%s=%s", k, v))
	}

	return environ, nil
}

func (dependency *Dependency) GetSetting(key string) interface{} {
	if dependency.Settings == nil {
		return nil
	}
	for k, v := range dependency.Settings {
		if strings.ToLower(k) == strings.ToLower(key) {
			return v
		}
	}
	return nil
}
