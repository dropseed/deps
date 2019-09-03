package ci

import (
	"os"
	"os/exec"
	"strings"

	"github.com/dropseed/deps/internal/ci/circleci"
	"github.com/dropseed/deps/internal/ci/generic"
	"github.com/dropseed/deps/internal/ci/githubactions"
	"github.com/dropseed/deps/internal/ci/gitlabci"
	"github.com/dropseed/deps/internal/ci/travisci"
	"github.com/dropseed/deps/internal/output"
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
	if gitlabci.Is() {
		return &gitlabci.GitLabCI{}
	}
	return &generic.GenericCI{}
}

func BaseAutoconfigure() {
	if cmd := exec.Command("git", "config", "user.name", "deps"); cmd != nil {
		output.Event("Autoconfigure: %s", strings.Join(cmd.Args, " "))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	if cmd := exec.Command("git", "config", "user.email", "bot@dependencies.io"); cmd != nil {
		output.Event("Autoconfigure: %s", strings.Join(cmd.Args, " "))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
}
