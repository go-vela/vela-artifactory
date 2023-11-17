// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/http/httptest"
	"testing"

	"github.com/go-vela/vela-artifactory/cmd/vela-artifactory/mock"
)

func TestArtifactory_DockerPromote_Exec(t *testing.T) {
	// setup types
	s := httptest.NewServer(mock.Handlers())

	p := &Plugin{
		Config: &Config{
			Action:   "docker-promote",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      s.URL,
			Username: mock.Username,
			Password: mock.Password,
		},
		Copy:   &Copy{},
		Delete: &Delete{},
		DockerPromote: &DockerPromote{
			TargetRepo:     "docker",
			DockerRegistry: "github/octocat",
		},
		SetProp: &SetProp{},
		Upload:  &Upload{},
	}

	err := p.Exec()
	if err != nil {
		t.Errorf("Exec returned err %v", err)
	}
}

func TestArtifactory_DockerPromote_Exec_Error(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "docker-promote",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      mock.InvalidArtifactoryServerURL,
			Username: mock.Username,
			Password: mock.Password,
		},
		Copy:   &Copy{},
		Delete: &Delete{},
		DockerPromote: &DockerPromote{
			TargetRepo:     "docker",
			DockerRegistry: "github/octocat",
		},
		SetProp: &SetProp{},
		Upload:  &Upload{},
	}

	err := p.Exec()
	if err != nil {
		t.Errorf("Exec returned err %v", err)
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
