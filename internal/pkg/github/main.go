package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/dependencies-io/pullrequest/internal/app/pullrequest"
)

// PullRequest stores additional GitHub specific data
type PullRequest struct {
	// directly use the properties of base Pullrequest
	*pullrequest.Pullrequest
	RepoFullName string
	APIToken     string
}

// NewPullRequestFromEnv creates a PullRequest using env variables
func NewPullRequestFromEnv(prBase *pullrequest.Pullrequest) *PullRequest {
	return &PullRequest{
		Pullrequest:  prBase,
		RepoFullName: os.Getenv("GITHUB_REPO_FULL_NAME"),
		APIToken:     os.Getenv("GITHUB_API_TOKEN"),
	}
}

func (pr PullRequest) getCreateJSONData() []byte {
	var base string
	if base = os.Getenv("SETTING_GITHUB_BASE_BRANCH"); base == "" {
		base = pr.DefaultBaseBranch
	}

	pullrequestMap := map[string]string{
		"title": pr.Title,
		"head":  pr.Branch,
		"base":  base,
		"body":  pr.Body,
	}
	fmt.Printf("%+v\n", pullrequestMap)
	pullrequestData, _ := json.Marshal(pullrequestMap)
	return pullrequestData
}

func (pr PullRequest) createPR() (map[string]interface{}, error) {
	// open the actual PR, first of two API calls

	pullrequestData := pr.getCreateJSONData()
	pullrequestsURL := fmt.Sprintf("https://api.github.com/repos/%v/pulls", pr.RepoFullName)

	client := &http.Client{}

	req, err := http.NewRequest("POST", pullrequestsURL, bytes.NewBuffer(pullrequestData))
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "token "+pr.APIToken)
	req.Header.Add("User-Agent", "dependencies.io pullrequest")
	req.Header.Set("Content-Type", "application/json")

	if !pr.Config.IsProduction() {
		fmt.Printf("Skipping GitHub API call due to \"%v\" env\n", pr.Config.Env)
		return nil, nil
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("failed to create pull request: %+v", resp)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// Create performs the creation of the PR on GitHub
func (pr PullRequest) Create() error {
	// check the optional settings now, before actually creating the PR (which we'll have to update)
	var labels []string
	if labelsEnv := os.Getenv("SETTING_GITHUB_LABELS"); labelsEnv != "" {
		if err := json.Unmarshal([]byte(labelsEnv), &labels); err != nil {
			return err
		}
	}

	var assignees []string
	if assigneesEnv := os.Getenv("SETTING_GITHUB_ASSIGNEES"); assigneesEnv != "" {
		if err := json.Unmarshal([]byte(assigneesEnv), &assignees); err != nil {
			return err
		}
	}

	var milestone int64
	if milestoneEnv := os.Getenv("SETTING_GITHUB_MILESTONE"); milestoneEnv != "" {
		var err error
		milestone, err = strconv.ParseInt(milestoneEnv, 10, 32)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Preparing to open GitHub pull request for %v\n", pr.RepoFullName)
	data, err := pr.createPR()
	if err != nil {
		return err
	}

	// pr has been created at this point, now have to add meta fields in
	// another request
	issueURL, _ := data["issue_url"].(string)
	htmlURL, _ := data["html_url"].(string)
	fmt.Printf("Successfully created %v\n", htmlURL)

	if len(labels) > 0 || len(assignees) > 0 || milestone > 0 {
		issueMap := make(map[string]interface{})

		if len(labels) > 0 {
			issueMap["labels"] = labels
		}
		if len(assignees) > 0 {
			issueMap["assignees"] = assignees
		}
		if milestone > 0 {
			issueMap["milestone"] = milestone
		}

		fmt.Printf("%+v\n", issueMap)
		issueData, _ := json.Marshal(issueMap)

		req, err := http.NewRequest("PATCH", issueURL, bytes.NewBuffer(issueData))
		if err != nil {
			return err
		}
		req.Header.Add("Authorization", "token "+pr.APIToken)
		req.Header.Add("User-Agent", "dependencies.io pullrequest")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("failed to update pull request: %+v", resp)
		}

		fmt.Printf("Successfully updated PR fields on %v\n", htmlURL)
	}

	return nil
}
