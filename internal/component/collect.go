package component

import (
	"os"

	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/pkg/schema"
)

func (r *Runner) Collect(inputPath string) (*schema.Dependencies, error) {
	if output.Verbosity > 0 {
		output.Event("Collecting with %s...", r.Given)
	} else {
		output.Event("Collecting with %s...", r.GetName())
	}
	output.Debug("Input path: %s", inputPath)

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
