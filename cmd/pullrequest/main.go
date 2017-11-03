package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dependencies-io/pullrequest/internal/app/config"
	"github.com/dependencies-io/pullrequest/internal/app/pullrequest"
	"github.com/dependencies-io/pullrequest/internal/pkg/github"
	"github.com/dependencies-io/pullrequest/internal/pkg/gitlab"
)

const maxBodyLength = 65535

func main() {

	configFlags := parseFlags()

	if configFlags.branch == "" {
		fmt.Printf("\"branch\" is required")
		os.Exit(1)
	}

	title, err := configFlags.titleFromConfigFlags()
	if err != nil {
		panic(err)
	}

	body, err := configFlags.bodyFromConfigFlags()
	if err != nil {
		panic(err)
	}

	// look for additional user content to add to the body
	if pullrequestNotes := os.Getenv("SETTING_PULLREQUEST_NOTES"); pullrequestNotes != "" {
		body = strings.TrimSpace(pullrequestNotes) + "\n\n---\n\n" + body
	}

	// trim the pr body string to a max of this size,
	// should rarely happen but this way API call should still be success
	if len(body) > maxBodyLength {
		body = body[:maxBodyLength]
	}

	config := config.NewConfigFromEnv()

	prBase := pullrequest.NewPullrequestFromEnv(configFlags.branch, title, body, config)

	switch gitHost := os.Getenv("GIT_HOST"); strings.ToLower(gitHost) {
	case "github":
		pr := github.NewPullRequestFromEnv(prBase)
		err := pr.Create()
		if err != nil {
			panic(err)
		}
	case "gitlab":
		pr := gitlab.NewMergeRequestFromEnv(prBase)
		err := pr.Create()
		if err != nil {
			panic(err)
		}
	default:
		fmt.Printf("Unknown GIT_HOST \"%v\"\n", gitHost)
		if config.IsProduction() {
			os.Exit(1)
		}
	}

	if !config.IsProduction() {
		fmt.Printf("pullrequest exiting successfullly in \"%v\" environment\n", config.Env)
	}
}
