package main

import "testing"

func TestArtifactory_Copy_Validate(t *testing.T) {
	// setup types
	c := &Copy{
		Flat:      false,
		Recursive: false,
		Source:    "foo/bar",
		Target:    "bar/foo",
	}

	err := c.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Copy_Validate_NoSource(t *testing.T) {
	// setup types
	c := &Copy{
		Flat:      false,
		Recursive: false,
		Target:    "bar/foo",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Copy_Validate_NoTarget(t *testing.T) {
	// setup types
	c := &Copy{
		Flat:      false,
		Recursive: false,
		Source:    "foo/bar",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
