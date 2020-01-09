// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import "testing"

func TestArtifactory_Config_Validate(t *testing.T) {
	// setup types
	c := &Config{
		APIKey:   "superSecretAPIKey",
		Password: "superSecretPassword",
		URL:      "https://myarti.com/artifactory",
		Username: "octocat",
	}

	err := c.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Config_Validate_NoAPIKeyOrPassword(t *testing.T) {
	// setup types
	c := &Config{
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
		APIKey:   "superSecretAPIKey",
		Password: "superSecretPassword",
		URL:      "https://myarti.com/artifactory",
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
