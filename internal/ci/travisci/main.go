package travisci

import (
	"os"
)

type TravisCI struct {
}

func Is() bool {
	return os.Getenv("TRAVIS") != ""
}

func (travis *TravisCI) Autoconfigure() error {
	return nil
}

func (travis *TravisCI) Branch() string {
	if b := os.Getenv("TRAVIS_PULL_REQUEST_BRANCH"); b != "" {
		return b
	}
	if b := os.Getenv("TRAVIS_BRANCH"); b != "" {
		return b
	}
	return ""
}
