package component

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/pkg/schema"
)

func (r *Runner) Act(inputDependencies *schema.Dependencies) (*schema.Dependencies, error) {
	output.Event("Updating with %s", r.Given)

	inputFilename, err := inputTempFile(inputDependencies)
	if err != nil {
		return nil, err
	}
	if !output.IsDebug() {
		defer os.Remove(inputFilename)
	}

	// TODO hooks??

	outputPath, err := r.run(r.getCommand(r.Config.Act, "act"), inputFilename)
	if err != nil {
		return nil, err
	}
	if !output.IsDebug() {
		defer os.Remove(outputPath)
	}

	outputDependencies, err := schema.NewDependenciesFromJSONPath(outputPath)
	if err != nil {
		return nil, err
	}

	return outputDependencies, nil
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
