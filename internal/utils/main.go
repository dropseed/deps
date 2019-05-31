package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"

	"github.com/dropseed/deps/internal/output"
	"github.com/imdario/mergo"
)

func TemporaryCopyOfDir(dirToCopy string) (string, error) {
	tmpPath := "" // use the default

	if runtime.GOOS == "darwin" {
		// the default is not shared with docker, by default
		tmpPath = "/tmp"
	}

	dir, err := ioutil.TempDir(tmpPath, "deps-")
	if err != nil {
		return "", err
	}

	repoDir := path.Join(dir, "repo")

	// TODO better cross platform way to do this? things I tried all had an issue
	cmd := exec.Command("cp", "-a", dirToCopy, repoDir)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	output.Debug("Made temporary copy of %s into %s\n", dirToCopy, repoDir)

	return repoDir, nil
}

func ParseTagFromString(tag, output string) (string, error) {
	pattern := fmt.Sprintf("\\<%s\\>.*?<\\/%s\\>", tag, tag)
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return "", err
	}

	matches := regex.FindAllString(output, -1)
	if len(matches) == 0 {
		return "", nil
	}

	mergedMatches := map[string]interface{}{}

	// Account for <Tag> and </Tag> when trimming string
	beginningOffset := len(tag) + 2
	endingOffset := beginningOffset + 1

	for _, match := range matches {
		match = match[beginningOffset : len(match)-endingOffset]
		matchData := map[string]interface{}{}
		if err = json.Unmarshal([]byte(match), &matchData); err != nil {
			return "", err
		}

		if err = mergo.Merge(&mergedMatches, matchData); err != nil {
			return "", err
		}
	}

	merged, err := json.Marshal(mergedMatches)
	if err != nil {
		return "", err
	}

	return string(merged), nil
}

func StringSliceToMap(slice []string) map[string]bool {
	ret := map[string]bool{}
	for _, s := range slice {
		ret[s] = true
	}
	return ret
}
