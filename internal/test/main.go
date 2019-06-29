package test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"

	"github.com/dropseed/deps/internal/schema"

	"github.com/dropseed/deps/internal/component"
	"github.com/dropseed/deps/internal/output"
)

var UpdateOutputData = false
var LooseOutputDataComparison = false
var ExitEarly = false
var FilterName = ""

func Run() error {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	configs, err := findTestConfigs(pwd)
	if err != nil {
		return err
	}

	if len(configs) < 1 {
		return errors.New("No test configs found")
	}

	output.Debug("%d configs loaded", len(configs))
	for _, config := range configs {
		output.Debug("- %s", config.path)
	}

	tests := testsMatchingFilters(configs)
	if len(tests) == 0 {
		return errors.New("No matching tests found")
	}

	// TODO save actual test errors and show them at the end

	testsTotal := len(tests)
	testsFailed := 0
	testsPassed := 0
	output.Event("Tests found: %d", testsTotal)

	runner, err := component.NewRunnerFromPath(pwd)
	if err != nil {
		return err
	}

	if err := runner.Install(); err != nil {
		return err
	}

	for _, test := range tests {
		output.Event("Starting: %s", test.displayName())
		if err := runTest(runner, test); err != nil {
			testsFailed++
			if ExitEarly {
				return err
			}
			output.Error(err.Error())
			output.Error("Failed: %s\n", test.displayName())
		} else {
			testsPassed++
			output.Success("Passed: %s\n", test.displayName())
		}
	}

	resultString := fmt.Sprintf("%d/%d tests passed", testsPassed, testsTotal)
	if testsFailed > 0 {
		return errors.New(resultString)
	}
	output.Success(resultString)
	return nil
}

func runTest(runner *component.Runner, test *Test) error {
	copyRepoPath, err := temporaryCopyOfDir(test.config.joinPath(test.Repo))
	if err != nil {
		return err
	}

	if err := os.Chdir(copyRepoPath); err != nil {
		return err
	}

	if env, err := test.UserConfig.Environ(); err != nil {
		return err
	} else {
		runner.Env = env
	}

	if output.IsDebug() {
		if wd, err := os.Getwd(); err != nil {
			return err
		} else {
			output.Debug("Running from %s", wd)
		}
	}

	if test.Collect.Output != "" {
		outputDeps, err := runner.Collect(test.UserConfig.Path)
		if err != nil {
			return err
		}
		if err := checkOutput(outputDeps, test.config.joinPath(test.Collect.Output)); err != nil {
			return err
		}
	}

	if test.Act.Input != "" && test.Act.Output != "" {
		inputDeps, err := schema.NewDependenciesFromJSONPath(test.config.joinPath(test.Act.Input))
		if err != nil {
			return err
		}
		outputDeps, err := runner.Act(inputDeps)
		if err != nil {
			return err
		}
		if err := checkOutput(outputDeps, test.config.joinPath(test.Act.Output)); err != nil {
			return err
		}
	}

	if test.Diff != "" {

		diffRepo := test.config.joinPath(test.Diff)

		if UpdateOutputData {
			output.Event("Updating test diff repo with actual results")
			if err := os.RemoveAll(diffRepo); err != nil {
				panic(err)
			}
			if err := copyDir(copyRepoPath, diffRepo); err != nil {
				panic(err)
			}
		}

		diff := getDiff(diffRepo, copyRepoPath, test.DiffArgs...)
		if diff != "" {
			println(diff)
			return errors.New("Diff does not match expected")
		}
	}

	return nil
}

func checkOutput(outputDeps *schema.Dependencies, outputPath string) error {
	if UpdateOutputData {
		output.Event("Writing parsed output to %s", outputPath)
		out, err := json.MarshalIndent(outputDeps, "", "  ")
		if err != nil {
			return err
		}
		out = append(out, "\n"...)
		if err := ioutil.WriteFile(outputPath, out, 0644); err != nil {
			panic(err)
		}
	}

	expectedOutputDeps, err := schema.NewDependenciesFromJSONPath(outputPath)
	if err != nil {
		return err
	}

	if err := compare(outputDeps, expectedOutputDeps); err != nil {
		return err
	}

	return nil
}

func getDiff(a string, b string, args ...string) string {
	cmdArgs := []string{
		"-Naur",
		a,
		b,
	}
	cmdArgs = append(cmdArgs, args...)
	cmd := exec.Command("diff", cmdArgs...)
	out, _ := cmd.CombinedOutput()
	// if err != nil {
	// 	panic(err)
	// }
	return strings.TrimSpace(string(out))
}

func compare(given, expected *schema.Dependencies) error {
	if LooseOutputDataComparison {
		if given.UpdateID != expected.UpdateID {
			return fmt.Errorf("Update IDs don't match: %s != %s", given.UpdateID, expected.UpdateID)
		}
	} else {
		if match, err := schemasMatchExactly(given, expected); err != nil {
			return err
		} else if !match {
			return errors.New("Output doesn't match")
		}
	}

	return nil
}

func testsMatchingFilters(configs []*Config) []*Test {

	output.Debug("Name filter: \"%s\"", FilterName)

	tests := []*Test{}
	for _, cfg := range configs {
		for _, c := range cfg.Tests {
			match := true

			match = match && (FilterName == "" || strings.Contains(c.Name, FilterName))

			if match {
				tests = append(tests, c)
			} else {
				output.Debug("%s does not match filters", c.displayName())
			}
		}
	}
	return tests
}

func temporaryCopyOfDir(dirToCopy string) (string, error) {
	tmpPath := "" // use the default

	if runtime.GOOS == "darwin" {
		// the default is not shared with docker, by default
		tmpPath = "/tmp"
	}

	dir, err := ioutil.TempDir(tmpPath, "deps-")
	if err != nil {
		return "", err
	}

	repoDir := path.Join(dir, "repo")
	if err := copyDir(dirToCopy, repoDir); err != nil {
		panic(err)
	}
	output.Debug("Made temporary copy of %s into %s\n", dirToCopy, repoDir)

	return repoDir, nil
}

func copyDir(from, to string) error {
	cmd := exec.Command("cp", "-a", from, to)
	if output.IsDebug() {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
