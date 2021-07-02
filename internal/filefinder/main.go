package filefinder

import (
	"io/ioutil"
	"os"
	"path"
	"regexp"
)

const maxDepth = 2

var directoryNamesToSkip = map[string]bool{
	".git":         true,
	"node_modules": true,
	"env":          true,
	"vendor":       true,
	".venv":        true,
}

// FindInDir takes a set of named patterns,
// and returns a map of matched paths and the pattern name they matched
func FindInDir(dir string, patterns map[string]*regexp.Regexp) map[string]string {
	return findInDir(dir, patterns, 1)
}

// Same as FindInDir with a specific depth
func DeepFindInDir(dir string, patterns map[string]*regexp.Regexp, depth int) map[string]string {
	return findInDir(dir, patterns, maxDepth-depth+1)
}

func findInDir(dir string, patterns map[string]*regexp.Regexp, depth int) map[string]string {
	if depth > maxDepth {
		return nil
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil
	}

	matches := map[string]string{}

	for _, f := range files {
		name := f.Name()
		p := path.Join(dir, name)

		fileInfo, err := os.Stat(p)
		if err != nil {
			continue
		}

		if fileInfo.IsDir() {
			if directoryNamesToSkip[name] {
				continue
			}
			for k, v := range findInDir(p, patterns, depth+1) {
				matches[k] = v
			}
		} else if match := patternMatchingPath(p, patterns); match != "" {
			matches[p] = match
		}
	}

	return matches
}

func patternMatchingPath(p string, patterns map[string]*regexp.Regexp) string {
	basename := path.Base(p)
	for patternName, pattern := range patterns {
		if pattern.MatchString(basename) {
			return patternName
		}
	}
	return ""
}
