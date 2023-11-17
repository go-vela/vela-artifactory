// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"

	"github.com/sirupsen/logrus"
)

const copyAction = "copy"

// Copy represents the plugin configuration for copy information.
type Copy struct {
	// Flat is a flag that enables removing source file directory hierarchy
	Flat bool
	// Recursive is a flag that enables copying sub-directories from source
	Recursive bool
	// Path is the source path to artifact(s) to copy
	Path string
	// Target is the path to copy artifact(s) to
	Target string
}

// Exec formats and runs the commands for copying artifacts in Artifactory.
func (c *Copy) Exec(cli artifactory.ArtifactoryServicesManager) error {
	logrus.Trace("running copy with provided configuration")

	// create new copy parameters
	p := services.NewMoveCopyParams()

	// add copy configuration to copy parameters
	p.CommonParams = &utils.CommonParams{
		Pattern:   c.Path,
		Recursive: c.Recursive,
		Target:    c.Target,
	}
	p.Flat = c.Flat

	// send API call to copy artifacts in Artifactory
	_, _, err := cli.Copy(p)
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the Copy is properly configured.
func (c *Copy) Validate() error {
	logrus.Trace("validating copy plugin configuration")

	// verify path is provided
	if len(c.Path) == 0 {
		return fmt.Errorf("no copy path provided")
	}

	// verify target is provided
	if len(c.Target) == 0 {
		return fmt.Errorf("no copy target provided")
	}

	return nil
}
