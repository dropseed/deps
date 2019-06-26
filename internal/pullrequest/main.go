package pullrequest

import (
	"github.com/dropseed/deps/internal/schema"
)

// Pullrequest stores the basic data
type Pullrequest struct {
	Branch            string
	Title             string
	Body              string
	DefaultBaseBranch string
	Dependencies      *schema.Dependencies
	Action            *schema.Action
}

// NewPullrequestFromEnv creates a Pullrequest using env variables
func NewPullrequestFromEnv(deps *schema.Dependencies) (*Pullrequest, error) {
	branch := deps.GetBranchName()

	title, err := deps.GenerateTitle()
	if err != nil {
		return nil, err
	}

	body, err := deps.GenerateBody()
	if err != nil {
		return nil, err
	}

	return &Pullrequest{
		Branch:            branch,
		Title:             title,
		Body:              body,
		DefaultBaseBranch: "",
		Dependencies:      deps,
		Action:            &schema.Action{Metadata: map[string]interface{}{}},
	}, nil
}
