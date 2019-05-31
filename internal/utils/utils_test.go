package utils

import "testing"

func TestParseDependencies(t *testing.T) {
	match, err := ParseTagFromString("Dependencies", "Some string output <Dependencies>{\"testing\": true}</Dependencies> and more.")
	if err != nil {
		t.Error(err)
	}
	if match != "{\"testing\":true}" {
		t.FailNow()
	}
}

func TestParseActions(t *testing.T) {
	s := `<Actions>{"PR #0":{"dependencies":{"manifests":{"Dockerfile":{"current":{"dependencies":{"ubuntu":{"source":"dockerhub","constraint":"latest"}}},"updated":{"dependencies":{"ubuntu":{"source":"dockerhub","constraint":"16.04"}}}}}},"metadata":{"foo":"bar"}}}</Actions>`
	_, err := ParseTagFromString("Actions", s)
	if err != nil {
		t.Error(err)
	}
	// expected := "{\"PR #0\":{\"dependencies\":{\"manifests\":{\"Dockerfile\":{\"current\":{\"dependencies\":{\"ubuntu\":{\"source\":\"dockerhub\",\"constraint\":\"latest\"}}},\"updated\":{\"dependencies\":{\"ubuntu\":{\"source\":\"dockerhub\",\"constraint\":\"16.04\"}}}}}},\"metadata\":{\"foo\":\"bar\"}}}"
	// if match != expected {
	// 	println(match)
	// 	t.FailNow()
	// }
}
