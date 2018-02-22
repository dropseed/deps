package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/dependencies-io/pullrequest/internal/env"
)

// Manifest contains manifest data
type Manifest struct {
	// TODO path string
	Current *ManifestVersion `json:"current"`
	Updated *ManifestVersion `json:"updated,omitempty"`
}

// ManifestVersion constains data for a manifest at a specific point in time
type ManifestVersion struct {
	Dependencies map[string]ManifestDependency `json:"dependencies"`
}

type ManifestDependency struct {
	Source     string    `json:"source"`
	Constraint string    `json:"constraint"`
	Available  []Version `json:"available,omitempty"`
}

// GetSummaryLineForDependencyName returns a bulleted list item string
func (manifest *Manifest) GetSummaryLineForDependencyName(name, manifestPath string) (string, error) {
	currentDependency := manifest.Current.Dependencies[name]
	updatedDependency := manifest.Updated.Dependencies[name]
	return fmt.Sprintf("- `%v` in `%v` from \"%v\" to \"%v\"", name, manifestPath, currentDependency.Constraint, updatedDependency.Constraint), nil
}

// GetBodyContentForDependencyName compiles the markdown content for this dependency update
func (manifest *Manifest) GetBodyContentForDependencyName(name, manifestPath string) (string, error) {
	// TODO add notes

	currentDependency := manifest.Current.Dependencies[name]
	updatedDependency := manifest.Updated.Dependencies[name]

	subject := fmt.Sprintf(
		"[Dependencies.io](https://www.dependencies.io) has updated `%v` (a %v dependency in `%v`) from \"%v\" to \"%v\".",
		name,
		updatedDependency.Source,
		manifestPath,
		currentDependency.Constraint,
		updatedDependency.Constraint,
	)

	content := ""

	for _, v := range updatedDependency.Available {
		vContent := updatedDependency.GetContentForVersion(name, &v)
		content += fmt.Sprintf("\n\n<details>\n<summary>%v</summary>\n\n%v\n\n</details>", v.Name, vContent)
	}

	return subject + content, nil
}

// GetContentForVersion finds the content for a given version, optionally from the remote API
func (dependency *ManifestDependency) GetContentForVersion(dependencyName string, version *Version) string {
	if version.Content != "" {
		return version.Content
	}

	attemptRemote := env.IsProduction() && env.GetSetting("PULLREQUEST_VERSIONS_API_DISABLED", "") == ""
	if attemptRemote && dependency.Source != "" && dependencyName != "" && version.Name != "" {
		apiURL := fmt.Sprintf("https://versions.dependencies.io/%s/%s/%s", dependency.Source, dependencyName, version.Name)
		tr := &http.Transport{
			IdleConnTimeout: 45 * time.Second,
		}
		client := &http.Client{Transport: tr}
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			panic(err)
		}
		req.Header.Add("User-Agent", "deps")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("GitHub-API-Token", os.Getenv("GITHUB_API_TOKEN"))
		if resp, err := client.Do(req); err == nil && resp.StatusCode >= 200 && resp.StatusCode < 400 {
			body, err := ioutil.ReadAll(resp.Body)
			resp.Body.Close()
			if err == nil {
				var data map[string]interface{}
				if err := json.Unmarshal(body, &data); err == nil {
					if content, ok := data["content"]; ok && content != nil {
						if content, ok := content.(string); ok && content != "" {
							return content
						}
					}
				}
			}
		}
	}

	return "_No content found. Please open an issue at https://github.com/dependencies-io/support if you think this content could have been found._"
}
