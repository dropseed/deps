package git

import (
	"bytes"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/dropseed/deps/internal/output"
)

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

func GetBranchName(suffix string) string {
	prefix := getBranchPrefix()
	return prefix + suffix
}

func IsDepsBranch(branchName string) bool {
	prefix := getBranchPrefix()
	return strings.HasPrefix(branchName, prefix)
}

func GetDepsBranches() []string {
	depsBranches := []string{}
	branches := listBranches()
	for _, branch := range branches {
		if IsDepsBranch(branch) {
			depsBranches = append(depsBranches, branch)
		}
	}
	return depsBranches
}

func listBranches() []string {
	cmd := exec.Command("git", "branch", "--list", "--all", "--no-column")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	s := string(out)

	lines := strings.Split(s, "\n")

	branches := []string{}
	for _, line := range lines {
		branch := strings.TrimSpace(line)
		if strings.HasPrefix(branch, "remotes/origin/") {
			branch = branch[15:]
		}
		// if IsDepsBranch(branch) {
		branches = append(branches, branch)
		// }
	}
	return branches
}

func MergeWouldConflict(branch string) bool {
	mergeCmd := exec.Command("git", "merge", branch, "--no-commit")
	mergeOutput, mergeErr := mergeCmd.CombinedOutput()

	abortMergeCmd := exec.Command("git", "merge", "--abort")
	abortMergeCmd.CombinedOutput()

	if mergeErr != nil {
		output.Warning((string(mergeOutput)))
	}

	return mergeErr != nil
}

func Merge(branch string) bool {
	cmd := exec.Command("git", "merge", branch)
	out, err := cmd.CombinedOutput()
	outS := string(out)

	if err != nil {
		output.Error(outS)
		exec.Command("git", "merge", "--abort").Run()
		return false
	}

	if strings.Contains(outS, "Already up-to-date") {
		return false
	}

	return true
}

// MergeAvailable returns true if there are changes that can be merged
// whether or not there would be a conflict
func MergeAvailable(branch string) bool {
	cmd := exec.Command("git", "merge", branch, "--no-commit")
	out, _ := cmd.CombinedOutput()
	outS := string(out)

	if strings.Contains(outS, "fatal: refusing to merge unrelated histories") {
		output.Debug("Backfilling repo history to allow merging")
		if err := run("pull", "--unshallow"); err != nil {
			panic(err)
		}

		output.Debug("Retrying merge check")
		cmd := exec.Command("git", "merge", branch, "--no-commit")
		out, _ := cmd.CombinedOutput()
		outS = string(out)
	}

	fmt.Println(outS)

	// Clean up the merge no matter what
	exec.Command("git", "merge", "--abort").Run()

	return !strings.Contains(outS, "Already up-to-date")
}

func getBranchPrefix() string {
	branchPrefix := ""
	branchSeparator := "/"
	return fmt.Sprintf("%sdeps%s", branchPrefix, branchSeparator)
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

func GitRemoteHostname() string {
	remote := GitRemote()
	parsed, err := url.Parse(remote)
	if err != nil {
		panic(err)
	}
	return parsed.Hostname()
}

func GitRemoteToHTTPS(original string) string {
	re := regexp.MustCompile("^git@([^:]+):(.+)")
	updated := re.ReplaceAllString(original, "https://$1/$2")
	return updated
}

func Clone(url, path string) {
	if err := run("clone", url, path); err != nil {
		panic(err)
	}
}

func BranchMatching(startsWith string) string {
	branches := listBranches()
	for _, b := range branches {
		if strings.HasPrefix(b, startsWith) {
			return b
		}
	}
	return ""
}

func CurrentRef() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	out, err := cmd.CombinedOutput()

	if err != nil {
		panic(err)
	}

	return strings.TrimSpace(string(out))
}

func Add() {
	if err := run("add", "."); err != nil {
		panic(err)
	}
}

func Unstage() {
	if err := run("reset", "."); err != nil {
		panic(err)
	}
}

func Commit(message string) {
	if err := run("commit", "-m", message); err != nil {
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
	cmd := exec.Command("git", "stash", "save", "--include-untracked", message)
	out, err := cmd.CombinedOutput()
	outS := string(out)
	fmt.Println(outS)
	if err != nil {
		panic(err)
	}
	if strings.Contains(outS, "No local changes to save") {
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

func DeleteRemoteBranch(branch string) error {
	return run("push", "origin", "--delete", branch)
}

func Fetch() {
	if err := run("fetch"); err != nil {
		panic(err)
	}
}

func HasStagedChanges() bool {
	cmd := exec.Command("git", "diff", "--name-only", "--staged")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	cleaned := strings.TrimSpace(string(out))
	return cleaned != ""
}

func IsDirty() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	out, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	cleaned := strings.TrimSpace(string(out))
	return cleaned != ""
}

func Status() string {
	cmd := exec.Command("git", "status")
	status, err := cmd.CombinedOutput()
	if err != nil {
		panic(err)
	}
	s := string(status)
	return s
}

func RepoRoot() string {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}
	return string(out)
}

func run(args ...string) error {
	cmdString := fmt.Sprintf("git %s", strings.Join(args, " "))
	output.Debug(cmdString)
	cmd := exec.Command("git", args...)

	out := bytes.Buffer{}

	if output.IsDebug() {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	} else {
		cmd.Stdout = &out
		cmd.Stderr = &out
	}

	if err := cmd.Run(); err != nil {

		if !output.IsDebug() {
			// Show more output if it wasn't showing already
			fmt.Println(cmdString)
			fmt.Println(out.String())
		}

		return err
	}

	return nil
}
