package install

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"

	"github.com/dropseed/deps/internal/output"
)

type lockfilePattern struct {
	pattern    *regexp.Regexp
	installCmd string
}

var lockfilePatterns = []lockfilePattern{
	lockfilePattern{
		pattern:    regexp.MustCompile("^yarn.lock$"),
		installCmd: "yarn install",
	},
	lockfilePattern{
		pattern:    regexp.MustCompile("^package-lock.json$"),
		installCmd: "npm ci",
	},
	lockfilePattern{
		pattern:    regexp.MustCompile("^Pipfile.lock$"),
		installCmd: "pipenv sync --dev",
	},
	lockfilePattern{
		pattern:    regexp.MustCompile("^poetry.lock$"),
		installCmd: "poetry install",
	},
}

const maxInferenceDepth = 2

var directoryNamesToSkip = map[string]bool{
	".git":         true,
	"node_modules": true,
	"env":          true,
	"vendor":       true,
}

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
	return findLockfiles(dir, 1)
}

func findLockfiles(dir string, depth int) []*lockfile {
	if depth > maxInferenceDepth {
		return nil
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}

	lockfiles := []*lockfile{}

	for _, f := range files {
		name := f.Name()
		p := path.Join(dir, name)

		fileInfo, err := os.Stat(p)
		if err != nil {
			output.Debug("Error os.Stat: %v", err)
			continue
		}

		if fileInfo.IsDir() {
			if directoryNamesToSkip[name] {
				continue
			}

			if dirDeps := findLockfiles(p, depth+1); dirDeps != nil {
				lockfiles = append(lockfiles, dirDeps...)
			}
		} else if lockfile := lockfileMatchingPath(p); lockfile != nil {
			lockfiles = append(lockfiles, lockfile)
		}
	}

	return lockfiles
}

func lockfileMatchingPath(p string) *lockfile {
	basename := path.Base(p)
	for _, lockfilePattern := range lockfilePatterns {
		if lockfilePattern.pattern.MatchString(basename) {
			return &lockfile{
				Path:       p,
				installCmd: lockfilePattern.installCmd,
			}
		}
	}
	return nil
}
