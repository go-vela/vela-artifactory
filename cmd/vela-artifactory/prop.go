// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"
	"strings"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"

	"github.com/sirupsen/logrus"

	json "github.com/ghodss/yaml"
	yaml "gopkg.in/yaml.v2"
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

// String formats and returns a query string for the property.
func (p *Prop) String() string {
	// check if property value is provided
	if len(p.Value) > 0 {
		return fmt.Sprintf("%s=%s", p.Name, p.Value)
	}

	return fmt.Sprintf("%s=%s", p.Name, strings.Join(p.Values, ","))
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
	// raw input of properties provided for plugin
	RawProps interface{}
	// enables setting properties on sub-directories for the artifact(s) in the path
	Recursive bool
}

// Exec formats and runs the commands for setting properties on artifacts in Artifactory.
func (s *SetProp) Exec(cli *artifactory.ArtifactoryServicesManager) error {
	logrus.Trace("running set-prop with provided configuration")

	// create new search parameters
	searchParams := services.NewSearchParams()

	// add search configuration to search parameters
	searchParams.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{
		Pattern:   s.Path,
		Recursive: s.Recursive,
	}

	// send API call to search path for artifacts in Artifactory
	files, err := cli.SearchFiles(searchParams)
	if err != nil {
		return err
	}

	// create new property parameters
	p := services.NewPropsParams()

	// add property configuration to property parameters
	p.Items = files
	p.Props = s.String()

	// send API call to set properties for artifacts in Artifactory
	_, err = cli.SetProps(p)
	if err != nil {
		return err
	}

	return nil
}

// String formats and returns a query string for the properties.
func (s *SetProp) String() string {
	// variable to store properties
	var properties []string

	// iterate through all properties
	for _, prop := range s.Props {
		// add string for property from provided properties
		properties = append(properties, prop.String())
	}

	return strings.Join(properties, ";")
}

// Unmarshal captures the provided properties and
// serializes them into their expected form.
func (s *SetProp) Unmarshal() error {
	// capture raw properties provided to plugin
	raw, err := yaml.Marshal(s.RawProps)
	if err != nil {
		return err
	}

	// convert YAML to JSON for compatiblity purposes
	bytes, err := json.YAMLToJSON(raw)
	if err != nil {
		return err
	}

	// serialize raw properties into expected Props type
	err = json.Unmarshal(bytes, &s.Props)
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the SetProp is properly configured.
func (s *SetProp) Validate() error {
	logrus.Trace("validating set prop plugin configuration")

	if len(s.Path) == 0 {
		return fmt.Errorf("no set-prop path provided")
	}

	// serialize provided properties into expected type
	err := s.Unmarshal()
	if err != nil {
		return fmt.Errorf("unable to unmarshal set-prop props: %v", err)
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
