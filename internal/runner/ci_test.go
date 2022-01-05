package runner

import (
	"fmt"
	"testing"

	"github.com/dropseed/deps/pkg/schema"
)

func checkRender(jsonName, template, expected string) error {
	deps, err := schema.NewDependenciesFromJSONPath(fmt.Sprintf("./testdata/%s.json", jsonName))
	if err != nil {
		panic(err)
	}
	message, err := renderCommitMessage(deps, template)
	if err != nil {
		return err
	}
	if message != expected {
		return fmt.Errorf("Message does not match: %s\n---\n%s\n", message, expected)
	}
	return nil
}

func TestCommitMessageSubjectTemplates(t *testing.T) {
	if err := checkRender("single_dependency", "{{.Subject}}", "Update pullrequest from 0.1.0 to 0.3.0"); err != nil {
		t.Error(err)
	}
	if err := checkRender("single_dependency", "[deps] {{.Subject}}", "[deps] Update pullrequest from 0.1.0 to 0.3.0"); err != nil {
		t.Error(err)
	}
	if err := checkRender("single_dependency", "{{.Subject}} (skip ci)", "Update pullrequest from 0.1.0 to 0.3.0 (skip ci)"); err != nil {
		t.Error(err)
	}
	if err := checkRender("single_dependency", "    {{.Subject}}   ", "Update pullrequest from 0.1.0 to 0.3.0"); err != nil {
		t.Error(err)
	}
	if err := checkRender("single_dependency", "{{.Subject}}\n\nChangelog: updated", "Update pullrequest from 0.1.0 to 0.3.0\n\nChangelog: updated"); err != nil {
		t.Error(err)
	}
}

func TestCommitMessageSubjectBodyTemplates(t *testing.T) {
	if err := checkRender("single_dependency", "{{.SubjectAndBody}}", "Update pullrequest from 0.1.0 to 0.3.0"); err != nil {
		t.Error(err)
	}
	if err := checkRender("two_dependencies", "{{.SubjectAndBody}}", "Update requirements.txt (pullrequest and requests)\n\n- `pullrequest` in `requirements.txt` from 0.1.0 to 0.3.0\n- `requests` in `requirements.txt` from 0.1.0 to 0.3.0"); err != nil {
		t.Error(err)
	}
	if err := checkRender("single_lockfile", "{{.SubjectAndBody}}", "Update yarn.lock (postcss-cli and tailwindcss)\n\n- `yarn.lock` was updated (including 2 direct and 44 transitive dependencies)\n  - `postcss-cli` was updated from 6.1.2 to 6.1.3\n  - `tailwindcss` was updated from 1.0.1 to 1.1.2"); err != nil {
		t.Error(err)
	}
}
