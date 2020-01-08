// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import "testing"

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

func TestArtifactory_SetProp_Validate(t *testing.T) {
	// setup types
	s := &SetProp{
		Path: "foo/bar",
		Props: []*Prop{
			{
				Name:  "foo",
				Value: "bar",
			},
		},
	}

	err := s.Validate()
	if err != nil {
		t.Errorf("Validate returned err: %v", err)
	}
}

func TestArtifactory_SetProp_Validate_InvalidProp(t *testing.T) {
	// setup types
	s := &SetProp{
		Path: "foo/bar",
		Props: []*Prop{
			{
				Name: "foo",
			},
		},
	}

	err := s.Validate()
	if err == nil {
		t.Errorf("Validate should have returned err")
	}
}

func TestArtifactory_SetProp_Validate_NoPath(t *testing.T) {
	// setup types
	s := &SetProp{
		Props: []*Prop{
			{
				Name:  "foo",
				Value: "bar",
			},
		},
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
