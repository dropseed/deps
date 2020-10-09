package schemaext

import (
	"fmt"
	"sort"
	"strings"

	"github.com/dropseed/deps/pkg/schema"
)

const maxBodyLength = 65535

func TitleForDeps(s *schema.Dependencies) string {

	lockfiles := map[string]*schema.Lockfile{}
	manifests := map[string]*schema.Manifest{}

	for name, lockfile := range s.Lockfiles {
		if lockfile.HasUpdates() {
			lockfiles[name] = lockfile
		}
	}
	for name, manifest := range s.Manifests {
		if manifest.HasUpdates() {
			manifests[name] = manifest
		}
	}

	foundLockfiles := len(lockfiles) > 0
	foundManifests := len(manifests) > 0

	if foundLockfiles && foundManifests {
		lfPlural := "lockfiles"
		if len(lockfiles) == 1 {
			lfPlural = "lockfile"
		}
		mPlural := "manifests"
		if len(manifests) == 1 {
			mPlural = "manifest"
		}
		return fmt.Sprintf("Update %v %v and %v %v", len(lockfiles), lfPlural, len(manifests), mPlural)
	}

	if foundLockfiles {
		lockfilePaths := []string{}
		for k := range lockfiles {
			lockfilePaths = append(lockfilePaths, k)
		}
		sort.Strings(lockfilePaths)

		if len(lockfilePaths) == 1 {
			return fmt.Sprintf("Update %v", lockfilePaths[0])
		}
		return fmt.Sprintf("Update lockfiles: %v", strings.Join(lockfilePaths, ", "))
	}

	if foundManifests {

		manifestPaths := make([]string, 0, len(manifests))
		for k := range manifests {
			manifestPaths = append(manifestPaths, k)
		}

		if len(manifestPaths) == 1 {

			manifestPath := manifestPaths[0]
			manifest := manifests[manifestPath]
			dependencies := manifest.Updated.Dependencies
			dependencyNames := make([]string, 0, len(dependencies))
			for k := range dependencies {
				dependencyNames = append(dependencyNames, k)
			}

			if len(dependencyNames) == 1 {
				name := dependencyNames[0]
				dep := dependencies[name]
				installed := manifest.Current.Dependencies[name].Constraint
				updated := dep.Constraint
				inManifest := ""
				if manifestPath != "" && manifestPath != "." && manifestPath != "/" {
					inManifest = fmt.Sprintf(" in %s", manifestPath)
				}
				return fmt.Sprintf("Update %s%s from %s to %s", dependencyNameForDisplay(name), inManifest, installed, updated)
			}

			// more than 1 dependency
			// create a "set" of sources
			sources := make(map[string]bool)
			for _, dep := range dependencies {
				source := dep.Source
				sources[source] = true
			}

			// get the keys remaining
			sourceNames := []string{}
			for k := range sources {
				sourceNames = append(sourceNames, k)
			}

			sort.Strings(sourceNames)

			// TODO if > 2 items, put an "and " in front of the last one

			return fmt.Sprintf("Update %v dependencies from %v", len(dependencies), strings.Join(sourceNames, ", "))

		}

		// More than 1 manifest
		return fmt.Sprintf("Update dependencies in %v", strings.Join(manifestPaths, ", "))
	}

	return ""
}

func DescriptionForDeps(s *schema.Dependencies) string {
	lockfiles := map[string]*schema.Lockfile{}
	manifests := map[string]*schema.Manifest{}

	for name, lockfile := range s.Lockfiles {
		if lockfile.HasUpdates() {
			lockfiles[name] = lockfile
		}
	}
	for name, manifest := range s.Manifests {
		if manifest.HasUpdates() {
			manifests[name] = manifest
		}
	}

	foundLockfiles := len(lockfiles) > 0
	foundManifests := len(manifests) > 0

	if !foundLockfiles && !foundManifests {
		return ""
	}

	summaryLines := []string{}

	if foundLockfiles {
		lines, err := getSummaryLinesForLockfiles(lockfiles)
		if err != nil {
			panic(err)
		}
		summaryLines = append(summaryLines, lines...)
	}

	if foundManifests {
		lines, err := getSummaryLinesForManifests(manifests)
		if err != nil {
			panic(err)
		}
		summaryLines = append(summaryLines, lines...)
	}

	summaryHeader := "The following dependencies have been updated by [dependencies.io](https://www.dependencies.io/):\n\n"

	notes := "" // TODO use go template instead
	// notes := env.GetSetting("pullrequest_notes", "")
	// if notes != "" {
	// 	notes = notes + "\n\n---\n\n"
	// }

	final := notes + summaryHeader + strings.Join(summaryLines, "\n") + "\n"

	if len(final) > maxBodyLength {
		final = final[:maxBodyLength]
	}

	return final
}
