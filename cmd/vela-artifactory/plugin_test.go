// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"
)

func TestArtifactory_Plugin_Exec_Copy(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "copy",
			APIKey:   "superSecretAPIKey",
			Password: "superSecretPassword",
			URL:      "https://myarti.com/artifactory",
			Username: "octocat",
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
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestArtifactory_Plugin_Exec_Delete(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "delete",
			APIKey:   "superSecretAPIKey",
			Password: "superSecretPassword",
			URL:      "https://myarti.com/artifactory",
			Username: "octocat",
		},
		Copy: &Copy{},
		Delete: &Delete{
			ArgsFile:  "",
			DryRun:    false,
			Recursive: false,
			Path:      "foo/bar",
		},
		SetProp: &SetProp{},
		Upload:  &Upload{},
	}

	err := p.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestArtifactory_Plugin_Exec_InvalidAction(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "foobar",
			APIKey:   "superSecretAPIKey",
			Password: "superSecretPassword",
			URL:      "https://myarti.com/artifactory",
			Username: "octocat",
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

func TestArtifactory_Plugin_Exec_SetProp(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "set-prop",
			APIKey:   "superSecretAPIKey",
			Password: "superSecretPassword",
			URL:      "https://myarti.com/artifactory",
			Username: "octocat",
		},
		Copy:   &Copy{},
		Delete: &Delete{},
		SetProp: &SetProp{
			Path: "foo/bar",
			Props: []*Prop{
				{
					Name:  "foo",
					Value: "bar",
				},
			},
		},
		Upload: &Upload{},
	}

	err := p.Exec()
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestArtifactory_Plugin_Exec_Upload(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "upload",
			APIKey:   "superSecretAPIKey",
			Password: "superSecretPassword",
			URL:      "https://myarti.com/artifactory",
			Username: "octocat",
		},
		Copy:    &Copy{},
		Delete:  &Delete{},
		SetProp: &SetProp{},
		Upload: &Upload{
			ArgsFile:    "",
			DryRun:      false,
			Flat:        true,
			IncludeDirs: false,
			Recursive:   true,
			Regexp:      false,
			Path:        "foo/bar",
			Sources:     []string{"baz.txt"},
		},
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
			APIKey:   "superSecretAPIKey",
			Password: "superSecretPassword",
			URL:      "https://myarti.com/artifactory",
			Username: "octocat",
		},
		Copy: &Copy{
			Flat:      false,
			Recursive: false,
			Path:      "foo/bar",
			Target:    "bar/foo",
		},
		Delete: &Delete{
			ArgsFile:  "",
			DryRun:    false,
			Recursive: false,
			Path:      "foo/bar",
		},
		SetProp: &SetProp{
			Path: "foo/bar",
			Props: []*Prop{
				{
					Name:  "foo",
					Value: "bar",
				},
			},
		},
		Upload: &Upload{
			ArgsFile:    "",
			DryRun:      false,
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
			APIKey:   "superSecretAPIKey",
			Password: "superSecretPassword",
			URL:      "https://myarti.com/artifactory",
			Username: "octocat",
		},
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
			APIKey:   "superSecretAPIKey",
			Password: "superSecretPassword",
			URL:      "https://myarti.com/artifactory",
			Username: "octocat",
		},
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

func TestArtifactory_Plugin_Validate_NoDelete(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "delete",
			APIKey:   "superSecretAPIKey",
			Password: "superSecretPassword",
			URL:      "https://myarti.com/artifactory",
			Username: "octocat",
		},
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

func TestArtifactory_Plugin_Validate_NoSetProp(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "set-prop",
			APIKey:   "superSecretAPIKey",
			Password: "superSecretPassword",
			URL:      "https://myarti.com/artifactory",
			Username: "octocat",
		},
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

func TestArtifactory_Plugin_Validate_NoUpload(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
			Action:   "upload",
			APIKey:   "superSecretAPIKey",
			Password: "superSecretPassword",
			URL:      "https://myarti.com/artifactory",
			Username: "octocat",
		},
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
