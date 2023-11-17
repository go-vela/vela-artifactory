// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"

	"github.com/go-vela/vela-artifactory/cmd/vela-artifactory/mock"
)

func TestArtifactory_Config_New(t *testing.T) {
	// setup types
	c := &Config{
		Action:   "copy",
		APIKey:   mock.APIKey,
		DryRun:   false,
		Password: mock.Password,
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
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
		APIKey:   mock.APIKey,
		DryRun:   false,
		Password: mock.Password,
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
	}

	err := c.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Config_Validate_NoAction(t *testing.T) {
	// setup types
	c := &Config{
		APIKey:   mock.APIKey,
		DryRun:   false,
		Password: mock.Password,
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
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
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
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
		Password: mock.Password,
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
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
		APIKey:   mock.APIKey,
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
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
		APIKey:   mock.APIKey,
		Password: mock.Password,
		Username: mock.Username,
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
		Password: mock.Password,
		URL:      mock.InvalidArtifactoryServerURL,
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Config_Validate_NoAuth(t *testing.T) {
	// setup types
	c := &Config{
		Action: "copy",
		DryRun: false,
		URL:    mock.InvalidArtifactoryServerURL,
	}

	err := c.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
