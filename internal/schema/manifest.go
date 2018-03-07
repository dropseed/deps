package schema

import (
	"fmt"
)

// Manifest contains manifest data
type Manifest struct {
	// TODO path string
	Current *ManifestVersion `json:"current"`
	Updated *ManifestVersion `json:"updated,omitempty"`
}

// ManifestVersion constains data for a manifest at a specific point in time
type ManifestVersion struct {
	Dependencies map[string]ManifestDependency `json:"dependencies"`
}

type ManifestDependency struct {
	*Dependency
	Constraint string    `json:"constraint"`
	Available  []Version `json:"available,omitempty"`
}

// GetSummaryLineForDependencyName returns a bulleted list item string
func (manifest *Manifest) GetSummaryLineForDependencyName(name, manifestPath string) (string, error) {
	currentDependency := manifest.Current.Dependencies[name]
	updatedDependency := manifest.Updated.Dependencies[name]
	return fmt.Sprintf("- `%v` in `%v` from \"%v\" to \"%v\"", name, manifestPath, currentDependency.Constraint, updatedDependency.Constraint), nil
}

// GetBodyContentForDependencyName compiles the markdown content for this dependency update
func (manifest *Manifest) GetBodyContentForDependencyName(name, manifestPath string) (string, error) {
	// TODO add notes

	currentDependency := manifest.Current.Dependencies[name]
	updatedDependency := manifest.Updated.Dependencies[name]

	subject := fmt.Sprintf(
		"### `%s`\n\nThis dependency is located in `%s` and was updated from \"%s\" to \"%s\".",
		name,
		manifestPath,
		currentDependency.Constraint,
		updatedDependency.Constraint,
	)

	content := "\n\n" + updatedDependency.GetMarkdownContentForVersions(name, updatedDependency.Available)

	return subject + content, nil
}
