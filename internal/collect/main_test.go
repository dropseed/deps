package collect

import "testing"

func TestGetOutput(t *testing.T) {
	s, err := getOutputForDependenciesJSONPath("../schema/testdata/single_dependency.json")
	if err != nil {
		t.Error(err)
	}
	if s != "<Dependencies>{\"manifests\":{\"/\":{\"current\":{\"dependencies\":{\"pullrequest\":{\"source\":\"go\",\"constraint\":\"0.1.0\"}}},\"updated\":{\"dependencies\":{\"pullrequest\":{\"source\":\"go\",\"constraint\":\"0.3.0\",\"available\":[{\"name\":\"0.3.0\",\"content\":\"Third release\"},{\"name\":\"0.2.0\",\"content\":\"Second release\"}]}}}}}}</Dependencies>" {
		t.Error(s)
	}
}
