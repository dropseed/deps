package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/algobardo/yaml"
	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/schema"
)

var directoryNamesToSkip = map[string]bool{
	".git":         true,
	"node_modules": true,
	"env":          true,
	"vendor":       true,
}

type Config struct {
	Cases []*Case `yaml:"cases"`
	path  string
}

func (c *Config) relPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	rel, err := filepath.Rel(wd, c.path)
	if err != nil {
		panic(err)
	}
	return rel
}

type Case struct {
	Name           string      `yaml:"name"`
	Type           string      `yaml:"type"`
	RepoContents   string      `yaml:"repo_contents"`
	InputDataPath  string      `yaml:"input_data_path,omitempty"`
	OutputDataPath string      `yaml:"output_data_path"`
	Tests          []string    `yaml:"tests,omitempty"`
	UserConfig     *UserConfig `yaml:"user_config,omitempty"`
}

func (c *Case) displayName() string {
	return fmt.Sprintf("\"%s\" %s", c.Name, c.Type)
}

type UserConfig struct {
	Path     string                 `yaml:"path"`
	Settings map[string]interface{} `yaml:"settings"`
}

func (c *Case) asConfigDependency(imageName string) *config.Dependency {
	depConfig := &config.Dependency{
		Type: imageName,
	}
	if c.UserConfig != nil {
		depConfig.Path = c.UserConfig.Path
		depConfig.Settings = c.UserConfig.Settings
	}
	depConfig.Compile()
	return depConfig
}

func (c *Case) inputSchema() (*schema.Dependencies, error) {
	if c.InputDataPath != "" {
		inputSchema, err := schema.NewDependenciesFromJSONPath(c.InputDataPath)
		if err != nil {
			return nil, err
		}
		return inputSchema, nil
	}
	return nil, nil
}

func (c *Case) outputSchema() (interface{}, error) {
	if c.OutputDataPath != "" {
		// Try to return as Dependencies first
		outputSchema, err := schema.NewDependenciesFromJSONPath(c.OutputDataPath)
		if err == nil && (len(outputSchema.Manifests) > 0 || len(outputSchema.Lockfiles) > 0) {
			return outputSchema, nil
		}

		// Fall back to basic json parse
		fileContent, err := ioutil.ReadFile(c.OutputDataPath)
		if err != nil {
			return nil, err
		}
		var output interface{}
		if err := json.Unmarshal(fileContent, &output); err != nil {
			return nil, err
		}
		return output, nil
	}
	return nil, nil
}

// NewConfigFromPath loads a Config from a file
func NewConfigFromPath(path string) (*Config, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	cfg, err := NewConfigFromReader(f)
	if cfg != nil {
		cfg.path = path
	}
	return cfg, err
}

func NewConfigFromReader(reader io.Reader) (*Config, error) {
	config := &Config{}
	decoder := yaml.NewDecoder(reader)
	decoder.SetDefaultMapType(reflect.TypeOf(map[string]interface{}{}))
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

func findTestConfigs(dir string) ([]*Config, error) {
	configPaths := findTestConfigPaths(dir, 0)

	if len(configPaths) < 1 {
		return nil, errors.New("no test config files found")
	}

	configs := []*Config{}
	for _, p := range configPaths {
		config, err := NewConfigFromPath(p)
		if err != nil {
			return nil, err
		}
		configs = append(configs, config)
	}

	return configs, nil
}

func findTestConfigPaths(dir string, depth int) []string {
	if depth > 2 {
		return []string{}
	}

	paths := []string{}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

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

			paths = append(paths, findTestConfigPaths(p, depth+1)...)
		} else if isConfigFile(name) {
			paths = append(paths, p)
		}
	}

	return paths
}

func isConfigFile(name string) bool {
	return strings.HasPrefix(name, "dependencies_test") && (strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml"))
}
