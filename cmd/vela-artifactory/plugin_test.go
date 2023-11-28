// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"

	"github.com/go-vela/vela-artifactory/cmd/vela-artifactory/mock"
)

func TestArtifactory_Plugin_Exec_InvalidAction(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "foobar",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      mock.InvalidArtifactoryServerURL,
			Username: mock.Username,
			Password: mock.Password,
		},
		Copy:    &Copy{},
		Delete:  &Delete{},
		SetProp: &SetProp{},
		Upload:  &Upload{},
	}

	err := p.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestArtifactory_Plugin_Validate(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "copy",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      mock.InvalidArtifactoryServerURL,
			Username: mock.Username,
			Password: mock.Password,
		},
		Copy: &Copy{
			Flat:      false,
			Recursive: false,
			Path:      "foo/bar",
			Target:    "bar/foo",
		},
		Delete: &Delete{
			Recursive: false,
			Path:      "foo/bar",
		},
		DockerPromote: &DockerPromote{
			TargetRepo:     "docker",
			DockerRegistry: "github/octocat",
		},
		SetProp: &SetProp{
			Path: "foo/bar",
			Props: []*Prop{
				{
					Name:  "foo",
					Value: "bar",
				},
			},
			RawProps: `[{"name": "single", "value": "foo"}]`,
		},
		Upload: &Upload{
			Flat:        true,
			IncludeDirs: false,
			Recursive:   true,
			Regexp:      false,
			Path:        "foo/bar",
			Sources:     []string{"baz.txt"},
		},
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Plugin_Validate_InvalidAction(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "foobar",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      mock.InvalidArtifactoryServerURL,
			Username: mock.Username,
			Password: mock.Password,
		},
		Copy:          &Copy{},
		Delete:        &Delete{},
		DockerPromote: &DockerPromote{},
		SetProp:       &SetProp{},
		Upload:        &Upload{},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Plugin_Validate_NoConfig(t *testing.T) {
	// setup types
	p := &Plugin{
		Config:  &Config{},
		Copy:    &Copy{},
		Delete:  &Delete{},
		SetProp: &SetProp{},
		Upload:  &Upload{},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Plugin_Validate_NoCopy(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "copy",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      mock.InvalidArtifactoryServerURL,
			Username: mock.Username,
			Password: mock.Password,
		},
		Copy:          &Copy{},
		Delete:        &Delete{},
		DockerPromote: &DockerPromote{},
		SetProp:       &SetProp{},
		Upload:        &Upload{},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Plugin_Validate_NoDelete(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "delete",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      mock.InvalidArtifactoryServerURL,
			Username: mock.Username,
			Password: mock.Password,
		},
		Copy:          &Copy{},
		Delete:        &Delete{},
		DockerPromote: &DockerPromote{},
		SetProp:       &SetProp{},
		Upload:        &Upload{},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Plugin_Validate_NoDockerPromote(t *testing.T) {
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
		Copy:          &Copy{},
		Delete:        &Delete{},
		DockerPromote: &DockerPromote{},
		SetProp:       &SetProp{},
		Upload:        &Upload{},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Plugin_Validate_NoSetProp(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "set-prop",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      mock.InvalidArtifactoryServerURL,
			Username: mock.Username,
			Password: mock.Password,
		},
		Copy:          &Copy{},
		Delete:        &Delete{},
		DockerPromote: &DockerPromote{},
		SetProp:       &SetProp{},
		Upload:        &Upload{},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Plugin_Validate_NoUpload(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "upload",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      mock.InvalidArtifactoryServerURL,
			Username: mock.Username,
			Password: mock.Password,
		},
		Copy:          &Copy{},
		Delete:        &Delete{},
		DockerPromote: &DockerPromote{},
		SetProp:       &SetProp{},
		Upload:        &Upload{},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
