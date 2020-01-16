// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"
)

func TestArtifactory_Config_New(t *testing.T) {
	// setup types
	c := &Config{
		Action:   "copy",
		APIKey:   "superSecretAPIKey",
		DryRun:   false,
		Password: "superSecretPassword",
		URL:      "https://myarti.com/artifactory",
		Username: "octocat",
	}

	got, err := c.New()
	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if got == nil {
		t.Errorf("New is nil")
	}
}

func TestArtifactory_Config_Validate(t *testing.T) {
	// setup types
	c := &Config{
		Action:   "copy",
		APIKey:   "superSecretAPIKey",
		DryRun:   false,
		Password: "superSecretPassword",
		URL:      "https://myarti.com/artifactory",
		Username: "octocat",
	}

	err := c.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Config_Validate_NoAction(t *testing.T) {
	// setup types
	c := &Config{
		APIKey:   "superSecretAPIKey",
		DryRun:   false,
		Password: "superSecretPassword",
		URL:      "https://myarti.com/artifactory",
		Username: "octocat",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Config_Validate_NoAPIKeyOrPassword(t *testing.T) {
	// setup types
	c := &Config{
		Action:   "copy",
		DryRun:   false,
		URL:      "https://myarti.com/artifactory",
		Username: "octocat",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Config_Validate_NoAPIKey(t *testing.T) {
	// setup types
	c := &Config{
		Action:   "copy",
		DryRun:   false,
		Password: "superSecretPassword",
		URL:      "https://myarti.com/artifactory",
		Username: "octocat",
	}

	err := c.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Config_Validate_NoPassword(t *testing.T) {
	// setup types
	c := &Config{
		Action:   "copy",
		DryRun:   false,
		APIKey:   "superSecretAPIKey",
		URL:      "https://myarti.com/artifactory",
		Username: "octocat",
	}

	err := c.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Config_Validate_NoUrl(t *testing.T) {
	// setup types
	c := &Config{
		Action:   "copy",
		DryRun:   false,
		APIKey:   "superSecretAPIKey",
		Password: "superSecretPassword",
		Username: "octocat",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Config_Validate_NoUsername(t *testing.T) {
	// setup types
	c := &Config{
		Action:   "copy",
		DryRun:   false,
		APIKey:   "superSecretAPIKey",
		Password: "superSecretPassword",
		URL:      "https://myarti.com/artifactory",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
