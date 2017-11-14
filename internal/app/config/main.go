package config

import (
	"errors"
)

// Config stores program-wide setttings
type Config struct {
	EnvSettings *EnvSettings
	Flags       *Flags
}

// LoadEnvSettings parses the flags and stores them on Config
func (config *Config) LoadEnvSettings() error {
	e := NewEnvSettingsFromEnv()
	config.EnvSettings = e
	return nil
}

// LoadFlags parses the flags and stores them on Config
func (config *Config) LoadFlags() error {
	f := ParseFlags()
	config.Flags = f
	return nil
}

// Validate will return any errors in the Config or its Flags
func (config *Config) Validate() error {
	if config.Flags.Branch == "" {
		return errors.New("the \"branch\" flag is required")
	}

	_, err := config.TitleFromConfig()
	if err != nil {
		return err
	}

	_, err = config.BodyFromConfig()
	if err != nil {
		return err
	}

	return nil
}
