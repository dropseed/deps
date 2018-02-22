package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/dependencies-io/pullrequest/internal/env"
)

// BranchForJob branches off of GIT_SHA
func BranchForJob() (string, error) {
	gitSha := os.Getenv("GIT_SHA")
	if gitSha == "" {
		return "", errors.New("GIT_SHA not found in env, required to branch")
	}

	branchName, err := GetJobBranchName()
	if err != nil {
		return "", err
	}

	checkoutShaOutput, err := exec.Command("git", "checkout", gitSha).CombinedOutput()
	fmt.Println(string(checkoutShaOutput))
	if err != nil {
		return "", err
	}

	checkoutBranchOutput, err := exec.Command("git", "checkout", "-b", branchName).CombinedOutput()
	fmt.Println(string(checkoutBranchOutput))
	if err != nil {
		return "", err
	}

	return branchName, nil
}

// AddCommit a set of paths
func AddCommit(message string, paths []string) error {
	args := append([]string{"add"}, paths...)
	addOutput, err := exec.Command("git", args...).CombinedOutput()
	fmt.Println(string(addOutput))
	if err != nil {
		return err
	}

	commitMessagePrefix := env.GetSetting("commit_message_prefix", "")
	commitMessage := commitMessagePrefix + message

	commitOutput, err := exec.Command("git", "commit", "-m", commitMessage).CombinedOutput()
	fmt.Println(string(commitOutput))
	if err != nil {
		return err
	}

	return nil
}

// Push a given branch to the origin
func Push(branchName string) error {
	if !env.IsProduction() {
		fmt.Printf("Not pushing git branch in '%s' env\n", env.GetCurrentEnv())
		return nil
	}

	pushOutput, err := exec.Command("git", "push", "--set-upstream", "origin", branchName).CombinedOutput()
	fmt.Println(string(pushOutput))
	if err != nil {
		return err
	}

	return nil
}

func GetJobBranchName() (string, error) {
	jobID := os.Getenv("JOB_ID")
	if jobID == "" {
		return "", errors.New("JOB_ID not found in env, required to generate branch name")
	}

	// expected to be a UUID4, and we'll simply get the first part of the string
	// which isn't guaranteed to be unique, but within the open branches on 1
	// repo, we'll take the chance so that it's easier to use/type for a user
	shortenedJobIDParts := strings.SplitN(jobID, "-", 2)
	shortenedJobID := shortenedJobIDParts[0]

	branchPrefix := env.GetSetting("branch_prefix", "")

	return fmt.Sprintf("%sdeps/update-%s", branchPrefix, shortenedJobID), nil
}
