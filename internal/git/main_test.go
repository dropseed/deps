package git

import (
	"testing"
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

// func TestBranchExists(t *testing.T) {
// 	output.Verbosity = 1
// 	if !BranchExists("master") {
// 		t.FailNow()
// 	}

// 	if BranchExists("foo") {
// 		t.FailNow()
// 	}
// }

func TestGitRemoteToHTTPS(t *testing.T) {
	tests := []struct {
		input  string
		output string
	}{
		{
			input:  "git@github.com:dropseed/test.git",
			output: "https://github.com/dropseed/test.git",
		},
		{
			input:  "git@github.com:dropseed/test",
			output: "https://github.com/dropseed/test",
		},
		{
			input:  "git@gitlab.com:dropseed/test.git",
			output: "https://gitlab.com/dropseed/test.git",
		},
		{
			input:  "git@gitlab.com:dropseed/test/two.git",
			output: "https://gitlab.com/dropseed/test/two.git",
		},
	}

	for _, test := range tests {
		actual := GitRemoteToHTTPS(test.input)
		if actual != test.output {
			t.Errorf("%s\n%s != %s", test.input, test.output, actual)
		}
	}
}

func TestGitRemoteHostname(t *testing.T) {
	if GitRemoteHostname() != "github.com" {
		t.Fail()
	}
}

// func TestGetDepsBranches(t *testing.T) {
// 	branches := GetDepsBranches()
// 	for _, b := range branches {
// 		println(b)
// 	}
// 	t.Fail()
// }
