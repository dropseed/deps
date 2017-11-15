package config

import "os"

// EnvSettings stores pullrequest settings available via SETTING_.*
type EnvSettings struct {
	PullrequestNotes  string
	RelatedPRBehavior string
	RelatedPRAuthor   string
}

// NewEnvSettingsFromEnv gets EnvSettings using env variables
func NewEnvSettingsFromEnv() *EnvSettings {
	// for debugging we can change the author that is used for the search
	prAuthor := os.Getenv("DEBUG_SETTING_RELATED_PR_AUTHOR")
	if prAuthor == "" {
		prAuthor = "app/dependencies"
	}
	return &EnvSettings{
		PullrequestNotes:  os.Getenv("SETTING_PULLREQUEST_NOTES"),
		RelatedPRBehavior: os.Getenv("SETTING_RELATED_PR_BEHAVIOR"),
		RelatedPRAuthor:   prAuthor,
	}
}
