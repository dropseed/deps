package config

import (
	"fmt"
	"regexp"

	"github.com/dependencies-io/deps/internal/schema"
	"github.com/dependencies-io/deps/internal/versionfilter"
)

type ManifestUpdates struct {
	Enabled *bool     `mapstructure:"enabled,omitempty" yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Filters []*Filter `mapstructure:"filters,omitempty" yaml:"filters,omitempty" json:"filters,omitempty"`
	// ConstraintPrefix string    `mapstructure:"constraint_prefix,omitempty" yaml:"constraint_prefix,omitempty" json:"constraint_prefix,omitempty"`
}

type Filter struct {
	Name     string `mapstructure:"name" yaml:"name" json:"name"`
	Enabled  *bool  `mapstructure:"enabled,omitempty" yaml:"enabled,omitempty" json:"enabled,omitempty"`
	Versions string `mapstructure:"versions,omitempty" yaml:"versions,omitempty" json:"versions,omitempty"`
	Group    *bool  `mapstructure:"group,omitempty" yaml:"group,omitempty" json:"group,omitempty"`
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

// func (manifestUpdates *ManifestUpdates) FilteredVersions(name string, current string, availableVersions []schema.Version) ([]schema.Version, error) {
// 	if len(availableVersions) < 1 {
// 		return []schema.Version{}, nil
// 	}
// 	for _, filter := range manifestUpdates.Filters {
// 		if filter.MatchesName(name) {
// 			if *filter.Enabled {

// 				available := []string{}
// 				for _, v := range availableVersions {
// 					available = append(available, v.Name)
// 				}
// 				availableStrings, err := filter.FilterVersions(name, current, available)
// 				if err != nil {
// 					return nil, err
// 				}

// 				// convert the strings back to their Version objects with content
// 				matchMap := utils.StringSliceToMap(availableStrings)
// 				filteredVersions := []schema.Version{}
// 				for _, v := range availableVersions {
// 					if ok := matchMap[v.Name]; ok {
// 						filteredVersions = append(filteredVersions, v)
// 					}
// 				}

// 				return filteredVersions, nil

// 			} else {
// 				return []schema.Version{}, nil
// 			}
// 		}
// 	}

// 	return []schema.Version{}, nil
// }

func (filter *Filter) MatchesName(name string) bool {
	nameRegex := regexp.MustCompile(filter.Name)
	return nameRegex.MatchString(name)
}

func (filter *Filter) FilterVersions(name string, current string, available []string) ([]string, error) {
	f := versionfilter.NewVersionFilter(filter.Versions)
	if f == nil {
		return nil, fmt.Errorf("unable to parse filter %s", filter.Versions)
	}

	return f.Matching(available, current)
}
