package pullrequest

// func TestNewPullrequestFromEnv(t *testing.T) {
// 	config := config.Config{}
// 	// config.LoadFlags()
// 	config.LoadEnvSettings()
//
// 	os.Setenv("GIT_BRANCH", "tester")
//
// 	pr := NewPullrequestFromEnv("branch-name", "pr title", "pr body", &config)
//
// 	if pr.DefaultBaseBranch != "tester" {
// 		t.Error("DefaultBaseBranch value incorrect")
// 	}
//
// 	if pr.Branch != "branch-name" {
// 		t.Error("Branch value incorrect")
// 	}
//
// 	if pr.Title != "pr title" {
// 		t.Error("Title value incorrect")
// 	}
//
// 	if pr.Body != "pr body" {
// 		t.Error("Body value incorrect")
// 	}
//
// 	if pr.Config != &config {
// 		t.Error("Config value incorrect")
// 	}
// }
