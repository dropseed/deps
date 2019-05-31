package component

import (
	"io"
	"os"
	"reflect"

	"github.com/algobardo/yaml"
	"github.com/mitchellh/mapstructure"
)

const DefaultFilename = "dependencies_component.yml"

type Config struct {
	Install string `mapstructure:"install" yaml:"install" json:"install"`
	Collect string `mapstructure:"collect" yaml:"collect" json:"collect"`
	Act     string `mapstructure:"act" yaml:"act" json:"act"`
}

func NewConfigFromPath(path string) (*Config, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	return NewConfigFromReader(f)
}

func NewConfigFromReader(reader io.Reader) (*Config, error) {
	temp := map[string]interface{}{}
	decoder := yaml.NewDecoder(reader)
	decoder.SetDefaultMapType(reflect.TypeOf(map[string]interface{}{}))
	if err := decoder.Decode(&temp); err != nil {
		return nil, err
	}

	return newConfigFromMap(temp)
}

func newConfigFromMap(m map[string]interface{}) (*Config, error) {
	config := &Config{}

	mapDecoderConfig := mapstructure.DecoderConfig{
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

	return config, nil
}
