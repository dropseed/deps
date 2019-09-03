package github

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/bradleyfalzon/ghinstallation"
	"github.com/dropseed/deps/internal/git"
)

func dereferenceGitHubIssueLinks(body string) (string, error) {
	r, err := regexp.Compile("https://github.com/([^/]+/[^/]+/(issues|pull)/\\d+)")
	if err != nil {
		return "", err
	}
	sanitized := r.ReplaceAllString(body, "https://www.dependencies.io/github-redirect/$1")
	return sanitized, nil
}

func getRepoFullName() (string, error) {

	// Custom override
	if s := os.Getenv("DEPS_GITHUB_REPOSITORY"); s != "" {
		return s, nil
	}

	// GitHub Actions
	if s := os.Getenv("GITHUB_REPOSITORY"); s != "" {
		return s, nil
	}

	if s := os.Getenv("TRAVIS_REPO_SLUG"); s != "" {
		return s, nil
	}

	if s := os.Getenv("CIRCLE_PROJECT_USERNAME"); s != "" {
		return fmt.Sprintf("%s/%s", s, os.Getenv("CIRCLE_PROJECT_REPONAME")), nil
	}

	if s := getRepoFullNameFromRemote(git.GitRemote()); s != "" {
		return s, nil
	}

	return "", errors.New("Unable to find GitHub repo full name")
}

func getRepoFullNameFromRemote(remote string) string {
	pattern := regexp.MustCompile("([a-zA-Z0-9_-]+\\/[a-zA-Z0-9_-]+)(\\.git)?\\/?$")
	matches := pattern.FindStringSubmatch(remote)
	if len(matches) > 0 {
		return matches[1]
	}
	return ""
}

func getAPIToken() string {
	if s := os.Getenv("DEPS_GITHUB_TOKEN"); s != "" {
		return s
	}

	if key := os.Getenv("DEPS_GITHUB_APP_KEY"); key != "" {
		// key path
		// key raw
		// key b64 - only this for now

		keyBytes, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			panic("Could not decode base64 DEPS_GITHUB_APP_KEY")
		}

		appID, err := strconv.Atoi(os.Getenv("DEPS_GITHUB_APP_ID"))
		if err != nil {
			panic("Invalid DEPS_GITHUB_APP_ID")
		}
		installationID, err := strconv.Atoi(os.Getenv("DEPS_GITHUB_APP_INSTALLATION_ID"))
		if err != nil {
			panic("Invalid DEPS_GITHUB_APP_INSTALLATION_ID")
		}

		tr := http.DefaultTransport
		itr, err := ghinstallation.New(tr, appID, installationID, keyBytes)
		if err != nil {
			panic(err)
		}
		token, err := itr.Token()
		if err != nil {
			panic(err)
		}
		return token
	}

	// Used by GitHub Actions
	if s := os.Getenv("GITHUB_TOKEN"); s != "" {
		return s
	}

	return ""
}
