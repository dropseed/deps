package schema

import (
	"errors"
)

type Lockfile struct {
	Current *LockfileVersion `json:"current"`
	Updated *LockfileVersion `json:"updated,omitempty"`
}

type LockfileVersion struct {
	Dependencies map[string]*LockfileDependency `json:"dependencies"`
	Fingerprint  string                         `json:"fingerprint"`
}

// Dependency stores data for a manifest or lockfile dependency (some fields will be empty)
type LockfileDependency struct {
	// Constraint   string   `json:"constraint,omitempty"`
	Version      *Version `json:"version"`
	IsTransitive bool     `json:"is_transitive,omitempty"`
	*Dependency
}

func (lockfile *Lockfile) Validate() error {
	if lockfile.Current != nil {
		if err := lockfile.Current.Validate(); err != nil {
			return err
		}
	} else {
		return errors.New("lockfile.current is required")
	}

	if lockfile.Updated != nil {
		if err := lockfile.Updated.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (lockfile *Lockfile) HasUpdates() bool {
	return lockfile.Updated != nil && len(lockfile.Updated.Dependencies) > 0
}

func (lv *LockfileVersion) Validate() error {
	if lv.Fingerprint == "" {
		return errors.New("lockfile fingerprint is required")
	}

	for _, dependency := range lv.Dependencies {
		if err := dependency.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (ld *LockfileDependency) Validate() error {
	if ld.Version != nil {
		if err := ld.Version.Validate(); err != nil {
			return err
		}
	} else {
		return errors.New("lockfile dependency.version is required")
	}
	return nil
}
