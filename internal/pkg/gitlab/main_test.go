package gitlab

import (
	"os"
	"testing"

	"github.com/dependencies-io/pullrequest/internal/app/config"
	"github.com/dependencies-io/pullrequest/internal/app/pullrequest"
)

func getMR() (*MergeRequest, *pullrequest.Pullrequest) {
	config := config.NewConfigFromEnv()
	prBase := pullrequest.NewPullrequestFromEnv("branch", "title", "body", config)
	return NewMergeRequestFromEnv(prBase), prBase
}

func TestNewMergeRequestFromEnv(t *testing.T) {
	os.Setenv("GITLAB_API_URL", "testurl")
	os.Setenv("GITLAB_API_TOKEN", "testtoken")

	pr, prBase := getMR()

	if pr.Pullrequest != prBase {
		t.Error("Pullrequest value incorrect")
	}

	if pr.ProjectAPIURL != "testurl" {
		t.Error("ProjectAPIURL value incorrect")
	}

	if pr.APIToken != "testtoken" {
		t.Error("APIToken value incorrect")
	}
}

func TestCreateTestEnv(t *testing.T) {
	pr, _ := getMR()
	err := pr.Create()
	if err != nil {
		t.Fail()
	}
}

func TestCreateProductionEnv(t *testing.T) {
	// this will try to send an actual API call to github.com, and fail
	os.Setenv("DEPENDENCIES_ENV", "production")
	pr, _ := getMR()
	err := pr.Create()
	if err == nil {
		t.Fail()
	}
}
