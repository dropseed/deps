package gitlab

import (
	"errors"
	"os"
)

type GitLabRepo struct {
	apiToken string
}

func NewRepoFromEnv() (*GitLabRepo, error) {
	token := getAPIToken()

	if token == "" {
		return nil, errors.New("Unable to find GitLab API token.\n\nVisit https://docs.dependencies.io/gitlab for more information.")
	}

	return &GitLabRepo{
		apiToken: token,
	}, nil
}

func (repo *GitLabRepo) CheckRequirements() error {
	if repo.apiToken == "" {
		return errors.New("GitLab API token not found")
	}
	return nil
}

func (repo *GitLabRepo) Autoconfigure() {
}

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
