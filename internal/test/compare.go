package test

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/dependencies-io/deps/internal/output"
	"github.com/dependencies-io/deps/internal/schema"
	"github.com/sergi/go-diff/diffmatchpatch"
)

const LOOSEHOLDER = "testloose"

func schemasMatchExactly(given, expected interface{}) (bool, error) {
	givenOut, err := json.MarshalIndent(given, "", "  ")
	if err != nil {
		return false, err
	}
	expectedOut, err := json.MarshalIndent(expected, "", "  ")
	if err != nil {
		return false, err
	}

	givenOutString := string(givenOut)
	expectedOutString := string(expectedOut)

	if givenOutString != expectedOutString {
		output.Error(strings.Repeat("=", 80))
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(expectedOutString, givenOutString, false)
		fmt.Println(dmp.DiffPrettyText(diffs))
		output.Error(strings.Repeat("=", 80))
		return false, nil
	}

	return true, nil
}

func schemasMatchLoosely(given, expected interface{}) (bool, error) {
	if gd, gok := given.(*schema.Dependencies); gok {
		if ed, ok := expected.(*schema.Dependencies); ok {
			cleanSchemaDependencies(gd, ed)
		}
	}
	return schemasMatchExactly(given, expected)
}

func cleanSchemaDependencies(given, expected *schema.Dependencies) {
	output.Debug("Cleaning parsed manifest dependencies for loose comparison")
	// for p, m := range expected.Manifests {
	// 	for depName, dep := range m.Current.Dependencies {
	// 		// Truncate the parsed available versions to the same length as the expected
	// 		// so that new versions can be added onto the end, but nothing else
	// 		if givenDep, ok := given.Manifests[p].Current.Dependencies[depName]; ok {
	// 			if execute.Verbosity > 0 {
	// 				output.Debug("Truncating %s availble from %d to %d", depName, len(givenDep.Available), len(dep.Available))
	// 			}
	// 			givenDep.Available = givenDep.Available[:len(dep.Available)]
	// 		}
	// 	}
	// }

	output.Debug("Cleaning lockfile dependencies for loose comparison")

	for _, l := range expected.Lockfiles {
		if lfv := l.Current; lfv != nil {
			cleanSchemaLockfileVersion(lfv)
		}
		if lfv := l.Updated; lfv != nil {
			cleanSchemaLockfileVersion(lfv)
		}
	}
	for _, l := range given.Lockfiles {
		if lfv := l.Current; lfv != nil {
			cleanSchemaLockfileVersion(lfv)
		}
		if lfv := l.Updated; lfv != nil {
			cleanSchemaLockfileVersion(lfv)
		}
	}
}

func cleanSchemaLockfileVersion(lfv *schema.LockfileVersion) {
	lfv.Fingerprint = LOOSEHOLDER
	for depName, dep := range lfv.Dependencies {
		if dep.IsTransitive {
			delete(lfv.Dependencies, depName)
		} else {
			dep.Version.Name = LOOSEHOLDER
			// dep.Version.Constraint = LOOSEHOLDER
		}
	}
}
