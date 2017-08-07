package main

import (
    "os"
    "bytes"
    "fmt"
    "strings"
    "net/http"
    "io/ioutil"
    "flag"
    "encoding/json"
    "strconv"
)


type pullrequest struct {
    branch string
    title string
    body string
    gitHost string
    baseBranch string
}

func createGitHubPullRequest(pr pullrequest) {
    apiToken := os.Getenv("GITHUB_API_TOKEN")
    repoFullName := os.Getenv("GITHUB_REPO_FULL_NAME")
    pullrequestsUrl := fmt.Sprintf("https://api.github.com/repos/%v/pulls", repoFullName)

    fmt.Printf("Preparing to open GitHub pull request for %v\n", repoFullName)

    client := &http.Client{}

    pullrequestMap := map[string]string{
        "title": pr.title,
        "head": pr.branch,
        "base": pr.baseBranch,
        "body": pr.body,
    }
    fmt.Printf("%+v\n", pullrequestMap)
    pullrequestData, _ := json.Marshal(pullrequestMap)

    req, err := http.NewRequest("POST", pullrequestsUrl, bytes.NewBuffer(pullrequestData))
    req.Header.Add("Authorization", "token " + apiToken)
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
        req.Header.Add("Authorization", "token " + apiToken)
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

    pr := pullrequest{
        branch: *branch,
        title: *title,
        body: *body,
        gitHost: os.Getenv("GIT_HOST"),
        baseBranch: os.Getenv("GIT_BRANCH"),
    }

    fmt.Printf("Creating pull request for %v\n", pr.branch)

    switch strings.ToLower(pr.gitHost) {
        case "github":
            createGitHubPullRequest(pr)
        case "gitlab":
            fmt.Printf("gitlab\n")
        default:
            fmt.Printf("Unknown GIT_HOST \"%v\"\n", pr.gitHost)
            os.Exit(1)
    }
}
