package gitlab

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

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
	remote := git.GitRemote()
	if strings.HasPrefix(remote, "https://gitlab-ci-token:") {
		parts := strings.SplitN(remote, "@", 2)
		keep := parts[1]
		updatedRemote := fmt.Sprintf("https://%s:%s@%s", repo.apiUsername, repo.apiToken, keep)
		maskedRemote := strings.Replace(updatedRemote, repo.apiToken, "*****", 1)
		if cmd := exec.Command("git", "remote", "set-url", "origin", updatedRemote); cmd != nil {
			output.Event("Autoconfigure: git remote set-url origin %s", maskedRemote)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				panic(err)
			}
		}
	}
}
