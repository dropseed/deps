package component

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/dropseed/deps/internal/git"

	"github.com/dropseed/deps/internal/output"
)

type Runner struct {
	Index  int
	Given  string
	Config *Config
	Path   string
	Env    []string
}

const DefaultRemotePrefix = "dropseed/deps-"
const DefaultCacheDirName = "deps"

func NewRunnerFromString(s string) (*Runner, error) {
	runner, err := NewRunnerFromPath(s)

	if os.IsNotExist(err) {
		runner, err = newRunnerFromRemote(s)
	}

	return runner, err
}

func NewRunnerFromPath(s string) (*Runner, error) {
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
		git.Clone(url, clonePath)
	} else if err != nil {
		return nil, err
	} else {
		output.Event("FIXME: removing previous clone of component and recloning")
		if err := os.RemoveAll(clonePath); err != nil {
			panic(err)
		}
		git.Clone(url, clonePath)
	}

	// run git pull - need to be able to specify version somehow though
	// "version" optional on deps config? anything that can be git checkout in this case
	// so maybe sharing these across repos isn't bad... checkout happens every time

	// TODO need to point it to the repo, not run in current dir
	// if err := git.Pull(); err != nil {
	// 	return nil, err
	// }

	return NewRunnerFromPath(clonePath)
}

func (r *Runner) GetName() string {
	return path.Base(r.Path)
}

func (r *Runner) getCommand(defaultCmd, cmdName string) string {
	command := defaultCmd
	if override := r.getOverrideFromEnv(cmdName); override != "" {
		command = override
	}
	return command
}

func (r *Runner) getOverrideFromEnv(name string) string {
	override := os.Getenv(fmt.Sprintf("DEPS_%d_%s", r.Index, strings.ToUpper(name)))
	if override != "" {
		output.Event("Overriding %s command from env", name)
	}
	return override
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

	if output.IsDebug() {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

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
