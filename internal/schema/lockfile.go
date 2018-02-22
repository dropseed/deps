package schema

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

type Lockfile struct {
	Current *LockfileVersion `json:"current"`
	Updated *LockfileVersion `json:"updated,omitempty"`
}

type LockfileVersion struct {
	Fingerprint  string                        `json:"fingerprint"`
	Dependencies map[string]LockfileDependency `json:"dependencies"`
}

// Dependency stores data for a manifest or lockfile dependency (some fields will be empty)
type LockfileDependency struct {
	Source       string  `json:"source"`
	Installed    Version `json:"installed"`
	IsTransitive bool    `json:"is_transitive,omitempty"`
	Constraint   string  `json:"constraint,omitempty"`
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
	Updated int
	Added   int
	Removed int
}

func (lc *LockfileChanges) getTableLine(depType string) string {
	updatedString, addedString, removedString := "-", "-", "-"

	if lc.Updated > 0 {
		updatedString = strconv.Itoa(lc.Updated)
	}
	if lc.Added > 0 {
		addedString = strconv.Itoa(lc.Added)
	}
	if lc.Removed > 0 {
		removedString = strconv.Itoa(lc.Removed)
	}

	return fmt.Sprintf("| %v | %v | %v | %v |", strings.Title(depType), updatedString, addedString, removedString)
}

// GetSummaryLine returns a summary line for a bulleted markdown list
func (lockfile *Lockfile) GetSummaryLine(lockfilePath string) (string, error) {
	return fmt.Sprintf("- `%v` was updated", lockfilePath), nil
}

// GetBodyContent compiles the long-form content for changes to the lockfile
func (lockfile *Lockfile) GetBodyContent(lockfilePath string) (string, error) {
	changesByType := map[string]*LockfileChanges{}

	for name, dep := range lockfile.Current.Dependencies {
		depType := dep.GetDependencyTypeString()

		_, ok := changesByType[depType]
		if !ok {
			changesByType[depType] = &LockfileChanges{}
		}
		changesForType := changesByType[depType]

		if updatedDep, found := lockfile.Updated.Dependencies[name]; !found {
			changesForType.Removed++
		} else {
			if dep.Installed.Name != updatedDep.Installed.Name {
				changesForType.Updated++
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

			changesForType.Added++
		}
	}

	beforeTable := fmt.Sprintf("The following changes were made in the `%v` update:\n\n", lockfilePath)
	afterTable := "\n\nView the git diff for more details about exactly what changed."

	tableHeader := "| Type | Updated | Added | Removed |\n|---|---|---|---|\n"
	lines := []string{}

	keys := []string{}
	for k := range changesByType {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, depType := range keys {
		changes := changesByType[depType]
		lines = append(lines, changes.getTableLine(depType))
	}
	tableBody := strings.Join(lines, "\n")

	return beforeTable + tableHeader + tableBody + afterTable, nil
}
