package runner

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/pullrequest/github"
)

func CI(updateLimit int) error {
	branch := getCurrentBranch()

	startingRef := git.CurrentRef()

	output.Debug("Fetching all branches so we can check for existing updates")
	git.FetchAllBranches()

	output.Debug("Checking out the tip of the branch (no point in looking at old commits for updates)")
	git.Checkout(branch)

	if !git.CanPush() {
		preparePush()
	}

	if git.IsDepsBranch(branch) {
		output.Event("Deps branch detected: running lockfile updates directly on this branch")
		branch = ""
		manifestUpdatesDisabled = true
	}

	newUpdates, _, _, err := collectUpdates(updateLimit)
	if err != nil {
		return err
	}

	if len(newUpdates) > 0 {
		output.Event("Performing updates")
		if err := newUpdates.run(branch, true); err != nil {
			return err
		}
	} else {
		output.Success("No new updates")
	}

	output.Debug("Attempting to put git back in the original state")
	git.ResetAndClean()
	git.Checkout(startingRef)

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

func (updates Updates) run(branch string, commitPush bool) error {
	for _, update := range updates {
		if err := update.runner.Act(update.dependencies, branch, commitPush); err != nil {
			return err
		}
		update.completed = true
	}
	return nil
}

func preparePush() {
	gitHost := git.GitHost()

	if gitHost == "github" {
		token := github.GetAPIToken()
		output.Debug("Writing GitHub token to ~/.netrc")
		echo := fmt.Sprintf("echo -e \"machine github.com\n  login x-access-token\n  password %s\" >> ~/.netrc", token)
		cmd := exec.Command("sh", "-c", echo)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			panic(err)
		}
	}
}
