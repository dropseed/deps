package component

import (
	"os"

	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/schema"
)

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
