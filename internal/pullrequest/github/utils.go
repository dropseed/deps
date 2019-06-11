package github

import (
	"fmt"
	"os"
	"regexp"
)

func dereferenceGitHubIssueLinks(body string) (string, error) {
	r, err := regexp.Compile("https://github.com/([^/]+/[^/]+/(issues|pull)/\\d+)")
	if err != nil {
		return "", err
	}
	sanitized := r.ReplaceAllString(body, "https://www.dependencies.io/github-redirect/$1")
	return sanitized, nil
}

func getRepoFullName() string {

	if s := os.Getenv("DEPS_GITHUB_REPO_FULL_NAME"); s != "" {
		return s
	}

	if s := os.Getenv("TRAVIS_REPO_SLUG"); s != "" {
		return s
	}

	if s := os.Getenv("CIRCLE_PROJECT_USERNAME"); s != "" {
		return fmt.Sprintf("%s/%s", s, os.Getenv("CIRCLE_PROJECT_REPONAME"))
	}

	// git remote

	return ""
}

func getAPIToken() string {
	if s := os.Getenv("DEPS_GITHUB_API_TOKEN"); s != "" {
		return s
	}

	// Used by GitHub Actions
	if s := os.Getenv("GITHUB_TOKEN"); s != "" {
		return s
	}

	return ""
}
