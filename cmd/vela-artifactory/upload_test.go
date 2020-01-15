// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import "testing"

func TestArtifactory_Upload_Exec_Error(t *testing.T) {
	// setup types
	config := &Config{
		Action:   "upload",
		APIKey:   "superSecretAPIKey",
		Password: "superSecretPassword",
		URL:      "http://localhost:8081/artifactory",
		Username: "octocat",
	}

	cli, err := config.New()
	if err != nil {
		t.Errorf("Unable to create Artifactory client: %v", err)
	}

	u := &Upload{
		ArgsFile:    "",
		DryRun:      false,
		Flat:        true,
		IncludeDirs: false,
		Recursive:   true,
		Regexp:      false,
		Path:        "foo/bar",
		Sources:     []string{"baz.txt"},
	}

	err = u.Exec(cli)
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestArtifactory_Upload_Validate(t *testing.T) {
	// setup types
	u := &Upload{
		ArgsFile:    "",
		DryRun:      false,
		Flat:        true,
		IncludeDirs: false,
		Recursive:   true,
		Regexp:      false,
		Path:        "foo/bar",
		Sources:     []string{"baz.txt"},
	}

	err := u.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Upload_Validate_NoPath(t *testing.T) {
	// setup types
	u := &Upload{
		ArgsFile:    "",
		DryRun:      false,
		Flat:        true,
		IncludeDirs: false,
		Recursive:   true,
		Regexp:      false,
		Sources:     []string{"baz.txt"},
	}

	err := u.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Upload_Validate_NoSources(t *testing.T) {
	// setup types
	u := &Upload{
		ArgsFile:    "",
		DryRun:      false,
		Flat:        true,
		IncludeDirs: false,
		Recursive:   true,
		Regexp:      false,
		Path:        "foo/bar",
	}

	err := u.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
