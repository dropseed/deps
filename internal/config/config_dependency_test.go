package config

import "testing"

func TestEmptyCompile(t *testing.T) {
	dep := Dependency{}
	dep.Compile()
	if dep.Type != "" {
		t.Error("wrong type")
	}
	if dep.Path != "." {
		t.Error("wrong path")
	}
	if dep.Settings == nil {
		t.Error("settings nil")
	}
	if !*dep.LockfileUpdates.Enabled {
		t.Error("lockfile updates disabled")
	}
	if dep.LockfileUpdates.Settings == nil {
		t.Error("lockfile settings nil")
	}
	if !*dep.ManifestUpdates.Enabled {
		t.Error("manifest updates disabled")
	}
	if dep.ManifestUpdates.Settings == nil {
		t.Error("manifest settings nil")
	}
	if len(dep.ManifestUpdates.Filters) != 1 {
		t.Error("manifest filters wrong")
	}
	filter := dep.ManifestUpdates.Filters[0]
	if filter.Name != ".*" {
		t.Error("filter name wrong")
	}
	if !*filter.Enabled {
		t.Error("filter disabled")
	}
	if *filter.Group {
		t.Error("filter group wrong")
	}
}
