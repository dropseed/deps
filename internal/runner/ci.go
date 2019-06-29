package runner

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/pullrequest/adapter"
)

type updateError struct {
	update *Update
	err    error
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

	startingBranch := getCurrentBranch()

	git.Checkout(startingBranch)

	// // TODO does this belong? or user responsibility (we can give instruction)
	if !git.CanPush() {
		repo.PreparePush()
	}

	updateErrors := []*updateError{}

	isDepsBranch := git.IsDepsBranch(startingBranch)

	// TODO better to pass this through collect or something
	manifestUpdatesDisabled = isDepsBranch

	newUpdates, existingUpdates, err := collectUpdates()
	if err != nil {
		return err
	}

	// TODO put specific update limit on new? or disable
	// same with existing
	if updateLimit > -1 {
		// newUpdates = newUpdates[:updateLimit]
	}

	// TODO this is also because collectors may have done some crap and not cleaned up
	output.Event("Temporarily saving your uncommitted changes in a git stash")
	stashed := git.Stash(fmt.Sprintf("Deps save before update"))

	// Stash pop needs to happen last (so be added first)
	defer func() {
		if stashed {
			output.Event("Putting original uncommitted changes back")
			if err := git.StashPop(); err != nil {
				output.Error("Error putting stash back: %v", err)
			}
		}
	}()

	if !isDepsBranch {
		output.Event("Performing %d new updates on %s", len(newUpdates), startingBranch)

		for _, update := range newUpdates {
			output.Event("Running update: %s", update.title)
			if err := runUpdate(update, startingBranch, update.branch); err != nil {
				updateErrors = append(updateErrors, &updateError{
					update: update,
					err:    err,
				})
				output.Error("Update failed: %v", err)
			} else {
				update.completed = true
			}
		}

		for _, update := range existingUpdates {
			output.Event("Checking existing update: %s", update.title)
			if err := runUpdate(update, update.branch, update.branch); err != nil {
				updateErrors = append(updateErrors, &updateError{
					update: update,
					err:    err,
				})
				output.Error("Update failed: %v", err)
			} else {
				update.completed = true
			}
		}
	} else {
		output.Event("Checking for updates to existing deps branch %s", startingBranch)
		var matchingExistingUpdate *Update
		for _, update := range existingUpdates {
			if update.branch == startingBranch {
				matchingExistingUpdate = update
			}
		}
		if matchingExistingUpdate != nil {
			output.Event("Applying latest matching update to this branch")
			if err := runUpdate(matchingExistingUpdate, matchingExistingUpdate.branch, matchingExistingUpdate.branch); err != nil {
				updateErrors = append(updateErrors, &updateError{
					update: matchingExistingUpdate,
					err:    err,
				})
				output.Error("Update failed: %v", err)
			} else {
				matchingExistingUpdate.completed = true
			}
		}
	}

	if len(updateErrors) > 0 {
		output.Error("There were %d errors making the updates", len(updateErrors))
		for _, ue := range updateErrors {
			output.Error("- [%s] %s\n  %v", ue.update.id, ue.update.title, ue.err)
		}
		return fmt.Errorf("%d errors", len(updateErrors))
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

func runUpdate(update *Update, base, head string) error {
	if base == head {
		output.Event("Running changes directly (no branches)")
		git.Checkout(head)
	} else {
		git.Checkout(base)
		git.Branch(head)
	}

	defer func() {
		// Theres should only be uncommitted changes if we're bailing early
		git.ResetAndClean()
		git.CheckoutLast()
	}()

	outputDeps, err := update.runner.Act(update.dependencies)
	if err != nil {
		return err
	}

	var pr adapter.PullrequestAdapter
	gitHost := git.GitHost()

	if base != head {
		// pr = repo.NewPullrequest(outputDeps, pullrequestToBranch)
		// TODO should pass head here too
		pr, err = adapter.PullrequestAdapterFromDependenciesAndHost(outputDeps, gitHost, base)
		if err != nil {
			return err
		}
	}

	title, err := outputDeps.GenerateTitle()
	if err != nil {
		return err
	}

	// if nothing to commit, don't worry about it
	if git.IsDirty() {
		git.AddCommit(title)
		// TODO try adding more lines for dependency breakdown,
		// especially on lockfiles
	}

	git.PushBranch(head)

	// TODO hooks or what do you do otherwise?

	if pr != nil {
		output.Debug("Waiting a second for the push to be processed by the host")
		time.Sleep(2 * time.Second)

		if err := pr.CreateOrUpdate(); err != nil {
			return err
		}
	}

	return nil
}
