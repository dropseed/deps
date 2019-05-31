package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/dropseed/deps/internal/env"
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

	attemptRemote := env.GetSetting("PULLREQUEST_VERSIONS_API_DISABLED", "") == ""
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

	return ""
}
