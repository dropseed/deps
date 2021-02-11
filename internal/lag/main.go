package lag

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/dropseed/deps/internal/cache"
	"github.com/dropseed/deps/internal/git"
	"github.com/dropseed/deps/internal/install"
	"github.com/dropseed/deps/internal/output"
)

const lagFilename = "lag.json"

type lagData struct {
	Lockfiles map[string]string `json:"lockfiles"`
}

type lagManager struct {
	dataPath string
	data     *lagData
}

func NewLagManager() (*lagManager, error) {
	depsCache := cache.GetCachePath()
	lagDataPath := path.Join(depsCache, lagFilename)
	manager := lagManager{
		dataPath: lagDataPath,
	}

	if err := manager.loadData(); err != nil {
		return nil, err
	}

	return &manager, nil
}

func (manager *lagManager) loadData() error {
	jsonFile, err := os.Open(manager.dataPath)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var data lagData
	if err := json.Unmarshal(byteValue, &data); err != nil {
		return err
	}

	manager.data = &data

	return nil
}

func (manager *lagManager) saveData() error {
	out, err := json.MarshalIndent(manager.data, "", "  ")
	if err != nil {
		return err
	}
	out = append(out, "\n"...)
	if err := ioutil.WriteFile(manager.dataPath, out, 0644); err != nil {
		panic(err)
	}
	return nil
}

func (manager *lagManager) getSavedLockfileIdentifier(path, id string) string {
	if manager.data == nil {
		return ""
	}
	existingID, ok := manager.data.Lockfiles[path]
	if !ok {
		return ""
	}
	return existingID
}

func (manager *lagManager) SaveLockfileIdentifier(path, id string) error {
	// Reload before saving to prevent overwrites in another process
	if err := manager.loadData(); err != nil {
		return err
	}

	if manager.data == nil {
		manager.data = &lagData{
			Lockfiles: map[string]string{},
		}
	} else if manager.data.Lockfiles == nil {
		manager.data.Lockfiles = map[string]string{}
	}

	manager.data.Lockfiles[path] = id

	return manager.saveData()
}

func IdentifierForFile(p string) string {
	// Can't use mtime because git changes it on checkout etc.
	cmd := exec.Command("git", "hash-object", p)
	out, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(string(out))
}

func Run(dir string) error {
	if git.RepoRoot() == "" {
		return nil
	}

	// TODO if components are responsible for this, then we do it this way
	// Disable all output (FindOrInfer prints some stuff right now...)
	// tempVerbosity := output.Verbosity
	// output.Verbosity = -1
	// cfg, err := config.FindOrInfer()
	// if err != nil {
	// 	return err
	// }
	// output.Verbosity = tempVerbosity

	manager, err := NewLagManager()
	if err != nil {
		return err
	}

	laggingLockfiles := map[string]string{}

	for _, lockfile := range install.FindLockfiles(dir) {
		lockfileID := IdentifierForFile(lockfile.Path)

		existingID := manager.getSavedLockfileIdentifier(lockfile.Path, lockfileID)

		if lockfileID == existingID { //} || existingID == "ignore" {
			continue
		}

		// Use existing ID so we can tell what to do
		// "" would be not installed
		// "ignore" would be skip (not shown)
		// "anything else" would be re-install
		laggingLockfiles[lockfile.RelPath()] = existingID
	}

	if len(laggingLockfiles) > 0 {
		output.Warning("\nDeps: You have some dependencies that should be installed!\n")

		for lockfilePath, lockfileID := range laggingLockfiles {
			if lockfileID == "" {
				output.Unstyled("- %s has not been installed with deps, so can't be tracked", lockfilePath)
			} else {
				output.Unstyled("- %s has been updated but not installed", lockfilePath)
			}
		}

		output.Warning("\nRun `deps install`")
	}

	return nil
}
