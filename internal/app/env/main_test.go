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
