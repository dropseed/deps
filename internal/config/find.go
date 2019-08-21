package config

import (
	"os"
	"path"
)

func FindFilename(dir string, filenames ...string) string {
	for _, f := range filenames {
		p := path.Join(dir, f)
		if fileExists(p) {
			return p
		}
	}
	return ""
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
