package schema

import (
	"errors"
	"fmt"
)

// Manifest contains manifest data
type Manifest struct {
	// TODO path string
	Current *ManifestVersion `json:"current"`
	// TODO remove?
	LockfilePath string           `json:"lockfile_path,omitempty"`
	Updated      *ManifestVersion `json:"updated,omitempty"`
}

// ManifestVersion constains data for a manifest at a specific point in time
type ManifestVersion struct {
	Dependencies map[string]*ManifestDependency `json:"dependencies"`
}

type ManifestDependency struct {
	// Latest     *Version `json:"latest"`
	Constraint string `json:"constraint"`
	// Version    *Version `json:"version"`
	*Dependency
}

func (manifest *Manifest) Validate() error {
	if manifest.Current != nil {
		if err := manifest.Current.Validate(); err != nil {
			return err
		}
	} else {
		return errors.New("manifest.current is requried")
	}

	if manifest.Updated != nil {
		if err := manifest.Updated.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (manifest *Manifest) HasUpdates() bool {
	return manifest.Updated != nil && len(manifest.Updated.Dependencies) > 0
}

func (mv *ManifestVersion) Validate() error {
	for _, dependency := range mv.Dependencies {
		if err := dependency.Validate(); err != nil {
			return err
		}
	}
	return nil
}
func (md *ManifestDependency) Validate() error {
	if md.Constraint == "" {
		return errors.New("manifest dependency constraint is required")
	}
	// if md.Latest != nil {
	// 	if err := md.Latest.Validate(); err != nil {
	// 		return err
	// 	}
	// } else {
	// 	return errors.New("manifest dependency latest is required")
	// }
	return nil
}

// GetSummaryLineForDependencyName returns a bulleted list item string
func (manifest *Manifest) GetSummaryLineForDependencyName(name, manifestPath string) (string, error) {
	currentDependency := manifest.Current.Dependencies[name]
	updatedDependency := manifest.Updated.Dependencies[name]
	return fmt.Sprintf("- `%v` in `%v` from \"%v\" to \"%v\"", name, manifestPath, currentDependency.Constraint, updatedDependency.Constraint), nil
}

// func (manifest *Manifest) GetSummaryLine(manifestPath string) (string, error) {
// 	if len(manifest.Updated.Dependencies) == 1 {
// 		// for nam
// 		// currentDependency := manifest.Current.Dependencies[name]
// 		// updatedDependency := manifest.Updated.Dependencies[name]
// 		// return fmt.Sprintf("- `%v` in `%v` from \"%v\" to \"%v\"", name, manifestPath, currentDependency.Constraint, updatedDependency.Constraint), nil
// 		return "update single manifest dep", nil
// 	}
// 	// update x, y , z in manifestPath
// 	return "update multiple manifest deps", nil
// }

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

	// TODO figure out
	// content := "\n\n" + updatedDependency.GetMarkdownContentForVersion(name, updatedDependency.Constraint)
	content := ""

	return subject + content, nil
}
