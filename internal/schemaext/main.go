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
			if shortOverview := getShortOverviewForLockfile(lockfiles[lockfilePaths[0]]); shortOverview != "" {
				return fmt.Sprintf("Update %v (%s)", lockfilePaths[0], shortOverview)
			} else {
				return fmt.Sprintf("Update %v", lockfilePaths[0])
			}
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

			if shortOverview := getShortOverviewForManifest(manifest); shortOverview != "" {
				return fmt.Sprintf("Update %s (%s)", manifestPath, shortOverview)
			} else {
				return fmt.Sprintf("Update %s", manifestPath)
			}
		}

		// More than 1 manifest
		return fmt.Sprintf("Update dependencies in %v", strings.Join(manifestPaths, ", "))
	}

	return ""
}

func DescriptionForDeps(s *schema.Dependencies) string {
	summaryHeader := "The following dependencies have been updated by [deps](https://www.dependencies.io/):"
	summaryLines := DescriptionItemsForDeps(s)
	final := summaryHeader + "\n\n" + summaryLines + "\n"

	if len(final) > maxBodyLength {
		final = final[:maxBodyLength]
	}

	return final
}

func DescriptionItemsForDeps(s *schema.Dependencies) string {
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

	return strings.Join(summaryLines, "\n")
}
