package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/schemaext"
	"github.com/dropseed/deps/pkg/schema"
)

type PullRequest struct {
	Base         string
	Head         string
	Title        string
	Body         string
	Dependencies *schema.Dependencies
	Config       *config.Dependency

	ProjectAPIURL string
	APIUsername   string
	APIPassword   string
}

func NewPullRequest(base string, head string, deps *schema.Dependencies, cfg *config.Dependency) (*PullRequest, error) {
	apiURL, err := getProjectAPIURL()
	if err != nil {
		return nil, err
	}

	return &PullRequest{
		Base:          base,
		Head:          head,
		Title:         schemaext.TitleForDeps(deps),
		Body:          schemaext.DescriptionForDeps(deps),
		Dependencies:  deps,
		Config:        cfg,
		ProjectAPIURL: apiURL,
		APIUsername:   getAPIUsername(),
		APIPassword:   getAPIPassword(),
	}, nil
}

func (pr *PullRequest) request(verb string, url string, input []byte) (*http.Response, string, error) {
	client := &http.Client{}

	req, err := http.NewRequest(verb, url, bytes.NewBuffer(input))
	if err != nil {
		return nil, "", err
	}

	req.SetBasicAuth(pr.APIUsername, pr.APIPassword)
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

// // Create will create the pull request on Bitbucket
func (pr *PullRequest) CreateOrUpdate() error {
	output.Debug("Preparing to open Bitbucket pull request for %v\n", pr.ProjectAPIURL)

	pullrequestMap := pr.getPullRequestOptions()
	output.Debug("%+v\n", pullrequestMap)
	pullrequestData, _ := json.Marshal(pullrequestMap)

	url := pr.ProjectAPIURL + "/pullrequests"
	output.Debug("Creating pull request at %s", url)

	resp, body, err := pr.request("POST", url, pullrequestData)
	if err != nil {
		return err
	}

	if resp.StatusCode == 201 {
		output.Event("Successfully created Bitbucket pull request for %v\n", pr.ProjectAPIURL)
		return nil
	}

	return fmt.Errorf("Failed to create pull request: %s", body)
}

func (pr *PullRequest) getPullRequestOptions() map[string]interface{} {
	base := pr.Base
	if target := pr.Config.GetSetting("bitbucket_destination"); target != nil {
		base = target.(string)
	}

	pullrequestMap := make(map[string]interface{})
	pullrequestMap["title"] = pr.Title
	pullrequestMap["source"] = map[string]interface{}{
		"branch": map[string]string{
			"name": pr.Head,
		},
	}
	pullrequestMap["destination"] = map[string]interface{}{
		"branch": map[string]string{
			"name": base,
		},
	}
	pullrequestMap["description"] = pr.Body

	otherFields := []string{
		"close_source_branch",
		"reviewers",
	}

	for _, f := range otherFields {
		if s := pr.Config.GetSetting(fmt.Sprintf("bitbucket_%s", f)); s != nil {
			pullrequestMap[f] = s
		}
	}

	return pullrequestMap
}
