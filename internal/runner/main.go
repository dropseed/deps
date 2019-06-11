package runner

import (
	"errors"
	"os"

	"github.com/dropseed/deps/internal/component"
	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/output"
)

const COLLECTOR = "collector"
const ACTOR = "actor"

// Run a full interactive update process
func Local() error {

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

	branch := false
	if err := availableUpdates.Prompt(branch); err != nil {
		return err
	}

	return nil
}

func CI() error {
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

	branch := true
	if err := availableUpdates.Run(branch); err != nil {
		return err
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

		// overriding this is really just a debug/testing thing... shouldn't need to commit it?
		// so what's the best way to accomplish that...
		// cmd line flag or env var...?
		// or a separate patch type config file? so you can ignore it/trash it etc.
		// but it can have complex types?
		// dependencies_components.yml - still not easy in CI really

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
			availableUpdates[runner] = updates
		}
	}

	return availableUpdates, nil
}
