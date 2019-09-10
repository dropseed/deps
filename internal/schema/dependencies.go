package schema

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"
)

const maxBodyLength = 65535

type Dependencies struct {
	Lockfiles map[string]*Lockfile `json:"lockfiles,omitempty"`
	Manifests map[string]*Manifest `json:"manifests,omitempty"`

	// Not included in JSON
	Title       string `json:"-"`
	Description string `json:"-"`
	UpdateID    string `json:"-"`
	UniqueID    string `json:"-"`
}

// NewDependenciesFromJSONPath loads Dependencies from a JSON file path
func NewDependenciesFromJSONPath(path string) (*Dependencies, error) {
	fileContent, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return NewDependenciesFromJSONContent(fileContent)
}

// NewDependenciesFromJSONContent creates a Dependencies instance with Unmarshalled JSON data
func NewDependenciesFromJSONContent(content []byte) (*Dependencies, error) {
	deps := Dependencies{}
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&deps); err != nil {
		return nil, err
	}

	if err := deps.ValidateAndCompile(); err != nil {
		return nil, err
	}

	return &deps, nil
}

func (s *Dependencies) ValidateAndCompile() error {
	for _, lockfile := range s.Lockfiles {
		if err := lockfile.Validate(); err != nil {
			return err
		}
	}
	for _, manifest := range s.Manifests {
		if err := manifest.Validate(); err != nil {
			return err
		}
	}

	s.Title = s.generateTitle()
	s.Description = s.generateDescription()
	s.UpdateID = s.getUpdateID()
	s.UniqueID = s.getUniqueID()

	return nil
}

// generateTitle generates a title string from the dependencies dependencies, optinally for the related PR search
func (s *Dependencies) generateTitle() string {

	lockfiles := map[string]*Lockfile{}
	manifests := map[string]*Manifest{}

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
				return fmt.Sprintf("Update %v in %v from %v to %v", name, manifestPath, installed, updated)
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

// GenerateBody generates a body string from the dependencies schema
func (s *Dependencies) generateDescription() string {
	lockfiles := map[string]*Lockfile{}
	manifests := map[string]*Manifest{}

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

func getSummaryLinesForLockfiles(lockfiles map[string]*Lockfile) ([]string, error) {
	summaries := make([]string, 0, len(lockfiles))

	// iterate using the sorted keys instead of unpredictable map
	keys := []string{}
	for path := range lockfiles {
		keys = append(keys, path)
	}
	sort.Strings(keys)

	for _, lockfilePath := range keys {
		lockfile := lockfiles[lockfilePath]
		s, err := lockfile.GetSummaryLine(lockfilePath)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, s)
	}
	return summaries, nil
}

func getSummaryLinesForManifests(manifests map[string]*Manifest) ([]string, error) {
	summaries := make([]string, 0, len(manifests))

	// iterate using the sorted keys instead of unpredictable map
	keys := []string{}
	for path := range manifests {
		keys = append(keys, path)
	}
	sort.Strings(keys)

	for _, manifestPath := range keys {
		manifest := manifests[manifestPath]
		// summary, err := manifest.GetSummaryLine(manifestPath)
		// if err != nil {
		// 	return nil, err
		// }
		// summaries = append(summaries, summary)
		// iterate using the sorted keys instead of unpredictable map
		keys := []string{}
		for name := range manifest.Updated.Dependencies {
			keys = append(keys, name)
		}
		sort.Strings(keys)

		for _, dependencyName := range keys {
			s, err := manifest.GetSummaryLineForDependencyName(dependencyName, manifestPath)
			if err != nil {
				return nil, err
			}
			summaries = append(summaries, s)
		}
	}
	return summaries, nil
}

func (dependencies *Dependencies) getUpdateID() string {
	truncated := Dependencies{
		// TODO if type is important to separate updates between components,
		// then can add Dependencies.Type and use that too
		Lockfiles: map[string]*Lockfile{},
		Manifests: map[string]*Manifest{},
	}

	if dependencies.Lockfiles != nil {
		for name := range dependencies.Lockfiles {
			// Only care about the filename
			truncated.Lockfiles[name] = nil
		}
	}

	if dependencies.Manifests != nil {
		for name, manifest := range dependencies.Manifests {
			if !manifest.HasUpdates() {
				continue
			}

			// Only care about the filename + dependency names
			truncatedManifest := &Manifest{
				Updated: &ManifestVersion{
					Dependencies: map[string]*ManifestDependency{},
				},
			}
			for dep := range manifest.Updated.Dependencies {
				truncatedManifest.Updated.Dependencies[dep] = nil
			}

			truncated.Manifests[name] = truncatedManifest
		}
	}

	return getShortMD5(truncated)
}

func (dependencies *Dependencies) getUniqueID() string {
	return getShortMD5(dependencies)
}

func getShortMD5(i interface{}) string {
	out, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	sum := md5.Sum(out)
	str := hex.EncodeToString(sum[:])
	short := str[:7]
	return short
}
