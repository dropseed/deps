package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dropseed/deps/internal/schemaext"

	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/pkg/schema"

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
	nodeID        string
	apiURL        string
	apiBaseURL    string
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
		Title:         schemaext.TitleForDeps(deps),
		Body:          schemaext.DescriptionForDeps(deps),
		Dependencies:  deps,
		Config:        cfg,
		RepoOwnerName: owner,
		RepoName:      repo,
		RepoFullName:  fullName,
		APIToken:      getAPIToken(),
		apiBaseURL:    getAPIBaseURL(),
	}, nil
}

func (pr *PullRequest) GetSetting(name string) interface{} {
	return pr.Config.GetSettingForSchema(name, pr.Dependencies)
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
	return fmt.Sprintf("%s/repos/%s/pulls", pr.apiBaseURL, pr.RepoFullName)
}

func (pr *PullRequest) getCreateJSONData() ([]byte, error) {
	base := pr.Base
	if override := pr.GetSetting("github_base_branch"); override != nil {
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

	labels := pr.GetSetting("github_labels")
	assignees := pr.GetSetting("github_assignees")
	milestone := pr.GetSetting("github_milestone")

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

	// Store these for future API calls
	pr.nodeID = data["node_id"].(string)
	pr.apiURL = data["url"].(string)

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

	if automerge := pr.GetSetting("github_automerge"); automerge != nil {

		currentAutomerge := data["auto_merge"]

		automergeMethod, ok := automerge.(string)
		if !ok {
			automergeBool, ok := automerge.(bool)
			if !ok {
				return fmt.Errorf("github_automerge must be a string (\"merge\", \"rebase\", or \"squash\") or a boolean")
			}
			if automergeBool {
				automergeMethod = "squash"
			} else {
				automergeMethod = "none"
			}
		}

		if currentAutomerge == nil && automergeMethod != "none" {
			output.Event("Enabling \"%s\" automerge on %v\n", automergeMethod, htmlURL)
			if err := pr.enableAutomerge(automergeMethod); err != nil {
				return err
			}
		} else if currentAutomerge != nil && automergeMethod == "none" {
			output.Event("Disabling automerge on %v\n", htmlURL)
			if err := pr.disableAutomerge(); err != nil {
				return err
			}
		} else {
			output.Event("No automerge change to make on %v\nRequested: %s\nExisting: %+v", automergeMethod, currentAutomerge, htmlURL)
		}
	}

	return nil
}

func (pr *PullRequest) graphqlRequest(body map[string]interface{}) error {
	bodyData, _ := json.Marshal(body)
	resp, respBody, err := pr.request("POST", pr.apiBaseURL+"/graphql", bodyData)

	output.Debug("%+v", resp)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("GraphQL API returned %d\n%s", resp.StatusCode, respBody)
	}

	var respData map[string]interface{}
	if err := json.Unmarshal([]byte(respBody), &respData); err != nil {
		return err
	}

	if respData["errors"] != nil {
		return fmt.Errorf("GraphQL API returned errors:\n%s", respBody)
	}

	return err
}

func (pr *PullRequest) enableAutomerge(mergeMethod string) error {
	err := pr.graphqlRequest(map[string]interface{}{
		"query": `mutation($input: EnablePullRequestAutoMergeInput!) {
			enablePullRequestAutoMerge(input: $input) {
				clientMutationId
			}
		}`,
		"variables": map[string]interface{}{
			"input": map[string]interface{}{
				"pullRequestId": pr.nodeID,
				"mergeMethod":   strings.ToUpper(mergeMethod),
			},
		},
	})

	if err != nil && strings.Contains(err.Error(), "Pull request is in clean status") {
		output.Event("Pull request is in clean status, automerge can't be enabled so we will merge manually")
		return pr.merge(mergeMethod)
	} else if err != nil {
		return err
	}

	return nil
}

func (pr *PullRequest) disableAutomerge() error {
	return pr.graphqlRequest(map[string]interface{}{
		"query": `mutation {
			disablePullRequestAutoMerge(input: {pullRequestId: "%s"}) {
				clientMutationId
			}
		}`,
		"variables": map[string]interface{}{
			"pullRequestId": pr.nodeID,
		},
	})
}

func (pr *PullRequest) merge(mergeMethod string) error {
	data := map[string]interface{}{
		"merge_method": strings.ToLower(mergeMethod),
	}
	dataStr, _ := json.Marshal(data)
	resp, body, err := pr.request("PUT", pr.apiURL+"/merge", dataStr)

	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to merge PR:\n%s", body)
	}

	return nil
}
