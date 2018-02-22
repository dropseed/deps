package env

import "testing"
import "os"

func TestProductionEnv(t *testing.T) {
	os.Setenv("DEPENDENCIES_ENV", "production")
	if !IsProduction() {
		t.Fail()
	}

	if GetCurrentEnv() != "production" {
		t.Fail()
	}
}

func TestTestEnv(t *testing.T) {
	os.Setenv("DEPENDENCIES_ENV", "test")
	if IsProduction() {
		t.Fail()
	}

	if GetCurrentEnv() != "test" {
		t.Fail()
	}
}

func TestGetSetting(t *testing.T) {
	os.Setenv("SETTING_FOO_BAR", "test")
	if s := GetSetting("foo_bar", ""); s != "test" {
		t.Fail()
	}

	if s := GetSetting("foo_barred", ""); s != "" {
		t.Fail()
	}
}
