package schema

import (
	"errors"
	"fmt"
	"sort"
	"strings"
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

// GetDependencyTypeString returns a string representation of the dependencies relationship to the repo
func (dep *LockfileDependency) GetDependencyTypeString() string {
	if dep.IsTransitive {
		return "transitive"
	}

	return "direct"
}

// LockfileChanges stores data about what changes were made to a lockfile
type LockfileChanges struct {
	Updated []string
	Added   []string
	Removed []string
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

func (lockfile *Lockfile) changesByType() map[string]*LockfileChanges {
	changesByType := map[string]*LockfileChanges{}

	for name, dep := range lockfile.Current.Dependencies {
		depType := dep.GetDependencyTypeString()

		_, ok := changesByType[depType]
		if !ok {
			changesByType[depType] = &LockfileChanges{}
		}
		changesForType := changesByType[depType]

		if updatedDep, found := lockfile.Updated.Dependencies[name]; !found {
			changesForType.Removed = append(changesForType.Removed, name)
		} else {
			if dep.Version.Name != updatedDep.Version.Name {
				changesForType.Updated = append(changesForType.Updated, name)
			}
		}
	}

	for name, dep := range lockfile.Updated.Dependencies {
		if _, found := lockfile.Current.Dependencies[name]; !found {
			depType := dep.GetDependencyTypeString()

			_, ok := changesByType[depType]
			if !ok {
				changesByType[depType] = &LockfileChanges{}
			}
			changesForType := changesByType[depType]

			changesForType.Added = append(changesForType.Added, name)
		}
	}

	return changesByType
}

// GetSummaryLine returns a summary line for a bulleted markdown list
func (lockfile *Lockfile) GetSummaryLine(lockfilePath string) (string, error) {
	changesByType := lockfile.changesByType()
	additional := ""
	if direct, found := changesByType["direct"]; found && len(direct.Updated) > 0 {
		additional = fmt.Sprintf(" (including %d updated direct dependencies)", len(direct.Updated))
	}
	return fmt.Sprintf("- `%v` was updated%v", lockfilePath, additional), nil
}

// GetBodyContent compiles the long-form content for changes to the lockfile
func (lockfile *Lockfile) GetBodyContent(lockfilePath string) (string, error) {
	changesByType := lockfile.changesByType()

	contentParts := []string{}

	contentParts = append(contentParts, "### "+lockfilePath)

	if transitive, found := changesByType["transitive"]; found {
		contentParts = append(contentParts, fmt.Sprintf("%d transitive dependencies were updated, %d were added, and %d removed. View the git diff for more details about exactly what changed.", len(transitive.Updated), len(transitive.Added), len(transitive.Removed)))
	}

	if direct, found := changesByType["direct"]; found && len(direct.Updated) > 0 {
		contentParts = append(contentParts, fmt.Sprintf("The following %d direct dependencies were updated:", len(direct.Updated)))

		sort.Strings(direct.Updated) // sort first to get predictable order
		for _, name := range direct.Updated {
			currentDep := lockfile.Current.Dependencies[name]
			dep := lockfile.Updated.Dependencies[name]
			versionContent := dep.GetMarkdownContentForVersion(name, dep.Version)
			contentParts = append(contentParts, fmt.Sprintf("#### `%s` was updated from %s to %s\n\n%s", name, currentDep.Version.Name, dep.Version.Name, versionContent))
		}
	}

	return strings.Join(contentParts, "\n\n"), nil
}
