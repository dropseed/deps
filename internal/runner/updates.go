package runner

import (
	"fmt"

	"github.com/dropseed/deps/internal/component"
	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/schema"
	"github.com/manifoldco/promptui"
)

// Update contains relevant data for a potential dependency update
type Update struct {
	dependencies     *schema.Dependencies
	dependencyConfig *config.Dependency
	completed        bool
}

func (update *Update) shouldSkip() bool {
	branch, err := update.dependencies.GetBranchName()
	if err != nil {
		panic(err)
	}
	if git.BranchExists(branch) {
		output.Debug("Skipping update: branch %s already exists", branch)
		return true
	}
	return false
}

type Updates map[*component.Runner][]Update

func (updates Updates) PrintOverview() {
	if len(updates) < 1 {
		output.Success("No updates found")
	}

	for runner, runnerUpdates := range updates {
		plural := ""
		if len(runnerUpdates) != 1 {
			plural = "s"
		}
		output.Success("Found %d update%s for %s", len(runnerUpdates), plural, runner.GetName())
		for _, update := range runnerUpdates {
			id, err := update.dependencies.GetID()
			if err != nil {
				panic(err)
			}
			title, err := update.dependencies.GenerateTitle()
			if err != nil {
				panic(err)
			}
			fmt.Printf("[%s] %s\n", id, title)
		}
	}
}

func (availableUpdates Updates) Prompt(branch bool) error {
	for runner, updates := range availableUpdates {
		for {
			items := []string{}
			for _, update := range updates {
				if !update.completed {
					title, err := update.dependencies.GenerateTitle()
					if err != nil {
						return err
					}
					items = append(items, title)
				}
			}

			if len(items) < 1 {
				// No updates left
				break
			}

			items = append(items, "Skip")

			prompt := promptui.Select{
				Label: fmt.Sprintf("Choose an update to make on %s", runner.GetName()),
				Items: items,
			}

			i, _, err := prompt.Run()
			if err != nil {
				return err
			}

			if i < len(updates) {
				update := updates[i]
				if err := runner.Act(update.dependencies, branch); err != nil {
					return err
				}
				update.completed = true
			} else {
				// Chose skip
				break
			}
		}
	}

	return nil
}

func (availableUpdates Updates) Run(branch bool) error {
	for runner, updates := range availableUpdates {
		for _, update := range updates {
			if err := runner.Act(update.dependencies, branch); err != nil {
				return err
			}
			update.completed = true
		}
	}
	return nil
}

func NewUpdatesFromDependencies(dependencies *schema.Dependencies, dependencyConfig *config.Dependency) ([]Update, error) {
	updates := []Update{}

	if *dependencyConfig.LockfileUpdates.Enabled {
		for path, lockfile := range dependencies.Lockfiles {
			if lockfile.Updated == nil || len(lockfile.Updated.Dependencies) < 1 {
				continue
			}

			updateDependencies := schema.Dependencies{
				Lockfiles: map[string]*schema.Lockfile{
					path: lockfile,
				},
			}

			update := Update{
				dependencies:     &updateDependencies,
				dependencyConfig: dependencyConfig,
			}

			// TODO move this into runner later
			if update.shouldSkip() {
				continue
			}

			updates = append(updates, update)
		}
	}

	if *dependencyConfig.ManifestUpdates.Enabled {
		for path, manifest := range dependencies.Manifests {

			if manifest.Updated == nil || len(manifest.Updated.Dependencies) < 1 {
				continue
			}

			filteredGroups, err := dependencyConfig.ManifestUpdates.FilteredDependencyGroups(manifest.Updated.Dependencies)
			if err != nil {
				return nil, err
			}

			for _, groupDeps := range filteredGroups {

				updateDependencies := schema.Dependencies{
					Manifests: map[string]*schema.Manifest{
						path: &schema.Manifest{
							LockfilePath: manifest.LockfilePath,
							Current: &schema.ManifestVersion{
								Dependencies: map[string]*schema.ManifestDependency{},
							},
							Updated: &schema.ManifestVersion{
								Dependencies: map[string]*schema.ManifestDependency{},
							},
						},
					},
				}

				for name, dep := range groupDeps {
					updateDependencies.Manifests[path].Current.Dependencies[name] = manifest.Current.Dependencies[name]
					updateDependencies.Manifests[path].Updated.Dependencies[name] = dep
				}

				update := Update{
					dependencies:     &updateDependencies,
					dependencyConfig: dependencyConfig,
				}

				// TODO move this into runner later
				if update.shouldSkip() {
					continue
				}

				updates = append(updates, update)
			}
		}
	}

	return updates, nil
}
