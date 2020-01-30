// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"net/http/httptest"
	"testing"

	"github.com/target/go-arty/artifactory/fixtures/docker"
)

func TestArtifactory_DockerPromote_Exec(t *testing.T) {
	// Create http test server from our fake API handler
	s := httptest.NewServer(docker.FakeHandler())

	// setup types
	config := &Config{
		Action:   "docker-promote",
		APIKey:   "superSecretAPIKey",
		DryRun:   false,
		Password: "superSecretPassword",
		URL:      s.URL,
		Username: "octocat",
	}

	// setup types
	p := &DockerPromote{
		TargetRepo:     "docker",
		DockerRegistry: "github/octocat",
	}

	err := p.Exec(config)
	if err != nil {
		t.Errorf("Exec should have returned err: %w", err)
	}
}

func TestArtifactory_DockerPromote_Validate(t *testing.T) {
	// setup types
	p := &DockerPromote{
		TargetRepo:     "docker",
		DockerRegistry: "github/octocat",
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_DockerPromote_Validate_TargetRepo(t *testing.T) {
	// setup types
	p := &DockerPromote{
		DockerRegistry: "github/octocat",
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_DockerPromote_Validate_DockerRepo(t *testing.T) {
	// setup types
	p := &DockerPromote{
		TargetRepo: "docker",
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
