package config

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"

	yaml "github.com/algobardo/yaml"
	"github.com/dropseed/deps/internal/output"
	"github.com/mitchellh/mapstructure"
)

// DefaultFilename is where a config is placed by default
const DefaultFilename = "dependencies.yml"
const Version = 3

// Config stores a dependencies.yml config
type Config struct {
	Version      int           `mapstructure:"version" yaml:"version" json:"version"`
	Dependencies []*Dependency `mapstructure:"dependencies" yaml:"dependencies" json:"dependencies"`
}

func (config *Config) Compile() {
	for _, dependency := range config.Dependencies {
		dependency.Compile()
	}
}

func LoadOrInferConfigFromPath(configpath string, variables map[string]interface{}) (*Config, error) {
	config, err := NewConfigFromPath(configpath, nil)
	if os.IsNotExist(err) {
		output.Event("No local config found, inferring one from your files")
		config, err = InferredConfigFromDir(filepath.Dir(configpath))
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return config, nil
}

// NewConfigFromPath loads a Config from a file
func NewConfigFromPath(path string, variables map[string]interface{}) (*Config, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	return NewConfigFromReader(f, variables)
}

func NewConfigFromReader(reader io.Reader, variables map[string]interface{}) (*Config, error) {
	temp := map[string]interface{}{}
	decoder := yaml.NewDecoder(reader)
	decoder.SetDefaultMapType(reflect.TypeOf(map[string]interface{}{}))
	if err := decoder.Decode(&temp); err != nil {
		return nil, err
	}

	return newConfigFromMap(temp, variables)
}

func newConfigFromMap(m map[string]interface{}, variables map[string]interface{}) (*Config, error) {
	// set this so it is accessible in the decode func
	tempConfigVariables = variables
	defer resetVariables()

	config := &Config{}

	mapDecoderConfig := mapstructure.DecoderConfig{
		DecodeHook:  variableMapDecode,
		Result:      config,
		ErrorUnused: true,
	}
	mapDecoder, err := mapstructure.NewDecoder(&mapDecoderConfig)
	if err != nil {
		return nil, err
	}

	if err = mapDecoder.Decode(m); err != nil {
		return nil, err
	}

	if config.Version != Version {
		return nil, fmt.Errorf("Config must be version %d", Version)
	}

	return config, nil
}

func (config *Config) DumpYAML() (string, error) {
	out, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
