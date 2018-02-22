package git

import (
	"os"
	"testing"
)

func TestGetJobBranchName(t *testing.T) {
	os.Setenv("JOB_ID", "1111-2222-3333")
	if s, err := GetJobBranchName(); s != "deps/update-1111" || err != nil {
		t.Error(s)
	}

	os.Setenv("JOB_ID", "")
	if s, err := GetJobBranchName(); s != "" || err == nil {
		t.Error(err)
	}

	os.Setenv("JOB_ID", "abc")
	if s, err := GetJobBranchName(); s != "deps/update-abc" || err != nil {
		t.Error(s)
	}

	os.Setenv("JOB_ID", "1111-2222-3333")
	os.Setenv("SETTING_BRANCH_PREFIX", "foo/")
	if s, err := GetJobBranchName(); s != "foo/deps/update-1111" || err != nil {
		t.Error(s)
	}
}
