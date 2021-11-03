package runner

import (
	"fmt"

	"github.com/dropseed/deps/internal/schemaext"

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
	if err := deps.Validate(); err != nil {
		panic(err)
	}

	updateID := schemaext.UpdateIDForDeps(deps)
	uniqueID := schemaext.UniqueIDForDeps(deps)
	branch := git.GetBranchName(fmt.Sprintf("%s-%s", updateID, uniqueID))

	update := Update{
		dependencies:     deps,
		dependencyConfig: cfg,
		id:               updateID,
		title:            schemaext.TitleForDeps(deps),
		branch:           branch,
	}

	return &update
}

func (update *Update) exists() bool {
	b := git.BranchMatching(update.branch)
	return b != ""
}

func (update *Update) outdatedBranch() string {
	// Assumes exists() was already checked for an exact match,
	// so this checks ANY branches that match the prefix ID
	return git.BranchMatching(update.branchPrefix())
}

func (update *Update) branchPrefix() string {
	return git.GetBranchName(update.id)
}
