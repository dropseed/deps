package gitlab

import (
	"errors"
	"os"

	"github.com/dropseed/deps/internal/ci/gitlabci"
)

func getAPIToken() string {
	if s := os.Getenv("DEPS_GITLAB_TOKEN"); s != "" {
		return s
	}

	return ""
}

func getAPIUsername() string {
	if s := os.Getenv("DEPS_GITLAB_USERNAME"); s != "" {
		return s
	}

	return "gitlab-ci-token"
}

func getProjectAPIURL() (string, error) {
	if s := os.Getenv("DEPS_GITLAB_PROJECT_API_URL"); s != "" {
		return s, nil
	}

	if gitlabciURL := gitlabci.GetProjectAPIURL(); gitlabciURL != "" {
		return gitlabciURL, nil
	}

	// TODO otherwise from git remote?

	return "", errors.New("Unable to determine GitLab API url for this project")
}
