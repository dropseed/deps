package runner

import (
	"errors"
	"os"
	"strings"

	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/pullrequest/adapter"
)

type branchUpdate struct {
	base                    string
	checkout                string
	manifestUpdatesDisabled bool
}

func CI(updateLimit int) error {
	if git.IsDirty() {
		return errors.New("git status must be clean to run deps ci")
	}

	repo := adapter.NewRepoFromEnv()
	if repo == nil {
		return errors.New("Repo not found or not supported")
	}
	if err := repo.CheckRequirements(); err != nil {
		return err
	}

	output.Debug("Fetching all branches so we can check for existing updates")
	git.FetchAllBranches()

	branches := []branchUpdate{}

	startingBranch := getCurrentBranch()
	startingRef := git.CurrentRef()

	git.Checkout(startingBranch)
	if !git.CanPush() {
		repo.PreparePush()
	}

	if git.IsDepsBranch(startingBranch) {
		output.Event("Deps branch detected: running lockfile updates directly on this branch")
		branches = append(branches, branchUpdate{
			base:                    "",
			checkout:                startingBranch,
			manifestUpdatesDisabled: true,
		})
	} else {
		branches = append(branches, branchUpdate{
			base:                    startingBranch,
			checkout:                startingBranch,
			manifestUpdatesDisabled: false,
		})
		// Run lockfile updates on any existing deps branches
		for _, branch := range git.GetDepsBranches() {
			branches = append(branches, branchUpdate{
				base:                    "",
				checkout:                branch,
				manifestUpdatesDisabled: true,
			})
			// TODO if lockfile is updated in a branch,
			// then the PR description may need to be updated too?
			// title should be good
		}
	}

	updateErrors := []struct {
		update *Update
		err    error
	}{}

	for _, branch := range branches {
		output.Debug("Checking out the tip of %s branch", branch.checkout)
		git.Checkout(branch.checkout)

		// TODO not a great pattern here?
		manifestUpdatesDisabled = branch.manifestUpdatesDisabled

		// Limit new updates on the main branch only
		limit := -1
		if branch.base == startingBranch {
			limit = updateLimit
		}

		newUpdates, _, _, err := collectUpdates(limit)
		if err != nil {
			return err
		}

		if len(newUpdates) > 0 {
			output.Event("Performing updates on %s", branch)
			for _, update := range newUpdates {
				output.Event("Running update: %s", update.title)
				if err := update.runner.Act(update.dependencies, branch.base, true); err != nil {
					updateErrors = append(updateErrors, struct {
						update *Update
						err    error
					}{
						update: update,
						err:    err,
					})
					output.Error("Update failed: %v", err)
				} else {
					update.completed = true
				}
			}
		} else {
			output.Success("No new updates on %s", branch)
		}

		output.Debug("Attempting to put git back in the original state")
		git.ResetAndClean()
	}

	git.Checkout(startingRef)

	if len(updateErrors) > 0 {
		output.Error("There were %d errors making the updates", len(updateErrors))
		for _, ue := range updateErrors {
			output.Error("- [%s] %s\n  %v", ue.update.id, ue.update.title, ue.err)
		}
	}

	return nil
}

func getCurrentBranch() string {
	branch := git.CurrentRef()

	// CI environments may be checking out a specific ref,
	// so use the variables they provide to see if we get a different branch name
	if b := os.Getenv("TRAVIS_PULL_REQUEST_BRANCH"); b != "" {
		branch = b
	}
	if b := os.Getenv("TRAVIS_BRANCH"); b != "" {
		branch = b
	}
	if b := os.Getenv("CIRCLE_BRANCH"); b != "" {
		branch = b
	}
	if b := os.Getenv("GITHUB_REF"); strings.HasPrefix(b, "refs/heads/") {
		branch = b[11:]
	}

	if branch == "" {
		panic(errors.New("Unable to determine base branch"))
	}

	return branch
}
