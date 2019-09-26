package schemaext

import (
	"fmt"
	"sort"

	"github.com/dropseed/deps/internal/changelogs"
	"github.com/dropseed/deps/pkg/schema"
)

func getSummaryLinesForManifests(manifests map[string]*schema.Manifest) ([]string, error) {
	summaries := make([]string, 0, len(manifests))

	// iterate using the sorted keys instead of unpredictable map
	keys := []string{}
	for path := range manifests {
		keys = append(keys, path)
	}
	sort.Strings(keys)

	for _, manifestPath := range keys {
		manifest := manifests[manifestPath]
		// summary, err := manifest.GetSummaryLine(manifestPath)
		// if err != nil {
		// 	return nil, err
		// }
		// summaries = append(summaries, summary)
		// iterate using the sorted keys instead of unpredictable map
		keys := []string{}
		for name := range manifest.Updated.Dependencies {
			keys = append(keys, name)
		}
		sort.Strings(keys)

		for _, dependencyName := range keys {
			s, err := getSummaryLineForDependencyName(manifest, dependencyName, manifestPath)
			if err != nil {
				return nil, err
			}
			summaries = append(summaries, s)
		}
	}
	return summaries, nil
}

func getSummaryLineForDependencyName(manifest *schema.Manifest, name, manifestPath string) (string, error) {
	currentDependency := manifest.Current.Dependencies[name]
	updatedDependency := manifest.Updated.Dependencies[name]
	inManifest := ""
	if manifestPath != "" {
		inManifest = fmt.Sprintf(" in `%s`", manifestPath)
	}

	cf := &changelogs.ChangelogFinder{
		Source:     updatedDependency.Source,
		Dependency: name,
		Repo:       updatedDependency.Repo,
		Version:    updatedDependency.Constraint,
	}
	updatedURL := cf.FindURL()

	return fmt.Sprintf("- `%s`%s from \"%s\" to \"%s\"", dependencyNameForDisplay(name), inManifest, currentDependency.Constraint, optionalMarkdownLink(updatedDependency.Constraint, updatedURL)), nil
}
