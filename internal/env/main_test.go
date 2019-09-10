package env

import "testing"

func TestSettingVal(t *testing.T) {
	if s, _ := settingValToEnv(2); s != "2" {
		t.FailNow()
	}
}
