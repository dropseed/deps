package component

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/pullrequest/adapter"
	"github.com/dropseed/deps/internal/schema"
)

func (r *Runner) Act(inputDependencies *schema.Dependencies, baseBranch string, commitPush bool) error {
	output.Event("Updating with %s", r.Given)

	predictedUpdateBranch := ""
	stashed := false

	if commitPush {
		output.Event("Temporarily saving your uncommitted changes in a git stash")
		stashed = git.Stash(fmt.Sprintf("Deps save before update"))

		// Stash pop needs to happen last (so be added first)
		defer func() {
			if stashed {
				output.Event("Putting original uncommitted changes back")
				if err := git.StashPop(); err != nil {
					output.Error("Error putting stash back: %v", err)
				}
			}
		}()
	}

	if baseBranch == "" {
		output.Event("Running changes directly (no branches)")
	} else {
		// If we're given a base branch then we'll be creating a new
		// branch for the update
		predictedUpdateBranch = inputDependencies.GetBranchName()
		git.Branch(predictedUpdateBranch)

		defer func() {
			// Theres should only be uncommitted changes if we're bailing early
			git.ResetAndClean()
			git.CheckoutLast()
		}()
	}

	inputFilename, err := inputTempFile(inputDependencies)
	if err != nil {
		return err
	}
	if !output.IsDebug() {
		defer os.Remove(inputFilename)
	}

	outputPath, err := r.run(r.getCommand(r.Config.Act, "act"), inputFilename)
	if err != nil {
		return err
	}
	if !output.IsDebug() {
		defer os.Remove(outputPath)
	}

	outputDependencies, err := schema.NewDependenciesFromJSONPath(outputPath)
	if err != nil {
		return err
	}

	// baseBranch
	// before_update / after_branch?
	// how would this work more naturally now in ci? try without it and find out

	// Empty by default -- current branch
	updateBranch := ""

	if baseBranch != "" {
		updateBranch = outputDependencies.GetBranchName()

		if updateBranch != predictedUpdateBranch {
			output.Debug("Actual update differed from expected, renaming git branch")

			if git.BranchExists(updateBranch) {
				output.Warning("Aborting update branch rename since the new branch should already exist")
				return nil
			}

			git.RenameBranch(predictedUpdateBranch, updateBranch)
		}
	}

	var pr adapter.PullrequestAdapter
	gitHost := git.GitHost()

	if baseBranch != "" {
		// pr = repo.NewPullrequest(outputDependencies, baseBranch)
		pr, err = adapter.PullrequestAdapterFromDependenciesJSONPathAndHost(outputPath, gitHost, baseBranch)
		if err != nil {
			return err
		}
	}

	if commitPush {
		title, err := inputDependencies.GenerateTitle()
		if err != nil {
			return err
		}

		git.AddCommit(title)
		git.PushBranch(updateBranch)

		output.Debug("Waiting a second for the push to be processed by the host")
		time.Sleep(2 * time.Second)
	}

	if pr != nil {
		// TODO hooks or what do you do otherwise?

		// CreateOrUpdate instead

		if err := pr.Create(); err != nil {
			return err
		}
		if err := pr.DoRelated(); err != nil {
			return err
		}
	}

	return nil
}

func inputTempFile(inputDependencies *schema.Dependencies) (string, error) {
	inputJSON, err := json.MarshalIndent(inputDependencies, "", "  ")
	if err != nil {
		return "", err
	}
	inputFile, err := ioutil.TempFile("", "deps-")
	if err != nil {
		return "", err
	}
	if _, err := inputFile.Write(inputJSON); err != nil {
		panic(err)
	}
	return inputFile.Name(), nil
}
