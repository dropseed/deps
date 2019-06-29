package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dropseed/deps/internal/schema"

	"github.com/dropseed/deps/internal/env"
	"github.com/dropseed/deps/internal/pullrequest"
)

// MergeRequest stores additional GitLab specific data
type MergeRequest struct {
	// directly use the properties of base Pullrequest
	*pullrequest.Pullrequest
	ProjectAPIURL string
	APIToken      string
}

// NewPullrequestFromDependenciesEnv creates a PullRequest
func NewPullrequestFromDependenciesEnv(deps *schema.Dependencies, branch string) (*MergeRequest, error) {
	prBase, err := pullrequest.NewPullrequestFromEnv(deps, branch)
	if err != nil {
		return nil, err
	}

	return &MergeRequest{
		Pullrequest:   prBase,
		ProjectAPIURL: os.Getenv("GITLAB_API_URL"),
		APIToken:      os.Getenv("GITLAB_API_TOKEN"),
	}, nil
}

// Create will create the merge request on GitLab
func (pr *MergeRequest) CreateOrUpdate() error {
	fmt.Printf("Preparing to open GitLab merge request for %v\n", pr.ProjectAPIURL)

	client := &http.Client{}

	var base string
	if base = env.GetSetting("GITLAB_TARGET_BRANCH", ""); base == "" {
		base = pr.DefaultBaseBranch
	}

	pullrequestMap := make(map[string]interface{})
	pullrequestMap["title"] = pr.Title
	pullrequestMap["source_branch"] = pr.Branch
	pullrequestMap["target_branch"] = base
	pullrequestMap["description"] = pr.Body

	if assigneeIIEnv := env.GetSetting("GITLAB_ASSIGNEE_ID", ""); assigneeIIEnv != "" {
		var err error
		pullrequestMap["assignee_id"], err = strconv.ParseInt(assigneeIIEnv, 10, 32)
		if err != nil {
			return err
		}
	}

	if labelsEnv := env.GetSetting("GITLAB_LABELS", ""); labelsEnv != "" {
		var labels *[]string
		if err := json.Unmarshal([]byte(labelsEnv), &labels); err != nil {
			return err
		}
		pullrequestMap["labels"] = strings.Join(*labels, ",")
	}

	// TODO is it really supposed to be milestone ID instead of IID? How are you supposed to know that?!
	// if milestoneIdEnv := env.GetSetting("GITLAB_MILESTONE_ID", ""); milestoneIdEnv != "" {
	//     var err error
	//     pullrequestMap["milestone_id"], err = strconv.ParseInt(milestoneIdEnv, 10, 32)
	//     if err != nil {
	//         return err
	//     }
	// }

	if targetProjectIDEnv := env.GetSetting("GITLAB_TARGET_PROJECT_ID", ""); targetProjectIDEnv != "" {
		var err error
		pullrequestMap["target_project_id"], err = strconv.ParseInt(targetProjectIDEnv, 10, 32)
		if err != nil {
			return err
		}
	}

	if removeSourceBranchEnv := env.GetSetting("GITLAB_REMOVE_SOURCE_BRANCH", ""); removeSourceBranchEnv != "" {
		var removeSourceBranch *bool
		if err := json.Unmarshal([]byte(removeSourceBranchEnv), &removeSourceBranch); err != nil {
			return err
		}
		pullrequestMap["remove_source_branch"] = *removeSourceBranch
	}

	fmt.Printf("%+v\n", pullrequestMap)
	pullrequestData, _ := json.Marshal(pullrequestMap)

	req, err := http.NewRequest("POST", pr.ProjectAPIURL+"/merge_requests", bytes.NewBuffer(pullrequestData))
	if err != nil {
		return err
	}
	req.Header.Add("PRIVATE-TOKEN", pr.APIToken)
	req.Header.Add("User-Agent", "dependencies.io pullrequest")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 201 {
		return fmt.Errorf("failed to create merge request: %+v", resp)
	}

	fmt.Printf("Successfully created GitLab merge request for %v\n", pr.ProjectAPIURL)

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return err
	}
	pr.Action.Name = fmt.Sprintf("MR !%v", int(data["iid"].(float64)))
	pr.Action.Metadata["gitlab_merge_request"] = data

	return nil
}
