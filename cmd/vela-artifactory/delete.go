// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package main

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
)

const deleteAction = "delete"

// Delete represents the plugin configuration for delete information.
type Delete struct {
	// file path to load arguments from
	ArgsFile string
	// enables pretending to remove the artifact(s) in the path
	DryRun bool
	// enables removing sub-directories for the artifact(s) in the path
	Recursive bool
	// target path to artifact(s) to remove
	Path string
}

// Exec formats and runs the commands for removing artifacts in Artifactory.
func (d *Delete) Exec(cli *artifactory.ArtifactoryServicesManager) error {
	logrus.Trace("running delete with provided configuration")

	// create new delete parameters
	p := services.NewDeleteParams()

	// add delete configuration to delete parameters
	p.ArtifactoryCommonParams = &utils.ArtifactoryCommonParams{
		Pattern:   d.Path,
		Recursive: d.Recursive,
	}

	// send API call to capture paths to artifacts in Artifactory
	paths, err := cli.GetPathsToDelete(p)
	if err != nil {
		return err
	}

	// send API call to delete artifacts in Artifactory
	_, err = cli.DeleteFiles(paths)
	if err != nil {
		return err
	}

	return nil
}

// Validate verifies the Delete is properly configured.
func (d *Delete) Validate() error {
	logrus.Trace("validating delete plugin configuration")

	if len(d.Path) == 0 {
		return fmt.Errorf("no delete path provided")
	}

	return nil
}
