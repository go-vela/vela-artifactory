// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"

	"github.com/sirupsen/logrus"

	json "github.com/ghodss/yaml"
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
	logrus.Tracef("creating prop string for %s", p.Name)

	// check if property value is provided
	if len(p.Value) > 0 {
		return fmt.Sprintf("%s=%s", p.Name, p.Value)
	}

	return fmt.Sprintf("%s=%s", p.Name, strings.Join(p.Values, ","))
}

// Validate verifies the Prop is properly configured.
func (p *Prop) Validate() error {
	logrus.Tracef("validating prop configuration for %s", p.Name)

	// verify name is provided
	if len(p.Name) == 0 {
		return fmt.Errorf("no prop name provided")
	}

	// verify value or values are provided
	if len(p.Value) == 0 && len(p.Values) == 0 {
		return fmt.Errorf("no prop value or values provided")
	}

	return nil
}

// SetProp represents the plugin configuration for setting property information.
type SetProp struct {
	// Path is the target path to artifact(s) to set properties
	Path string
	// Props are properties to set on the artifact(s)
	Props []*Prop
	// RawProps is raw input of properties provided for plugin
	RawProps string
	// Recursive is a flag that enables setting properties on sub-directories for the artifact(s) in the path
	Recursive bool
}

// Exec formats and runs the commands for setting properties on artifacts in Artifactory.
func (s *SetProp) Exec(cli artifactory.ArtifactoryServicesManager) error {
	logrus.Trace("running set-prop with provided configuration")

	// create new search parameters
	searchParams := services.NewSearchParams()

	// add search configuration to search parameters
	searchParams.CommonParams = &utils.CommonParams{
		Pattern:   s.Path,
		Recursive: s.Recursive,
	}

	retries := 3
	var retryErr error
	var files *content.ContentReader

	// send API call to search for files in Artifactory
	// retry a couple times to avoid intermittent authentication failures from Artifactory
	for i := 0; i < retries; i++ {
		backoff := i*2 + 1

		// send API call to search path for artifacts in Artifactory
		files, retryErr = cli.SearchFiles(searchParams)
		if retryErr != nil {
			logrus.Errorf("Error searching for files using pattern %s on attempt %d: %s",
				searchParams.CommonParams.Pattern, i+1, retryErr.Error())

			if i == retries-1 {
				return retryErr
			}

			logrus.Infof("Retrying in %d seconds...", backoff)

			time.Sleep(time.Duration(backoff) * time.Second)
		} else {
			break
		}
	}

	// create new property parameters
	p := services.NewPropsParams()

	// add property configuration to property parameters
	p.Reader = files
	p.Props = s.String()

	// send API call to upload files in Artifactory
	// retry a couple times to avoid intermittent authentication failures from Artifactory
	for i := 0; i < retries; i++ {
		backoff := i*2 + 1

		// send API call to set properties for artifacts in Artifactory
		_, retryErr = cli.SetProps(p)
		if retryErr != nil {
			logrus.Errorf("Error setting properties for files on attempt %d: %s",
				i+1, retryErr.Error())

			if i == retries-1 {
				return retryErr
			}

			logrus.Infof("Retrying in %d seconds...", backoff)

			time.Sleep(time.Duration(backoff) * time.Second)
		} else {
			break
		}
	}

	return nil
}

// String formats and returns a query string for the properties.
func (s *SetProp) String() string {
	logrus.Trace("creating string for props")

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
	logrus.Trace("unmarshaling raw props")

	// cast raw properties into bytes
	bytes := []byte(s.RawProps)

	// serialize raw properties into expected Props type
	err := json.Unmarshal(bytes, &s.Props)
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the SetProp is properly configured.
func (s *SetProp) Validate() error {
	logrus.Trace("validating set prop plugin configuration")

	// verify path is provided
	if len(s.Path) == 0 {
		return fmt.Errorf("no set-prop path provided")
	}

	// serialize provided properties into expected type
	err := s.Unmarshal()
	if err != nil {
		return fmt.Errorf("unable to unmarshal set-prop props: %w", err)
	}

	// verify properties are provided
	if len(s.Props) == 0 {
		return fmt.Errorf("no set-prop props provided")
	}

	// iterate through all properties
	for _, prop := range s.Props {
		// verify the property is valid
		err := prop.Validate()
		if err != nil {
			return fmt.Errorf("invalid set-prop prop provided: %w", err)
		}
	}

	return nil
}
