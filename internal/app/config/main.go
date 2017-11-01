package config

import "os"
import "strings"

// Config stores program-wide setttings
type Config struct {
	Env string
}

// NewConfigFromEnv creates a new Config using env variables
func NewConfigFromEnv() *Config {
	return &Config{
		Env: os.Getenv("DEPENDENCIES_ENV"),
	}
}

// IsProduction checks if this is a production env
func (config Config) IsProduction() bool {
	return strings.ToLower(config.Env) == "production"
}
