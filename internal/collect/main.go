package collect

import (
	"encoding/json"
	"fmt"

	"github.com/dependencies-io/deps/internal/schema"
)

func getOutputForDependenciesJSONPath(path string) (string, error) {
	err := schema.ValidateDependenciesJSONPath(path)
	if err != nil {
		return "", err
	}

	deps, err := schema.NewDependenciesFromJSONPath(path)
	if err != nil {
		return "", err
	}

	jsonOutput, err := json.Marshal(deps)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("<Dependencies>%s</Dependencies>", string(jsonOutput)), nil
}

// DependenciesJSONPath collects the data by sending it to stdout
func DependenciesJSONPath(path string) error {
	s, err := getOutputForDependenciesJSONPath(path)
	if err != nil {
		return err
	}
	fmt.Println(s)
	return nil
}
