package github

import "regexp"

func dereferenceGitHubIssueLinks(body string) (string, error) {
	r, err := regexp.Compile("https://github.com/([^/]+/[^/]+/(issue|pull)/\\d+)")
	if err != nil {
		return "", err
	}
	sanitized := r.ReplaceAllString(body, "https://www.dependencies.io/github-redirect/$1")
	return sanitized, nil
}
