package hooks

import (
	"os"
	"os/exec"
	"strings"

	"github.com/dropseed/deps/internal/output"
	"github.com/dropseed/deps/internal/pullrequest"
)

func RunPullrequestHook(pr pullrequest.PullrequestAdapter, hookName string) error {
	hookCmd := pr.GetSetting(hookName)

	if hookCmd == nil {
		return nil
	}

	output.Event("Executing %s hook", hookName)

	hookCmdLines := strings.Split(hookCmd.(string), "\n")

	for _, line := range hookCmdLines {
		cmd := exec.Command("sh", "-c", line)
		// specific env too? would be pr.Config.Env?
		cmd.Env = os.Environ()
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}

	return nil
}
