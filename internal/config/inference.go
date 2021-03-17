package config

import (
	"path/filepath"
	"regexp"
	"sort"

	"github.com/dropseed/deps/internal/filefinder"
)

// InferredConfigFromPath loads a Config object based on the dependency files present
func InferredConfigFromDir(dir string) (*Config, error) {
	patterns := map[string]*regexp.Regexp{
		"requirements.txt": regexp.MustCompile("^.*requirements.*\\.txt$"),
		"Pipfile":          regexp.MustCompile("^Pipfile$"),
		"poetry.lock":      regexp.MustCompile("^poetry\\.lock$"),
		"package.json":     regexp.MustCompile("^package\\.json$"),
		"composer.json":    regexp.MustCompile("^composer\\.json$"),
	}

	types := map[string]struct {
		Name   string
		UseDir bool
	}{
		"requirements.txt": {
			Name:   "python",
			UseDir: false,
		},
		"Pipfile": {
			Name:   "python",
			UseDir: false,
		},
		"poetry.lock": {
			Name:   "python",
			UseDir: false,
		},
		"package.json": {
			Name:   "js",
			UseDir: true,
		},
		"composer.json": {
			Name:   "php",
			UseDir: true,
		},
		// 	Pattern: regexp.MustCompile("^Dockerfile.*$"),
		// 	Type:    "dockerfile",
	}

	dependencies := []*Dependency{}

	// Have to sort these for YAML consistency
	found := filefinder.FindInDir(dir, patterns)
	foundPathsSorted := []string{}
	for p := range found {
		foundPathsSorted = append(foundPathsSorted, p)
	}
	sort.Strings(foundPathsSorted)

	for _, path := range foundPathsSorted {
		patternName := found[path]
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
