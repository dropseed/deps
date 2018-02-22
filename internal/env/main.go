package env

import (
	"fmt"
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

// GetSetting gets a string for a given setting name from the env
func GetSetting(name string, defaultValue string) string {
	v := os.Getenv(fmt.Sprintf("SETTING_%s", strings.ToUpper(name)))
	if v == "" && defaultValue != "" {
		return defaultValue
	}
	return v
}
