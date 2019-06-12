package git

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dropseed/deps/internal/env"
	"github.com/dropseed/deps/internal/output"
)

// BranchForJob branches off of GIT_SHA
func Branch(to, from string) {
	if err := run("checkout", "-b", to, from); err != nil {
		panic(err)
	}
}

// Push a given branch to the origin
func PushBranch(branchName string) error {
	return run("push", "--set-upstream", "origin", branchName)
}

func GetBranchName(id string) string {
	branchPrefix := env.GetSetting("branch_prefix", "")
	branchSeparator := env.GetSetting("branch_separator", "/")

	return fmt.Sprintf("%sdeps%s%s", branchPrefix, branchSeparator, id)
}

func GitHost() string {
	// or can maybe tell from github actions env var too or gitlab pipeline, but both should have remote as well
	if override := os.Getenv("DEPS_GIT_HOST"); override != "" {
		return override
	}

	remote := GitRemote()

	// TODO regex, ssh urls, etc.

	if strings.HasPrefix(remote, "https://github.com/") {
		return "github"
	}

	if strings.HasPrefix(remote, "https://gitlab.com/") {
		return "gitlab"
	}

	return ""
}

func GitRemote() string {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	remote, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	s := string(remote)
	s = strings.TrimSpace(s)
	return s
}

func Clone(url, path string) error {
	return run("clone", url, path)
}

func BranchExists(branch string) bool {
	// TODO need to check exit code or stderr? what about other failures
	err := run("rev-parse", "--verify", branch)
	if err != nil {
		return false
	}
	return true
}

func CurrentSHA() string {
	cmd := exec.Command("git", "rev-parse", "--verify", "HEAD")
	out, err := cmd.CombinedOutput()

	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(out))
}

func AddCommit(message string) error {
	if err := run("add", "."); err != nil {
		return err
	}
	if err := run("commit", "-m", message); err != nil {
		return err
	}
	return nil
}

func CheckoutLast() error {
	return run("checkout", "-")
}

func Stash(message string) (bool, error) {
	cmd := exec.Command("git", "stash", "push", "--include-untracked", "-m", message)
	out, err := cmd.CombinedOutput()
	println(out)
	if err != nil {
		return false, err
	}
	if strings.Contains(string(out), "No local changes to save") {
		return false, nil
	}
	return true, nil
}

func StashPop() error {
	return run("stash", "pop")
}

func Pull() error {
	return run("pull")
}

func run(args ...string) error {
	cmd := exec.Command("git", args...)

	if output.Verbosity > 0 {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
