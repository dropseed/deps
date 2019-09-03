package ci

import (
	"github.com/dropseed/deps/internal/ci/circleci"
	"github.com/dropseed/deps/internal/ci/githubactions"
	"github.com/dropseed/deps/internal/ci/travisci"
)

type CIProvider interface {
	Autoconfigure() error
	Branch() string
}

func NewCIProvider() CIProvider {
	if circleci.Is() {
		return &circleci.CircleCI{}
	}
	if travisci.Is() {
		return &travisci.TravisCI{}
	}
	if githubactions.Is() {
		return &githubactions.GitHubActions{}
	}
	return nil
}
