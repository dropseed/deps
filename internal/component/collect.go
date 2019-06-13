package component

import (
	"os"

	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/schema"
)

func (r *Runner) Collect(inputPath string) (*schema.Dependencies, error) {
	output.Event("Collecting with %s", r.Given)

	outputPath, err := r.run(r.getCommand(r.Config.Collect, "collect"), inputPath)
	if err != nil {
		return nil, err
	}
	if !output.IsDebug() {
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
