package config

import (
	"regexp"

	"github.com/dropseed/deps/internal/schema"
)

type ManifestUpdates struct {
	Enabled *bool     `mapstructure:"enabled,omitempty" yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Filters []*Filter `mapstructure:"filters,omitempty" yaml:"filters,omitempty" json:"filters,omitempty"`
	// ConstraintPrefix string    `mapstructure:"constraint_prefix,omitempty" yaml:"constraint_prefix,omitempty" json:"constraint_prefix,omitempty"`
}

type Filter struct {
	Name    string `mapstructure:"name" yaml:"name" json:"name"`
	Enabled *bool  `mapstructure:"enabled,omitempty" yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Group   *bool  `mapstructure:"group,omitempty" yaml:"group,omitempty" json:"group,omitempty"`
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
