package config

import (
	"path/filepath"
	"regexp"

	"github.com/dropseed/deps/internal/filefinder"
)

// InferredConfigFromPath loads a Config object based on the dependency files present
func InferredConfigFromDir(dir string) (*Config, error) {
	patterns := map[string]*regexp.Regexp{
		"requirements.txt": regexp.MustCompile("^.*requirements.*\\.txt$"),
		"Pipfile":          regexp.MustCompile("^Pipfile$"),
		"package.json":     regexp.MustCompile("^package\\.json$"),
		"composer.json":    regexp.MustCompile("^composer\\.json$"),
	}

	types := map[string]struct {
		Name   string
		UseDir bool
	}{
		"requirements.txt": struct {
			Name   string
			UseDir bool
		}{
			Name:   "python",
			UseDir: false,
		},
		"Pipfile": struct {
			Name   string
			UseDir bool
		}{
			Name:   "python",
			UseDir: false,
		},
		"package.json": struct {
			Name   string
			UseDir bool
		}{
			Name:   "js",
			UseDir: true,
		},
		"composer.json": struct {
			Name   string
			UseDir bool
		}{
			Name:   "php",
			UseDir: true,
		},
		// 	Pattern: regexp.MustCompile("^Dockerfile.*$"),
		// 	Type:    "dockerfile",
	}

	dependencies := []*Dependency{}
	for path, patternName := range filefinder.FindInDir(dir, patterns) {
		depType := types[patternName]
		dep := &Dependency{
			Path: path,
			Type: depType.Name,
		}
		if depType.UseDir {
			dep.Path = filepath.Dir(dep.Path)
		}
		dependencies = append(dependencies, dep)
	}

	config := &Config{
		Version:      Version,
		Dependencies: dependencies,
	}
	// make the dependency paths relative to the dir we were asked to look in
	for _, dep := range config.Dependencies {
		p, err := filepath.Rel(dir, dep.Path)
		if err != nil {
			panic(err)
		}
		dep.Path = p
	}
	return config, nil
}
