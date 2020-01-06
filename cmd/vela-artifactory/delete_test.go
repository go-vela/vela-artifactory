package main

import "testing"

func TestArtifactory_Delete_Validate(t *testing.T) {
	// setup types
	d := &Delete{
		ArgsFile:  "",
		DryRun:    false,
		Recursive: false,
		Path:      "foo/bar",
	}

	err := d.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Delete_Validate_NoPath(t *testing.T) {
	// setup types
	d := &Delete{
		ArgsFile:  "",
		DryRun:    false,
		Recursive: false,
	}

	err := d.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
