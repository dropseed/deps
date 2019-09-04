package github

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/output"
)

type GitHubRepo struct {
	apiToken string
}

func NewRepo() *GitHubRepo {
	return &GitHubRepo{
		apiToken: getAPIToken(),
	}
}

func (repo *GitHubRepo) CheckRequirements() error {
	if repo.apiToken == "" {
		return errors.New("Unable to find GitHub API token.\n\nVisit https://docs.dependencies.io/github for more information.")
	}
	return nil
}

func (repo *GitHubRepo) Autoconfigure() {
	output.Debug("Writing GitHub token to ~/.netrc")
	hostname := git.GitRemoteHostname()
	echo := fmt.Sprintf("echo -e \"machine %s\n  login x-access-token\n  password %s\" >> ~/.netrc", hostname, repo.apiToken)
	cmd := exec.Command("sh", "-c", echo)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		panic(err)
	}
}

// func (repo *GitHubRepo) NewPullrequest(deps *schema.Dependencies, baseBranch string) *PullRequest {
// 	prBase, err := pullrequest.NewPullrequest(deps)
// 	if err != nil {
// 		panic(err)
// 	}
// 	prBase.DefaultBaseBranch = baseBranch

// 	fullName, err := getRepoFullName()
// 	if err != nil {
// 		panic(err)
// 	}
// 	parts := strings.Split(fullName, "/")
// 	owner := parts[0]
// 	repoName := parts[1]

// 	return &PullRequest{
// 		Pullrequest:   prBase,
// 		RepoOwnerName: owner,
// 		RepoName:      repoName,
// 		RepoFullName:  fullName,
// 		APIToken:      GetAPIToken(),
// 	}
// }
