package pullrequest

import (
	"errors"

	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/pullrequest/github"
	"github.com/dropseed/deps/internal/schema"
)

// PullrequestAdapter implements the basic Pullrequest functions
type PullrequestAdapter interface {
	CreateOrUpdate() error
}

type RepoAdapter interface {
	CheckRequirements() error
	Autoconfigure()
	// NewPullrequest(*schema.Dependencies, string) PullrequestAdapter
}

func NewRepo() (RepoAdapter, error) {
	gitHost := gitHost()

	if gitHost == GITHUB {
		return github.NewRepoFromEnv()
	}

	return nil, errors.New("Repo not found or not supported")
}

func NewPullrequest(base string, head string, deps *schema.Dependencies, cfg *config.Dependency) (PullrequestAdapter, error) {
	gitHost := gitHost()

	if gitHost == GITHUB {
		return github.NewPullrequest(base, head, deps, cfg)
	}

	return nil, errors.New("Repo not found or not supported")
}
