package config

import (
	"regexp"

	"github.com/dropseed/deps/pkg/schema"
)

type ManifestUpdates struct {
	Enabled  *bool     `mapstructure:"enabled,omitempty" yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Settings Settings  `mapstructure:"settings,omitempty" yaml:"settings,omitempty" json:"settings"`
	Filters  []*Filter `mapstructure:"filters,omitempty" yaml:"filters,omitempty" json:"filters,omitempty"`
}

type Filter struct {
	Name     string   `mapstructure:"name" yaml:"name" json:"name"`
	Enabled  *bool    `mapstructure:"enabled,omitempty" yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Group    *bool    `mapstructure:"group,omitempty" yaml:"group,omitempty" json:"group,omitempty"`
	Settings Settings `mapstructure:"settings,omitempty" yaml:"settings,omitempty" json:"settings"`
}

func (manifestUpdates *ManifestUpdates) FilteredDependencyGroups(dependencies map[string]*schema.ManifestDependency) (map[string]map[string]*schema.ManifestDependency, error) {
	groups := map[string]map[string]*schema.ManifestDependency{}

	dependenciesSeen := map[string]bool{}

	for _, filter := range manifestUpdates.Filters {
		for name, dep := range dependencies {
			if filter.MatchesName(name) {

				if _, seen := dependenciesSeen[name]; seen {
					// dependency will only be grouped in
					// the first filter that matched it
					continue
				}

				dependenciesSeen[name] = true

				if *filter.Enabled {

					groupName := filter.Name

					if !*filter.Group {
						// make a new group just for this dependency
						// works because it is only seen once anyway
						// (won't work if multiple manifests in 1 collector though?)
						groupName = name
					}

					groupDeps := groups[groupName]
					if groupDeps == nil {
						groupDeps = map[string]*schema.ManifestDependency{}
					}
					groupDeps[name] = dep
					groups[groupName] = groupDeps

				}
			}
		}
	}

	return groups, nil
}

func (filter *Filter) MatchesName(name string) bool {
	nameRegex := regexp.MustCompile(filter.Name)
	return nameRegex.MatchString(name)
}

func (filter *Filter) MatchesEntireSchema(deps *schema.Dependencies) bool {
	// Manifest filters can't match if there are any lockfiles involved (unusual)
	if deps.HasLockfiles() {
		return false
	}

	// The filter has to match *all* deps in the schema to be applied
	if deps.HasManifests() {
		for _, manifest := range deps.Manifests {
			if manifest.HasUpdates() {
				for depName := range manifest.Updated.Dependencies {
					if !filter.MatchesName(depName) {
						return false
					}
				}
			}
		}

		return true
	}

	return false
}
