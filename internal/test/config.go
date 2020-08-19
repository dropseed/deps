package test

import (
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/dropseed/deps/internal/config"
	"gopkg.in/yaml.v2"
)

var directoryNamesToSkip = map[string]bool{
	".git":         true,
	"node_modules": true,
	"env":          true,
	"vendor":       true,
}

type Config struct {
	Tests []*Test `yaml:"tests"`
	path  string
}

func (c *Config) compile() {
	for _, test := range c.Tests {
		test.config = c
		if test.UserConfig == nil {
			test.UserConfig = &config.Dependency{}
		}
		// Set default data paths
		if test.Collect == nil {
			test.Collect = &TestPhase{
				Input:  test.Data,
				Output: test.Data,
			}
		}
		if test.Act == nil {
			test.Act = &TestPhase{
				Input:  test.Data,
				Output: test.Data,
			}
		}
		test.UserConfig.Compile()
	}
}

func (c *Config) joinPath(s string) string {
	return path.Join(path.Dir(c.path), s)
}

type Test struct {
	Name       string             `yaml:"name"`
	Repo       string             `yaml:"repo"`
	UserConfig *config.Dependency `yaml:"user_config,omitempty"`
	config     *Config
	Diff       string     `yaml:"diff,omitempty"`
	DiffArgs   []string   `yaml:"diff_args,omitempty"`
	Collect    *TestPhase `yaml:"collect,omitempty"`
	Act        *TestPhase `yaml:"act,omitempty"`
	Data       string     `yaml:"data"`
}

type TestPhase struct {
	Input  string `yaml:"input,omitempty"`
	Output string `yaml:"output"`
}

func (t *Test) displayName() string {
	return t.Name
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
	cfg := &Config{}
	decoder := yaml.NewDecoder(reader)
	decoder.SetDefaultMapType(reflect.TypeOf(map[string]interface{}{}))
	if err := decoder.Decode(cfg); err != nil {
		return nil, err
	}

	cfg.compile()

	return cfg, nil
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
	return strings.HasPrefix(name, "deps_test") && (strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml"))
}
