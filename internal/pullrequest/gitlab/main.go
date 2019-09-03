package gitlab

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/schema"
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

// // Create will create the merge request on GitLab
func (pr *MergeRequest) CreateOrUpdate() error {
	fmt.Printf("Preparing to open GitLab merge request for %v\n", pr.ProjectAPIURL)

	client := &http.Client{}

	pullrequestMap := pr.getMergeRequestOptions()
	fmt.Printf("%+v\n", pullrequestMap)
	pullrequestData, _ := json.Marshal(pullrequestMap)

	url := pr.ProjectAPIURL + "/merge_requests"
	output.Debug("Creating merge request at %s", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(pullrequestData))
	if err != nil {
		return err
	}
	req.Header.Add("PRIVATE-TOKEN", pr.APIToken)
	req.Header.Add("User-Agent", "deps")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	// TODO if it exists already, we need to update it

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

	fmt.Printf("%+v", data)

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

	if assignee := pr.Config.GetSetting("gitlab_assignee_id"); assignee != nil {
		pullrequestMap["assignee_id"] = assignee
	}

	if labels := pr.Config.GetSetting("gitlab_labels"); labels != nil {
		labelStrings := []string{}
		for _, l := range labels.([]interface{}) {
			labelStrings = append(labelStrings, l.(string))
		}
		pullrequestMap["labels"] = strings.Join(labelStrings, ",")
	}

	// TODO is it really supposed to be milestone ID instead of IID? How are you supposed to know that?!
	// if milestoneIdEnv := env.GetSetting("GITLAB_MILESTONE_ID", ""); milestoneIdEnv != nil {
	//     var err error
	//     pullrequestMap["milestone_id"], err = strconv.ParseInt(milestoneIdEnv, 10, 32)
	//     if err != nil {
	//         return err
	//     }
	// }

	if targetProjectID := pr.Config.GetSetting("gitlab_target_project_id"); targetProjectID != nil {
		pullrequestMap["target_project_id"] = targetProjectID
	}

	if removeSourceBranch := pr.Config.GetSetting("gitlab_remove_source_branch"); removeSourceBranch != nil {
		pullrequestMap["remove_source_branch"] = removeSourceBranch
	}

	return pullrequestMap
}
