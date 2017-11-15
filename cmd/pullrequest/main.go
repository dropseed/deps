package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dependencies-io/pullrequest/internal/app/config"
	"github.com/dependencies-io/pullrequest/internal/app/env"
	"github.com/dependencies-io/pullrequest/internal/app/pullrequest"
	"github.com/dependencies-io/pullrequest/internal/pkg/github"
	"github.com/dependencies-io/pullrequest/internal/pkg/gitlab"
)

func main() {

	config := config.Config{}
	config.LoadEnvSettings()
	config.LoadFlags()
	configErr := config.Validate()
	if configErr != nil {
		panic(configErr)
	}

	title, err := config.TitleFromConfig()
	if err != nil {
		panic(err)
	}

	body, err := config.BodyFromConfig()
	if err != nil {
		panic(err)
	}

	branch := config.Flags.Branch

	// The PR gets the final copy of the basic data sent to it, plus the config for additional options
	prBase := pullrequest.NewPullrequestFromEnv(branch, title, body, &config)

	switch gitHost := os.Getenv("GIT_HOST"); strings.ToLower(gitHost) {
	case "github":
		pr := github.NewPullRequestFromEnv(prBase)
		err := pr.Create()
		if err != nil {
			panic(err)
		}
		err = pr.DoRelated()
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
		if env.IsProduction() {
			os.Exit(1)
		}
	}

	if !env.IsProduction() {
		fmt.Printf("pullrequest exiting successfullly in \"%v\" environment\n", env.GetCurrentEnv())
	}
}
