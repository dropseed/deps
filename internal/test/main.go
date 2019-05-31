package test

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/dependencies-io/deps/internal/output"
)

var UpdateOutputData = false
var LooseOutputDataComparison = false
var ExitEarly = false
var FilterType = ""
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

	cases := casesMatchingFilters(configs)
	if len(cases) == 0 {
		return errors.New("no test cases found")
	}

	testsTotal := len(cases)
	testsFailed := 0
	testsPassed := 0
	output.Event("Test cases found: %d", testsTotal)

	output.Event("Building docker image %s", pwd)
	buildImage(pwd, getImageName(pwd))

	for _, testCase := range cases {
		output.Event("Starting: %s", testCase.displayName())
		if err := runTestCase(testCase, pwd); err != nil {
			testsFailed++
			if ExitEarly {
				return err
			}
			output.Error(err.Error())
			output.Error("Failed: %s\n", testCase.displayName())
		} else {
			testsPassed++
			output.Success("Passed: %s\n", testCase.displayName())
		}
	}

	resultString := fmt.Sprintf("%d passed and %d failed of %d total", testsPassed, testsFailed, testsTotal)
	if testsFailed > 0 {
		return errors.New(resultString)
	}
	output.Success(resultString)
	return nil
}

func runTestCase(testCase *Case, dir string) error {
	return errors.New("unimplemented")
	// 	depConfig := testCase.asConfigDependency(getImageName(dir))

	// 	repoPath := path.Join(dir, testCase.RepoContents)
	// 	var err error
	// 	repoPath, err = utils.TemporaryCopyOfDir(repoPath)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	r, err := runner.NewDockerRunner(repoPath, depConfig, testCase.Type)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if _, err = execute.Run("git", "-C", repoPath, "init"); err != nil {
	// 		panic(err)
	// 	}
	// 	if _, err = execute.Run("git", "-C", repoPath, "add", "."); err != nil {
	// 		panic(err)
	// 	}
	// 	if _, err = execute.Run("git", "-C", repoPath, "-c", "user.name=\"Example User\"", "-c", "user.email=\"testing@example.com\"", "commit", "-m", "First test commit"); err != nil {
	// 		panic(err)
	// 	}
	// 	sha, err := execute.Run("git", "-C", repoPath, "rev-parse", "HEAD")
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	r.SetEnv("DEPENDENCIES_ENV", "test")
	// 	r.SetEnv("GIT_SHA", strings.TrimSpace(string(sha)))
	// 	r.SetEnv("GIT_HOST", "test")
	// 	r.SetEnv("GIT_BRANCH", "master")
	// 	r.SetEnv("JOB_ID", "0")

	// 	inputSchema, err := testCase.inputSchema()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	expectedOutputSchema, err := testCase.outputSchema()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	outputSchema, err := r.Run(inputSchema)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	if UpdateOutputData {
	// 		output.Event("Writing parsed output to %s", testCase.OutputDataPath)
	// 		out, err := json.MarshalIndent(outputSchema, "", "  ")
	// 		if err != nil {
	// 			return err
	// 		}
	// 		out = append(out, "\n"...)
	// 		if err := ioutil.WriteFile(testCase.OutputDataPath, out, 0644); err != nil {
	// 			panic(err)
	// 		}
	// 	}

	// 	if LooseOutputDataComparison {
	// 		match, err := schemasMatchLoosely(outputSchema, expectedOutputSchema)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if !match {
	// 			return errors.New("output does not loosely match expected")
	// 		}
	// 	} else {
	// 		match, err := schemasMatchExactly(outputSchema, expectedOutputSchema)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		if !match {
	// 			return errors.New("output does not exactly match expected")
	// 		}
	// 	}

	// 	if len(testCase.Tests) > 0 {
	// 		if err := runTestCaseExtraTests(testCase, repoPath); err != nil {
	// 			return err
	// 		}
	// 	}

	// 	return nil
}

func runTestCaseExtraTests(testCase *Case, repoPath string) error {
	// wd, err := os.Getwd()
	// if err != nil {
	// 	panic(err)
	// }
	// if err := os.Chdir(repoPath); err != nil {
	// 	panic(err)
	// }
	// execute.Env = []string{"CWD=" + wd}

	// for _, test := range testCase.Tests {
	// 	cmd := strings.TrimSpace(test)
	// 	output.Event("- subtest: %s", cmd)
	// 	_, err := execute.Run("sh", "-c", cmd)
	// 	if err != nil {
	// 		return err
	// 	}
	// }

	// if err := os.Chdir(wd); err != nil {
	// 	panic(err)
	// }
	// execute.Env = []string{}

	return nil
}

func casesMatchingFilters(configs []*Config) []*Case {

	output.Debug("Type filter: \"%s\"", FilterType)
	output.Debug("Name filter: \"%s\"", FilterName)

	cases := []*Case{}
	for _, cfg := range configs {
		for _, c := range cfg.Cases {
			match := true

			match = match && (FilterType == "" || c.Type == FilterType)
			match = match && (FilterName == "" || strings.Contains(c.Name, FilterName))

			if match {
				cases = append(cases, c)
			} else {
				output.Debug("%s does not match filters", c.displayName())
			}
		}
	}
	return cases
}

func getImageName(dir string) string {
	name := "deps-test-" + path.Base(dir)
	name = strings.Replace(name, "/", "-", -1)
	return name
}

func buildImage(path, name string) error {
	// _, err := execute.Run("docker", "build", "-t", name, path)
	// if err != nil {
	// 	return err
	// }
	return nil
}
