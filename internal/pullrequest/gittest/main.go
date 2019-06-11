package gittest

import (
	"fmt"

	"github.com/dropseed/deps/internal/pullrequest"
)

// PullRequest stores additional GitLab specific data
type PullRequest struct {
	// directly use the properties of base Pullrequest
	*pullrequest.Pullrequest
}

// NewPullrequestFromDependenciesJSONPathAndEnv creates a PullRequest
func NewPullrequestFromDependenciesJSONPathAndEnv(dependenciesJSONPath string) (*PullRequest, error) {
	prBase, err := pullrequest.NewPullrequestFromJSONPathAndEnv(dependenciesJSONPath)
	if err != nil {
		return nil, err
	}

	return &PullRequest{
		Pullrequest: prBase,
	}, nil
}

func (pr *PullRequest) PreparePush() error {
	return nil
}

// Create will mock the create op
func (pr *PullRequest) Create() error {
	fmt.Println("Doing mock create PR behavior...success!")
	pr.Action.Name = "PR #0"
	pr.Action.Metadata["foo"] = "bar"
	return nil
}

// DoRelated will mock the related PR op
func (pr *PullRequest) DoRelated() error {
	fmt.Println("Doing mock related PR behavior...success!")
	return nil
}
