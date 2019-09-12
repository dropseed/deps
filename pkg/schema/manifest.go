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
	inManifest := ""
	if manifestPath != "" {
		inManifest = fmt.Sprintf(" in `%s`", manifestPath)
	}
	return fmt.Sprintf("- `%s`%s from \"%s\" to \"%s\"", dependencyNameForDisplay(name), inManifest, currentDependency.Constraint, updatedDependency.Constraint), nil
}
