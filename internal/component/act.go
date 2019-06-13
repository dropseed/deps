package component

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"
	"time"

	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/pullrequest/adapter"
	"github.com/dropseed/deps/internal/schema"
)

func (r *Runner) Act(inputDependencies *schema.Dependencies, baseBranch string) error {
	output.Event("Updating with %s", r.Given)

	predictedUpdateBranch := ""
	stashed := false

	if baseBranch != "" {
		// If we're given a base branch then we'll be creating a new
		// branch for the update
		id := inputDependencies.GetID()
		var err error
		err = nil
		output.Event("Temporarily saving your uncommitted changes in a git stash")
		stashed, err = git.Stash(fmt.Sprintf("Deps save before update %s", id))
		if err != nil {
			return err
		}
		predictedUpdateBranch = inputDependencies.GetBranchName()
		git.Branch(predictedUpdateBranch)
	} else {
		output.Event("Running changes directly (no branches)")
	}

	// Try to revert this stuff if something goes wrong
	defer func() {
		if r := recover(); r != nil {
			output.Error("Recovering from update error")

			if err := git.CheckoutLast(); err != nil {
				panic(err)
			}

			if stashed {
				if err := git.StashPop(); err != nil {
					panic(err)
				}
			}

			// TODO delete baseBranch?

			debug.PrintStack()

			panic(r)
		}
	}()

	inputFilename, err := inputTempFile(inputDependencies)
	if err != nil {
		return err
	}
	if !output.IsDebug() {
		defer os.Remove(inputFilename)
	}

	command := r.Config.Act
	if override := r.getOverrideFromEnv("act"); override != "" {
		command = override
	}

	outputPath, err := r.run(command, inputFilename)
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

	if baseBranch != "" {

		updateBranch := outputDependencies.GetBranchName()
		if updateBranch != predictedUpdateBranch {
			output.Debug("Actual update differed from expected, renaming git branch")
			git.RenameBranch(predictedUpdateBranch, updateBranch)
		}

		pr, err := adapter.PullrequestAdapterFromDependenciesJSONPathAndHost(outputPath, git.GitHost(), baseBranch)
		if err != nil {
			return err
		}

		// TODO run commit also, just commit all, use inputDependencies to get title, etc.?
		title, err := inputDependencies.GenerateTitle()
		if err != nil {
			return err
		}
		if err = git.AddCommit(title); err != nil {
			return err
		}

		if pr != nil {
			// TODO hooks or what do you do otherwise?

			if err := git.PushBranch(updateBranch); err != nil {
				// TODO better to check for "Authentication failed" in output?
				if err := pr.PreparePush(); err != nil {
					return err
				}

				if err := git.PushBranch(updateBranch); err != nil {
					return err
				}
			}

			output.Debug("Waiting a second for the push to be processed by the host")
			time.Sleep(2 * time.Second)

			if err := pr.Create(); err != nil {
				return err
			}
			if err := pr.DoRelated(); err != nil {
				return err
			}
		}

		if err := git.CheckoutLast(); err != nil {
			return err
		}

		if stashed {
			output.Event("Putting original uncommitted changes back")
			if err := git.StashPop(); err != nil {
				return err
			}
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
