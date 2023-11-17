// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/http/httptest"
	"testing"

	"github.com/go-vela/vela-artifactory/cmd/vela-artifactory/mock"
)

func TestArtifactory_Delete_Exec(t *testing.T) {
	// setup types
	s := httptest.NewServer(mock.Handlers())

	p := &Plugin{
		Config: &Config{
			Action:   "delete",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      s.URL,
			Username: mock.Username,
			Password: mock.Password,
		},
		Copy: &Copy{},
		Delete: &Delete{
			Recursive: false,
			Path:      "foo/bar",
		},
		SetProp: &SetProp{},
		Upload:  &Upload{},
	}

	err := p.Exec()
	if err != nil {
		t.Errorf("Exec returned err %v", err)
	}
}

func TestArtifactory_Delete_Exec_NotFound(t *testing.T) {
	// setup types
	s := httptest.NewServer(mock.Handlers())

	p := &Plugin{
		Config: &Config{
			Action:   "delete",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      s.URL,
			Username: mock.Username,
			Password: mock.Password,
		},
		Copy: &Copy{},
		Delete: &Delete{
			Recursive: false,
			Path:      "bar/bar/foo",
		},
		SetProp: &SetProp{},
		Upload:  &Upload{},
	}

	err := p.Exec()
	if err != nil {
		t.Errorf("Exec returned err %v", err)
	}
}

func TestArtifactory_Delete_Exec_Error(t *testing.T) {
	// setup types
	config := &Config{
		Action:   "copy",
		APIKey:   mock.APIKey,
		DryRun:   false,
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
		Password: mock.Password,
	}

	cli, err := config.New()
	if err != nil {
		t.Errorf("Unable to create Artifactory client: %v", err)
	}

	d := &Delete{
		Recursive: false,
		Path:      "foo/bar",
	}

	err = d.Exec(*cli)
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestArtifactory_Delete_Validate(t *testing.T) {
	// setup types
	d := &Delete{
		Recursive: false,
		Path:      "foo/bar",
	}

	err := d.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Delete_Validate_NoPath(t *testing.T) {
	// setup types
	d := &Delete{
		Recursive: false,
	}

	err := d.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
