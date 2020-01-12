// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"testing"
)

func TestArtifactory_Plugin_Exec(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{
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

	err := p.Exec()
	if err != nil {
		t.Errorf("Exec returned err: %v", err)
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

func TestArtifactory_Plugin_Validate_NoConfig(t *testing.T) {
	// setup types
	p := &Plugin{
		Config: &Config{},
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
		Copy: &Copy{},
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
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Plugin_Validate_NoDelete(t *testing.T) {
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
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Plugin_Validate_NoSetProp(t *testing.T) {
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

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Plugin_Validate_NoUpload(t *testing.T) {
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
		Upload: &Upload{},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}
