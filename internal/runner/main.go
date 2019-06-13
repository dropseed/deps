package runner

import (
	"errors"
	"fmt"
	"os"

	"github.com/dropseed/deps/internal/git"

	"github.com/dropseed/deps/internal/component"
	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/output"
)

const COLLECTOR = "collector"
const ACTOR = "actor"

// Run a full interactive update process
func Local() error {
	baseBranch := ""

	cfg, err := getConfig()
	if err != nil {
		return err
	}

	availableUpdates, err := getAvailableUpdates(cfg)
	if err != nil {
		return err
	}

	// TODO for the updates that have branches: try lockfile update on them?
	// if no branch, act on update

	availableUpdates.PrintOverview()

	if err := availableUpdates.Prompt(baseBranch); err != nil {
		return err
	}

	return nil
}

func CI(updateLimit int) error {
	baseBranch := git.CurrentBranch()

	cfg, err := getConfig()
	if err != nil {
		return err
	}

	// get Repo obj? required if running "CI" version
	// and can validate before proceeding?

	availableUpdates, err := getAvailableUpdates(cfg)
	if err != nil {
		return err
	}

	newUpdates := Updates{}      // PRs for these
	limitedUpdates := Updates{}  // nothing
	existingUpdates := Updates{} // lockfile update on these?

	for _, update := range availableUpdates {
		if update.branchExists() {
			existingUpdates = append(existingUpdates, update)
		} else if updateLimit > -1 && len(newUpdates) >= updateLimit {
			limitedUpdates = append(limitedUpdates, update)
		} else {
			newUpdates = append(newUpdates, update)
		}
	}

	if len(existingUpdates) > 0 {
		output.Event("%d existing updates", len(existingUpdates))
		existingUpdates.PrintOverview()
		fmt.Println()
	}

	if len(limitedUpdates) > 0 {
		output.Event("%d updates skipped based on limit", len(limitedUpdates))
		limitedUpdates.PrintOverview()
		fmt.Println()
	}

	if len(newUpdates) > 0 {
		output.Event("%d new updates to be made", len(newUpdates))
		newUpdates.PrintOverview()
		fmt.Println()
	}

	if len(newUpdates) > 0 {
		output.Event("Performing updates")
		if err := newUpdates.Run(baseBranch); err != nil {
			return err
		}
	} else {
		output.Success("No new updates")
	}

	return nil
}

func getConfig() (*config.Config, error) {
	cfg, err := config.NewConfigFromPath(config.DefaultFilename, nil)
	if os.IsNotExist(err) {
		output.Event("No local config found, detecting your dependencies automatically")
		// should we always check for inferred? and could let them know what they
		// don't have in theirs?
		// dump both to yaml, use regular diff tool and highlighting?
		cfg, err = config.InferredConfigFromDir(".")
		if err != nil {
			return nil, err
		}

		inferred, err := cfg.DumpYAML()
		if err != nil {
			return nil, err
		}
		println("---")
		println(inferred)
		println("---")
	} else if err != nil {
		return nil, err
	}

	if len(cfg.Dependencies) < 1 {
		return nil, errors.New("no dependencies found")
	}

	cfg.Compile()

	return cfg, nil
}

func getAvailableUpdates(cfg *config.Config) (Updates, error) {
	availableUpdates := Updates{}

	for index, dependencyConfig := range cfg.Dependencies {

		runner, err := component.NewRunnerFromString(dependencyConfig.Type)
		if err != nil {
			return nil, err
		}
		env, err := dependencyConfig.Environ()
		if err != nil {
			return nil, err
		}
		runner.Index = index
		runner.Env = env

		// add a .shouldInstall - true when local or ref changed?

		err = runner.Install()
		if err != nil {
			return nil, err
		}

		dependencies, err := runner.Collect(dependencyConfig.Path)
		if err != nil {
			return nil, err
		}

		updates, err := NewUpdatesFromDependencies(dependencies, dependencyConfig)
		if err != nil {
			return nil, err
		}

		if len(updates) > 0 {
			for _, update := range updates {
				// Store this for use later
				update.runner = runner
			}
			availableUpdates = append(availableUpdates, updates...)
		}
	}

	return availableUpdates, nil
}
