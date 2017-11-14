package config

import "os"

// EnvSettings stores pullrequest settings available via SETTING_.*
type EnvSettings struct {
	PullrequestNotes string
}

// NewEnvSettingsFromEnv gets EnvSettings using env variables
func NewEnvSettingsFromEnv() *EnvSettings {
	return &EnvSettings{
		PullrequestNotes: os.Getenv("SETTING_PULLREQUEST_NOTES"),
	}
}
