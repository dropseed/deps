package ci

import (
	"os"
	"os/exec"
	"strings"

	"github.com/dropseed/deps/internal/ci/bitbucketpipelines"
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
	if bitbucketpipelines.Is() {
		return &bitbucketpipelines.BitbucketPipelines{}
	}
	return &generic.GenericCI{}
}

func BaseAutoconfigure() {

	gitName := "deps"
	gitEmail := "bot@dependencies.io"

	if s := os.Getenv("DEPS_GIT_NAME"); s != "" {
		gitName = s
	}
	if s := os.Getenv("DEPS_GIT_EMAIL"); s != "" {
		gitEmail = s
	}

	if cmd := exec.Command("git", "config", "user.name", gitName); cmd != nil {
		output.Event("Autoconfigure: %s", strings.Join(cmd.Args, " "))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}

	if cmd := exec.Command("git", "config", "user.email", gitEmail); cmd != nil {
		output.Event("Autoconfigure: %s", strings.Join(cmd.Args, " "))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
}
