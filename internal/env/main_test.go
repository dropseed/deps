package env

import "testing"
import "os"

func TestProductionEnv(t *testing.T) {
	os.Setenv("DEPENDENCIES_ENV", "production")
	if !IsProduction() {
		t.FailNow()
	}

	if GetCurrentEnv() != "production" {
		t.FailNow()
	}
}

func TestTestEnv(t *testing.T) {
	os.Setenv("DEPENDENCIES_ENV", "test")
	if IsProduction() {
		t.FailNow()
	}

	if GetCurrentEnv() != "test" {
		t.FailNow()
	}
}

func TestSettingVal(t *testing.T) {
	if s, _ := settingValToEnv(2); s != "2" {
		t.FailNow()
	}
}
