// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/sirupsen/logrus"
)

const deleteAction = "delete"

// Delete represents the plugin configuration for delete information.
type Delete struct {
	// Recursive is a flag that enables removing sub-directories for the artifact(s) in the path
	Recursive bool
	// Path is the target path to artifact(s) to remove
	Path string
}

// Exec formats and runs the commands for removing artifacts in Artifactory.
func (d *Delete) Exec(cli artifactory.ArtifactoryServicesManager) error {
	logrus.Trace("running delete with provided configuration")

	// create new delete parameters
	p := services.NewDeleteParams()

	// add delete configuration to delete parameters
	p.CommonParams = &utils.CommonParams{
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

	// verify path is provided
	if len(d.Path) == 0 {
		return fmt.Errorf("no delete path provided")
	}

	return nil
}
