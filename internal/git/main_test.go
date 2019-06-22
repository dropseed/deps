package git

import (
	"testing"

	"github.com/dropseed/deps/internal/output"
)

// func TestGetJobBranchName(t *testing.T) {
// 	os.Setenv("JOB_ID", "1111-2222-3333")
// 	if s, err := GetJobBranchName(); s != "deps/update-1111" || err != nil {
// 		t.Error(s)
// 	}

// 	os.Setenv("JOB_ID", "")
// 	if s, err := GetJobBranchName(); s != "" || err == nil {
// 		t.Error(err)
// 	}

// 	os.Setenv("JOB_ID", "abc")
// 	if s, err := GetJobBranchName(); s != "deps/update-abc" || err != nil {
// 		t.Error(s)
// 	}

// 	os.Setenv("JOB_ID", "1111-2222-3333")
// 	os.Setenv("SETTING_BRANCH_PREFIX", "foo/")
// 	if s, err := GetJobBranchName(); s != "foo/deps/update-1111" || err != nil {
// 		t.Error(s)
// 	}
// }

func TestBranchExists(t *testing.T) {
	output.Verbosity = 1
	if !BranchExists("master") {
		t.FailNow()
	}

	if BranchExists("foo") {
		t.FailNow()
	}
}

// func TestGetDepsBranches(t *testing.T) {
// 	branches := GetDepsBranches()
// 	for _, b := range branches {
// 		println(b)
// 	}
// 	t.Fail()
// }
