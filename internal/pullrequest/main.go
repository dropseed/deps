package pullrequest

import (
	"errors"
	"os"
	"strings"

	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/pullrequest/bitbucket"
	"github.com/dropseed/deps/internal/pullrequest/github"
	"github.com/dropseed/deps/internal/pullrequest/gitlab"
	"github.com/dropseed/deps/pkg/schema"
)

const GITHUB = "github"
const GITLAB = "gitlab"
const BITBUCKET = "bitbucket"

// PullrequestAdapter implements the basic Pullrequest functions
type PullrequestAdapter interface {
	CreateOrUpdate() error
}

type RepoAdapter interface {
	CheckRequirements() error
	Autoconfigure()
	// NewPullrequest(*schema.Dependencies, string) PullrequestAdapter
}

func NewRepo() RepoAdapter {
	gitHost := gitHost()

	if gitHost == GITHUB {
		return github.NewRepo()
	}

	if gitHost == GITLAB {
		return gitlab.NewRepo()
	}

	if gitHost == BITBUCKET {
		return bitbucket.NewRepo()
	}

	return nil
}

func NewPullrequest(base string, head string, deps *schema.Dependencies, cfg *config.Dependency) (PullrequestAdapter, error) {
	gitHost := gitHost()

	if gitHost == GITHUB {
		return github.NewPullRequest(base, head, deps, cfg)
	}

	if gitHost == GITLAB {
		return gitlab.NewMergeRequest(base, head, deps, cfg)
	}

	if gitHost == BITBUCKET {
		return bitbucket.NewPullRequest(base, head, deps, cfg)
	}

	return nil, errors.New("Repo not found or not supported")
}

func gitHost() string {
	// or can maybe tell from github actions env var too or gitlab pipeline, but both should have remote as well
	if override := os.Getenv("DEPS_GIT_HOST"); override != "" {
		return override
	}

	remote := git.GitRemote()

	// TODO https://user:pass@

	if strings.HasPrefix(remote, "https://github.com/") || strings.HasPrefix(remote, "git@github.com:") {
		return GITHUB
	}

	if strings.HasPrefix(remote, "https://gitlab.com/") || strings.HasPrefix(remote, "git@gitlab.com:") {
		return GITLAB
	}

	if strings.HasPrefix(remote, "https://bitbucket.org") || strings.HasPrefix(remote, "git@bitbucket.org:") {
		return BITBUCKET
	}

	// More generic matching (github.example.com, etc. but could also accidently match gitlab.example.com/org/github-api)

	if strings.Contains(remote, "github") {
		return GITHUB
	}

	if strings.Contains(remote, "gitlab") {
		return GITLAB
	}

	if strings.Contains(remote, "bitbucket") {
		return BITBUCKET
	}

	return ""
}
