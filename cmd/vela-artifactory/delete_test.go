// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import "testing"

func TestArtifactory_Delete_Exec_Error(t *testing.T) {
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

	d := &Delete{
		Recursive: false,
		Path:      "foo/bar",
	}

	err = d.Exec(cli)
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestArtifactory_Delete_Validate(t *testing.T) {
	// setup types
	d := &Delete{
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
		Recursive: false,
	}

	err := d.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
