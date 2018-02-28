package pullrequest

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/dependencies-io/deps/internal/git"
	"github.com/dependencies-io/deps/internal/schema"
)

// Pullrequest stores the basic data
type Pullrequest struct {
	Branch            string
	Title             string
	Body              string
	DefaultBaseBranch string
	Dependencies      *schema.Dependencies
	Action            *schema.Action
}

// NewPullrequestFromEnv creates a Pullrequest using env variables
func NewPullrequestFromJSONPathAndEnv(dependenciesJSONPath string) (*Pullrequest, error) {
	branch, err := git.GetJobBranchName()
	if err != nil {
		return nil, err
	}

	dependencies, err := schema.NewDependenciesFromJSONPath(dependenciesJSONPath)
	if err != nil {
		return nil, err
	}

	title, err := dependencies.GenerateTitle()
	if err != nil {
		return nil, err
	}

	body, err := dependencies.GenerateBody()
	if err != nil {
		return nil, err
	}

	return &Pullrequest{
		Branch:            branch,
		Title:             title,
		Body:              body,
		DefaultBaseBranch: os.Getenv("GIT_BRANCH"),
		Dependencies:      dependencies,
		Action:            &schema.Action{Metadata: map[string]interface{}{}},
	}, nil
}

func (pr *Pullrequest) Create() error {
	return nil
}

func (pr *Pullrequest) DoRelated() error {
	return nil
}

func (pr *Pullrequest) GetActionsJSON() (string, error) {
	if pr.Action.Name == "" {
		return "", errors.New("Action name must not be empty")
	}

	output := map[string]interface{}{
		pr.Action.Name: map[string]interface{}{
			"metadata":     pr.Action.Metadata,
			"dependencies": pr.Dependencies,
		},
	}

	jsonOutput, err := json.Marshal(output)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("<Actions>%v</Actions>", string(jsonOutput)), nil
}

// OutputActions sends the Action to stdout
func (pr *Pullrequest) OutputActions() error {
	output, err := pr.GetActionsJSON()
	if err != nil {
		return err
	}

	if output == "" {
		return errors.New("no <Actions> to output")
	}

	fmt.Println(output)

	return nil
}
