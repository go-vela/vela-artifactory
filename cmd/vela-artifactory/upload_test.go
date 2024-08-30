// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/http/httptest"
	"testing"

	"github.com/go-vela/vela-artifactory/cmd/vela-artifactory/mock"
)

func TestArtifactory_Plugin_Exec_Upload(t *testing.T) {
	// setup types
	s := httptest.NewServer(mock.Handlers())

	p := &Plugin{
		Config: &Config{
			Action:   "upload",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      s.URL,
			Username: mock.Username,
			Password: mock.Password,
			Client: &Client{
				Retries:            3,
				RetryWaitMilliSecs: 1,
			},
		},
		Copy:    &Copy{},
		Delete:  &Delete{},
		SetProp: &SetProp{},
		Upload: &Upload{
			Flat:        true,
			IncludeDirs: false,
			Recursive:   true,
			Regexp:      false,
			Path:        "foo/bar",
			Sources:     []string{"mock/testdata/baz.txt"},
		},
	}

	err := p.Exec()
	if err != nil {
		t.Errorf("Exec returned err %v", err)
	}
}

func TestArtifactory_Plugin_Exec_UploadWithBuildProps(t *testing.T) {
	// setup types
	s := httptest.NewServer(mock.Handlers())

	p := &Plugin{
		Config: &Config{
			Action:   "upload",
			Token:    mock.Token,
			APIKey:   mock.APIKey,
			DryRun:   false,
			URL:      s.URL,
			Username: mock.Username,
			Password: mock.Password,
			Client: &Client{
				Retries:            3,
				RetryWaitMilliSecs: 1,
			},
		},
		Copy:    &Copy{},
		Delete:  &Delete{},
		SetProp: &SetProp{},
		Upload: &Upload{
			Flat:        true,
			IncludeDirs: false,
			Recursive:   true,
			Regexp:      false,
			Path:        "foo/bar",
			Sources:     []string{"mock/testdata/baz.txt"},
			BuildProps:  "build.name=buildName;build.number=17;build.timestamp=1600856623553",
		},
	}

	err := p.Exec()
	if err != nil {
		t.Errorf("Exec returned err %v", err)
	}
}

func TestArtifactory_Upload_Exec_Error(t *testing.T) {
	// setup types
	config := &Config{
		Action:   "upload",
		APIKey:   mock.APIKey,
		DryRun:   false,
		URL:      mock.InvalidArtifactoryServerURL,
		Username: mock.Username,
		Password: mock.Password,
		Client: &Client{
			Retries:            3,
			RetryWaitMilliSecs: 1,
		},
	}

	cli, err := config.New()
	if err != nil {
		t.Errorf("Unable to create Artifactory client: %v", err)
	}

	u := &Upload{
		Flat:        true,
		IncludeDirs: false,
		Recursive:   true,
		Regexp:      false,
		Path:        "foo/bar",
		Sources:     []string{"baz.txt"},
	}

	err = u.Exec(*cli)
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestArtifactory_Upload_Validate(t *testing.T) {
	// setup types
	u := &Upload{
		Flat:        true,
		IncludeDirs: false,
		Recursive:   true,
		Regexp:      false,
		Path:        "foo/bar",
		Sources:     []string{"baz.txt"},
	}

	err := u.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Upload_Validate_NoPath(t *testing.T) {
	// setup types
	u := &Upload{
		Flat:        true,
		IncludeDirs: false,
		Recursive:   true,
		Regexp:      false,
		Sources:     []string{"baz.txt"},
	}

	err := u.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Upload_Validate_NoSources(t *testing.T) {
	// setup types
	u := &Upload{
		Flat:        true,
		IncludeDirs: false,
		Recursive:   true,
		Regexp:      false,
		Path:        "foo/bar",
	}

	err := u.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
