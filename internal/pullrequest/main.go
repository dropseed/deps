package pullrequest

import (
	"os"
	"strings"

	"github.com/dropseed/deps/internal/git"
)

// Pullrequest stores the basic data
// type Pullrequest struct {
// 	Base         string
// 	Head         string
// 	Title        string
// 	Body         string
// 	Dependencies *schema.Dependencies
// 	Config       *config.Dependency
// }

// // NewPullrequest creates a Pullrequest using env variables
// func NewPullrequest(base string, head string, deps *schema.Dependencies, cfg *config.Dependency) (*Pullrequest, error) {
// 	return &Pullrequest{
// 		Base:         base,
// 		Head:         head,
// 		Title:        deps.Title,
// 		Body:         deps.Description,
// 		Dependencies: deps,
// 		Config:       cfg,
// 	}, nil
// }

const GITHUB = "github"
const GITLAB = "gitlab"

func gitHost() string {
	// or can maybe tell from github actions env var too or gitlab pipeline, but both should have remote as well
	if override := os.Getenv("DEPS_GIT_HOST"); override != "" {
		return override
	}

	remote := git.GitRemote()

	// TODO https://user:pass@

	if strings.HasPrefix(remote, "https://github.com/") || strings.HasPrefix(remote, "git@github.com:") {
		return GITHUB
	}

	if strings.HasPrefix(remote, "https://gitlab.com/") || strings.HasPrefix(remote, "git@gitlab.com:") {
		return GITLAB
	}

	// More generic matching (github.example.com, etc. but could also accidently match gitlab.example.com/org/github-api)

	if strings.Contains(remote, "github") {
		return GITHUB
	}

	if strings.Contains(remote, "gitlab") {
		return GITLAB
	}

	return ""
}
