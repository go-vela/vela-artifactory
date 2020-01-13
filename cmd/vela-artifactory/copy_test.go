// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import "testing"

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
