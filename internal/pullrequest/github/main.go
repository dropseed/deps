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
)

// PullRequest stores additional GitHub specific data
type PullRequest struct {
	Base         string
	Head         string
	Title        string
	Body         string
	Dependencies *schema.Dependencies
	Config       *config.Dependency

	RepoOwnerName string
	RepoName      string
	RepoFullName  string
	APIToken      string
}

func NewPullRequest(base string, head string, deps *schema.Dependencies, cfg *config.Dependency) (*PullRequest, error) {
	fullName, err := getRepoFullName()
	if err != nil {
		return nil, err
	}
	parts := strings.Split(fullName, "/")
	owner := parts[0]
	repo := parts[1]

	return &PullRequest{
		Base:          base,
		Head:          head,
		Title:         deps.Title,
		Body:          deps.Description,
		Dependencies:  deps,
		Config:        cfg,
		RepoOwnerName: owner,
		RepoName:      repo,
		RepoFullName:  fullName,
		APIToken:      getAPIToken(),
	}, nil
}

func (pr *PullRequest) request(verb string, url string, input []byte) (*http.Response, string, error) {
	client := &http.Client{}

	req, err := http.NewRequest(verb, url, bytes.NewBuffer(input))
	if err != nil {
		return nil, "", err
	}

	req.Header.Add("Authorization", "token "+pr.APIToken)
	req.Header.Add("User-Agent", "deps")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	return resp, string(body), err
}

func (pr *PullRequest) pullsURL() string {
	apiBase := "https://api.github.com" // or from setting/env
	return fmt.Sprintf("%s/repos/%s/pulls", apiBase, pr.RepoFullName)
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

func (pr *PullRequest) createPR() (map[string]interface{}, error) {
	pullrequestData, err := pr.getCreateJSONData()
	if err != nil {
		return nil, err
	}

	output.Debug("Creating pull request:\n%s", pullrequestData)

	resp, body, err := pr.request("POST", pr.pullsURL(), pullrequestData)
	if err != nil {
		return nil, err
	}

	if strings.Index(string(body), "pull request already exists") != -1 {
		output.Event("Pull request already exists")
		return nil, nil
	} else if resp.StatusCode != 201 {
		return nil, fmt.Errorf("Failed to create pull request:\\n\n%s\n\n%+v\n\n%+v", body, resp.Request, resp)
	}

	output.Event("Created pull request")
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return nil, err
	}

	return data, nil
}

func (pr *PullRequest) getExisting() (map[string]interface{}, error) {
	params := fmt.Sprintf("?head=%s:%s&base=%s", pr.RepoOwnerName, pr.Head, pr.Base)
	resp, body, err := pr.request("GET", pr.pullsURL()+params, nil)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("List pull requests API returned %d", resp.StatusCode)
	}

	var data []map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return nil, err
	}

	if len(data) != 1 {
		return nil, fmt.Errorf("Found %d matches for existing pull request", len(data))
	}

	return data[0], nil
}

// Create performs the creation of the PR on GitHub
func (pr *PullRequest) CreateOrUpdate() error {
	// check the optional settings now, before actually creating the PR (which we'll have to update)

	labels := pr.Config.GetSetting("github_labels")
	assignees := pr.Config.GetSetting("github_assignees")
	milestone := pr.Config.GetSetting("github_milestone")

	fmt.Printf("Preparing to open GitHub pull request for %v\n", pr.RepoFullName)

	// TODO if pr exists then update original comment (if diff from original)
	// check return code on status?
	data, err := pr.createPR()
	if err != nil {
		return err
	}
	if data == nil {
		data, err = pr.getExisting()
		if err != nil {
			return err
		}
	}

	// pr has been created at this point, now have to add meta fields in
	// another request
	issueURL, _ := data["issue_url"].(string)
	htmlURL, _ := data["html_url"].(string)
	issueTitle, _ := data["title"].(string)
	issueBody, _ := data["body"].(string)

	if labels != nil || assignees != nil || milestone != nil || pr.Title != issueTitle || pr.Body != issueBody {
		issueMap := make(map[string]interface{})

		// Make sure these are correct and up-to-date
		issueMap["title"] = pr.Title
		issueMap["body"] = pr.Body

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

		resp, _, err := pr.request("PATCH", issueURL, issueData)
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return fmt.Errorf("failed to update pull request: %+v", resp)
		}

		output.Event("Successfully updated PR fields on %v\n", htmlURL)
	}

	return nil
}
