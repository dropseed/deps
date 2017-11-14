package github

import (
	"os"
	"testing"

	"github.com/dependencies-io/pullrequest/internal/app/config"
	"github.com/dependencies-io/pullrequest/internal/app/pullrequest"
)

func getPR() (*PullRequest, *pullrequest.Pullrequest) {
	config := config.Config{}
	config.LoadEnvSettings()
	prBase := pullrequest.NewPullrequestFromEnv("branch", "title", "body", &config)
	return NewPullRequestFromEnv(prBase), prBase
}

func TestNewPullRequestFromEnv(t *testing.T) {
	os.Setenv("GITHUB_REPO_FULL_NAME", "dropseed/test")
	os.Setenv("GITHUB_API_TOKEN", "testtoken")

	pr, prBase := getPR()

	if pr.Pullrequest != prBase {
		t.Error("Pullrequest value incorrect")
	}

	if pr.RepoFullName != "dropseed/test" {
		t.Error("RepoFullName value incorrect")
	}

	if pr.APIToken != "testtoken" {
		t.Error("APIToken value incorrect")
	}
}

func TestCreateTestEnv(t *testing.T) {
	pr, _ := getPR()
	err := pr.Create()
	if err != nil {
		t.Fail()
	}
}

func TestCreateProductionEnv(t *testing.T) {
	// this will try to send an actual API call to github.com, and fail
	os.Setenv("DEPENDENCIES_ENV", "production")
	pr, _ := getPR()
	err := pr.Create()
	if err == nil {
		t.Fail()
	}
}
