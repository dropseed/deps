package schema

import "testing"

func TestInvalidSchema(t *testing.T) {
	err := ValidateDependenciesJSONPath("testdata/invalid.json")
	if err == nil {
		t.Fail()
	}
	if err.Error() != "The document is not valid. see errors:\n- bad: Additional property bad is not allowed\n" {
		t.Error("Error string doesn't match")
	}
}

func TestValidSchema(t *testing.T) {
	err := ValidateDependenciesJSONPath("testdata/single_dependency.json")
	if err != nil {
		t.Error(err)
	}
}
