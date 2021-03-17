package component

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/dropseed/deps/internal/cache"
	"github.com/dropseed/deps/internal/config"
	"github.com/dropseed/deps/internal/git"

	"github.com/dropseed/deps/internal/output"
)

type Runner struct {
	Index         int
	Given         string
	Config        *Config
	Path          string
	Env           []string
	shouldInstall bool
}

const DefaultRemotePrefix = "dropseed/deps-"

func NewRunnerFromString(s string) (*Runner, error) {
	runner, err := NewRunnerFromPath(s)

	if os.IsNotExist(err) {
		runner, err = newRunnerFromRemote(s)
	}

	return runner, err
}

func NewRunnerFromPath(s string) (*Runner, error) {
	componentPath := s

	configPath := config.FindFilename(componentPath, DefaultFilenames...)
	if configPath == "" {
		return nil, os.ErrNotExist
	}

	cfg, err := NewConfigFromPath(configPath)
	if err != nil {
		return nil, err
	}

	return &Runner{
		Given:         s,
		Config:        cfg,
		Path:          componentPath,
		shouldInstall: true,
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

	depsCache := cache.GetCachePath()

	cloneDirName := path.Base(url)
	cloneDirName = strings.Replace(cloneDirName, ".git", "", -1)
	clonePath := path.Join(depsCache, "components", cloneDirName)

	output.Debug("Storing component in %s", clonePath)

	cloned := false
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(clonePath); os.IsNotExist(err) {
		git.Clone(url, clonePath)
		cloned = true
	} else if err != nil {
		return nil, err
	}

	// run git commands from the new repo
	if err := os.Chdir(clonePath); err != nil {
		panic(err)
	}

	refBefore := ""

	if !cloned {
		refBefore = git.CurrentRef()
		if err := git.Pull(); err != nil {
			return nil, err
		}
	}

	// TODO checkout user specified Version
	// split @ from string?
	// git.Checkout

	refAfter := git.CurrentRef()

	if err := os.Chdir(cwd); err != nil {
		panic(err)
	}

	runner, err := NewRunnerFromPath(clonePath)
	if err != nil {
		return nil, err
	}

	runner.shouldInstall = refBefore != refAfter

	return runner, nil
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
