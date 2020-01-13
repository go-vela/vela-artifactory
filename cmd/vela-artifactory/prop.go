// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

const setPropAction = "set-prop"

// Prop represents the plugin configuration for setting a property.
type Prop struct {
	// name of the property to set on the artifact(s)
	Name string
	// value of the property to set on the artifact(s)
	Value string
	// values of the property to set on the artifact(s)
	Values []string
}

// Validate verifies the Prop is properly configured.
func (p *Prop) Validate() error {
	logrus.Trace("validating prop configuration")

	if len(p.Name) == 0 {
		return fmt.Errorf("no prop name provided")
	}

	if len(p.Value) == 0 && len(p.Values) == 0 {
		return fmt.Errorf("no prop value(s) provided")
	}

	return nil
}

// SetProp represents the plugin configuration for setting property information.
type SetProp struct {
	// target path to artifact(s) to set properties
	Path string
	// properties to set on the artifact(s)
	Props []*Prop
}

// Validate verifies the SetProp is properly configured.
func (s *SetProp) Validate() error {
	logrus.Trace("validating set prop plugin configuration")

	if len(s.Path) == 0 {
		return fmt.Errorf("no set-prop path provided")
	}

	if len(s.Props) == 0 {
		return fmt.Errorf("no set-prop props provided")
	}

	for _, prop := range s.Props {
		err := prop.Validate()
		if err != nil {
			return fmt.Errorf("invalid set-prop prop provided: %v", err)
		}
	}

	return nil
}
