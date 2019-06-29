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
func NewPullrequestFromEnv(deps *schema.Dependencies, branch string) (*Pullrequest, error) {
	return &Pullrequest{
		Branch:            branch,
		Title:             deps.Title,
		Body:              deps.Description,
		DefaultBaseBranch: "",
		Dependencies:      deps,
		Action:            &schema.Action{Metadata: map[string]interface{}{}},
	}, nil
}
