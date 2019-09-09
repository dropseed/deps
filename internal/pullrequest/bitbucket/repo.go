package bitbucket

import (
	"errors"
)

type BitbucketRepo struct {
	apiPassword string
	apiUsername string
}

func NewRepo() *BitbucketRepo {
	return &BitbucketRepo{
		apiUsername: getAPIUsername(),
		apiPassword: getAPIPassword(),
	}
}

func (repo *BitbucketRepo) CheckRequirements() error {
	if repo.apiPassword == "" {
		return errors.New("Unable to find Bitbucket API password.\n\nVisit https://docs.dependencies.io/bitbucket for more information.")
	}
	if repo.apiUsername == "" {
		return errors.New("Unable to find Bitbucket API username.\n\nVisit https://docs.dependencies.io/bitbucket for more information.")
	}
	return nil
}

func (repo *BitbucketRepo) Autoconfigure() {
}
