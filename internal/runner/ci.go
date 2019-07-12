package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/pullrequest"
	"github.com/dropseed/deps/internal/pullrequest/github"
)

type updateResult struct {
	update *Update
	err    error
}

func CI(autoconfigure bool, types []string, updateLimit int) error {
	if git.IsDirty() {
		return errors.New("git status must be clean to run deps ci")
	}

	gitHost := git.GitHost()
	var repo pullrequest.RepoAdapter
	if gitHost == "github" {
		repo = github.NewRepoFromEnv()
	} else {
		return errors.New("Repo not found or not supported")
	}

	if err := repo.CheckRequirements(); err != nil {
		return err
	}

	if autoconfigure {
		if err := autoconfigureRepo(repo); err != nil {
			return err
		}
	}

	output.Debug("Fetching all branches so we can check for existing updates")
	git.Fetch()

	startingBranch := getCurrentBranch()

	git.Checkout(startingBranch)

	updateErrors := []*updateResult{}

	isDepsBranch := git.IsDepsBranch(startingBranch)

	// TODO better to pass this through collect or something
	manifestUpdatesDisabled = isDepsBranch

	cfg, err := getConfig()
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
				updateErrors = append(updateErrors, &updateResult{
					update: update,
					err:    err,
				})
				output.Error("Update failed: %v", err)
			}
		}

		for _, update := range outdatedUpdates {
			output.Event("Updating outdated update: %s", update.title)
			if err := runUpdate(update, update.branch, update.branch); err != nil {
				updateErrors = append(updateErrors, &updateResult{
					update: update,
					err:    err,
				})
				output.Error("Update failed: %v", err)
			}
		}
	} else {
		output.Event("Checking for updates to existing deps branch %s", startingBranch)
		var outdatedMatch *Update
		for _, update := range outdatedUpdates {
			if update.branch == startingBranch {
				outdatedMatch = update
			}
		}
		if outdatedMatch != nil {
			output.Event("Applying latest matching update to this branch")
			if err := runUpdate(outdatedMatch, outdatedMatch.branch, outdatedMatch.branch); err != nil {
				updateErrors = append(updateErrors, &updateResult{
					update: outdatedMatch,
					err:    err,
				})
				output.Error("Update failed: %v", err)
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

	// TODO also updates successful and print green
	// so summary is successful, skipped, and errored?

	return nil
}

func autoconfigureRepo(repo pullrequest.RepoAdapter) error {

	if cmd := exec.Command("git", "config", "user.name", "deps"); cmd != nil {
		output.Event("Autoconfigure: %s", strings.Join(cmd.Args, " "))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	if cmd := exec.Command("git", "config", "user.email", "bot@dependencies.io"); cmd != nil {
		output.Event("Autoconfigure: %s", strings.Join(cmd.Args, " "))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return err
		}
	}

	if circleci := os.Getenv("CIRCLECI"); circleci != "" {
		// CircleCI uses ssh clones by default,
		// so try to switch to https

		if cmd := exec.Command("git", "config", "--global", "--remove-section", "url.ssh://git@github.com"); cmd != nil {
			output.Event("Autoconfigure: %s", strings.Join(cmd.Args, " "))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Run() // Don't worry about an error
		}

		originalOrigin := git.GitRemote()
		if updatedOrigin := git.GitRemoteToHTTPS(originalOrigin); originalOrigin != updatedOrigin {
			if cmd := exec.Command("git", "remote", "set-url", "origin", updatedOrigin); cmd != nil {
				output.Event("Autoconfigure: %s", strings.Join(cmd.Args, " "))
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					return err
				}
			}
		}
	}

	repo.PreparePush()

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

	var pr pullrequest.PullrequestAdapter
	gitHost := git.GitHost()

	if gitHost == "github" {
		pr, err = github.NewPullrequest(base, head, outputDeps, update.dependencyConfig)
		if err != nil {
			return err
		}
	}

	if !git.IsDirty() {
		output.Event("No changes to commit, exiting update early")
		return nil
	}

	git.AddCommit(outputDeps.Title)
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
