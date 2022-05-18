package schemaext

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"

	"github.com/dropseed/deps/pkg/schema"
)

func UpdateIDForDeps(dependencies *schema.Dependencies) string {
	truncated := schema.Dependencies{
		// TODO if type is important to separate updates between components,
		// then can add Dependencies.Type and use that too
		Lockfiles: map[string]*schema.Lockfile{},
		Manifests: map[string]*schema.Manifest{},
	}

	if dependencies.HasLockfiles() {
		for name := range dependencies.Lockfiles {
			// Only care about the filename
			truncated.Lockfiles[name] = nil
		}
	}

	if dependencies.HasManifests() {
		for name, manifest := range dependencies.Manifests {
			if !manifest.HasUpdates() {
				continue
			}

			// Only care about the filename + dependency names
			truncatedManifest := &schema.Manifest{
				Updated: &schema.ManifestVersion{
					Dependencies: map[string]*schema.ManifestDependency{},
				},
			}
			for dep := range manifest.Updated.Dependencies {
				truncatedManifest.Updated.Dependencies[dep] = nil
			}

			truncated.Manifests[name] = truncatedManifest
		}
	}

	return getShortMD5(truncated)
}

func UniqueIDForDeps(dependencies *schema.Dependencies) string {
	return getShortMD5(dependencies)
}

func getShortMD5(i interface{}) string {
	out, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	sum := md5.Sum(out)
	str := hex.EncodeToString(sum[:])
	short := str[:7]
	return short
}
