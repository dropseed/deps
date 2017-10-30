package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type pullrequest struct {
	branch            string
	title             string
	body              string
	gitHost           string
	defaultBaseBranch string
}

func createGitHubPullRequest(pr pullrequest) {
	apiToken := os.Getenv("GITHUB_API_TOKEN")
	repoFullName := os.Getenv("GITHUB_REPO_FULL_NAME")
	pullrequestsUrl := fmt.Sprintf("https://api.github.com/repos/%v/pulls", repoFullName)

	fmt.Printf("Preparing to open GitHub pull request for %v\n", repoFullName)

	client := &http.Client{}

	var base string
	if base = os.Getenv("SETTING_GITHUB_BASE_BRANCH"); base == "" {
		base = pr.defaultBaseBranch
	}

	pullrequestMap := map[string]string{
		"title": pr.title,
		"head":  pr.branch,
		"base":  base,
		"body":  pr.body,
	}
	fmt.Printf("%+v\n", pullrequestMap)
	pullrequestData, _ := json.Marshal(pullrequestMap)

	req, err := http.NewRequest("POST", pullrequestsUrl, bytes.NewBuffer(pullrequestData))
	req.Header.Add("Authorization", "token "+apiToken)
	req.Header.Add("User-Agent", "dependencies.io pullrequest")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 201 {
		fmt.Printf("Failed to create pull request: %+v\n", resp)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		panic(err)
	}

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		panic(err)
	}

	// pr has been created at this point, now have to add meta fields in
	// another request
	issueUrl, _ := data["issue_url"].(string)
	htmlUrl, _ := data["html_url"].(string)
	fmt.Printf("Successfully created %v\n", htmlUrl)

	var labels []string
	if labelsEnv := os.Getenv("SETTING_GITHUB_LABELS"); labelsEnv != "" {
		if err := json.Unmarshal([]byte(labelsEnv), &labels); err != nil {
			panic(err)
		}
	}

	var assignees []string
	if assigneesEnv := os.Getenv("SETTING_GITHUB_ASSIGNEES"); assigneesEnv != "" {
		if err := json.Unmarshal([]byte(assigneesEnv), &assignees); err != nil {
			panic(err)
		}
	}

	var milestone int64
	if milestoneEnv := os.Getenv("SETTING_GITHUB_MILESTONE"); milestoneEnv != "" {
		milestone, err = strconv.ParseInt(milestoneEnv, 10, 32)
		if err != nil {
			panic(err)
		}
	}

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

		req, err := http.NewRequest("PATCH", issueUrl, bytes.NewBuffer(issueData))
		req.Header.Add("Authorization", "token "+apiToken)
		req.Header.Add("User-Agent", "dependencies.io pullrequest")
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			panic(err)
		}

		if resp.StatusCode != 200 {
			fmt.Printf("Failed to create pull request: %+v\n", resp)
			os.Exit(1)
		}

		fmt.Printf("Successfully updated PR fields on %v\n", htmlUrl)
	}
}

func createGitLabPullRequest(pr pullrequest) {
	apiToken := os.Getenv("GITLAB_API_TOKEN")
	projectApiUrl := os.Getenv("GITLAB_API_URL")

	fmt.Printf("Preparing to open GitLab merge request for %v\n", projectApiUrl)

	client := &http.Client{}

	var base string
	if base = os.Getenv("SETTING_GITLAB_TARGET_BRANCH"); base == "" {
		base = pr.defaultBaseBranch
	}

	pullrequestMap := make(map[string]interface{})
	pullrequestMap["title"] = pr.title
	pullrequestMap["source_branch"] = pr.branch
	pullrequestMap["target_branch"] = base
	pullrequestMap["description"] = pr.body

	if assigneeIdEnv := os.Getenv("SETTING_GITLAB_ASSIGNEE_ID"); assigneeIdEnv != "" {
		var err interface{}
		pullrequestMap["assignee_id"], err = strconv.ParseInt(assigneeIdEnv, 10, 32)
		if err != nil {
			panic(err)
		}
	}

	if labelsEnv := os.Getenv("SETTING_GITLAB_LABELS"); labelsEnv != "" {
		var labels *[]string
		if err := json.Unmarshal([]byte(labelsEnv), &labels); err != nil {
			panic(err)
		}
		pullrequestMap["labels"] = strings.Join(*labels, ",")
	}

	// TODO is it really supposed to be milestone ID instead of IID? How are you supposed to know that?!
	// if milestoneIdEnv := os.Getenv("SETTING_GITLAB_MILESTONE_ID"); milestoneIdEnv != "" {
	//     var err interface{}
	//     pullrequestMap["milestone_id"], err = strconv.ParseInt(milestoneIdEnv, 10, 32)
	//     if err != nil {
	//         panic(err)
	//     }
	// }

	if targetProjectIdEnv := os.Getenv("SETTING_GITLAB_TARGET_PROJECT_ID"); targetProjectIdEnv != "" {
		var err interface{}
		pullrequestMap["target_project_id"], err = strconv.ParseInt(targetProjectIdEnv, 10, 32)
		if err != nil {
			panic(err)
		}
	}

	if removeSourceBranchEnv := os.Getenv("SETTING_GITLAB_REMOVE_SOURCE_BRANCH"); removeSourceBranchEnv != "" {
		var removeSourceBranch *bool
		if err := json.Unmarshal([]byte(removeSourceBranchEnv), &removeSourceBranch); err != nil {
			panic(err)
		}
		pullrequestMap["remove_source_branch"] = *removeSourceBranch
	}

	fmt.Printf("%+v\n", pullrequestMap)
	pullrequestData, _ := json.Marshal(pullrequestMap)

	req, err := http.NewRequest("POST", projectApiUrl+"/merge_requests", bytes.NewBuffer(pullrequestData))
	req.Header.Add("PRIVATE-TOKEN", apiToken)
	req.Header.Add("User-Agent", "dependencies.io pullrequest")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 201 {
		fmt.Printf("Failed to create merge request: %+v\n", resp)
		os.Exit(1)
	}

	fmt.Printf("Successfully created GitLab merge request for %v\n", projectApiUrl)
}

func main() {

	branch := flag.String("branch", "", "branch that pull request will be created from")
	title := flag.String("title", "", "pull request title")
	body := flag.String("body", "", "pull request body")
	flag.Parse()

	if *branch == "" {
		fmt.Printf("\"branch\" is required")
		os.Exit(1)
	}

	if *title == "" {
		fmt.Printf("\"title\" is required")
		os.Exit(1)
	}

	if *body == "" {
		fmt.Printf("\"body\" is required")
		os.Exit(1)
	}

	// look for additional user content to add to the body
	if pullrequestNotes := os.Getenv("SETTING_PULLREQUEST_NOTES"); pullrequestNotes != "" {
		*body = strings.TrimSpace(pullrequestNotes) + "\n\n---\n\n" + *body
	}

	pr := pullrequest{
		branch:            *branch,
		title:             *title,
		body:              *body,
		gitHost:           os.Getenv("GIT_HOST"),
		defaultBaseBranch: os.Getenv("GIT_BRANCH"),
	}

	fmt.Printf("Creating pull request for %v\n", pr.branch)

	switch strings.ToLower(pr.gitHost) {
	case "github":
		createGitHubPullRequest(pr)
	case "gitlab":
		createGitLabPullRequest(pr)
	default:
		fmt.Printf("Unknown GIT_HOST \"%v\"\n", pr.gitHost)
		os.Exit(1)
	}
}
