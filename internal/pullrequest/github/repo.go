package github

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/dropseed/deps/internal/output"
)

type GitHubRepo struct {
	apiToken string
}

func NewRepoFromEnv() *GitHubRepo {
	return &GitHubRepo{
		apiToken: GetAPIToken(),
	}
}

func (repo *GitHubRepo) CheckRequirements() error {
	if repo.apiToken == "" {
		return errors.New("GitHub API token not found")
	}
	return nil
}

func (repo *GitHubRepo) PreparePush() {
	// switch remote to https
	// remove global config
	// 	git config --global --remove-section url."ssh://git@github.com"
	// get-url then replace ssh with https
	// git remote set-url origin https://github.com/$CIRCLE_PROJECT_USERNAME/$CIRCLE_PROJECT_REPONAME

	output.Debug("Writing GitHub token to ~/.netrc")
	echo := fmt.Sprintf("echo -e \"machine github.com\n  login x-access-token\n  password %s\" >> ~/.netrc", repo.apiToken)
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
