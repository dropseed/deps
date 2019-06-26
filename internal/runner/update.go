package runner

import (
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
}

func NewUpdate(deps *schema.Dependencies, cfg *config.Dependency) *Update {
	id := deps.GetUpdateID()
	title, err := deps.GenerateTitle()
	if err != nil {
		panic(err)
	}

	update := Update{
		dependencies:     deps,
		dependencyConfig: cfg,
		id:               id,
		title:            title,
	}

	return &update
}

func (update *Update) branchExists() bool {
	branch := update.dependencies.GetBranchName()
	return git.BranchExists(branch)
}
