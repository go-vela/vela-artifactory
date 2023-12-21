// SPDX-License-Identifier: Apache-2.0

package main

import (
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/go-vela/vela-artifactory/cmd/vela-artifactory/mock"
)

func TestArtifactory_Prop_Exec(t *testing.T) {
	// setup types
	s := httptest.NewServer(mock.Handlers())

	p := &Plugin{
		Config: &Config{
			Action:                 "set-prop",
			Token:                  mock.Token,
			APIKey:                 mock.APIKey,
			DryRun:                 false,
			URL:                    s.URL,
			Username:               mock.Username,
			Password:               mock.Password,
			HTTPRetries:            3,
			HTTPRetryWaitMilliSecs: 1,
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
			RawProps: `[{"name": "single", "value": "foo"}]`,
		},
		Upload: &Upload{},
	}

	err := p.Exec()
	if err != nil {
		t.Errorf("Exec returned err %v", err)
	}
}

func TestArtifactory_Prop_String_Value(t *testing.T) {
	// setup types
	p := &Prop{
		Name:   "foo",
		Value:  "bar",
		Values: []string{},
	}

	want := "foo=bar"

	got := p.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

func TestArtifactory_Prop_String_Values(t *testing.T) {
	// setup types
	p := &Prop{
		Name:   "foo",
		Value:  "",
		Values: []string{"baz"},
	}

	want := "foo=baz"

	got := p.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

func TestArtifactory_Prop_Validate(t *testing.T) {
	// setup types
	p := &Prop{
		Name:   "foo",
		Value:  "bar",
		Values: []string{"baz"},
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Prop_Validate_NoName(t *testing.T) {
	// setup types
	p := &Prop{
		Value:  "bar",
		Values: []string{"baz"},
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_Prop_Validate_NoValue(t *testing.T) {
	// setup types
	p := &Prop{
		Name:   "foo",
		Values: []string{"baz"},
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Prop_Validate_NoValues(t *testing.T) {
	// setup types
	p := &Prop{
		Name:   "foo",
		Values: []string{"baz"},
	}

	err := p.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_Prop_Validate_NoValueOrValues(t *testing.T) {
	// setup types
	p := &Prop{
		Name: "foo",
	}

	err := p.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_SetProp_Exec_Error(t *testing.T) {
	// setup types
	config := &Config{
		Action:                 "set-prop",
		APIKey:                 mock.APIKey,
		DryRun:                 false,
		URL:                    mock.InvalidArtifactoryServerURL,
		Username:               mock.Username,
		Password:               mock.Password,
		HTTPRetries:            3,
		HTTPRetryWaitMilliSecs: 1,
	}

	cli, err := config.New()
	if err != nil {
		t.Errorf("Unable to create Artifactory client: %v", err)
	}

	s := &SetProp{
		Path: "foo/bar",
		Props: []*Prop{
			{
				Name:  "single",
				Value: "foo",
			},
		},
		RawProps: `[{"name": "single", "value": "foo"}]`,
	}

	err = s.Exec(*cli)
	if err == nil {
		t.Errorf("Exec should have returned err")
	}
}

func TestArtifactory_SetProp_String(t *testing.T) {
	// setup types
	s := &SetProp{
		Path: "foo/bar",
		Props: []*Prop{
			{
				Name:  "single",
				Value: "foo",
			},
		},
		RawProps: `[{"name": "single", "value": "foo"}]`,
	}

	want := "single=foo"

	got := s.String()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("String is %v, want %v", got, want)
	}
}

func TestArtifactory_SetProp_Validate(t *testing.T) {
	// setup types
	s := &SetProp{
		Path:     "foo/bar",
		RawProps: `[{"name": "single", "value": "foo"}]`,
	}

	err := s.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_SetProp_Validate_Invalid(t *testing.T) {
	// setup types
	s := &SetProp{
		Path:     "foo/bar",
		RawProps: "!@#$%^&*()",
	}

	err := s.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_SetProp_Validate_NoPath(t *testing.T) {
	// setup types
	s := &SetProp{
		RawProps: `[{"name": "single", "value": "foo"}]`,
	}

	err := s.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_SetProp_Validate_NoProps(t *testing.T) {
	// setup types
	s := &SetProp{
		Path: "foo/bar",
	}

	err := s.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_SetProp_Validate_NoPropValue(t *testing.T) {
	// setup types
	s := &SetProp{
		Path:     "foo/bar",
		RawProps: `[{"name": "single"}]`,
	}

	err := s.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_SetProp_Unmarshal(t *testing.T) {
	// setup types
	s := &SetProp{
		Path: "foo/bar",
		RawProps: `
[
  {"name": "single", "value": "foo"},
  {"name": "multiple", "values": ["bar", "baz"]}
]
`}

	want := []*Prop{
		{
			Name:  "single",
			Value: "foo",
		},
		{
			Name:   "multiple",
			Values: []string{"bar", "baz"},
		},
	}

	err := s.Unmarshal()
	if err != nil {
		t.Errorf("Unmarshal returned err: %v", err)
	}

	if !reflect.DeepEqual(s.Props, want) {
		t.Errorf("Unmarshal is %v, want %v", s.Props, want)
	}
}

func TestArtifactory_SetProp_Unmarshal_FailUnmarshal(t *testing.T) {
	// setup types
	s := &SetProp{
		Path:     "foo/bar",
		RawProps: "!@#$%^&*()",
	}

	err := s.Unmarshal()
	if err == nil {
		t.Errorf("Unmarshal should have returned err")
	}
}
