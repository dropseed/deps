package runner

import (
	"errors"
	"os"
	"strings"

	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/output"
)

func CI(updateLimit int) error {
	// get Repo obj? required if running "CI" version
	// and can validate before proceeding?

	newUpdates, _, _, err := collectUpdates(updateLimit)
	if err != nil {
		return err
	}

	if len(newUpdates) > 0 {
		output.Event("Performing updates")
		if err := newUpdates.run(getBaseBranch()); err != nil {
			return err
		}
	} else {
		output.Success("No new updates")
	}

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

func (updates Updates) run(branch string) error {
	for _, update := range updates {
		if err := update.runner.Act(update.dependencies, branch); err != nil {
			return err
		}
		update.completed = true
	}
	return nil
}
