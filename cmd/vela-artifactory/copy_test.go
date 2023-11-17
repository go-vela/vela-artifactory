// SPDX-License-Identifier: Apache-2.0

package main

import "testing"

func TestArtifactory_Copy_Exec_Error(t *testing.T) {
	// setup types
	config := &Config{
		Action:   "copy",
		APIKey:   "superSecretAPIKey",
		DryRun:   false,
		Password: "superSecretPassword",
		URL:      "http://localhost:8081/artifactory",
		Username: "octocat",
	}

	cli, err := config.New()
	if err != nil {
		t.Errorf("Unable to create Artifactory client: %v", err)
	}

	c := &Copy{
		Flat:      false,
		Recursive: false,
		Path:      "foo/bar",
		Target:    "bar/foo",
	}

	err = c.Exec(*cli)
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestArtifactory_Copy_Validate(t *testing.T) {
	// setup types
	c := &Copy{
		Flat:      false,
		Recursive: false,
		Path:      "foo/bar",
		Target:    "bar/foo",
	}

	err := c.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Copy_Validate_NoPath(t *testing.T) {
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
		Path:      "foo/bar",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
