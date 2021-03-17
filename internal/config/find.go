package config

import (
	"errors"
	"os"
	"path"

	"github.com/dropseed/deps/internal/output"
)

func FindFilename(dir string, filenames ...string) string {
	for _, f := range filenames {
		p := path.Join(dir, f)
		if fileExists(p) {
			return p
		}
	}
	return ""
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FindOrInfer() (*Config, error) {
	configPath := FindFilename("", DefaultFilenames...)

	if configPath != "" {
		cfg, err := NewConfigFromPath(configPath)
		if err != nil {
			return nil, err
		}

		cfg.Compile()

		return cfg, nil
	}

	output.Event("No local config found, detecting your dependencies automatically")
	// should we always check for inferred? and could let them know what they
	// don't have in theirs?
	// dump both to yaml, use regular diff tool and highlighting?
	cfg, err := InferredConfigFromDir(".")
	if err != nil {
		return nil, err
	}

	inferred, err := cfg.DumpYAML()
	if err != nil {
		return nil, err
	}
	output.Unstyled("---")
	output.Unstyled(inferred)
	output.Unstyled("---")

	cfg.Compile()

	if len(cfg.Dependencies) < 1 {
		return nil, errors.New("no dependencies found")
	}

	return cfg, nil
}
