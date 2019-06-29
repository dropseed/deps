package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/dropseed/deps/internal/schema"

	"github.com/dropseed/deps/internal/env"
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/pullrequest"
)

// PullRequest stores additional GitHub specific data
type PullRequest struct {
	// directly use the properties of base Pullrequest
	*pullrequest.Pullrequest
	RepoOwnerName string
	RepoName      string
	RepoFullName  string
	APIToken      string
	Number        int
	CreatedAt     string
}

// NewPullrequestFromDependenciesEnv creates a PullRequest
func NewPullrequestFromDependenciesEnv(deps *schema.Dependencies, branch string) (*PullRequest, error) {
	prBase, err := pullrequest.NewPullrequestFromEnv(deps, branch)
	if err != nil {
		return nil, err
	}

	fullName, err := getRepoFullName()
	if err != nil {
		return nil, err
	}
	parts := strings.Split(fullName, "/")
	owner := parts[0]
	repo := parts[1]

	return &PullRequest{
		Pullrequest:   prBase,
		RepoOwnerName: owner,
		RepoName:      repo,
		RepoFullName:  fullName,
		APIToken:      GetAPIToken(),
	}, nil
}

func (pr *PullRequest) getCreateJSONData() ([]byte, error) {
	var base string
	// TODO settings can probably be given directly now from config
	if base = env.GetSetting("GITHUB_BASE_BRANCH", ""); base == "" {
		base = pr.DefaultBaseBranch
	}

	// TODO does this need to change?
	body, err := dereferenceGitHubIssueLinks(pr.Body)
	if err != nil {
		return nil, err
	}

	pullrequestMap := map[string]string{
		"title": pr.Title,
		"head":  pr.Branch,
		"base":  base,
		"body":  body,
	}

	pullrequestData, _ := json.Marshal(pullrequestMap)
	return pullrequestData, nil
}

func (pr *PullRequest) addHeadersToRequest(req *http.Request) {
	req.Header.Add("Authorization", "token "+pr.APIToken)
	req.Header.Add("User-Agent", "deps")
	req.Header.Set("Content-Type", "application/json")
}

func (pr *PullRequest) createPR() (map[string]interface{}, error) {
	// open the actual PR, first of two API calls

	pullrequestData, err := pr.getCreateJSONData()
	if err != nil {
		return nil, err
	}

	output.Debug("Creating pull request:\n%s", pullrequestData)

	apiBase := "https://api.github.com"
	// TODO can maybe automatically get this from remote instead?
	// if base := env.GetSetting("DEPS_GITHUB_API_BASE_URL", ""); base == "" {
	// 	apiBase = base
	// }

	pullrequestsURL := fmt.Sprintf("%s/repos/%v/pulls", apiBase, pr.RepoFullName)

	client := &http.Client{}

	req, err := http.NewRequest("POST", pullrequestsURL, bytes.NewBuffer(pullrequestData))
	if err != nil {
		return nil, err
	}

	pr.addHeadersToRequest(req)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 201 {
		return nil, fmt.Errorf("Failed to create pull request:\n\n%s\n\n%+v", body, resp)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}

// Create performs the creation of the PR on GitHub
func (pr *PullRequest) CreateOrUpdate() error {
	// check the optional settings now, before actually creating the PR (which we'll have to update)
	var labels []string
	if labelsEnv := env.GetSetting("GITHUB_LABELS", ""); labelsEnv != "" {
		if err := json.Unmarshal([]byte(labelsEnv), &labels); err != nil {
			return err
		}
	}

	var assignees []string
	if assigneesEnv := env.GetSetting("GITHUB_ASSIGNEES", ""); assigneesEnv != "" {
		if err := json.Unmarshal([]byte(assigneesEnv), &assignees); err != nil {
			return err
		}
	}

	var milestone int64
	if milestoneEnv := env.GetSetting("GITHUB_MILESTONE", ""); milestoneEnv != "" {
		var err error
		milestone, err = strconv.ParseInt(milestoneEnv, 10, 32)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Preparing to open GitHub pull request for %v\n", pr.RepoFullName)

	// TODO if pr exists then update original comment (if diff from original)

	data, err := pr.createPR()
	if err != nil {
		return err
	}

	pr.Number = int(data["number"].(float64))
	pr.CreatedAt = data["created_at"].(string)

	// set the Action info for reporting back to dependencies.io
	pr.Action.Name = fmt.Sprintf("PR #%v", pr.Number)
	pr.Action.Metadata["github_pull_request"] = data

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

		pr.addHeadersToRequest(req)

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
