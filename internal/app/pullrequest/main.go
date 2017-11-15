package pullrequest

import "os"
import "github.com/dependencies-io/pullrequest/internal/app/config"

// Adapter implements the basic Pullrequest functions
type Adapter interface {
	Create()
	DoRelated()
}

// Pullrequest stores the basic data
type Pullrequest struct {
	Branch            string
	Title             string
	Body              string
	DefaultBaseBranch string
	Config            *config.Config
}

// NewPullrequestFromEnv creates a Pullrequest using env variables
func NewPullrequestFromEnv(branch string, title string, body string, config *config.Config) *Pullrequest {
	return &Pullrequest{
		Branch:            branch,
		Title:             title,
		Body:              body,
		DefaultBaseBranch: os.Getenv("GIT_BRANCH"),
		Config:            config,
	}
}
