package gitlabci

import (
	"os"
)

type GitLabCI struct {
}

func Is() bool {
	return os.Getenv("GITLAB_CI") != ""
}

func (gitlab *GitLabCI) Autoconfigure() error {
	return nil
}

func (gitlab *GitLabCI) Branch() string {
	if b := os.Getenv("CI_COMMIT_REF_NAME"); b != "" {
		return b
	}
	return ""
}
