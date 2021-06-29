package runner

import (
	"fmt"
	"testing"

	"github.com/dropseed/deps/pkg/schema"
)

func checkRender(template, expected string) error {
	deps, err := schema.NewDependenciesFromJSONPath("./testdata/single_dependency.json")
	if err != nil {
		panic(err)
	}
	message, err := renderCommitMessage(deps, template)
	if err != nil {
		return err
	}
	if message != expected {
		return fmt.Errorf("Message does not match: %s", message)
	}
	return nil
}

func TestCommitMessageTemplates(t *testing.T) {
	if err := checkRender("{{.Subject}}", "Update pullrequest from 0.1.0 to 0.3.0"); err != nil {
		t.Error(err)
	}
	if err := checkRender("[deps] {{.Subject}}", "[deps] Update pullrequest from 0.1.0 to 0.3.0"); err != nil {
		t.Error(err)
	}
	if err := checkRender("{{.Subject}} (skip ci)", "Update pullrequest from 0.1.0 to 0.3.0 (skip ci)"); err != nil {
		t.Error(err)
	}
	if err := checkRender("    {{.Subject}}   ", "Update pullrequest from 0.1.0 to 0.3.0"); err != nil {
		t.Error(err)
	}
	if err := checkRender("{{.Subject}}\n\nChangelog: updated", "Update pullrequest from 0.1.0 to 0.3.0\n\nChangelog: updated"); err != nil {
		t.Error(err)
	}
}
