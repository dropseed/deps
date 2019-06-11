package component

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime/debug"
	"strings"

	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/pullrequest/adapter"

	"github.com/dropseed/deps/internal/schema"

	"github.com/dropseed/deps/internal/output"
)

type Runner struct {
	Index  int
	Given  string
	Config *Config
	Path   string
	Env    []string
}

const DEBUG = true
const DefaultRemotePrefix = "dropseed/deps-"
const DefaultCacheDirName = "deps"

func NewRunnerFromString(s string) (*Runner, error) {
	runner, err := newRunnerFromPath(s)

	if os.IsNotExist(err) {
		runner, err = newRunnerFromRemote(s)
	}

	return runner, err
}

func newRunnerFromPath(s string) (*Runner, error) {
	componentPath := s

	configPath := path.Join(componentPath, DefaultFilename)
	config, err := NewConfigFromPath(configPath)
	if err != nil {
		return nil, err
	}

	return &Runner{
		Given:  s,
		Config: config,
		Path:   componentPath,
	}, nil
}

func newRunnerFromRemote(s string) (*Runner, error) {
	url := s

	if !strings.Contains(url, "/") {
		// shorthand for dropseed/deps-{}
		url = DefaultRemotePrefix + url
	}

	if !strings.HasPrefix(url, "http") {
		// automatically prefix owner/repo with github
		url = "https://github.com/" + url
	}

	output.Debug("Using component from %s", url)

	// get cache dir for the current dir
	userCache, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	depsCache := path.Join(userCache, DefaultCacheDirName)
	output.Debug("Making deps cache at %s", depsCache)
	if err := os.MkdirAll(depsCache, os.ModePerm); os.IsExist(err) {
		output.Debug("Deps cache already exists")
	} else if err != nil {
		output.Debug("Error making deps cache")
		return nil, err
	}

	// does it not have permission to do 777 on travis?
	// push another beta and test, but probably so -- what should the
	// perms be?

	cloneDirName := path.Base(url)
	cloneDirName = strings.Replace(cloneDirName, ".git", "", -1)
	clonePath := path.Join(depsCache, "components", cloneDirName)

	output.Debug("Storing component in %s", clonePath)

	// or clone into components specifically for this working repo?
	// basename-hash of path in user home dir?

	if _, err := os.Stat(clonePath); os.IsNotExist(err) {
		if err := git.Clone(url, clonePath); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	// run git pull - need to be able to specify version somehow though
	// "version" optional on deps config? anything that can be git checkout in this case
	// so maybe sharing these across repos isn't bad... checkout happens every time
	if err := git.Pull(); err != nil {
		return nil, err
	}

	return newRunnerFromPath(clonePath)
}

func (r *Runner) GetName() string {
	return path.Base(r.Path)
}

func (r *Runner) getOverrideFromEnv(name string) string {
	override := os.Getenv(fmt.Sprintf("DEPS_%d_%s", r.Index, strings.ToUpper(name)))
	if override != "" {
		output.Event("Overriding %s command from env", name)
	}
	return override
}

func (r *Runner) Install() error {
	// TODO install does not need to happen every time
	// - yes it does if local (assume could have changed)
	// - yes it does if new (just cloned)
	// - yes it does if ref chagned
	// so maybe not horrible if it runs every time right now

	output.Event("Installing %s", r.Given)

	command := r.Config.Install
	if override := r.getOverrideFromEnv("install"); override != "" {
		command = override
	}
	output.Debug(command)

	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = r.Path
	cmd.Env = r.Env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func (r *Runner) Collect(inputPath string) (*schema.Dependencies, error) {
	output.Event("Collecting with %s", r.Given)

	command := r.Config.Collect
	if override := r.getOverrideFromEnv("collect"); override != "" {
		command = override
	}

	outputPath, err := r.run(command, inputPath)
	if err != nil {
		return nil, err
	}
	if !DEBUG {
		defer os.Remove(outputPath)
	}

	output.Debug("Finished")

	dependencies, err := schema.NewDependenciesFromJSONPath(outputPath)
	if err != nil {
		output.Error("Unable to load output JSON from collector")
		return nil, err
	}

	return dependencies, nil
}

func (r *Runner) Act(inputDependencies *schema.Dependencies, branch bool) error {
	// In CI
	// - clear out git (using stash)
	// - branch, commit, push, PR
	// Local
	// - just update files

	output.Event("Updating with %s", r.Given)

	gitSha := git.CurrentSHA()

	branchName, err := inputDependencies.GetBranchName()
	if err != nil {
		return err
	}

	stashed := false

	if branch {
		id, err := inputDependencies.GetID()
		if err != nil {
			return err
		}
		stashed, err = git.Stash(fmt.Sprintf("Deps save before update %s", id))
		if err != nil {
			return err
		}
		git.Branch(branchName, gitSha)
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

			// TODO delete branch?

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

	// branch
	// before_update / after_branch?
	// how would this work more naturally now in ci? try without it and find out

	if branch {

		// TODO run commit also, just commit all, use inputDependencies to get title, etc.?
		title, err := inputDependencies.GenerateTitle()
		if err != nil {
			return err
		}
		if err = git.AddCommit(title); err != nil {
			return err
		}

		pr, err := adapter.PullrequestAdapterFromDependenciesJSONPathAndHost(outputPath, git.GitHost())
		if err != nil {
			return err
		}
		if pr != nil {
			// TODO hooks or what do you do otherwise?

			if err = git.PushBranch(branchName); err != nil {
				if strings.Contains(err.Error(), "Authentication failed") {
					if err = pr.PreparePush(); err != nil {
						return err
					}

					if err = git.PushBranch(branchName); err != nil {
						return err
					}
				} else {
					return err
				}
			}
			if err = pr.Create(); err != nil {
				return err
			}
			if err = pr.DoRelated(); err != nil {
				return err
			}
			// TODO remove this?
			// if err = pr.OutputActions(); err != nil {
			// 	return err
			// }
		}

		if err = git.CheckoutLast(); err != nil {
			return err
		}

		if stashed {
			if err := git.StashPop(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *Runner) run(command string, inputPath string) (string, error) {
	tmpfile, err := ioutil.TempFile("", "deps-")
	if err != nil {
		return "", err
	}
	outputPath := tmpfile.Name()

	commandString := fmt.Sprintf("%s %s %s", command, inputPath, outputPath)

	output.Debug(commandString)

	cmd := exec.Command(
		"sh",
		"-c",
		commandString,
	)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Env = r.Env
	cmd.Env = append(cmd.Env, fmt.Sprintf("DEPS_COMPONENT_PATH=%s", r.Path))
	if err != nil {
		return "", err
	}

	err = cmd.Run()
	if err != nil {
		return "", err
	}

	output.Debug(outputPath)

	return outputPath, nil
}
