package schema

// Dependency contains fields and functions common to lockfiles and manifests
type Dependency struct {
	Source string `json:"source"`
	Repo   string `json:"repo,omitempty"`
}
