// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/jfrog/jfrog-client-go/artifactory"
	"github.com/jfrog/jfrog-client-go/artifactory/services"
	"github.com/jfrog/jfrog-client-go/artifactory/services/utils"
	"github.com/jfrog/jfrog-client-go/utils/io/content"
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

	retries := 3
	var retryErr error
	var paths *content.ContentReader

	// send API call to search for paths in Artifactory
	// retry a couple times to avoid intermittent authentication failures from Artifactory
	for i := 0; i < retries; i++ {
		backoff := i*2 + 1

		// send API call to capture paths to artifacts in Artifactory
		paths, retryErr = cli.GetPathsToDelete(p)
		if retryErr != nil {
			logrus.Errorf("Error getting delete paths for files using pattern %s on attempt %d: %s",
				p.CommonParams.Pattern,
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

	// send API call to delete artifacts in Artifactory
	// retry a couple times to avoid intermittent authentication failures from Artifactory
	for i := 0; i < retries; i++ {
		backoff := i*2 + 1

		// send API call to delete artifacts in Artifactory
		_, retryErr = cli.DeleteFiles(paths)
		if retryErr != nil {
			logrus.Errorf("Error deleting paths at %s on attempt %d: %s",
				p.CommonParams.Pattern, i+1, retryErr.Error())

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

// Validate verifies the Delete is properly configured.
func (d *Delete) Validate() error {
	logrus.Trace("validating delete plugin configuration")

	// verify path is provided
	if len(d.Path) == 0 {
		return fmt.Errorf("no delete path provided")
	}

	return nil
}
