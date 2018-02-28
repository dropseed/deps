package github

import (
	"os"
	"testing"

	"github.com/dependencies-io/deps/internal/pullrequest"
)

func getPR(dependenciesJSONPath string) (*PullRequest, *pullrequest.Pullrequest) {
	os.Setenv("JOB_ID", "test")
	pr, err := NewPullrequestFromDependenciesJSONPathAndEnv(dependenciesJSONPath)
	if err != nil {
		panic(err)
	}
	return pr, pr.Pullrequest
}

func TestNewPullRequestFromEnv(t *testing.T) {
	os.Setenv("GITHUB_REPO_FULL_NAME", "dropseed/test")
	os.Setenv("GITHUB_API_TOKEN", "testtoken")

	pr, prBase := getPR("../../schema/testdata/single_dependency.json")

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
	pr, _ := getPR("../../schema/testdata/single_dependency.json")
	err := pr.Create()
	if err != nil {
		t.Fail()
	}
}

func TestCreateProductionEnv(t *testing.T) {
	// this will try to send an actual API call to github.com, and fail
	os.Setenv("DEPENDENCIES_ENV", "production")
	pr, _ := getPR("../../schema/testdata/single_dependency.json")
	err := pr.Create()
	if err == nil {
		t.Fail()
	}
}

// func TestGetActionsJSON(t *testing.T) {
// 	os.Setenv("DEPENDENCIES_ENV", "test")
// 	pr, _ := getPR("./testdata/action_dependencies.json")
// 	pr.Create()
// 	output, err := pr.GetActionsJSON()
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if output != "<Actions>{\"PR #0\":{\"dependencies\":{\"manifests\":{\"package.json\":{}}},\"metadata\":{}}}</Actions>" {
// 		t.Errorf("Output doesn't match expected: %v", output)
// 	}
// }
