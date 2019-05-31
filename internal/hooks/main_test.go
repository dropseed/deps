package hooks

import (
	"os"
	"testing"
)

func TestHooksValid(t *testing.T) {
	os.Setenv("SETTING_TEST_HOOK", "[\"ls -a -l\", \"pwd\"]")
	err := Run("test_hook")
	if err != nil {
		t.Error(err)
	}
}

func TestHooksInvalidJSON(t *testing.T) {
	os.Setenv("SETTING_TEST_HOOK", "foo")
	err := Run("test_hook")
	if err == nil {
		t.FailNow()
	}
}

func TestHooksInvalidHook(t *testing.T) {
	os.Setenv("SETTING_TEST_HOOK", "[\"ls -a -l\", \"heywhoashit\"]")
	err := Run("test_hook")
	if err == nil {
		t.FailNow()
	}
}
