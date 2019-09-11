package runner

// import (
// 	"strings"
// 	"testing"

// 	"github.com/dropseed/deps/internal/config"
// 	"github.com/dropseed/deps/pkg/schema"
// )

// func getDepConfig() *config.Dependency {
// 	content := `version: 3
// dependencies:
// - type: js
// `

// 	config, err := config.NewConfigFromReader(strings.NewReader(content), nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	config.Compile()
// 	return config.Dependencies[0]
// }

// func TestNoDependencies(t *testing.T) {
// 	dependencies, err := schema.NewDependenciesFromJSONPath("./testdata/no_dependencies.json")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	updates, err := newUpdatesFromDependencies(dependencies, getDepConfig())
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if len(updates) != 0 {
// 		t.FailNow()
// 	}
// }

// func TestLockfileUpdate(t *testing.T) {
// 	dependencies, err := schema.NewDependenciesFromJSONPath("./testdata/single_lockfile.json")
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	updates, err := newUpdatesFromDependencies(dependencies, getDepConfig())
// 	if err != nil {
// 		t.Error(err)
// 	}

// 	if len(updates) != 1 {
// 		t.FailNow()
// 	}

// 	if updates[0].summary != "- `yarn.lock` was updated (including 5 updated direct dependencies)" {
// 		t.FailNow()
// 	}
// }

// // func TestManifestUpdates(t *testing.T) {
// // 	dependencies, err := schema.NewDependenciesFromJSONPath("./testdata/two_dependencies.json")
// // 	if err != nil {
// // 		t.Error(err)
// // 	}

// // 	updates, err := newUpdatesFromDependencies(dependencies, getDepConfig())
// // 	if err != nil {
// // 		t.Error(err)
// // 	}

// // 	if len(updates) != 2 {
// // 		println(len(updates))
// // 		t.FailNow()
// // 	}

// // 	if updates[0].summary != "- `pullrequest` in `requirements.txt` from \"0.1.0\" to \"0.3.0\"" {
// // 		println(updates[0].summary)
// // 		t.FailNow()
// // 	}

// // 	if updates[1].summary != "- `requests` in `requirements.txt` from \"1.0.0\" to \"3.0.0\"" {
// // 		println(updates[1].summary)
// // 		t.FailNow()
// // 	}
// // }
