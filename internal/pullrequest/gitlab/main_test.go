package gitlab

import (
	"os"
	"testing"

	"github.com/dependencies-io/pullrequest/internal/pullrequest"
)

func getMR(dependenciesJSONPath string) (*MergeRequest, *pullrequest.Pullrequest) {
	os.Setenv("JOB_ID", "test")
	pr, err := NewPullrequestFromDependenciesJSONPathAndEnv(dependenciesJSONPath)
	if err != nil {
		panic(err)
	}
	return pr, pr.Pullrequest
}

func TestNewMergeRequestFromEnv(t *testing.T) {
	os.Setenv("GITLAB_API_URL", "testurl")
	os.Setenv("GITLAB_API_TOKEN", "testtoken")

	pr, prBase := getMR("../../schema/testdata/single_dependency.json")

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
	pr, _ := getMR("../../schema/testdata/single_dependency.json")
	err := pr.Create()
	if err != nil {
		t.Fail()
	}
}

func TestCreateProductionEnv(t *testing.T) {
	// this will try to send an actual API call to github.com, and fail
	os.Setenv("DEPENDENCIES_ENV", "production")
	pr, _ := getMR("../../schema/testdata/single_dependency.json")
	err := pr.Create()
	if err == nil {
		t.Fail()
	}
}

// func TestGetActionsJSON(t *testing.T) {
// 	os.Setenv("DEPENDENCIES_ENV", "test")
// 	pr, _ := getMR("./testdata/action_dependencies.json")
// 	pr.Create()
// 	output, err := pr.GetActionsJSON()
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if output != "<Actions>{\"MR !0\":{\"dependencies\":{\"manifests\":{\"package.json\":{}}},\"metadata\":{}}}</Actions>" {
// 		t.Errorf("Output doesn't match expected: %v", output)
// 	}
// }
