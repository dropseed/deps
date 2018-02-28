package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dependencies-io/pullrequest/internal/env"
)

// Dependency contains fields and functions common to lockfiles and manifests
type Dependency struct {
	Source string `json:"source"`
}

func (dependency *Dependency) GetMarkdownContentForVersions(dependencyName string, versions []Version) string {
	contentParts := []string{}
	versionsNotFound := []string{}

	for _, v := range versions {
		if vContent := dependency.getMarkdownContentForVersion(dependencyName, &v); vContent != "" {
			contentParts = append(contentParts, vContent)
		} else {
			versionsNotFound = append(versionsNotFound, v.Name)
		}
	}

	if len(versionsNotFound) > 0 {
		versions := ""
		if len(versionsNotFound) == 1 {
			versions = versionsNotFound[0]
		} else if len(versionsNotFound) == 2 {
			versions = versionsNotFound[0] + " or " + versionsNotFound[1]
		} else {
			versionsNotFound[len(versionsNotFound)-1] = "or " + versionsNotFound[len(versionsNotFound)-1]
			versions = strings.Join(versionsNotFound, ", ")
		}

		footnote := fmt.Sprintf("<small>*We didn't find any content for %s. Feel free to open an issue at https://github.com/dependencies-io/support to suggest any improvements.*</small>", versions)

		if len(contentParts) > 0 {
			// we need a manual break here to give spacing under the details sections
			footnote = "<br />" + footnote
		}

		contentParts = append(contentParts, footnote)
	}

	return strings.Join(contentParts, "\n\n")
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

	return ""
}
