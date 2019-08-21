package runner

import (
	"errors"
	"fmt"
	"strings"

	"github.com/dropseed/deps/internal/component"
	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/output"
)

const COLLECTOR = "collector"
const ACTOR = "actor"

func getConfig() (*config.Config, error) {
	configPath := config.FindFilename("", config.DefaultFilenames...)

	if configPath != "" {
		cfg, err := config.NewConfigFromPath(configPath)
		if err != nil {
			return nil, err
		}

		cfg.Compile()

		return cfg, nil
	}

	output.Event("No local config found, detecting your dependencies automatically")
	// should we always check for inferred? and could let them know what they
	// don't have in theirs?
	// dump both to yaml, use regular diff tool and highlighting?
	cfg, err := config.InferredConfigFromDir(".")
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

	cfg.Compile()

	if len(cfg.Dependencies) < 1 {
		return nil, errors.New("no dependencies found")
	}

	return cfg, nil
}

func organizeUpdates(updates Updates) (Updates, Updates, Updates, error) {
	newUpdates := Updates{}      // PRs for these
	outdatedUpdates := Updates{} // lockfile update on these?
	existingUpdates := Updates{}

	for _, update := range updates {
		if update.exists() {
			existingUpdates.addUpdate(update)
		} else if outdated := update.outdatedBranch(); outdated != "" {
			update.branch = outdated // change the branch to the existing match
			outdatedUpdates.addUpdate(update)
		} else {
			newUpdates.addUpdate(update)
		}
	}

	if len(outdatedUpdates) > 0 {
		fmt.Println()
		output.Event("%d outdated updates", len(outdatedUpdates))
		outdatedUpdates.printOverview()
	}

	if len(newUpdates) > 0 {
		fmt.Println()
		output.Event("%d new updates to be made", len(newUpdates))
		newUpdates.printOverview()
	}

	return newUpdates, outdatedUpdates, existingUpdates, nil
}

func collectUpdates(cfg *config.Config, types []string) (Updates, error) {
	if len(types) > 0 {
		output.Event("Only collecting types: %s", strings.Join(types, ", "))
	}
	typesMap := map[string]bool{}
	for _, t := range types {
		typesMap[t] = true
	}

	updates := Updates{}

	for index, dependencyConfig := range cfg.Dependencies {

		if _, ok := typesMap[dependencyConfig.Type]; len(typesMap) > 0 && !ok {
			continue
		}

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

		if err := runner.Install(); err != nil {
			return nil, err
		}

		dependencies, err := runner.Collect(dependencyConfig.Path)
		if err != nil {
			return nil, err
		}

		depUpdates, err := newUpdatesFromDependencies(dependencies, dependencyConfig)
		if err != nil {
			return nil, err
		}

		if len(depUpdates) > 0 {
			for _, update := range depUpdates {
				// Store this for use later
				update.runner = runner
				updates.addUpdate(update)
			}
		}
	}

	return updates, nil
}
