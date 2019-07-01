package pullrequest

// PullrequestAdapter implements the basic Pullrequest functions
type PullrequestAdapter interface {
	CreateOrUpdate() error
}

type RepoAdapter interface {
	CheckRequirements() error
	PreparePush()
	// NewPullrequest(*schema.Dependencies, string) PullrequestAdapter
}
