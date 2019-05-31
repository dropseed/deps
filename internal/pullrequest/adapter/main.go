package adapter

import (
	"fmt"
	"strings"

	"github.com/dropseed/deps/internal/pullrequest/github"
	"github.com/dropseed/deps/internal/pullrequest/gitlab"
	"github.com/dropseed/deps/internal/pullrequest/gittest"
)

// PullrequestAdapter implements the basic Pullrequest functions
type PullrequestAdapter interface {
	Create() error
	DoRelated() error
	OutputActions() error
}

// PullrequestAdapterFromDependenciesJSONPathAndHost returns a host-specific Pullrequest struct
func PullrequestAdapterFromDependenciesJSONPathAndHost(dependenciesJSONPath, host string) (PullrequestAdapter, error) {
	switch strings.ToLower(host) {
	case "github":
		pr, err := github.NewPullrequestFromDependenciesJSONPathAndEnv(dependenciesJSONPath)
		if err == nil {
			return pr, nil
		}
		return nil, err

	case "gitlab":
		pr, err := gitlab.NewPullrequestFromDependenciesJSONPathAndEnv(dependenciesJSONPath)
		if err == nil {
			return pr, nil
		}
		return nil, err

	case "test":
		// in test env we will run a mock version of PR
		// behavior, so that we can further test the interaction
		pr, err := gittest.NewPullrequestFromDependenciesJSONPathAndEnv(dependenciesJSONPath)
		if err == nil {
			return pr, nil
		}
		return nil, err

	}

	fmt.Printf("No pull request adapter for git host '%s'\n", host)
	return nil, nil
}
