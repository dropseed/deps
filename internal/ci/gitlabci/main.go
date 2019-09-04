package gitlabci

import (
	"fmt"
	"net/url"
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

func GetProjectAPIURL() string {
	if base := os.Getenv("CI_API_V4_URL"); base != "" {
		path := url.PathEscape(os.Getenv("CI_PROJECT_PATH"))
		return fmt.Sprintf("%s/projects/%s", base, path)
	}
	return ""
}
