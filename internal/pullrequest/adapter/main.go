package adapter

import (
	"fmt"
	"strings"

	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/pullrequest/github"
	"github.com/dropseed/deps/internal/pullrequest/gitlab"
	"github.com/dropseed/deps/internal/pullrequest/gittest"
)

// PullrequestAdapter implements the basic Pullrequest functions
type PullrequestAdapter interface {
	Create() error
	DoRelated() error
	// OutputActions() error
}

type RepoAdapter interface {
	CheckRequirements() error
	PreparePush()
	// NewPullrequest(*schema.Dependencies, string) PullrequestAdapter
}

func NewRepoFromEnv() RepoAdapter {
	gitHost := git.GitHost()
	if gitHost == "github" {
		return github.NewRepoFromEnv()
	}
	return nil
}

// func NewPullrequestFromRepo(repo RepoAdapter, deps *schema.Dependencies, baseBranch string) PullrequestAdapter {
// 	return repo.NewPullrequest(deps, baseBranch)
// }

// PullrequestAdapterFromDependenciesJSONPathAndHost returns a host-specific Pullrequest struct
func PullrequestAdapterFromDependenciesJSONPathAndHost(dependenciesJSONPath, host, baseBranch string) (PullrequestAdapter, error) {
	switch strings.ToLower(host) {
	case "github":
		pr, err := github.NewPullrequestFromDependenciesJSONPathAndEnv(dependenciesJSONPath)
		pr.DefaultBaseBranch = baseBranch
		if err == nil {
			return pr, nil
		}
		return nil, err

	case "gitlab":
		pr, err := gitlab.NewPullrequestFromDependenciesJSONPathAndEnv(dependenciesJSONPath)
		pr.DefaultBaseBranch = baseBranch
		if err == nil {
			return pr, nil
		}
		return nil, err

	case "test":
		// in test env we will run a mock version of PR
		// behavior, so that we can further test the interaction
		pr, err := gittest.NewPullrequestFromDependenciesJSONPathAndEnv(dependenciesJSONPath)
		pr.DefaultBaseBranch = baseBranch
		if err == nil {
			return pr, nil
		}
		return nil, err

	}

	fmt.Printf("No pull request adapter for git host '%s'\n", host)
	return nil, nil
}
