package githubactions

import (
	"os"
	"strings"
)

type GitHubActions struct {
}

func Is() bool {
	return os.Getenv("GITHUB_ACTION") != ""
}

func (actions *GitHubActions) Autoconfigure() error {
	return nil
}

func (actions *GitHubActions) Branch() string {
	if b := os.Getenv("GITHUB_REF"); strings.HasPrefix(b, "refs/heads/") {
		return b[11:]
	}
	return ""
}
