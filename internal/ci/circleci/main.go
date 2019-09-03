package circleci

import (
	"os"
	"os/exec"
	"strings"

	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/output"
)

type CircleCI struct {
}

func Is() bool {
	return os.Getenv("CIRCLECI") != ""
}

func (circle *CircleCI) Autoconfigure() error {
	// CircleCI uses ssh clones by default,
	// so try to switch to https
	if cmd := exec.Command("git", "config", "--global", "--remove-section", "url.ssh://git@github.com"); cmd != nil {
		output.Event("Autoconfigure: %s", strings.Join(cmd.Args, " "))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run() // Don't worry about an error
	}

	originalOrigin := git.GitRemote()
	if updatedOrigin := git.GitRemoteToHTTPS(originalOrigin); originalOrigin != updatedOrigin {
		if cmd := exec.Command("git", "remote", "set-url", "origin", updatedOrigin); cmd != nil {
			output.Event("Autoconfigure: %s", strings.Join(cmd.Args, " "))
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (circle *CircleCI) Branch() string {
	return os.Getenv("CIRCLE_BRANCH")
}
