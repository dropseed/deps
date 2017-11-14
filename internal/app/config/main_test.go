package config

import "testing"

func TestNotLoaded(t *testing.T) {
	config := Config{}
	if config.EnvSettings != nil {
		t.Fail()
	}
	if config.Flags != nil {
		t.Fail()
	}
}

func TestLoadEnvSettings(t *testing.T) {
	config := Config{}
	config.LoadEnvSettings()
	if config.EnvSettings == nil {
		t.Fail()
	}
	if config.Flags != nil {
		t.Fail()
	}
}

func TestLoadFlags(t *testing.T) {
	config := Config{}
	config.LoadFlags()
	if config.EnvSettings != nil {
		t.Fail()
	}
	if config.Flags == nil {
		t.Fail()
	}
}
