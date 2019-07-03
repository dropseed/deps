package changelogs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type ChangelogFetcher interface {
	fetch() string
}

type Dependency struct {
	source  string
	name    string
	version string
	repo    string
	url     string
}

func GetChangelog(depSource, depName, depVersion string) string {

	// a cache obj?

	dep := &Dependency{
		source:  depSource,
		name:    depName,
		version: depVersion,
	}

	// first check remote db for overrides
	// api/source/name
	// - repo (url) maybe
	// - url (just start with this)
	// - versions.1.0.0.url - or version url pattern?? or url_pattern

	var fetcher ChangelogFetcher

	if dep.source == "pypi" {
		// pypi fetcher
		fetcher = &Pypi{
			dependency: dep,
		}
	}

	dep.url = fetcher.fetch()

	attemptRemote := true //env.GetSetting("PULLREQUEST_VERSIONS_API_DISABLED", "") == ""

	if attemptRemote && depSource != "" && depName != "" && depVersion != "" {
		apiURL := fmt.Sprintf("https://versions.dependencies.io/%s/%s/%s", depSource, depName, depVersion)
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
		req.Header.Add("GitHub-API-Token", os.Getenv("DEPS_GITHUB_TOKEN")) // TODO abstract this?
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
