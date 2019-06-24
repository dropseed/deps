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
func Branch(to string) {
	if err := run("checkout", "-b", to); err != nil {
		panic(err)
	}
}

// Push a given branch to the origin
func PushBranch(branchName string) {
	if branchName == "" {
		if err := run("push"); err != nil {
			panic(err)
		}
	}
	if err := run("push", "--set-upstream", "origin", branchName); err != nil {
		panic(err)
	}
}

func CanPush() bool {
	if err := run("push", "--dry-run"); err != nil {
		return false
	}
	return true
}

func GetBranchName(id string) string {
	prefix := getBranchPrefix()
	return fmt.Sprintf("%s%s", prefix, id)
}

func IsDepsBranch(branchName string) bool {
	prefix := getBranchPrefix()
	return strings.HasPrefix(branchName, prefix)
}

func GetDepsBranches() []string {
	cmd := exec.Command("git", "branch", "--list", "--no-column")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	s := string(out)

	lines := strings.Split(s, "\n")

	branches := []string{}
	for _, line := range lines {
		branch := strings.TrimSpace(line)
		if IsDepsBranch(branch) {
			branches = append(branches, branch)
		}
	}
	return branches
}

func getBranchPrefix() string {
	branchPrefix := env.GetSetting("branch_prefix", "")
	branchSeparator := env.GetSetting("branch_separator", "/")

	return fmt.Sprintf("%sdeps%s", branchPrefix, branchSeparator)
}

func GitHost() string {
	// or can maybe tell from github actions env var too or gitlab pipeline, but both should have remote as well
	if override := os.Getenv("DEPS_GIT_HOST"); override != "" {
		return override
	}

	remote := GitRemote()

	// TODO https://user:pass@

	if strings.HasPrefix(remote, "https://github.com/") || strings.HasPrefix(remote, "git@github.com:") {
		return "github"
	}

	if strings.HasPrefix(remote, "https://gitlab.com/") || strings.HasPrefix(remote, "git@gitlab.com:") {
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
	// See if we have a matching branch locally
	if err := run("rev-parse", "--verify", branch); err == nil {
		return true
	}

	// Also need to check remote, in case the branch is cloned locally
	if err := run("ls-remote", "--exit-code", "--heads", "origin", branch); err == nil {
		return true
	}

	return false
}

func CurrentRef() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.CombinedOutput()

	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(out))
}

func AddCommit(message string) {
	if err := run("add", "."); err != nil {
		panic(err)
	}
	if err := run("-c", "user.name=deps", "-c", "user.email=bot@dependencies.io", "commit", "-m", message); err != nil {
		panic(err)
	}
}

func Checkout(s string) {
	if err := run("checkout", s); err != nil {
		panic(err)
	}
}

func CheckoutLast() {
	Checkout("-")
}

func ResetAndClean() {
	if err := run("reset", "--hard"); err != nil {
		panic(err)
	}
	if err := run("clean", "-df"); err != nil {
		panic(err)
	}
}

func Stash(message string) bool {
	cmd := exec.Command("git", "stash", "push", "--include-untracked", "-m", message)
	out, err := cmd.CombinedOutput()
	println(out)
	if err != nil {
		panic(err)
	}
	if strings.Contains(string(out), "No local changes to save") {
		return false
	}
	return true
}

func StashPop() error {
	return run("stash", "pop")
}

func Pull() error {
	return run("pull")
}

func RenameBranch(old, new string) {
	if err := run("branch", "-m", old, new); err != nil {
		panic(err)
	}
}

func FetchAllBranches() {
	if err := run("fetch", "--all"); err != nil {
		panic(err)
	}
}

func run(args ...string) error {
	output.Debug("git %s", strings.Join(args, " "))
	cmd := exec.Command("git", args...)

	if output.IsDebug() {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
