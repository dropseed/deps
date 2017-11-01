package config

import "testing"
import "os"

func TestNewConfigFromEnv(t *testing.T) {
	os.Setenv("DEPENDENCIES_ENV", "production")
	config := NewConfigFromEnv()
	if !config.IsProduction() {
		t.Fail()
	}
}
