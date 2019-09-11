package runner

import (
	"fmt"

	"github.com/dropseed/deps/internal/component"
	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/pkg/schema"
)

// Update contains relevant data for a potential dependency update
type Update struct {
	dependencies     *schema.Dependencies
	dependencyConfig *config.Dependency
	completed        bool
	runner           *component.Runner
	id               string
	title            string
	branch           string
}

func NewUpdate(deps *schema.Dependencies, cfg *config.Dependency) *Update {
	if err := deps.ValidateAndCompile(); err != nil {
		panic(err)
	}

	branch := git.GetBranchName(fmt.Sprintf("%s-%s", deps.UpdateID, deps.UniqueID))

	update := Update{
		dependencies:     deps,
		dependencyConfig: cfg,
		id:               deps.UpdateID,
		title:            deps.Title,
		branch:           branch,
	}

	return &update
}

func (update *Update) exists() bool {
	b := git.BranchMatching(update.branch)
	return b != ""
}

func (update *Update) outdatedBranch() string {
	// update id match only
	return git.BranchMatching(git.GetBranchName(update.id))
}
