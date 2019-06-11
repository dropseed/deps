package github

import "testing"

func TestNoopDereference(t *testing.T) {
	body := "hey this is normal\n\nwith newlines"
	cleaned, err := dereferenceGitHubIssueLinks(body)
	if err != nil {
		t.Error(err)
	}
	if body != cleaned {
		t.FailNow()
	}
}

func TestDereference(t *testing.T) {
	body := "hey this is normal\n\n[with](https://github.com/test-org/repo/issues/45) newlines"
	cleaned, err := dereferenceGitHubIssueLinks(body)
	if err != nil {
		t.Error(err)
	}
	if cleaned != "hey this is normal\n\n[with](https://www.dependencies.io/github-redirect/test-org/repo/issues/45) newlines" {
		t.Error(cleaned)
	}
}

func TestRepoNameFromRemote(t *testing.T) {
	remote := "https://github.com/dropseed/test.git/"
	name := getRepoFullNameFromRemote(remote)
	if name != "dropseed/test" {
		t.Error(name)
	}
}
