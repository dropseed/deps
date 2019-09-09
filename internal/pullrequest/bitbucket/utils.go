package bitbucket

import (
	"errors"
	"os"

	"github.com/dropseed/deps/internal/ci/bitbucketpipelines"
)

func getAPIPassword() string {
	if s := os.Getenv("DEPS_BITBUCKET_PASSWORD"); s != "" {
		return s
	}

	return ""
}

func getAPIUsername() string {
	if s := os.Getenv("DEPS_BITBUCKET_USERNAME"); s != "" {
		return s
	}

	return ""
}

func getProjectAPIURL() (string, error) {
	if s := os.Getenv("DEPS_BITBUCKET_REPO_API_URL"); s != "" {
		return s, nil
	}

	if ciURL := bitbucketpipelines.GetProjectAPIURL(); ciURL != "" {
		return ciURL, nil
	}

	// TODO otherwise from git remote?

	return "", errors.New("Unable to determine Bitbucket API url for this project")
}
