// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"
)

func TestArtifactory_DockerPromote_Exec(t *testing.T) {
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
