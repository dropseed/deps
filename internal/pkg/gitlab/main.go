package gitlab

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/dependencies-io/pullrequest/internal/app/env"
	"github.com/dependencies-io/pullrequest/internal/app/pullrequest"
)

// MergeRequest stores additional GitLab specific data
type MergeRequest struct {
	// directly use the properties of base Pullrequest
	*pullrequest.Pullrequest
	ProjectAPIURL string
	APIToken      string
}

// NewMergeRequestFromEnv creates a MergeRequest using env variables
func NewMergeRequestFromEnv(prBase *pullrequest.Pullrequest) *MergeRequest {
	return &MergeRequest{
		Pullrequest:   prBase,
		ProjectAPIURL: os.Getenv("GITLAB_API_URL"),
		APIToken:      os.Getenv("GITLAB_API_TOKEN"),
	}
}

// Create will create the merge request on GitLab
func (pr *MergeRequest) Create() error {
	fmt.Printf("Preparing to open GitLab merge request for %v\n", pr.ProjectAPIURL)

	client := &http.Client{}

	var base string
	if base = os.Getenv("SETTING_GITLAB_TARGET_BRANCH"); base == "" {
		base = pr.DefaultBaseBranch
	}

	pullrequestMap := make(map[string]interface{})
	pullrequestMap["title"] = pr.Title
	pullrequestMap["source_branch"] = pr.Branch
	pullrequestMap["target_branch"] = base
	pullrequestMap["description"] = pr.Body

	if assigneeIIEnv := os.Getenv("SETTING_GITLAB_ASSIGNEE_ID"); assigneeIIEnv != "" {
		var err error
		pullrequestMap["assignee_id"], err = strconv.ParseInt(assigneeIIEnv, 10, 32)
		if err != nil {
			return err
		}
	}

	if labelsEnv := os.Getenv("SETTING_GITLAB_LABELS"); labelsEnv != "" {
		var labels *[]string
		if err := json.Unmarshal([]byte(labelsEnv), &labels); err != nil {
			return err
		}
		pullrequestMap["labels"] = strings.Join(*labels, ",")
	}

	// TODO is it really supposed to be milestone ID instead of IID? How are you supposed to know that?!
	// if milestoneIdEnv := os.Getenv("SETTING_GITLAB_MILESTONE_ID"); milestoneIdEnv != "" {
	//     var err error
	//     pullrequestMap["milestone_id"], err = strconv.ParseInt(milestoneIdEnv, 10, 32)
	//     if err != nil {
	//         return err
	//     }
	// }

	if targetProjectIDEnv := os.Getenv("SETTING_GITLAB_TARGET_PROJECT_ID"); targetProjectIDEnv != "" {
		var err error
		pullrequestMap["target_project_id"], err = strconv.ParseInt(targetProjectIDEnv, 10, 32)
		if err != nil {
			return err
		}
	}

	if removeSourceBranchEnv := os.Getenv("SETTING_GITLAB_REMOVE_SOURCE_BRANCH"); removeSourceBranchEnv != "" {
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

	if env.IsProduction() {
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != 201 {
			return fmt.Errorf("failed to create merge request: %+v", resp)
		}
		fmt.Printf("Successfully created GitLab merge request for %v\n", pr.ProjectAPIURL)
	} else {
		fmt.Printf("Skipping GitLab API call due to \"%v\" env\n", env.GetCurrentEnv())
	}

	return nil
}

// DoRelated for GitLab is not yet implemented
func (pr *MergeRequest) DoRelated() error {
	return errors.New("related PR behavior is not yet supported for GitLab")
}
