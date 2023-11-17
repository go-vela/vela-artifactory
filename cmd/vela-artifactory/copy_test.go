// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/http/httptest"
	"testing"

	"github.com/go-vela/vela-artifactory/cmd/vela-artifactory/mock"
)

func TestArtifactory_Copy_Exec(t *testing.T) {
	// setup types
	s := httptest.NewServer(mock.Handlers())

	p := &Plugin{
		Config: &Config{
			Action:   "copy",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			Password: mock.Password,
			URL:      s.URL,
			Username: mock.Username,
		},
		Copy: &Copy{
			Flat:      false,
			Recursive: false,
			Path:      "foo/bar",
			Target:    "bar/foo",
		},
		Delete:  &Delete{},
		SetProp: &SetProp{},
		Upload:  &Upload{},
	}

	err := p.Exec()
	if err != nil {
		t.Errorf("Exec returned err %v", err)
	}
}

func TestArtifactory_Copy_Exec_Error(t *testing.T) {
	// setup types
	config := &Config{
		Action:   "copy",
		APIKey:   mock.APIKey,
		DryRun:   false,
		Password: mock.Password,
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
	}

	cli, err := config.New()
	if err != nil {
		t.Errorf("Unable to create Artifactory client: %v", err)
	}

	c := &Copy{
		Flat:      false,
		Recursive: false,
		Path:      "foo/bar",
		Target:    "bar/foo",
	}

	err = c.Exec(*cli)
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

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
