package schema

import (
	"fmt"

	"github.com/dropseed/deps/internal/changelogs"
)

// Dependency contains fields and functions common to lockfiles and manifests
type Dependency struct {
	Source string `json:"source"`
	Repo   string `json:"repo,omitempty"`
}

func (dependency *Dependency) GetMarkdownContentForVersion(dependencyName string, version *Version) string {
	if vContent := dependency.getMarkdownContentForVersion(dependencyName, version); vContent != "" {
		return vContent
	} else {
		return fmt.Sprintf("<small>*We didn't find any content for %s. Feel free to open an issue at https://github.com/dropseed/support to suggest any improvements.*</small>", version.Name)
	}
}

func (dependency *Dependency) getMarkdownContentForVersion(dependencyName string, version *Version) string {
	content := dependency.getContentForVersion(dependencyName, version)
	if content == "" {
		return ""
	}
	return fmt.Sprintf("<details>\n<summary>%v</summary>\n\n%v\n\n</details>", version.Name, content)
}

// getContentForVersion finds the content for a given version, optionally from the remote API
func (dependency *Dependency) getContentForVersion(dependencyName string, version *Version) string {
	if version.Content != "" {
		return version.Content
	}

	return changelogs.GetChangelog(dependency.Source, dependencyName, version.Name)
}
