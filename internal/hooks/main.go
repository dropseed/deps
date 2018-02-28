package hooks

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/dependencies-io/deps/internal/env"
)

// Run a specified hook by name
func Run(name string) error {
	hooksString := env.GetSetting(name, "")

	if hooksString == "" {
		return nil
	}

	fmt.Printf("Running hooks for '%s'\n", name)

	var hooks []string
	if err := json.Unmarshal([]byte(hooksString), &hooks); err != nil {
		return err
	}

	totalHooks := len(hooks)

	for index, hook := range hooks {
		fmt.Printf("Running '%s' %d/%d\n", name, index+1, totalHooks)
		fmt.Printf("-> %s\n", hook)
		cmdOutput, err := exec.Command("sh", "-c", hook).CombinedOutput()
		fmt.Println(string(cmdOutput))
		if err != nil {
			return err
		}
	}

	return nil
}
