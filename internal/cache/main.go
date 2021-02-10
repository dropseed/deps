package cache

import (
	"os"
	"path"
)

const DefaultCacheDirName = "deps"

func GetCachePath() string {
	userCache, err := os.UserCacheDir()
	if err != nil {
		panic(err)
	}

	cachePath := path.Join(userCache, DefaultCacheDirName)
	if err := os.MkdirAll(cachePath, os.ModePerm); err != nil {
		panic(err)
	}
	return cachePath
}
