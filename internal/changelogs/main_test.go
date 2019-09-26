package changelogs

import "testing"

func TestPythonRequests(t *testing.T) {
	cf := &ChangelogFinder{
		Source:     "pypi",
		Dependency: "requests",
		Repo:       "",
		Version:    "==2.11.1",
	}
	if url := cf.FindURL(); url != "test" {
		t.Error(url)
	}
}
