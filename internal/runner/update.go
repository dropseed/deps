package runner

import (
	"fmt"

	"github.com/dropseed/deps/internal/component"
	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/schema"
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
	return git.BranchExists(update.branch)
}

func (update *Update) outdated() bool {
	// update id match only
	return git.BranchExists(git.GetBranchName(update.id))
}
