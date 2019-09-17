package runner

import (
	"errors"
	"fmt"
	"time"

	"github.com/dropseed/deps/internal/billing"
	"github.com/dropseed/deps/internal/ci"
	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/pullrequest"
)

type updateResult struct {
	update *Update
	err    error
}

func CI(autoconfigure bool, types []string) error {

	api, err := billing.NewAPI()
	if err != nil {
		return err
	}

	if err := api.Validate(); err != nil {
		return err
	}

	if git.IsDirty() {
		return errors.New("git status must be clean to run deps ci")
	}

	repo := pullrequest.NewRepo()
	if repo == nil {
		return errors.New("Repo not found or not supported")
	}
	ciProvider := ci.NewCIProvider()

	if err := repo.CheckRequirements(); err != nil {
		return err
	}

	if autoconfigure {
		ci.BaseAutoconfigure()

		if err := ciProvider.Autoconfigure(); err != nil {
			return err
		}

		repo.Autoconfigure()
	}

	output.Debug("Fetching all branches so we can check for existing updates")
	git.Fetch()

	startingBranch := getCurrentBranch(ciProvider)

	git.Checkout(startingBranch)

	successfulUpdates := []*updateResult{}
	failedUpdates := []*updateResult{}

	if isDepsBranch := git.IsDepsBranch(startingBranch); isDepsBranch {
		return errors.New("You cannot run deps ci on a deps branch")
	}

	cfg, err := config.FindOrInfer()
	if err != nil {
		return err
	}

	allUpdates, err := collectUpdates(cfg, types)
	if err != nil {
		return err
	}

	output.Debug("%d collected updates", len(allUpdates))

	newUpdates, outdatedUpdates, existingUpdates, err := organizeUpdates(allUpdates)
	if err != nil {
		return err
	}

	output.Event("%d new updates", len(newUpdates))
	output.Event("%d outdated updates", len(outdatedUpdates))
	output.Event("%d existing updates", len(existingUpdates))

	// TODO this is also because collectors may have done some crap and not cleaned up
	if git.IsDirty() {
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
	}

	output.Event("Performing %d new updates on %s", len(newUpdates), startingBranch)

	for _, update := range newUpdates {
		output.Event("Running update: %s", update.title)
		if err := runUpdate(update, startingBranch, update.branch); err != nil {
			failedUpdates = append(failedUpdates, &updateResult{
				update: update,
				err:    err,
			})
			output.Error("Update failed: %v", err)
		} else {
			successfulUpdates = append(successfulUpdates, &updateResult{
				update: update,
				err:    err,
			})
			output.Success("Update succeeded: %v", update.title)
		}
	}

	for _, update := range outdatedUpdates {
		output.Event("Updating outdated update: %s", update.title)
		if err := runUpdate(update, update.branch, update.branch); err != nil {
			failedUpdates = append(failedUpdates, &updateResult{
				update: update,
				err:    err,
			})
			output.Error("Update failed: %v", err)
		} else {
			successfulUpdates = append(successfulUpdates, &updateResult{
				update: update,
				err:    err,
			})
			output.Success("Update succeeded: %v", update.title)
		}
	}

	if len(successfulUpdates) > 0 {
		output.Success("%d updates made successfully!", len(successfulUpdates))
		for _, ue := range successfulUpdates {
			output.Error("- [%s] %s", ue.update.id, ue.update.title)
		}
	}

	if len(failedUpdates) > 0 {
		output.Error("There were %d errors making the updates", len(failedUpdates))
		for _, ue := range failedUpdates {
			output.Error("- [%s] %s\n  %v", ue.update.id, ue.update.title, ue.err)
		}
	}

	if len(successfulUpdates) > 0 {
		if err := api.IncrementUsage(len(successfulUpdates)); err != nil {
			return err
		}
	}

	if len(failedUpdates) > 0 {
		return fmt.Errorf("%d errors", len(failedUpdates))
	}

	return nil
}

func getCurrentBranch(ci ci.CIProvider) string {
	branch := git.CurrentRef()

	// CI environments may be checking out a specific ref,
	// so use the variables they provide to see if we get a different branch name

	if b := ci.Branch(); b != "" {
		branch = b
	}

	if branch == "" {
		panic(errors.New("Unable to determine base branch"))
	}

	if branch == "HEAD" {
		panic(errors.New("Unable to determine base branch, only got HEAD"))
	}

	return branch
}

func runUpdate(update *Update, base, head string) error {
	git.Checkout(base)

	if base == head {
		// PR back to the main branch
		// (setting or env var?)
		base = "master"
	} else {
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

	pr, err := pullrequest.NewPullrequest(base, head, outputDeps, update.dependencyConfig)
	if err != nil {
		return err
	}

	if !git.IsDirty() {
		return errors.New("Update didn't generate any changes to commit")
	}

	git.Add()
	git.Commit(outputDeps.Title)
	// TODO try adding more lines for dependency breakdown,
	// especially on lockfiles

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
