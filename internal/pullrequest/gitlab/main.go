package gitlab

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/pkg/schema"
)

// MergeRequest stores additional GitLab specific data
type MergeRequest struct {
	Base         string
	Head         string
	Title        string
	Body         string
	Dependencies *schema.Dependencies
	Config       *config.Dependency

	ProjectAPIURL string
	APIToken      string
}

func NewMergeRequest(base string, head string, deps *schema.Dependencies, cfg *config.Dependency) (*MergeRequest, error) {
	apiURL, err := getProjectAPIURL()
	if err != nil {
		return nil, err
	}

	return &MergeRequest{
		Base:          base,
		Head:          head,
		Title:         deps.Title,
		Body:          deps.Description,
		Dependencies:  deps,
		Config:        cfg,
		ProjectAPIURL: apiURL,
		APIToken:      getAPIToken(),
	}, nil
}

func (pr *MergeRequest) request(verb string, url string, input []byte) (*http.Response, string, error) {
	client := &http.Client{}

	req, err := http.NewRequest(verb, url, bytes.NewBuffer(input))
	if err != nil {
		return nil, "", err
	}

	req.Header.Add("PRIVATE-TOKEN", pr.APIToken)
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

// // Create will create the merge request on GitLab
func (pr *MergeRequest) CreateOrUpdate() error {
	output.Debug("Preparing to open GitLab merge request for %v\n", pr.ProjectAPIURL)

	pullrequestMap := pr.getMergeRequestOptions()
	output.Debug("%+v\n", pullrequestMap)
	pullrequestData, _ := json.Marshal(pullrequestMap)

	url := pr.ProjectAPIURL + "/merge_requests"
	output.Debug("Creating merge request at %s", url)

	resp, body, err := pr.request("POST", url, pullrequestData)
	if err != nil {
		return err
	}

	if resp.StatusCode == 201 {
		output.Event("Successfully created GitLab merge request for %v\n", pr.ProjectAPIURL)
		return nil
	} else if resp.StatusCode == 409 {
		output.Event("Merge request already exists")
		var data map[string][]string
		if err := json.Unmarshal([]byte(body), &data); err != nil {
			return err
		}

		if message, hasMessage := data["message"]; hasMessage {
			pattern := regexp.MustCompile("!(\\d+)")
			matches := pattern.FindStringSubmatch(message[0])
			// finds !18 and 18...
			if len(matches) != 2 {
				return errors.New("Unable to find ID for existing merge request to update")
			}
			mrID := matches[1]
			return pr.update(mrID, pullrequestData)
		}
	}

	return fmt.Errorf("Failed to create merge request: %s", body)
}

func (pr *MergeRequest) update(iid string, data []byte) error {
	url := pr.ProjectAPIURL + "/merge_requests/" + iid
	resp, body, err := pr.request("PUT", url, data)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return fmt.Errorf("Error updating merge request:\n\n%s", body)
	}
	output.Success("Updated merge request %s", iid)
	return nil
}

func (pr *MergeRequest) getMergeRequestOptions() map[string]interface{} {
	base := pr.Base
	if target := pr.Config.GetSetting("gitlab_target_branch"); target != nil {
		base = target.(string)
	}

	pullrequestMap := make(map[string]interface{})
	pullrequestMap["title"] = pr.Title
	pullrequestMap["source_branch"] = pr.Head
	pullrequestMap["target_branch"] = base
	pullrequestMap["description"] = pr.Body

	if labels := pr.Config.GetSetting("gitlab_labels"); labels != nil {
		labelStrings := []string{}
		for _, l := range labels.([]interface{}) {
			labelStrings = append(labelStrings, l.(string))
		}
		pullrequestMap["labels"] = strings.Join(labelStrings, ",")
	}

	otherFields := []string{
		"assignee_id",
		"assignee_ids",
		"target_project_id",
		"milestone_id",
		"remove_source_branch",
		"allow_collaboration",
		"allow_maintainer_to_push",
		"squash",
	}

	for _, f := range otherFields {
		if s := pr.Config.GetSetting(fmt.Sprintf("gitlab_%s", f)); s != nil {
			pullrequestMap[f] = s
		}
	}

	return pullrequestMap
}
