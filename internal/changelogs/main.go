package changelogs

type ChangelogFinder struct {
	Source     string
	Dependency string
	Repo       string
	Version    string
}

func (cf *ChangelogFinder) FindURL() string {

	if cf.Repo != "" {
		// go straight there? or only if didn't find it elsewhere?
		return "repo"
	}

	if cf.Source == "git" {
		// use "dependency" as repo url
		return "repo"
	}

	if cf.Source == "pypi" {
		return findPypiURL(cf.Dependency, cf.Version)
	}

	return ""
}

func findPypiURL(dependency, version) string {

}
