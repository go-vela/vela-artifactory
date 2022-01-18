// Copyright (c) 2022 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"net/http/httptest"
	"testing"

	"github.com/target/go-arty/v2/artifactory/fixtures/docker"
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

	tests := []struct {
		dockerPromote *DockerPromote
		want          *error
	}{
		{ // dockerPromote that does not have a promoteProperty
			dockerPromote: &DockerPromote{
				TargetRepo:      "docker",
				DockerRegistry:  "github/octocat",
				TargetTags:      []string{"test"},
				PromoteProperty: false,
			},
			want: nil,
		},
		//TODO: investigate ways of combining handlers
		//
		// Go-Arty API docs for handlers:
		//	"github.com/target/go-arty/artifactory/fixtures/docker"
		//	"github.com/target/go-arty/artifactory/fixtures/storage"
		//
		// 	{ // dockerPromote that does have a promoteProperty
		// 		dockerPromote: &DockerPromote{
		// 			TargetRepo:      "docker",
		// 			DockerRegistry:  "github/octocat",
		// 			TargetTags:      []string{"test"},
		// 			PromoteProperty: true,
		// 		},
		// 		want: nil,
		// 	},
	}

	for _, test := range tests {
		err := test.dockerPromote.Exec(config)
		if err != nil {
			t.Errorf("Exec should have returned err: %v", err)
		}
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
