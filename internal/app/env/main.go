package env

import (
	"os"
	"strings"
)

// GetCurrentEnv gets the current environment this is running in ("production", "test", etc.)
func GetCurrentEnv() string {
	return os.Getenv("DEPENDENCIES_ENV")
}

// IsProduction checks if this is a production env
func IsProduction() bool {
	return strings.ToLower(GetCurrentEnv()) == "production"
}
