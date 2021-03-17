package install

import (
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"

	"github.com/dropseed/deps/internal/filefinder"
	"github.com/dropseed/deps/internal/output"
)

type lockfile struct {
	Path       string
	installCmd string
}

func (lf *lockfile) Install() error {
	output.Event("Installing %s", lf.RelPath())
	cmd := exec.Command("sh", "-c", lf.installCmd)
	cmd.Dir = path.Dir(lf.Path) // Could be down a directory, so make sure we run the command as if we're there
	cmd.Env = os.Environ()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func (lf *lockfile) RelPath() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	relPath, err := filepath.Rel(cwd, lf.Path)
	if err != nil {
		panic(err)
	}
	return relPath
}

func FindLockfiles(dir string) []*lockfile {
	patterns := map[string]*regexp.Regexp{
		"yarn.lock":         regexp.MustCompile("^yarn\\.lock$"),
		"package-lock.json": regexp.MustCompile("^package-lock\\.json$"),
		"Pipfile.lock":      regexp.MustCompile("^Pipfile\\.lock$"),
		"poetry.lock":       regexp.MustCompile("^poetry\\.lock$"),
	}
	commands := map[string]string{
		"yarn.lock":         "yarn install",
		"package-lock.json": "npm ci",
		"Pipfile.lock":      "pipenv sync --dev",
		"poetry.lock":       "poetry install",
	}
	lockfiles := []*lockfile{}
	for path, patternName := range filefinder.FindInDir(dir, patterns) {
		lockfiles = append(lockfiles, &lockfile{
			Path:       path,
			installCmd: commands[patternName],
		})
	}
	return lockfiles
}
