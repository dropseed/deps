package component

import (
	"os"
	"os/exec"

	"github.com/dropseed/deps/internal/output"
)

func (r *Runner) Install() error {
	if !r.shouldInstall {
		output.Debug("Skipping install of %s", r.Given)
		return nil
	}

	output.Event("Installing %s", r.Given)

	command := r.getCommand(r.Config.Install, "install")
	output.Debug(command)

	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = r.Path
	cmd.Env = append(os.Environ(), r.Env...)
	if output.IsDebug() {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
