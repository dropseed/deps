package runner

import (
	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/pkg/schema"
)

type Updates map[string]*Update

func (updates Updates) add(deps *schema.Dependencies, cfg *config.Dependency) {
	update := NewUpdate(deps, cfg)
	updates.addUpdate(update)
}

func (updates Updates) addUpdate(update *Update) {
	updates[update.id] = update
}

func (updates Updates) removeUpdate(update *Update) {
	delete(updates, update.id)
}

func (updates Updates) printOverview() {
	if len(updates) < 1 {
		output.Success("No updates found")
	}

	for _, update := range updates {
		output.Event("[%s] %s", update.id, update.title)
	}
}

func newUpdatesFromDependencies(dependencies *schema.Dependencies, dependencyConfig *config.Dependency) (Updates, error) {
	updates := Updates{}

	if *dependencyConfig.LockfileUpdates.Enabled {
		output.Debug("Filtering lockfile updates")
		for path, lockfile := range dependencies.Lockfiles {
			// output.Debug("%s has updates: %t", path, lockfile.HasUpdates())
			if !lockfile.HasUpdates() {
				continue
			}

			// All lockfile updates are split out individually

			updateDependencies := schema.Dependencies{
				Lockfiles: map[string]*schema.Lockfile{
					path: lockfile,
				},
			}

			updates.add(&updateDependencies, dependencyConfig)
		}
	} else {
		output.Event("Lockfile updates disbled")
	}

	if *dependencyConfig.ManifestUpdates.Enabled {
		output.Debug("Filtering manifest updates")
		for path, manifest := range dependencies.Manifests {
			// output.Debug("%s has updates: %t", path, manifest.HasUpdates())
			if !manifest.HasUpdates() {
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

				updates.add(&updateDependencies, dependencyConfig)
			}
		}
	} else {
		output.Event("Manifest updates disbled")
	}

	return updates, nil
}
