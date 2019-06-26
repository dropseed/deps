package config

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"regexp"
)

// MaxInferenceDepth determines how deep in the repo to look
const MaxInferenceDepth = 2

var directoryNamesToSkip = map[string]bool{
	".git":         true,
	"node_modules": true,
	"env":          true,
	"vendor":       true,
}

type inferencePattern struct {
	Pattern *regexp.Regexp
	Type    string
	UseDir  bool
}

var patterns = []inferencePattern{
	inferencePattern{
		Pattern: regexp.MustCompile("^.*requirements.*\\.txt$"),
		Type:    "python",
		UseDir:  false,
	},
	inferencePattern{
		Pattern: regexp.MustCompile("^Pipfile$"),
		Type:    "python",
		UseDir:  false,
	},
	inferencePattern{
		Pattern: regexp.MustCompile("^package.json$"),
		Type:    "js",
		UseDir:  true,
	},
	// inferencePattern{
	// 	Pattern: regexp.MustCompile("^composer.json$"),
	// 	Type:    "php",
	// 	UseDir:  true,
	// },
	// inferencePattern{
	// 	Pattern: regexp.MustCompile("^Dockerfile.*$"),
	// 	Type:    "dockerfile",
	// 	UseDir:  false,
	// },
}

// InferredConfigFromPath loads a Config object based on the dependency files present
func InferredConfigFromDir(dir string) (*Config, error) {
	dependencies := inferredDependenciesFromDir(dir, 1)
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

func inferredDependenciesFromDir(dir string, depth int) []*Dependency {
	if depth > MaxInferenceDepth {
		return nil
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	dependencies := []*Dependency{}

	for _, f := range files {
		name := f.Name()
		p := path.Join(dir, name)

		fileInfo, err := os.Stat(p)
		if err != nil {
			panic(err)
		}

		if fileInfo.IsDir() {
			if directoryNamesToSkip[name] {
				continue
			}

			if dirDeps := inferredDependenciesFromDir(p, depth+1); dirDeps != nil {
				dependencies = append(dependencies, dirDeps...)
			}
		} else if dependency := inferredDependencyFromPath(p); dependency != nil {
			dependencies = append(dependencies, dependency)
		}
	}

	return dependencies
}

func inferredDependencyFromPath(p string) *Dependency {
	basename := path.Base(p)
	for _, ip := range patterns {
		if ip.Pattern.MatchString(basename) {
			dep := &Dependency{
				Path: p,
				Type: ip.Type,
			}
			if ip.UseDir {
				dep.Path = filepath.Dir(dep.Path)
			}
			return dep
		}
	}
	return nil
}
