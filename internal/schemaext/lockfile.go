package schemaext

import (
	"fmt"
	"sort"

	"github.com/dropseed/deps/pkg/schema"
)

// LockfileChanges stores data about what changes were made to a lockfile
type LockfileChanges struct {
	Updated []string
	Added   []string
	Removed []string
}

func getSummaryLinesForLockfiles(lockfiles map[string]*schema.Lockfile) ([]string, error) {
	summaries := make([]string, 0, len(lockfiles))

	// iterate using the sorted keys instead of unpredictable map
	keys := []string{}
	for path := range lockfiles {
		keys = append(keys, path)
	}
	sort.Strings(keys)

	for _, lockfilePath := range keys {
		lockfile := lockfiles[lockfilePath]
		s, err := getSummaryLineForLockfile(lockfile, lockfilePath)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, nil
}

func lockfileChangesByType(lockfile *schema.Lockfile) map[string]*LockfileChanges {
	changesByType := map[string]*LockfileChanges{}

	for name, dep := range lockfile.Current.Dependencies {
		depType := "direct"
		if dep.IsTransitive {
			depType = "transitive"
		}

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
			depType := "direct"
			if dep.IsTransitive {
				depType = "transitive"
			}

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

func getSummaryLineForLockfile(lockfile *schema.Lockfile, lockfilePath string) (string, error) {
	changesByType := lockfileChangesByType(lockfile)

	subitems := ""

	numTransitive := 0
	numDirect := 0

	if transitive, found := changesByType["transitive"]; found && len(transitive.Updated) > 0 {
		numTransitive = len(transitive.Updated)
	}

	if direct, found := changesByType["direct"]; found && len(direct.Updated) > 0 {
		numDirect = len(direct.Updated)

		sort.Strings(direct.Updated) // sort first to get predictable order
		for _, name := range direct.Updated {
			currentDep := lockfile.Current.Dependencies[name]
			dep := lockfile.Updated.Dependencies[name]
			subitems += fmt.Sprintf("\n  - `%s` was updated from %s to %s", name, currentDep.Version.Name, dep.Version.Name)
		}
	}

	parens := fmt.Sprintf(" (including %d direct and %d transitive dependencies)", numDirect, numTransitive)

	return fmt.Sprintf("- `%s` was updated%s%s", lockfilePath, parens, subitems), nil
}
