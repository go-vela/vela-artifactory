// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

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
	// enables removing source file directory hierarchy
	Flat bool
	// enables copying sub-directories from source
	Recursive bool
	// source path to artifact(s) to copy
	Path string
	// target path to copy artifact(s) to
	Target string
}

// Exec formats and runs the commands for copying artifacts in Artifactory.
func (c *Copy) Exec(cli *artifactory.ArtifactoryServicesManager) error {
	logrus.Trace("running copy with provided configuration")

	// create new copy parameters
	p := services.NewMoveCopyParams()

	// add copy configuration to copy parameters
	p.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{
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

	if len(c.Path) == 0 {
		return fmt.Errorf("no copy path provided")
	}

	if len(c.Target) == 0 {
		return fmt.Errorf("no copy target provided")
	}

	return nil
}
