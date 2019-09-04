package gitlab

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/output"
)

type GitLabRepo struct {
	apiToken    string
	apiUsername string
}

func NewRepo() *GitLabRepo {
	return &GitLabRepo{
		apiToken:    getAPIToken(),
		apiUsername: getAPIUsername(),
	}
}

func (repo *GitLabRepo) CheckRequirements() error {
	if repo.apiToken == "" {
		return errors.New("Unable to find GitLab API token.\n\nVisit https://docs.dependencies.io/gitlab for more information.")
	}
	if repo.apiUsername == "" {
		return errors.New("Unable to find GitLab API username.\n\nVisit https://docs.dependencies.io/gitlab for more information.")
	}
	return nil
}

func (repo *GitLabRepo) Autoconfigure() {
	output.Debug("Writing GitLab token to ~/.netrc")
	hostname := git.GitRemoteHostname()
	echo := fmt.Sprintf("echo -e \"machine %s\n  login %s\n  password %s\" >> ~/.netrc", hostname, repo.apiUsername, repo.apiToken)
	cmd := exec.Command("sh", "-c", echo)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}
