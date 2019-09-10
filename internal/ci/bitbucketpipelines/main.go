package bitbucketpipelines

import (
	"fmt"
	"os"
)

type BitbucketPipelines struct {
}

func Is() bool {
	return os.Getenv("BITBUCKET_BUILD_NUMBER") != ""
}

func (gitlab *BitbucketPipelines) Autoconfigure() error {
	return nil
}

func (gitlab *BitbucketPipelines) Branch() string {
	if s := os.Getenv("BITBUCKET_BRANCH"); s != "" {
		return s
	}
	return ""
}

func GetProjectAPIURL() string {
	if slug := os.Getenv("BITBUCKET_REPO_FULL_NAME"); slug != "" {
		return fmt.Sprintf("https://api.bitbucket.org/2.0/repositories/%s", slug)
	}
	return ""
}
