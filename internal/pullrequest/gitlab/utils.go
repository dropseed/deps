package gitlab

import (
	"errors"
	"fmt"
	"net/url"
	"os"
)

func getAPIToken() string {
	if s := os.Getenv("DEPS_GITLAB_TOKEN"); s != "" {
		return s
	}

	if s := os.Getenv("GITLAB_TOKEN"); s != "" {
		return s
	}

	// if s := os.Getenv("CI_JOB_TOKEN"); s != "" {
	// 	return s
	// }

	return ""
}

func getProjectAPIURL() (string, error) {
	if s := os.Getenv("DEPS_GITLAB_PROJECT_API_URL"); s != "" {
		return s, nil
	}

	// GitLab CI support
	if base := os.Getenv("CI_API_V4_URL"); base != "" {
		slug := url.PathEscape(os.Getenv("CI_PROJECT_PATH_SLUG"))
		return fmt.Sprintf("%s/projects/%s", base, slug), nil
	}

	// TODO otherwise from git remote?

	return "", errors.New("Unable to determine GitLab API url for this project")
}
