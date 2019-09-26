package schema

import (
	"errors"
)

// Manifest contains manifest data
type Manifest struct {
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
	Constraint string `json:"constraint"`
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
	return nil
}
