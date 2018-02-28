package gittest

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
	pr, prBase := getPR("../../schema/testdata/single_dependency.json")

	if pr.Pullrequest != prBase {
		t.Error("Pullrequest value incorrect")
	}
}

func TestCreateTestEnv(t *testing.T) {
	pr, _ := getPR("../../schema/testdata/single_dependency.json")
	err := pr.Create()
	if err != nil {
		t.Error(err)
	}
}

func TestCreateProductionEnv(t *testing.T) {
	os.Setenv("DEPENDENCIES_ENV", "production")
	pr, _ := getPR("../../schema/testdata/single_dependency.json")
	err := pr.Create()
	if err != nil {
		t.Error(err)
	}
}

// func TestGetActionsJSON(t *testing.T) {
// 	pr, _ := getPR("../../schema/testdata/single_dependency.json")
// 	pr.Create()
// 	pr.Pullrequest.Config.DependenciesJSON = "./testdata/action_dependencies.json"
// 	output, err := pr.GetActionsJSON()
// 	if err != nil {
// 		t.Error(err)
// 		return
// 	}
// 	if output != "<Actions>{\"PR #0\":{\"dependencies\":{\"manifests\":{\"test\":{}}},\"metadata\":{\"foo\":\"bar\"}}}</Actions>" {
// 		t.Errorf("Output doesn't match expected: %v", output)
// 	}
// }
