package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/schema"

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
func NewPullrequest(base string, head string, deps *schema.Dependencies, cfg *config.Dependency) (*PullRequest, error) {
	prBase, err := pullrequest.NewPullrequest(base, head, deps, cfg)
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
	base := pr.Base
	if override := pr.Config.GetSetting("github_base_branch"); override != nil {
		base = override.(string)
	}

	// TODO does this need to change?
	body, err := dereferenceGitHubIssueLinks(pr.Body)
	if err != nil {
		return nil, err
	}

	pullrequestMap := map[string]string{
		"title": pr.Title,
		"head":  pr.Head,
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

	labels := pr.Config.GetSetting("github_labels")
	assignees := pr.Config.GetSetting("github_assignees")
	milestone := pr.Config.GetSetting("github_assignees")

	fmt.Printf("Preparing to open GitHub pull request for %v\n", pr.RepoFullName)

	// TODO if pr exists then update original comment (if diff from original)
	// check return code on status?
	data, err := pr.createPR()
	if err != nil {
		return err
	}

	pr.Number = int(data["number"].(float64))
	pr.CreatedAt = data["created_at"].(string)

	// pr has been created at this point, now have to add meta fields in
	// another request
	issueURL, _ := data["issue_url"].(string)
	htmlURL, _ := data["html_url"].(string)
	fmt.Printf("Successfully created %v\n", htmlURL)

	if labels != nil || assignees != nil || milestone != nil {
		issueMap := make(map[string]interface{})

		if labels != nil {
			issueMap["labels"] = labels
		}
		if assignees != nil {
			issueMap["assignees"] = assignees
		}
		if milestone != nil {
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
