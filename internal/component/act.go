package component

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"

	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/pullrequest/adapter"
	"github.com/dropseed/deps/internal/schema"
)

func (r *Runner) Act(inputDependencies *schema.Dependencies, baseBranch string) error {
	// In CI
	// - clear out git (using stash)
	// - baseBranch, commit, push, PR
	// Local
	// - just update files

	output.Event("Updating with %s", r.Given)

	gitSha := git.CurrentSHA()

	updateBranch := inputDependencies.GetBranchName()

	stashed := false

	if baseBranch != "" {
		id := inputDependencies.GetID()
		var err error
		err = nil
		output.Event("Temporarily saving your uncommitted changes in a git stash")
		stashed, err = git.Stash(fmt.Sprintf("Deps save before update %s", id))
		if err != nil {
			return err
		}
		git.Branch(updateBranch, gitSha)
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

	// Put the input in a file
	inputJSON, err := json.MarshalIndent(inputDependencies, "", "  ")
	if err != nil {
		return err
	}
	inputFile, err := ioutil.TempFile("", "deps-")
	if !DEBUG {
		defer os.Remove(inputFile.Name())
	}
	if err != nil {
		return err
	}
	if _, err := inputFile.Write(inputJSON); err != nil {
		panic(err)
	}

	// Run it

	command := r.Config.Act
	if override := r.getOverrideFromEnv("act"); override != "" {
		command = override
	}

	outputPath, err := r.run(command, inputFile.Name())
	if err != nil {
		return err
	}
	if !DEBUG {
		defer os.Remove(outputPath)
	}

	// baseBranch
	// before_update / after_branch?
	// how would this work more naturally now in ci? try without it and find out

	if baseBranch != "" {

		// TODO run commit also, just commit all, use inputDependencies to get title, etc.?
		title, err := inputDependencies.GenerateTitle()
		if err != nil {
			return err
		}
		if err = git.AddCommit(title); err != nil {
			return err
		}

		pr, err := adapter.PullrequestAdapterFromDependenciesJSONPathAndHost(outputPath, git.GitHost(), baseBranch)
		if err != nil {
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
