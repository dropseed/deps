package schema

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/dependencies-io/pullrequest/internal/env"
)

const maxBodyLength = 65535

type Dependencies struct {
	Lockfiles map[string]Lockfile `json:"lockfiles,omitempty"`
	Manifests map[string]Manifest `json:"manifests,omitempty"`
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
	if err := json.Unmarshal(content, &deps); err != nil {
		return nil, err
	}
	return &deps, nil
}

// GenerateTitle generates a title string from the dependencies schema
func (s *Dependencies) GenerateTitle() (string, error) {
	return s.generateTitle(false)
}

// GenerateRelatedPRTitleSearch generates a title search query string from the dependencies schema
func (s *Dependencies) GenerateRelatedPRTitleSearch() (string, error) {
	return s.generateTitle(true)
}

// generateTitle generates a title string from the dependencies dependencies, optinally for the related PR search
func (s *Dependencies) generateTitle(forRelatedPRSearch bool) (string, error) {
	lockfiles := s.Lockfiles
	manifests := s.Manifests
	foundLockfiles := len(lockfiles) > 0
	foundManifests := len(manifests) > 0

	if foundLockfiles && foundManifests {
		if forRelatedPRSearch {
			return "", errors.New("Can't generate a title to search for in this scenario")
		}
		lfPlural := "lockfiles"
		if len(lockfiles) == 1 {
			lfPlural = "lockfile"
		}
		mPlural := "manifests"
		if len(manifests) == 1 {
			mPlural = "manifest"
		}
		return fmt.Sprintf("Update %v %v and %v %v", len(lockfiles), lfPlural, len(manifests), mPlural), nil
	}

	if foundLockfiles {
		lockfilePaths := []string{}
		for k := range lockfiles {
			lockfilePaths = append(lockfilePaths, k)
		}
		sort.Strings(lockfilePaths)

		if len(lockfilePaths) == 1 {
			return fmt.Sprintf("Update %v", lockfilePaths[0]), nil
		}
		return fmt.Sprintf("Update lockfiles: %v", strings.Join(lockfilePaths, ", ")), nil
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
				if forRelatedPRSearch {
					// remove the last version name from the title returned
					updated = ""
				}
				return fmt.Sprintf("Update %v in %v from %v to %v", name, manifestPath, installed, updated), nil
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

			if forRelatedPRSearch {
				return "", errors.New("Can't generate a title to search for in this scenario")
			}

			return fmt.Sprintf("Update %v dependencies from %v", len(dependencies), strings.Join(sourceNames, ", ")), nil

		}

		// More than 1 manifest
		return fmt.Sprintf("Update dependencies in %v", strings.Join(manifestPaths, ", ")), nil
	}

	return "", errors.New("Unable to determine PR title")
}

// GenerateBody generates a body string from the dependencies schema
func (s *Dependencies) GenerateBody() (string, error) {
	lockfiles := s.Lockfiles
	manifests := s.Manifests
	foundLockfiles := len(lockfiles) > 0
	foundManifests := len(manifests) > 0

	if !foundLockfiles && !foundManifests {
		return "", errors.New("Dependencies must contain either lockfiles or manifests")
	}

	summaryLines := []string{}
	contentParts := []string{}

	if foundLockfiles {
		lines, err := getSummaryLinesForLockfiles(lockfiles)
		if err != nil {
			return "", err
		}
		summaryLines = append(summaryLines, lines...)

		parts, err := getBodyPartsForLockfiles(lockfiles)
		if err != nil {
			return "", err
		}
		contentParts = append(contentParts, parts...)
	}

	if foundManifests {
		lines, err := getSummaryLinesForManifests(manifests)
		if err != nil {
			return "", err
		}
		summaryLines = append(summaryLines, lines...)

		parts, err := getBodyPartsForManifests(manifests)
		if err != nil {
			return "", err
		}
		contentParts = append(contentParts, parts...)
	}

	summaryHeader := "## Overview\n\nThe following dependencies have been updated by [dependencies.io](https://www.dependencies.io/):\n\n"
	bodyHeader := "\n\n## Details\n\n"

	notes := env.GetSetting("pullrequest_notes", "")
	if notes != "" {
		notes = notes + "\n\n---\n\n"
	}

	final := notes + summaryHeader + strings.Join(summaryLines, "\n") + bodyHeader + strings.Join(contentParts, "\n\n---\n\n") + "\n"

	if len(final) > maxBodyLength {
		final = final[:maxBodyLength]
	}

	return final, nil
}

func getSummaryLinesForLockfiles(lockfiles map[string]Lockfile) ([]string, error) {
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

func getBodyPartsForLockfiles(lockfiles map[string]Lockfile) ([]string, error) {
	parts := make([]string, 0, len(lockfiles))

	// iterate using the sorted keys instead of unpredictable map
	keys := []string{}
	for path := range lockfiles {
		keys = append(keys, path)
	}
	sort.Strings(keys)

	for _, lockfilePath := range keys {
		lockfile := lockfiles[lockfilePath]
		s, err := lockfile.GetBodyContent(lockfilePath)
		if err != nil {
			return nil, err
		}
		parts = append(parts, s)
	}
	return parts, nil
}

func getSummaryLinesForManifests(manifests map[string]Manifest) ([]string, error) {
	summaries := make([]string, 0, len(manifests))

	// iterate using the sorted keys instead of unpredictable map
	keys := []string{}
	for path := range manifests {
		keys = append(keys, path)
	}
	sort.Strings(keys)

	for _, manifestPath := range keys {
		manifest := manifests[manifestPath]

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

func getBodyPartsForManifests(manifests map[string]Manifest) ([]string, error) {
	parts := []string{}

	// iterate using the sorted keys instead of unpredictable map
	keys := []string{}
	for path := range manifests {
		keys = append(keys, path)
	}
	sort.Strings(keys)

	for _, manifestPath := range keys {
		manifest := manifests[manifestPath]

		// iterate using the sorted keys instead of unpredictable map
		keys := []string{}
		for name := range manifest.Updated.Dependencies {
			keys = append(keys, name)
		}
		sort.Strings(keys)

		for _, dependencyName := range keys {
			depBody, err := manifest.GetBodyContentForDependencyName(dependencyName, manifestPath)
			if err != nil {
				return nil, err
			}
			parts = append(parts, depBody)
		}
	}

	return parts, nil
}
