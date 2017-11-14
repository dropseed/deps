package config

import (
	"os"
	"testing"
)

func TestPullrequestNotes(t *testing.T) {
	os.Setenv("SETTING_PULLREQUEST_NOTES", "notes test")
	e := NewEnvSettingsFromEnv()
	if e.PullrequestNotes != "notes test" {
		t.Error("Pullrequest notes don't match expected")
	}
}
