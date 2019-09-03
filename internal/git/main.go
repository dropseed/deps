package git

import (
	"bytes"
	"fmt"
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

func AddCommit(message string) {
	if err := run("add", "."); err != nil {
		panic(err)
	}
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
	println(outS)
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

func Fetch() {
	if err := run("fetch"); err != nil {
		panic(err)
	}
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
			println(cmdString)
			println(out.String())
		}

		return err
	}

	return nil
}
