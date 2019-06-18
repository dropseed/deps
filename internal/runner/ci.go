package runner

import (
	"errors"
	"os"
	"strings"

	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/output"
)

func CI(updateLimit int, ifEnv string) error {
	branch := getBaseBranch()
	if branch == "master" {
		if ifEnv != "" && !checkEnvCondition(ifEnv) {
			output.Success("Skipping deps on master branch based on if env condition %s", ifEnv)
			return nil
		}
	} else if git.IsDepsBranch(branch) {
		output.Success("Running deps on an existing deps update branch")
	} else {
		output.Success("Skipping deps on branch %s", branch)
		return nil
	}

	startingRef := git.CurrentRef()

	output.Debug("Fetching all branches so we can check for existing updates")
	git.FetchAllBranches()

	output.Debug("Checking out the tip of the branch (no point in looking at old commits for updates)")
	git.Checkout(branch)

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

func getBaseBranch() string {
	baseBranch := git.CurrentRef()

	// CI environments may be checking out a specific ref,
	// so use the variables they provide to see if we get a different branch name
	if b := os.Getenv("TRAVIS_BRANCH"); b != "" {
		baseBranch = b
	}
	if b := os.Getenv("CIRCLE_BRANCH"); b != "" {
		baseBranch = b
	}
	if b := os.Getenv("GITHUB_REF"); strings.HasPrefix(b, "refs/heads/") {
		baseBranch = b[11:]
	}

	if baseBranch == "" {
		panic(errors.New("Unable to determine base branch"))
	}

	return baseBranch
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

func checkEnvCondition(s string) bool {
	parts := strings.SplitN(s, "=", 2)
	envKey := parts[0]
	envVal := ""
	if len(parts) > 1 {
		envVal = parts[1]
	}

	actualVal := os.Getenv(envKey)

	if envVal == "" && actualVal != "" {
		return false
	}

	if envVal != actualVal {
		return false
	}

	return true
}
