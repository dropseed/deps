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

func TestRelatedPRBehavior(t *testing.T) {
	os.Setenv("SETTING_RELATED_PR_BEHAVIOR", "close")
	e := NewEnvSettingsFromEnv()
	if e.RelatedPRBehavior != "close" {
		t.Error("RelatedPRBehavior notes don't match expected")
	}
}

func TestRelatedPRAuthorDefault(t *testing.T) {
	e := NewEnvSettingsFromEnv()
	if e.RelatedPRAuthor != "app/dependencies" {
		t.Error("RelatedPRAuthor notes don't match expected")
	}
}

func TestRelatedPRAuthorDebug(t *testing.T) {
	os.Setenv("DEBUG_SETTING_RELATED_PR_AUTHOR", "davegaeddert")
	e := NewEnvSettingsFromEnv()
	if e.RelatedPRAuthor != "davegaeddert" {
		t.Error("RelatedPRAuthor notes don't match expected")
	}
}
