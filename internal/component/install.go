package component

import (
	"os"
	"os/exec"

	"github.com/dropseed/deps/internal/output"
)

func (r *Runner) Install() error {
	// TODO install does not need to happen every time
	// - yes it does if local (assume could have changed)
	// - yes it does if new (just cloned)
	// - yes it does if ref chagned
	// so maybe not horrible if it runs every time right now

	output.Event("Installing %s", r.Given)

	command := r.Config.Install
	if override := r.getOverrideFromEnv("install"); override != "" {
		command = override
	}
	output.Debug(command)

	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = r.Path
	cmd.Env = r.Env
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
