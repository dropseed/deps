package github

import (
	"errors"
	"fmt"
	"os"
	"regexp"

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

	// Used by GitHub Actions
	if s := os.Getenv("GITHUB_TOKEN"); s != "" {
		return s
	}

	panic("Unable to find GitHub API token")
}
